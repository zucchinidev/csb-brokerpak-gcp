package csbpg

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
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
			"shared_role": {
				Type:     schema.TypeString,
				Required: true,
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
	if username == "" {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "invalid 'username'",
		}}
	}
	password := d.Get("password").(string)
	if password == "" {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "invalid 'password'",
		}}
	}
	sharedRole := d.Get("shared_role").(string)
	if sharedRole == "" {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "invalid 'shared_role'",
		}}
	}

	id := fmt.Sprintf("bindinguser/%s", username)

	cf := m.(connectionFactory)

	db, err := cf.Connect()
	if err != nil {
		return diag.FromErr(err)
	}
	defer db.Close()
	log.Println("[DEBUG] connected")

	exists, err := roleExists(ctx, db, sharedRole)
	if err != nil {
		return diag.FromErr(err)
	}

	if !exists {
		log.Println("[DEBUG] role does not exist")
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("shared_role %s does not exist", sharedRole),
		}}
	}

	log.Println("[DEBUG] create user")
	// TODO: can't use $1 because this statement can't be prepared, but using %s looks unsafe
	_, err = db.Exec(fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD '%s' INHERIT IN ROLE %s", username, password, sharedRole))
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
