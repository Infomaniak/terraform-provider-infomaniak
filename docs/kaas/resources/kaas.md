---
page_title: "infomaniak_kaas"
subcategory: ""
description: |-
  The Kaas resource allows the user to manage a Kaas project
---

# infomaniak_kaas

The Kaas resource allows the user to manage a Kaas project.

## Example

```hcl
resource "infomaniak_kaas" "kluster" {
  pcp_id = "xxxxx"

  region = "yyyyy"
}
```

## Schema

### Required

- `pcp_id` (String) The id of the public cloud project where KaaS is installed.
- `region` (String) Region where the instance live.

### Read-Only

- `id` (String) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `kubeconfig` (String, Sensitive) The Kubeconfig to access the Kluster.
