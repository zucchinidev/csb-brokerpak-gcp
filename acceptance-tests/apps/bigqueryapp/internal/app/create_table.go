package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/gorilla/mux"
)

func handleCreateTable(client *bigquery.Client, datasetID string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling create table.")

		tableName, ok := mux.Vars(r)["tableName"]
		if !ok {
			fail(w, http.StatusBadRequest, "tableName missing.")
			return
		}

		appSchema := bigquery.Schema{
			{Name: "key", Type: bigquery.StringFieldType},
			{Name: "value", Type: bigquery.StringFieldType},
		}

		metaData := &bigquery.TableMetadata{
			Schema:         appSchema,
			ExpirationTime: time.Now().AddDate(0, 1, 0), // Table will be automatically deleted in 1 year.
		}
		tableRef := client.Dataset(datasetID).Table(tableName)
		if err := tableRef.Create(context.Background(), metaData); err != nil {
			fail(w, http.StatusFailedDependency, "Error creating table %s: %s", tableName, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Printf("created table %s", tableName)
	}
}
