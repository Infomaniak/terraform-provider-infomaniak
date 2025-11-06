---
page_title: "infomaniak_dbaas_packs"
subcategory: "DBaaS"
description: |-
  The DBaas Packs Data Source allows the user to read information about DBaaS packs
---

# infomaniak_dbaas_packs

The DBaas Packs Data Source allows the user to read information about DBaaS packs.

## Example

```hcl
data "infomaniak_dbaas_packs" "db-packs-data" {
  type = "mysql"
}
```

## Schema

### Required

- `type` (String) Database engine type name, available in `infomaniak_dbaas_constants` [data source](./dbaas_constants.md#read-only).

### Read-Only

- `packs` (List of Object) Available DBaaS pre-configured packages with their specifications and pricing.
  - `id` (Number) Unique identifier for the package.
  - `type` (String) Database engine type (e.g., "mysql", ...).
  - `group` (String) Package group category (e.g., "essential", "business", "enterprise").
  - `name` (String) Package name identifier (e.g., "essential-db-4", "business-db-16").
  - `instances` (Number) Number of database instances included in the package.
  - `cpu` (Number) Number of CPU cores allocated to the database instance.
  - `ram` (Number) Amount of RAM in GB allocated to the database instance.
  - `storage` (Number) Storage capacity in GB allocated to the database instance.
  - `rates` (Object) Pricing information for the package in different currencies.
    - `chf` (Object) Pricing in Swiss Francs.
      - `hour_excl_tax` (Number) Hourly price excluding tax.
      - `hour_incl_tax` (Number) Hourly price including tax.
    - `eur` (Object) Pricing in Euros.
      - `hour_excl_tax` (Number) Hourly price excluding tax.
      - `hour_incl_tax` (Number) Hourly price including tax.
