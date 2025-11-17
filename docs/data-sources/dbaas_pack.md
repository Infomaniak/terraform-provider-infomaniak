---
page_title: "infomaniak_dbaas_pack"
subcategory: "DBaaS"
description: |-
  The DBaas Pack Data Source allows the user to read information about DBaaS packs
---

# infomaniak_dbaas_packs

The DBaas Packs Data Source allows the user to read information about DBaaS packs.

## Example using pack name

```hcl
data "infomaniak_dbaas_packs" "db-packs-data" {
  type = "mysql"
  name = "business-db-4"
}
```

## Example using resources

``hcl
data "infomaniak_dbaas_packs" "db-packs-data" {
  type = "mysql"
  
  instances = 2
  cpu       = 2
  ram       = 8
  storage   = 160
}
```

## Schema

### Required

- `type` (String) Database engine type name, available in `infomaniak_dbaas_constants` [data source](./dbaas_constants.md#read-only).

At least one of the following is required:
- `name` (String) Package name identifier (e.g., "essential-db-4", "business-db-16").
- `instances` (Number) Number of database instances included in the package.
- `cpu` (Number) Number of CPU cores allocated to the database instance.
- `ram` (Number) Amount of RAM in GB allocated to the database instance.
- `storage` (Number) Storage capacity in GB allocated to the database instance.

### Read-Only

- `id` (Number) Unique identifier for the package.
- `type` (String) Database engine type (e.g., "mysql", ...).
- `group` (String) Package group category (e.g., "essential", "business", "enterprise").
- `rates` (Object) Pricing information for the package in different currencies.
  - `chf` (Object) Pricing in Swiss Francs.
    - `hour_excl_tax` (Number) Hourly price excluding tax.
    - `hour_incl_tax` (Number) Hourly price including tax.
  - `eur` (Object) Pricing in Euros.
    - `hour_excl_tax` (Number) Hourly price excluding tax.
    - `hour_incl_tax` (Number) Hourly price including tax.