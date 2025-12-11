---
page_title: "infomaniak_kaas_pack"
subcategory: "KaaS"
description: |-
  The KaaS Pack Data Source allows the user to retrieve information about a specific KaaS pack.
---

# infomaniak_kaas_pack

The KaaS Pack Data Source allows the user to retrieve information about a specific KaaS pack based on its name.

## Example Usage

```hcl
data "infomaniak_kaas_pack" "kaas_pack" {
  name = "shared"
}
```

## Schema
## Required
- `name` (String) The name of the KaaS pack to retrieve.

## Read-Only
- `id` (Number) The unique identifier of the KaaS pack.
- `description` (String) A description of what the KaaS pack offers.
- `price_per_hour` (Object) Pricing details per hour, broken down by currency.
    - `chf` (Number) Price per hour in Swiss Francs (CHF).
    - `eur` (Number) Price per hour in Euros (EUR).
- `limit_per_project` (Number) Maximum number of clusters allowed per project using this pack.
- `is_active` (Boolean) Indicates whether the pack is currently active and available for use.
