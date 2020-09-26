---
layout: "extended"
page_title: "Extended: extended_resource"
sidebar_current: "docs-extended-resource"
description: |-
  Sample resource in the Terraform provider extended.
---

# extended_resource

Sample resource in the Terraform provider extended.

## Example Usage

```hcl
resource "extended_resource" "example" {
  sample_attribute = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `sample_attribute` - Sample attribute.

