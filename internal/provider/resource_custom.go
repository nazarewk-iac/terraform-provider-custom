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
				Optional: true,
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
				Computed: true,
			},
			"output": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_sensitive": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCustomCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	p := Program(ctx, data, config)
	p.name = "create"
	diags = append(diags, p.openDir()...)
	defer func() { diags = append(diags, p.closeDir()...) }()
	if diags.HasError() {
		return
	}
	if _, ok := data.GetOk("program_create"); ok {
		p.name += "create>create"
		diags = append(diags, p.executeCommand("program_create")...)
		if diags.HasError() {
			return
		}
	} else {
		p.name += "create>update"
		diags = append(diags, p.executeCommand("program_update")...)
	}
	p.name = "create"

	diags = append(diags, p.storeId()...)
	if diags.HasError() {
		return
	}

	diags = append(diags, p.storeAttributes("state", "output", "output_sensitive")...)
	if diags.HasError() {
		return
	}
	return
}

func resourceCustomRead(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	p := Program(ctx, data, config)
	p.name = "read"
	diags = append(diags, p.openDir()...)
	defer func() { diags = append(diags, p.closeDir()...) }()
	if diags.HasError() {
		return
	}

	diags = append(diags, p.executeCommand("program_read")...)
	if diags.HasError() {
		return
	}

	diags = append(diags, p.storeId()...)
	if diags.HasError() {
		return
	}

	diags = append(diags, p.storeAttributes("state", "output", "output_sensitive")...)
	if diags.HasError() {
		return
	}
	return
}

func resourceCustomUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	p := Program(ctx, data, config)
	name := "update"
	p.name = name
	diags = append(diags, p.openDir()...)
	defer func() { diags = append(diags, p.closeDir()...) }()
	if diags.HasError() {
		return
	}
	diags = append(diags, p.executeCommand("program_update")...)
	if diags.HasError() {
		return
	}

	diags = append(diags, p.storeId()...)
	if diags.HasError() {
		return
	}

	diags = append(diags, p.storeAttributes("state", "output", "output_sensitive")...)
	if diags.HasError() {
		return
	}

	return
}

func resourceCustomDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	diags = append(diags, runProgram(ctx, data, config, "delete", "program_delete")...)

	data.SetId("")
	return
}
