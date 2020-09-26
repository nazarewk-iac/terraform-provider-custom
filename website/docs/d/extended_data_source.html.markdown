---
layout: "extended"
page_title: "Extended: extended_data_source"
sidebar_current: "docs-extended-data-source"
description: |-
  Sample data source in the Terraform provider extended.
---

# extended_data_source

Sample data source in the Terraform provider extended.

## Example Usage

```hcl
data "extended_data_source" "example" {
  sample_attribute = "foo"
}
```

## Attributes Reference

* `sample_attribute` - Sample attribute.
