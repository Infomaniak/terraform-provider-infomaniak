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
  public_cloud_id = xxxxx
  public_cloud_project_id = yyyyy
  
  kubernetes_version = "1.31"
  region = "zzzzz"
}
```

## Schema

### Required

- `public_cloud_id` (Integer) The id of the Public Cloud where KaaS is installed.
- `public_cloud_project_id` (String) The id of the public cloud project where KaaS is installed.
- `region` (String) Region where the instance live.

### Read-Only

- `id` (String) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `kubeconfig` (String, Sensitive) The Kubeconfig to access the Kluster.
