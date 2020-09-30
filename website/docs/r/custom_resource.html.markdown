---
layout: "custom"
page_title: "Custom: custom_resource"
sidebar_current: "docs-custom-resource"
description: |-
  Sample resource in the Terraform provider custom.
---

# custom_resource

Sample resource in the Terraform provider custom.

Following files will be created in directory defined by `${TF_CUSTOM_DIR}`:
- `provider_input` - (read-only) - inputs passed down from `provider` configuration
- `input` - (read-only) - inputs passed down from the resource's attribute, see arguments reference
- `input_sensitive` - (read-only) - inputs passed down from the resource's attribute, see arguments reference
- `id` - (read-write) - holds the current (and future if changed) value of `id` attribute, resource will be considered non-existing if empty,
- `state` - (read-write) - reflects current (and future if changed) managed state of the world, see arguments reference
- `old_state` - (read-only) - reflects previous state of the world,
- `output` - (write-only) - holds current unmanaged state of the world, see arguments reference
- `output_sensitive` - (write-only) - holds current unmanaged state of the world, see arguments reference

## Example Usage

```hcl
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

* `program_create` - (Optional) program to run on `create` operation, runs `program_update` instead if not provided. Be sure to write to `${TF_CUSTOM_DIR}/id`.
* `program_read` - program to run on `read` operation. It should update `${TF_CUSTOM_DIR}/state` to reflect world's state. Emptying `${TF_CUSTOM_DIR}/id` will inform Terraform that resource does not exist anymore.
* `program_update` - program to run on `update` (and optionally `create`) operations.  It should update `${TF_CUSTOM_DIR}/state` to reflect world's state.
* `program_delete` - program to run on `destroy` operation.
* `state` - (`${TF_CUSTOM_DIR}/state` read-write, `${TF_CUSTOM_DIR}/old_state` read-only) managed parts of resource's real state, it should be written to (during course of `create`, `read` and `update` commands) to reflect current state of the world.
* `input` - (`${TF_CUSTOM_DIR}/input` read-only) unmanaged/to-be-interpolated parts of resource's desired state
* `input_sensitive` - (`${TF_CUSTOM_DIR}/input_sensitive` read-only) same as `input`, but the content won't be printed during planning.
* `output` - (`${TF_CUSTOM_DIR}/output` write-only) additional (relative to `state` attribute) data the resource is providing.
* `output_sensitive` - (`${TF_CUSTOM_DIR}/output_sensitive` write-only) same as `output`, but content won't be printed during planning.
