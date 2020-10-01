# terraform-provider-custom

Escape hatch for defining Terraform 0.13+ resources with arbitrary commands in as unopinionated and universal way as possible.

Distributed through Terraform Registry at [nazarewk-iac/custom](https://registry.terraform.io/providers/nazarewk-iac/custom/latest)

See the Terraform Registry's [Documentation tab](https://registry.terraform.io/providers/nazarewk-iac/custom/latest/docs)

Source code available at [github.com/nazarewk-iac/terraform-provider-custom](https://github.com/nazarewk-iac/terraform-provider-custom)

# Provider design

see [Provider design](https://registry.terraform.io/providers/nazarewk-iac/custom/latest/docs#provider-design)

## Program interface and guidelines

see [Program interface and guidelines](https://registry.terraform.io/providers/nazarewk-iac/custom/latest/docs#program-interface-and-guidelines)

# Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x - binaries required in `PATH` for both development (tests) and usage
-	[Go](https://golang.org/doc/install) >= 1.15 - for development

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
