---
page_title: "infomaniak_kaas_instance_pool_flavor"
subcategory: "KaaS"
description: |-
  The KaaS Instance Pool Flavor Data Source allows users to read informations about a instance pool flavor.
---

# infomaniak_kaas_instance_pool_flavor

The KaaS Instance Pool Flavor Data Source allows users to read informations about a instance pool flavor.

## Example Usage
### Using the Name Argument
```hcl
data "infomaniak_kaas_instance_pool_flavor" "example_flavor_by_name" {
  public_cloud_id         = xxxxx
  public_cloud_project_id = yyyyy
  region                  = d4-a

  name = "a4-ram8-disk20-perf1"
}
```

### Searching by Hardware Specifications
If multiple matching flavors are found, the provider will prompt the user to refine their search criteria.
```hcl
data "infomaniak_kaas_instance_pool_flavor" "example_flavor_by_specs" {
  public_cloud_id         = xxxxx
  public_cloud_project_id = yyyyy
  region                  = d4-a

  cpu               = 4
  ram               = 8
  storage           = 20
  is_iops_optimized = false
}
```

## Schema
### Required
- public_cloud_id (Integer) The ID of the Public Cloud where the instance pool flavor is located.
- public_cloud_project_id (Integer) The ID of the Public Cloud project associated with the instance pool flavor.
- region (String) Region where the instance pool live.

### Optional
At least one of the following must be provided:
- name (String) The unique identifier name of the instance pool flavor.
- cpu (Integer) Number of vCPUs of the desired flavor.
- ram (Integer) Amount of RAM (in GB) of the desired flavor.
- storage (Integer) Storage size (in GB) of the desired flavor.
Additionally:
- is_memory_optimized (Boolean) Whether the flavor is optimized for memory usage.
- is_iops_optimized (Boolean) Whether the flavor is optimized for high disk usage.
- is_gpu_optimized (Boolean) Whether the flavor supports GPU acceleration.

### Read-Only
- is_available (Boolean) Indicates if the flavor is currently available for use.  
- rates (Object) Pricing details for the flavor:
    - hour_excl_tax (Float) Hourly price excluding taxes.
    - hour_incl_tax (Float) Hourly price including taxes.
