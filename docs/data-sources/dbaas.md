---
page_title: "infomaniak_dbaas"
subcategory: "DBaaS"
description: |-
  The DBaas Data Source allows the user to read information about a DBaaS project
---

# infomaniak_dbaas

The DBaas Data Source allows the user to read information about a DBaaS project.

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
data "infomaniak_dbaas" "db-0" {
  public_cloud_id         = xxxxx
  public_cloud_project_id = yyyyy
  id                      = zzzzz
}
```

## Schema

### Required

- `public_cloud_id` (Integer) The id of the Public Cloud where DBaaS is installed.
- `public_cloud_project_id` (Integer) The id of the public cloud project where DBaaS is installed.
- `id` (Integer) The id of the DBaaS project.

### Read-Only

- `region` (String) Region where the instance live.
- `kube_identifier` (String) A computed value that gives the kubernetes identifier of the DbaaS
- `pack_name` (String) The name of the pack corresponding the DBaaS project.
- `type` (String) The type of the database to use.
- `version` (String) The version of the database to use.
- `name` (String) The name of the DBaaS shown on the manager.
- `host` (String) The host to access the Database.
- `port` (String) The port to access the Database.
- `user` (String) The user to access the Database.
- `password` (String, Sensitive) The Kubeconfig to access the Database.
- `effective_configuration` (DynamicObject) The MySQL engine parameters on the API side.