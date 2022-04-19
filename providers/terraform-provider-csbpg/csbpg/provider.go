package csbpg

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ConfigureContextFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"csbpg_shared_role":  resourceSharedRole(),
			"csbpg_binding_user": resourceBindingUser(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	var diags diag.Diagnostics

	host := d.Get("host").(string)
	if host == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "invalid 'host'",
		})
	}

	port := d.Get("port").(int)
	if port == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "invalid 'port'",
		})
	}

	username := d.Get("username").(string)
	if username == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "invalid 'username'",
		})
	}

	password := d.Get("password").(string)
	if password == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "invalid 'password'",
		})
	}

	database := d.Get("database").(string)
	if password == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "invalid 'password'",
		})
	}

	return connectionFactory{
		host:     host,
		port:     port,
		username: username,
		password: password,
		database: database,
	}, diags
}
