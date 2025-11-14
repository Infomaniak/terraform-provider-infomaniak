---
page_title: "infomaniak_kaas_versions"
subcategory: "KaaS"
description: |-
  The KaaS Versions Data Source allows the user to retrieve available versions of Kubernetes.
---

# infomaniak_kaas_versions

The KaaS Versions Data Source allows the user to retrieve available versions of Kubernetes.  
The returned list is ordered such that the first element corresponds to the latest version.

## Example Usage

```hcl
data "infomaniak_kaas_versions" "kaas_versions" {}
```

## Schema
## Read-Only
- `items` (List of String) A list of available Kubernetes versions, sorted in descending order (latest version first).