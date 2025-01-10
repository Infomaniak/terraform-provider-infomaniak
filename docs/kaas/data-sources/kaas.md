---
page_title: "infomaniak_kaas"
subcategory: ""
description: |-
  The Kaas Data Source allows the user to read information about a Kaas project
---

# infomaniak_kaas (Data Source)

The Kaas Data Source allows the user to read information about a Kaas project.


## Example

```hcl
data "infomaniak_kaas" "kluster" {
    pcp_id = "xxxxx"
    id     = "yyyyy"
}
```

## Schema

### Required

- `id` (String) The id of the KaaS project.
- `pcp_id` (String) The id of the Public Cloud project where KaaS is installed.

### Read-Only

- `kubeconfig` (String, Sensitive) The Kubeconfig to access the Kluster.
- `region` (String) Region where the instance live.
