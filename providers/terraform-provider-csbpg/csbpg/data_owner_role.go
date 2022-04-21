package csbpg

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"sync"
)

var sharedDataOwnerRoleCreateMutex sync.Mutex

func createDataOwnerRole(db *sql.DB, cf connectionFactory) error {
	sharedDataOwnerRoleCreateMutex.Lock()
	defer sharedDataOwnerRoleCreateMutex.Unlock()

	exists, err := roleExists(db, cf.dataOwnerRole)
	if err != nil {
		return err
	}

	if !exists {
		log.Println("[DEBUG] data owner role does not exist - creating")
		_, err = db.Exec(fmt.Sprintf("CREATE ROLE %s WITH NOLOGIN", pq.QuoteIdentifier(cf.dataOwnerRole)))

		if err != nil {
			return err
		}
	}

	log.Println("[DEBUG] granting data owner role")
	_, err = db.Exec(fmt.Sprintf("GRANT CREATE ON DATABASE %s TO %s", pq.QuoteIdentifier(cf.database), pq.QuoteIdentifier(cf.dataOwnerRole)))

	return err
}
