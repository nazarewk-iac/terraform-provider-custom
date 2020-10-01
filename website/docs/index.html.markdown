---
layout: "custom"
page_title: "Provider: Custom"
sidebar_current: "docs-custom-index"
description: |-
  Escape hatch for defining Terraform 0.13+ resources with arbitrary commands in as unopinionated and universal way as possible.
---

# Custom provider

Escape hatch for defining Terraform 0.13+ resources with arbitrary commands in as unopinionated and universal way as possible.

Distributed through Terraform Registry at [nazarewk-iac/custom](https://registry.terraform.io/providers/nazarewk-iac/custom/latest)

See the Terraform Registry's [Documentation tab](https://registry.terraform.io/providers/nazarewk-iac/custom/latest/docs)

Source code available at [github.com/nazarewk-iac/terraform-provider-custom](https://github.com/nazarewk-iac/terraform-provider-custom)

## Provider design

The provider's primary goal is to be *quickly finished* as a feature-complete MVP (Minimal Viable Product) enabling Terraform developers to define custom logic.
Most of the efforts will go to thoroughly testing the code.

Below design decisions are supposed to help with above:

1. don't introduce new nice to have features:
    - less code to maintain in the provider,
    - most of those can be provided by wrapping resources in modules using new Terraform 0.12/0.13+ features

1. don't impose any code structure on the scripts:
    - just run the user-provided arguments list as-is
 
1. plain-text (`string` type) resource attributes:
    - use Terraform/HCL features to handle them as structured data, eg: `jsonencode()` and `jsondecode()` pair,
    - provider users are free to handle the data as they see fit,
    - only plumbing required to share the data with the `program`,

1. interface with the `program` through files in a temporary directory:
    - temporary directory path is exposed to the program through `TF_CUSTOM_DIR` environment variable,
    - filesystem permissions reflect what can be done with them,

1. attribute names map directly to file names:
    - so we don't need to pass anything other than temporary directory location to the `program`,

1. only one way to pass the data down to `program`:
    - through string attributes mapped to files,
    - environment variables are NOT configurable, if you really need them you can `source ${TF_CUSTOM_DIR}/input_sensitive` in shell program,

1. well defined `program` interface files

### Program interface and guidelines

1. Program receives temporary directory path in `${TF_CUSTOM_DIR}` environment variable.

1. Program interacts with files named after resource attributes and `id`, `old_state` and `provider_input` files.

    1. Program reads data from `provider_input`/`input`/`input_sensitive`:
        1. `provider_input` is the `input` attribute of `provider` configuration block,
    
    1. Program must fill-in `id` during create and update representing current and future values of `custom_resource.*.id`:
    
        1. Program empties `id` (eg. `echo -n > "${TF_CUSTOM_DIR}/id"`) during read when the resource disappeared (externally).
    
    1. Program stores the managed data in `state`:
    
        1. Program can read previous version of managed data from `old_state`,
        2. Program should write to `state` during read,
    
    1. Program writes to `output`/`output_sensitive` to expose additional data.


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
