package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceExtended(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"extended": func() (*schema.Provider, error) {
				return New(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceExtended,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"extended_resource.foo", "state", regexp.MustCompile("^qwe")),
				),
			},
		},
	})
}

const testAccResourceExtended = `
resource "extended_resource" "foo" {
  input = "qwe"
  program_create = ["sh", "-c", "cat \"$EXT_FILE_input\" > \"$EXT_FILE_state\" && sha256sum $EXT_FILE_input > $EXT_FILE_id"]
  program_read = ["sh", "-c", "cat $EXT_FILE_state"]
  program_update = ["sh", "-c", "cat \"$EXT_FILE_input\" > \"$EXT_FILE_state\" && sha256sum $EXT_FILE_input > $EXT_FILE_id"]
  program_delete = ["sh", "-c", "echo -n > $EXT_FILE_state"]
}
`
