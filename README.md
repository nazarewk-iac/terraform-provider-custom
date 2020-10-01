# terraform-provider-custom
Escape hatch for defining Terraform 0.13+ resources with arbitrary commands in as unopinionated and universal way as possible.

Distributed through Terraform Registry at [nazarewk-iac/custom](https://registry.terraform.io/providers/nazarewk-iac/custom/latest)

See the Terraform Registry's [Documentation tab](https://registry.terraform.io/providers/nazarewk-iac/custom/latest/docs/resources/custom_resource)

# Provider design

The provider's primary goal is to be *quickly finished* as a feature-complete MVP (Minimal Viable Product) enabling Terraform developers to define custom logic.
Most of the efforts will go to thoroughly testing the code.

Below design decisions are supposed to help with above:

- don't introduce new nice to have features:
    - less code to maintain in the provider,
    - most of those can be provided by wrapping resources in modules using new Terraform 0.12/0.13+ features

- don't impose any code structure on the scripts:
    - just run the user-provided arguments list as-is
 
- plain-text (`string` type) resource attributes:
    - use Terraform/HCL features to handle them as structured data, eg: `jsonencode()` and `jsondecode()` pair,
    - provider users are free to handle the data as they see fit,
    - only plumbing required to share the data with the `program`,

- interface with the `program` through files in a temporary directory:
    - temporary directory path is exposed to the program through `TF_CUSTOM_DIR` environment variable,
    - filesystem permissions reflect what can be done with them,

- attribute names map directly to file names:
    - so we don't need to pass anything other than temporary directory location to the `program`,

- only one way to pass the data down to `program`:
    - through string attributes mapped to files,
    - environment variables are NOT configurable, if you really need them you can `source ${TF_CUSTOM_DIR}/input_sensitive` in shell program,

- well defined `program` interface files

## Program interface and guidelines

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

# Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

# Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

# Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

# Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
