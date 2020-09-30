---
layout: "custom"
page_title: "Custom: custom_resource"
sidebar_current: "docs-custom-resource"
description: |-
  Sample resource in the Terraform provider custom.
---

# custom_resource

Sample resource in the Terraform provider custom.

## Example Usage

```hcl
resource "custom_resource" "example" {
  sample_attribute = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `sample_attribute` - Sample attribute.

