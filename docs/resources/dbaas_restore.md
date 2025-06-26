---
page_title: "infomaniak_dbaas_restore"
subcategory: "DBaaS"
description: |-
  The DBaas Restore resource allows the user to restore a database to a certain backup.
---

# infomaniak_dbaas_restore

The DBaas Restore resource allows the user to restore a database to a certain backup.

## No-Ops

Deleting this resource will effectively delete it from your Terraform state but will not delete the restore from showing up on your manager.

## Example

```hcl
resource "infomaniak_dbaas_restore" "db-0" {
  public_cloud_id         = xxxxx
  public_cloud_project_id = yyyyy
  dbaas_id                = zzzzz
  backup_id               = ttttt
}
```

## Schema

### Required

- `public_cloud_id` (Integer) The id of the Public Cloud where DBaaS is installed.
- `public_cloud_project_id` (Integer) The id of the public cloud project where DBaaS is installed.
- `dbaas_id` (Integer) The id of the dbaas project.
- `backup_id` (String) The id of the backup to recover from.

### Read-Only

- `id` (String) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `status` (String) The status of the restore.
