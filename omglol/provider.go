package omglol

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type auth struct {
	username string
	apiKey   string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OMGLOL_USERNAME", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OMGLOL_APIKEY", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"omglol_dns": dataSourceDns(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"omglol_dns": resourceDns(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	username := d.Get("username").(string)
	apiKey := d.Get("api_key").(string)

	return auth{username: username, apiKey: apiKey}, diags
}
