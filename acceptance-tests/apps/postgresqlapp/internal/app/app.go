package app

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"postgresqlapp/internal/credentials"
	"regexp"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	tableName   = "test"
	keyColumn   = "keyname"
	valueColumn = "valuedata"
)

func App(config *credentials.Config) *mux.Router {
	db, err := connect(config)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", aliveness).Methods(http.MethodHead, http.MethodGet)
	r.HandleFunc("/{schema}", handleCreateSchema(db)).Methods(http.MethodPut)
	r.HandleFunc("/{schema}", handleDropSchema(db)).Methods(http.MethodDelete)
	r.HandleFunc("/{schema}/{key}", handleSet(db)).Methods(http.MethodPut)
	r.HandleFunc("/{schema}/{key}", handleGet(db)).Methods(http.MethodGet)

	return r
}

func aliveness(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handled aliveness test.")
	w.WriteHeader(http.StatusNoContent)
}

func connect(config *credentials.Config) (*sql.DB, error) {
	connStr, err := createCon(config)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", connStr)

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}
	db.SetMaxIdleConns(0)

	_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (%s VARCHAR(255) NOT NULL, %s VARCHAR(255) NOT NULL)`, tableName, keyColumn, valueColumn))
	if err != nil {
		return nil, errors.Wrap(err, "error creating table")
	}

	_, err = db.Exec(fmt.Sprintf(`GRANT ALL ON TABLE public.%s TO PUBLIC`, tableName))
	if err != nil {
		return nil, errors.Wrap(err, "error granting table permissions")
	}

	return db, nil
}

func createCon(config *credentials.Config) (string, error) {
	// Create a TLS config with the CA/client key both configured
	parseConfig, err := pgx.ParseConfig(config.URI)
	if err != nil {
		return "", err
	}
	pair, err := tls.X509KeyPair([]byte(config.ClientCACert), []byte(config.ClientPrivateKey))
	if err != nil {
		return "", err
	}
	certPool := x509.NewCertPool()
	//if ok := certPool.AppendCertsFromPEM([]byte(caCert)); !ok {
	//	log.Fatal("Failed to append CA to cert pool")
	//}

	parseConfig.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{pair},
		RootCAs:            certPool,
	}
	return stdlib.RegisterConnConfig(parseConfig), nil
}

func fail(w http.ResponseWriter, code int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	log.Println(msg)
	http.Error(w, msg, code)
}

func schemaName(r *http.Request) (string, error) {
	schema, ok := mux.Vars(r)["schema"]

	switch {
	case !ok:
		return "", fmt.Errorf("schema missing")
	case len(schema) > 50:
		return "", fmt.Errorf("schema name too long")
	case len(schema) == 0:
		return "", fmt.Errorf("schema name cannot be zero length")
	case !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(schema):
		return "", fmt.Errorf("schema name contains invalid characters")
	default:
		return schema, nil
	}
}
