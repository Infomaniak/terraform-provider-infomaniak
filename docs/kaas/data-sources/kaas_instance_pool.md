---
page_title: "infomaniak_kaas_instance_pool"
subcategory: ""
description: |-
  The Kaas Instance Pool Data Source allows the user to read information about a Kaas project
---

# infomaniak_kaas_instance_pool (Data Source)

The Kaas Instance Pool Data Source allows the user to read information about a Kaas project.

## Example

```hcl
data "infomaniak_kaas_instance_pool" "instance_pool" {
  public_cloud_id = wwwwww
  public_cloud_project_id  = xxxxx
  kaas_id = yyyyy
  id      = zzzzz
}
```

## Schema

### Required

- `id` (Integer) The id of the Instance Pool inside the KaaS project.
- `kaas_id` (Integer) The id of the KaaS project.
- `public_cloud_project_id` (Integer) The id of the Public Cloud Project where KaaS is installed.
- `public_cloud_id` (Integer) The id of the Public Cloud where KaaS is installed.

### Read-Only

- `name` (String) The name of the Instance Pool.
- `flavor_name` (String) The flavor name
<!-- - `max_instances` (Number) The maximum amount of instances in the pool. -->
- `min_instances` (Number) The minimum amount of instances in the pool.
