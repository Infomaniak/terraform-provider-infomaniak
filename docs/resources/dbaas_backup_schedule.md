---
page_title: "infomaniak_dbaas_backup_schedule"
subcategory: "DBaaS"
description: |-
  The DBaas backup schedule resource allows the user to manage a specific DBaas backups schedules
---

# infomaniak_dbaas_backup_schedule

The DBaas backup schedule resource allows the user to manage a specific DBaas backups schedules.  
You can only create one and only one backup schedule for the moment.

To get your `public_cloud_id`:
```sh
account_id=$(curl -s -H "Authorization: Bearer $INFOMANIAK_TOKEN" https://api.infomaniak.com/2/profile | jq '.data.preferences.account.current_account_id')
curl -s -H "Authorization: Bearer $INFOMANIAK_TOKEN" https://api.infomaniak.com/1/public_clouds?account_id=$account_id | jq '.data[] | {"name": .customer_name, "cloud_id": .id}'
```

To get your `public_cloud_project_id`:
```sh
public_cloud_id=1234  # use the ID retrieved from the step above
curl -s -H "Authorization: Bearer $INFOMANIAK_TOKEN" https://api.infomaniak.com/1/public_clouds/$public_cloud_id/projects | jq '.data[] | {"name": .name, "project_id": .public_cloud_project_id}'
```

## Example

```hcl
resource "infomaniak_dbaas_backup_schedule" "db-0-backup-0" {
  public_cloud_id         = local.public_cloud_id
  public_cloud_project_id = local.public_cloud_project_id
  dbaas_id = infomaniak_dbaas.db-0.id

  scheduled_at = "12:00"
  retention = 3
  is_pitr_enabled = true
}

```

## Schema

### Required

- `public_cloud_id` (Integer) The id of the Public Cloud where DBaaS is installed.
- `public_cloud_project_id` (Integer) The id of the public cloud project where DBaaS is installed.
- `dbaas_id` (Integer) Id of the DBaaS.
- `scheduled_at` (Date) Backup hour in UTC TZ format (`HH:MM`).
- `retention` (Integer) The number of backups to keep.
- `is_pitr_enabled` (Boolean) Enable / disable point in time recovery. 

### Read-Only

- `id` (Integer) The identifier of the scheduled backup.
- `name` (String) The backup schedule generated name.
