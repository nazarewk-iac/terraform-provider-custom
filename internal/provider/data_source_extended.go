package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceExtended() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceExtendedRead,

		Schema: map[string]*schema.Schema{
			"sample_attribute": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceExtendedRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
