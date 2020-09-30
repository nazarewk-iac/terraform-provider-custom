package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCustom(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCustom,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"custom_resource.foo", "state", regexp.MustCompile("^.*\nqwe$")),
				),
			},
		},
	})
}

const testAccResourceCustom = `locals {
  create_update = <<EOT
	file_name="$(cat "$EXT_FILE_input" | tee "$EXT_FILE_id")"
	echo "$file_name" > "$EXT_FILE_state"
	cat "$EXT_FILE_input_sensitive" | tee -a "$EXT_FILE_state" > "$file_name"
  EOT
}

resource "custom_resource" "foo" {
  input = "/tmp/terraform-provider-custom_resource_test"
  input_sensitive = "qwe"
  program_create = [
    "sh",
    "-xeuo",
    "pipefail",
    "-c",
    local.create_update,
  ]
  program_read = [
    "sh",
    "-xeuo",
    "pipefail",
    "-c",
    <<EOT
	file_name="$(cat "$EXT_FILE_input")"
	cat "$file_name"
	cat "$EXT_FILE_state"
	cat "$file_name" > "$EXT_FILE_state"
  	EOT
  ]
  program_update = [
    "sh",
    "-xeuo",
    "pipefail",
    "-c",
    local.create_update,
  ]
  program_delete = [
    "sh",
    "-xeuo",
    "pipefail",
    "-c",
    <<EOT
	rm "$(cat "$EXT_FILE_input")"
  	EOT
  ]
}
`
