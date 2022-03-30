package app

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	tableName   = "test"
	keyColumn   = "keyname"
	valueColumn = "valuedata"
)

func App(uri string) *mux.Router {
	db := connect(uri)

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

func connect(uri string) *sql.DB {
	rootCertPool := x509.NewCertPool()
	caCert := os.Getenv(`TODO`)

	if ok := rootCertPool.AppendCertsFromPEM([]byte(caCert)); !ok {
		log.Fatal("Failed to append CA to cert pool")
	}

	certs, err := tls.LoadX509KeyPair(os.Getenv(`ca_cert`), os.Getenv(`server_key`))
	if err != nil {
		log.Fatalf("failed to create key pair: %s", err)
	}

	// Create a TLS config with the CA/client key both configured

	db, err := sql.Open("pgx", uri)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	db.SetMaxIdleConns(0)

	_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (%s VARCHAR(255) NOT NULL, %s VARCHAR(255) NOT NULL)`, tableName, keyColumn, valueColumn))
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	}

	_, err = db.Exec(fmt.Sprintf(`GRANT ALL ON TABLE public.%s TO PUBLIC`, tableName))
	if err != nil {
		log.Fatalf("Error granting table permissions: %s", err)
	}

	return db
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
