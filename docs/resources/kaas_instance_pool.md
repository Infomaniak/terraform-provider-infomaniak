---
page_title: "infomaniak_kaas_instance_pool"
subcategory: "KaaS"
description: |-
  The Kaas Instance Pools resource is used to manage Instance Pools inside a Kaas project
---

# infomaniak_kaas_instance_pool

The Kaas InstancePool resource is used to manage Instance Pools inside a Kaas project.

Typically, it will come after the creation of a Kaas project.

## Example

```hcl
resource "infomaniak_kaas" "kluster" {
  public_cloud_id         = xxxxx
  public_cloud_project_id = yyyyy

  name               = "kaastor"
  pack_name          = "shared"
  kubernetes_version = "1.31"
  region             = "dc-3"
}

resource "infomaniak_kaas_instance_pool" "instance_pool" {
  public_cloud_id         = infomaniak_kaas.kluster.public_cloud_id
  public_cloud_project_id = infomaniak_kaas.kluster.public_cloud_project_id
  kaas_id                 = infomaniak_kaas.kluster.id

  name              = "instance-pool-1"
  flavor_name       = "a1-ram2-disk20-perf1"
  min_instances     = 4
  availability_zone = "az-1"
}
```

## Schema

### Required

- `public_cloud_id` (Integer) The id of the Public Cloud where KaaS is installed.
- `public_cloud_project_id` (Integer) The id of the Public Cloud Project where KaaS is installed.
- `kaas_id` (Integer) The id of the KaaS project.
- `name` (String) The name of the KaaS shown on the manager.
- `availability_zone` (String) The availability zone where the instances will be populated.
- `flavor_name` (String) The flavor for the instances.
- `min_instances` (Integer) The minimum amount of instances in the pool.

### Read-Only

- `id` (Integer) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
