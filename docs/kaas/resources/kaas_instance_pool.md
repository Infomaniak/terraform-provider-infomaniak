---
page_title: "infomaniak_kaas_instance_pool"
subcategory: ""
description: |-
  The Kaas Instance Pools resource is used to manage Instance Pools inside a Kaas project
---

# infomaniak_kaas_instance_pool

The Kaas InstancePool resource is used to manage Instance Pools inside a Kaas project.

Typically, it will come after the creation of a Kaas project.

Setting `min_instances` = `max_instances` will disable autoscaling.

## Example

```hcl
resource "infomaniak_kaas" "kluster" {
  pcp_id = "xxxxx"

  region = "yyyyy"
}

resource "infomaniak_kaas_instance_pool" "instance_pool" {
  pcp_id  = infomaniak_kaas.kluster.pcp_id
  kaas_id = infomaniak_kaas.kluster.id

  name        = "instance-pool-1"
  flavor_name = "a1_ram2_disk20_perf1"
  min_instances   = 4
  max_instances   = 6
}
```

## Schema

### Required

- `pcp_id` (String) The id of the public cloud project where KaaS is installed.
- `kaas_id` (String) The id of the kaas project.
- `name` (String) The name of the instance pool.
- `flavor_name` (String) The flavor name.
- `max_instances` (Number) The maximum amount of instances in the pool.
- `min_instances` (Number) The minimum amount of instances in the pool.

### Read-Only

- `id` (String) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
