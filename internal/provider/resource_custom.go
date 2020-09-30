package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomCreate,
		ReadContext:   resourceCustomRead,
		UpdateContext: resourceCustomUpdate,
		DeleteContext: resourceCustomDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"program_create": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				MinItems: 1,
			},
			"program_read": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				MinItems: 1,
			},
			"program_update": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				MinItems: 1,
			},
			"program_delete": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				MinItems: 1,
			},
			"input": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"input_sensitive": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceCustomCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	diags = append(diags, runProgram(ctx, data, config, "program_create")...)
	return
}

func resourceCustomRead(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	diags = append(diags, runProgram(ctx, data, config, "program_read")...)
	return
}

func resourceCustomUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	diags = append(diags, runProgram(ctx, data, config, "program_update")...)
	return
}

func resourceCustomDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	diags = append(diags, runProgram(ctx, data, config, "program_delete")...)
	return
}
