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
					resource.TestMatchResourceAttr("custom_resource.foo", "state", regexp.MustCompile("^qwe$")),
					resource.TestMatchResourceAttr("custom_resource.foo", "output", regexp.MustCompile("^/tmp/terraform-provider-custom_resource_test$")),
				),
			},
		},
	})
}

const testAccResourceCustom = `locals {
  script = <<EOT
  set -xeuo pipefail

  main() {
	"cmd_$@"
  }

  cmd_update() {
	file_name="$(cat "$EXT_FILE_input" | tee "$EXT_FILE_id" "$EXT_FILE_output")"
	cat "$EXT_FILE_state" | tee "$EXT_FILE_state" > "$file_name"
  }

  cmd_read() {
	file_name="$(cat "$EXT_FILE_input")"
	cat "$file_name"
	cat "$EXT_FILE_state"
	echo -n "$file_name" > "$EXT_FILE_output"
	cat "$file_name" > "$EXT_FILE_state"
  }
  
  cmd_delete() {
	rm "$(cat "$EXT_FILE_input")"
  }

  main "$@"
  EOT

  program = ["sh", "-c", local.script, "command_string"]
}

resource "custom_resource" "foo" {
  input = "/tmp/terraform-provider-custom_resource_test"
  state = "qwe"

  program_create = concat(local.program, ["update"])
  program_read = concat(local.program, ["read"])
  program_update = concat(local.program, ["update"])
  program_delete = concat(local.program, ["delete"])
}
`
