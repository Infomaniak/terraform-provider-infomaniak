---
page_title: "infomaniak_dbaas"
subcategory: "DBaaS"
description: |-
  The DBaas resource allows the user to manage a DBaas project
---

# infomaniak_dbaas

The DBaas resource allows the user to manage a DBaas project.

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
resource "infomaniak_dbaas" "db-0" {
  public_cloud_id = xxxxx
  public_cloud_project_id = yyyyy
  
  name      = "db-0"
  pack_name = "pro-4"
  type      = "mysql"
  version   = "8.0.42"
  region    = "dc4-a"
}
```

## Schema

### Required

- `public_cloud_id` (Integer) The id of the Public Cloud where DBaaS is installed.
- `public_cloud_project_id` (Integer) The id of the public cloud project where DBaaS is installed.
- `region` (String) Region where the instance live.
- `pack_name` (String) The name of the pack corresponding the DBaaS project.
- `type` (String) The type of the database to use.
- `version` (String) The version of the database to use.
- `name` (String) The name of the DBaaS shown on the manager.

### Read-Only

- `id` (Integer) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `host` (String) The host to access the Database.
- `port` (String) The port to access the Database.
- `user` (String) The user to access the Database.
- `password` (String, Sensitive) The password to access the Database.
- `ca` (String) The database CA certificate.
