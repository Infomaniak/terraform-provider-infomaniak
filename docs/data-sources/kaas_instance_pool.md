---
page_title: "infomaniak_kaas_instance_pool"
subcategory: "KaaS"
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
- `min_instances` (Integer) The minimum amount of instances in the pool.
- `max_instances` (Integer) The maximum amount of instances in the pool.
- `availability_zone` (String) The availability zone where the instances will be populated.
- `flavor_name` (String) The flavor for the instances.
