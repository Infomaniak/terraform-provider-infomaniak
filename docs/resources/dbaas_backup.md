---
page_title: "infomaniak_dbaas_backup"
subcategory: "DBaaS"
description: |-
  The DBaas Backup resource allows the user to manage a Backup of a database.
---

# infomaniak_dbaas_backup

The DBaas Backup resource allows the user to manage a Backup of a database.

## Example

```hcl
resource "infomaniak_dbaas_backup" "db-0" {
  public_cloud_id         = xxxxx
  public_cloud_project_id = yyyyy
  dbaas_id                = zzzzz
}
```

## Schema

### Required

- `public_cloud_id` (Integer) The id of the Public Cloud where DBaaS is installed.
- `public_cloud_project_id` (Integer) The id of the public cloud project where DBaaS is installed.
- `dbaas_id` (Integer) The id of the dbaas project.

### Read-Only

- `id` (String) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `status` (String) The status of the backup.
