---
page_title: "infomaniak_dbaas_constants"
subcategory: "DBaaS"
description: |-
  The DBaas Constants Data Source allows the user to read information about DBaaS constants
---

# infomaniak_dbaas_constants

The DBaas Constants Data Source allows the user to read information about DBaaS constants.

## Example

```hcl
data "infomaniak_dbaas_constants" "db-consts-data" {}
```

## Schema

### Read-Only

- `regions` (List of String) Available DBaaS regions.
- `types` (List of Object) Available DBaaS
  - `name` (String) Database engine name
  - 'versions' (List of String) Available versions for `name` Database engine