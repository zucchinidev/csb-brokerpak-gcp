package app

import (
	"bigqueryapp/internal/credentials"
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"

	"github.com/gorilla/mux"
)

func App(creds credentials.BigQueryCredentials) *mux.Router {
	client, _ := bigquery.NewClient(context.Background(), creds.ProjectID, option.WithCredentialsJSON([]byte(creds.Credentials)))

	r := mux.NewRouter()

	r.HandleFunc("/", aliveness).Methods("HEAD", "GET")
	r.HandleFunc("/{tableName}", handleCreateTable(client, creds.DatasetID)).Methods(http.MethodPut)
	r.HandleFunc("/{tableName}", handleDropTable(client)).Methods(http.MethodDelete)
	r.HandleFunc("/{tableName}/{key}", handleSet(client, creds.ProjectID, creds.DatasetID)).Methods("PUT")
	r.HandleFunc("/{tableName}/{key}", handleGet(client, creds.ProjectID, creds.DatasetID)).Methods("GET")

	return r
}

func aliveness(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handled aliveness test.")
	w.WriteHeader(http.StatusNoContent)
}

func fail(w http.ResponseWriter, code int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	log.Println(msg)
	http.Error(w, msg, code)
}
