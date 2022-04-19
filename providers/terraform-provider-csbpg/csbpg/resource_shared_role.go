package csbpg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lib/pq"
	"log"
)

func resourceSharedRole() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		CreateContext: resourceSharedRoleCreate,
		ReadContext:   resourceSharedRoleRead,
		UpdateContext: resourceSharedRoleUpdate,
		DeleteContext: resourceSharedRoleDelete,
		Description:   "TODO",
		UseJSONNumber: true,
	}
}

func resourceSharedRoleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] ENTRY resourceSharedRoleCreate()")
	defer log.Println("[DEBUG] EXIT resourceSharedRoleCreate()")

	name := d.Get("name").(string)
	if name == "" {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "invalid 'name'",
		}}
	}
	id := fmt.Sprintf("sharedrole/%s", name)

	cf := m.(connectionFactory)

	db, err := cf.Connect()
	if err != nil {
		return diag.FromErr(err)
	}
	defer db.Close()
	log.Println("[DEBUG] connected")

	exists, err := roleExists(db, name)
	if err != nil {
		return diag.FromErr(err)
	}

	if !exists {
		log.Println("[DEBUG] creating role")
		// TODO: can't use $1 because this statement can't be prepared, but using %s looks unsafe
		_, err = db.Exec(fmt.Sprintf("CREATE ROLE %s WITH NOLOGIN", pq.QuoteIdentifier(name)))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Println("[DEBUG] granting")
	// TODO: can't use $1 because this statement can't be prepared, but using %s looks unsafe
	_, err = db.Exec(fmt.Sprintf("GRANT CREATE ON DATABASE %s TO %s", pq.QuoteIdentifier(cf.database), pq.QuoteIdentifier(name)))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] setting ID %s\n", id)
	d.SetId(id)
	return nil
}

func resourceSharedRoleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func resourceSharedRoleUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func resourceSharedRoleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func roleExists(db *sql.DB, name string) (bool, error) {
	log.Println("[DEBUG] ENTRY roleExists()")
	defer log.Println("[DEBUG] EXIT roleExists()")

	// TODO: can't use $1 because this statement can't be prepared, but using %s looks unsafe
	rows, err := db.Query(fmt.Sprintf("SELECT FROM pg_catalog.pg_roles WHERE rolname = '%s'", name))
	if err != nil {
		return false, fmt.Errorf("error finding role %s: %w", name, err)
	}
	defer rows.Close()
	return rows.Next(), nil
}
