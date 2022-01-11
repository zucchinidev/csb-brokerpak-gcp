package app

import (
	"context"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/bigquery"
	"github.com/gorilla/mux"
)

func handleSet(client *bigquery.Client, projectID, datasetID string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling set.")

		tableName, ok := mux.Vars(r)["tableName"]
		if !ok {
			fail(w, http.StatusBadRequest, "tableName missing.")
			return
		}

		key, ok := mux.Vars(r)["key"]
		if !ok {
			fail(w, http.StatusBadRequest, "key missing.")
			return
		}

		rawValue, err := io.ReadAll(r.Body)
		if err != nil {
			fail(w, http.StatusBadRequest, "Error parsing value: %s", err)
			http.Error(w, "Failed to parse value.", http.StatusBadRequest)
			return
		}

		value := string(rawValue)

		query := client.Query("INSERT INTO " + tableName + " VALUES (" + key + "," + value + ")")
		query.DefaultProjectID = projectID
		query.DefaultDatasetID = datasetID

		job, err := query.Run(context.Background())
		if err != nil {
			fail(w, http.StatusFailedDependency, "error in creating job inserting into table %s: %s", tableName, err)
			return
		}

		_, err = job.Wait(context.Background())
		if err != nil {
			fail(w, http.StatusFailedDependency, "error inserting into table %s: %s", tableName, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Printf("Key %q set to value %q.", key, value)
	}
}
