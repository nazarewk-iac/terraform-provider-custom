package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"input": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Input that all scripts using this provider will take.",
			},
			"input_sensitive": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "",
				Description: "Sensitive input that all scripts using this provider will take.",
			},
		},
		ConfigureContextFunc: providerConfigure,
		DataSourcesMap:       map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"custom_resource": resourceCustom(),
		},
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (meta interface{}, diag diag.Diagnostics) {
	meta = &Config{
		Input:          data.Get("input").(string),
		InputSensitive: data.Get("input_sensitive").(string),
	}
	return
}
