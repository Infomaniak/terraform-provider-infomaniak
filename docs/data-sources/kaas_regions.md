---
page_title: "infomaniak_kaas_regions"
subcategory: "KaaS"
description: |-
  The KaaS Regions Data Source allows the user to list available regions for KaaS.
---

# infomaniak_kaas_regions

The KaaS Regions Data Source allows the user to retrieve a list of available regions where Kubernetes clusters can be deployed.

## Example Usage

```hcl
data "infomaniak_kaas_regions" "kaas_regions" {}
```

## Schema
### Read-Only
- `items` (List of String) A list of region identifiers where KaaS is available.
