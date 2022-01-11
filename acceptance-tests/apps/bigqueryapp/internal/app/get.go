package app

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"cloud.google.com/go/bigquery"
)

func handleGet(client *bigquery.Client, projectID, datasetID string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling get.")

		tableName, ok := mux.Vars(r)["tableName"]
		if !ok {
			fail(w, http.StatusBadRequest, "tableName missing.")
			return
		}

		key, ok := mux.Vars(r)["key"]
		if !ok {
			fail(w, http.StatusBadRequest, "key missing")
			return
		}

		query := client.Query("SELECT value FROM " + tableName + " WHERE key = '" + "key" + "'")
		query.DefaultProjectID = projectID
		query.DefaultDatasetID = datasetID

		rows, err := query.Read(context.Background())
		if err != nil {
			fail(w, http.StatusFailedDependency, "error querying table %s: %s", tableName, err)
			return
		}

		var result string
		rows.Next(result)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")

		_, err = w.Write([]byte(result))

		if err != nil {
			log.Printf("Error writing value: %s", err)
			return
		}

		log.Printf("Value %q retrived from key %q.", result, key)
	}
}
