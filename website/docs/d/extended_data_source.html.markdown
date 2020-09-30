---
layout: "custom"
page_title: "Custom: custom_data_source"
sidebar_current: "docs-custom-data-source"
description: |-
  Sample data source in the Terraform provider custom.
---

# custom_data_source

Sample data source in the Terraform provider custom.

## Example Usage

```hcl
data "custom_data_source" "example" {
  sample_attribute = "foo"
}
```

## Attributes Reference

* `sample_attribute` - Sample attribute.
