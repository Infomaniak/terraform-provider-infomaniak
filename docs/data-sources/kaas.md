---
page_title: "infomaniak_kaas"
subcategory: "KaaS"
description: |-
  The Kaas Data Source allows the user to read information about a Kaas project
---

# infomaniak_kaas (Data Source)

The Kaas Data Source allows the user to read information about a Kaas project.

## Example

```hcl
data "infomaniak_kaas" "kluster" {
  public_cloud_id = xxxxx
  public_cloud_project_id = yyyyy
  id     = zzzzz
}
```

## Schema

### Required

- `id` (Integer) The id of the KaaS project.
- `public_cloud_project_id` (Integer) The id of the Public Cloud Project where KaaS is installed.
- `public_cloud_id` (Integer) The id of the Public Cloud where KaaS is installed.

### Read-Only

- `kubeconfig` (String, Sensitive) The Kubeconfig to access the Kluster.
- `region` (String) Region where the instance live.
- `pack_name` (String) The name of the pack corresponding the KaaS project.
- `kubernetes_version` (String) The version of Kubernetes to use.
- `name` (String) The name of the KaaS shown on the manager.
