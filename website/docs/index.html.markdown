---
layout: "custom"
page_title: "Provider: Custom"
sidebar_current: "docs-custom-index"
description: |-
  Terraform provider custom.
---

# Custom Provider

Use this paragraph to give a high-level overview of your provider, and any configuration it requires.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
provider "custom" {
}

locals {
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
```

## Argument Reference

The following arguments are supported:

* `input` - (`${TF_CUSTOM_DIR}/provider_input`) input to be passed down to all resources, but not stored in the terraform state file
