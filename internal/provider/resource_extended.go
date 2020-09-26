package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExtended() *schema.Resource {
	return &schema.Resource{
		Create: resourceExtendedCreate,
		Read:   resourceExtendedRead,
		Update: resourceExtendedUpdate,
		Delete: resourceExtendedDelete,

		Schema: map[string]*schema.Schema{
			"sample_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceExtendedCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceExtendedRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceExtendedUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceExtendedDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
