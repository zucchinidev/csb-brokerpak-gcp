package csbpg

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lib/pq"
	"log"
	"strings"
)

func resourceBindingUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		CreateContext: resourceBindingUserCreate,
		ReadContext:   resourceBindingUserRead,
		UpdateContext: resourceBindingUserUpdate,
		DeleteContext: resourceBindingUserDelete,
		Description:   "TODO",
		UseJSONNumber: true,
	}
}

func resourceBindingUserCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] ENTRY resourceSharedRoleCreate()")
	defer log.Println("[DEBUG] EXIT resourceSharedRoleCreate()")

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	id := fmt.Sprintf("bindinguser/%s", username)

	cf := m.(connectionFactory)

	db, err := cf.Connect()
	if err != nil {
		return diag.FromErr(err)
	}
	defer db.Close()
	log.Println("[DEBUG] connected")

	exists, err := roleExists(db, cf.dataOwnerRole)
	if err != nil {
		return diag.FromErr(err)
	}

	if !exists {
		log.Println("[DEBUG] data owner role does not exist - creating")
		// TODO: can't use $1 because this statement can't be prepared, but using %s looks unsafe
		_, err = db.Exec(fmt.Sprintf("CREATE ROLE %s WITH NOLOGIN", pq.QuoteIdentifier(cf.dataOwnerRole)))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Println("[DEBUG] granting data owner role")
	// TODO: can't use $1 because this statement can't be prepared, but using %s looks unsafe
	_, err = db.Exec(fmt.Sprintf("GRANT CREATE ON DATABASE %s TO %s", pq.QuoteIdentifier(cf.database), pq.QuoteIdentifier(cf.dataOwnerRole)))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] create binding user")
	// TODO: can't use $1 because this statement can't be prepared, but using %s looks unsafe
	_, err = db.Exec(fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD %s INHERIT IN ROLE %s", pq.QuoteIdentifier(username), safeQuote(password), pq.QuoteIdentifier(cf.dataOwnerRole)))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] setting ID %s\n", id)
	d.SetId(id)

	return nil
}

func resourceBindingUserRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func resourceBindingUserUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func resourceBindingUserDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func safeQuote(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `'`, `\\`))
}
