package provider

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCustom(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: replacedState("stateQWE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("custom_resource.foo", "state", regexp.MustCompile("^stateQWE$")),
					resource.TestMatchResourceAttr("custom_resource.foo", "output", regexp.MustCompile("^/tmp/terraform-provider-custom_resource_test$")),
				),
			},
			{
				Config:             replacedState("stateQWE"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("custom_resource.foo", "state", regexp.MustCompile("^stateQWE$")),
					resource.TestMatchResourceAttr("custom_resource.foo", "output", regexp.MustCompile("^/tmp/terraform-provider-custom_resource_test$")),
				),
			},
			{
				Config:             replacedState("stateASD"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: replacedState("stateASD"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("custom_resource.foo", "state", regexp.MustCompile("^stateASD$")),
					resource.TestMatchResourceAttr("custom_resource.foo", "output", regexp.MustCompile("^/tmp/terraform-provider-custom_resource_test$")),
				),
			},
		},
	})
}

func replacedState(replacement string) string {
	return strings.ReplaceAll(testAccResourceCustom, "STATE_PLACEHOLDER", replacement)
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
  state = "STATE_PLACEHOLDER"

  program_create = concat(local.program, ["update"])
  program_read = concat(local.program, ["read"])
  program_update = concat(local.program, ["update"])
  program_delete = concat(local.program, ["delete"])
}
`
