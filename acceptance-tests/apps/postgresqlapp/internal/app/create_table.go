package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func handleCreateTable(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling create table.")

		schema, err := schemaName(r)
		if err != nil {
			fail(w, http.StatusInternalServerError, "Schema name error: %s", err)
			return
		}

		tableName, err := additionalTableName(r)
		if err != nil {
			fail(w, http.StatusInternalServerError, "Table name error: %s", err)
			return
		}

		_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (%s VARCHAR(255) NOT NULL, %s VARCHAR(255) NOT NULL)`, schema, tableName, keyColumn, valueColumn))
		if err != nil {
			fail(w, http.StatusBadRequest, "Error creating table: %s", err)
			return
		}

		// Maye test if we can alter a table?

		_, err = db.Exec(fmt.Sprintf(`GRANT ALL ON TABLE %s.%s TO PUBLIC`, schema, tableName))
		if err != nil {
			fail(w, http.StatusBadRequest, "Error granting table permissions: %s", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Printf("Table %q.%q created", schema, tableName)
	}
}
