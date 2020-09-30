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
	file_name="$(cat "$TF_CUSTOM_DIR/input" | tee "$TF_CUSTOM_DIR/id" "$TF_CUSTOM_DIR/output")"
	cat "$TF_CUSTOM_DIR/state" | tee "$TF_CUSTOM_DIR/state" > "$file_name"
  }

  cmd_read() {
	file_name="$(cat "$TF_CUSTOM_DIR/input")"
	cat "$file_name"
	cat "$TF_CUSTOM_DIR/state"
	echo -n "$file_name" > "$TF_CUSTOM_DIR/output"
	cat "$file_name" > "$TF_CUSTOM_DIR/state"
  }
  
  cmd_delete() {
	rm "$(cat "$TF_CUSTOM_DIR/input")"
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
