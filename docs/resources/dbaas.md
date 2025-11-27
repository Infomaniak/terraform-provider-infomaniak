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

  allowed_cidrs = [
    "162.1.15.122/32",
    "1.1.1.1",
    "2345:425:2CA1:0000:0000:567:5673:23b5/64",
  ]

  configuration = {
    "connect_timeout": "10",
  }
}
```

Be careful with `allowed_cidrs`:
```hcl
resource "infomaniak_dbaas" "db-0" {
  public_cloud_id = xxxxx
  public_cloud_project_id = yyyyy
  
  name      = "db-0"
  pack_name = "pro-4"
  type      = "mysql"
  version   = "8.0.42"
  region    = "dc4-a"

  allowed_cidrs = [] // If you set an empty list here, it means no one can access it ! Even you !

  configuration = {
    connect_timeout = 300,
    max_connections = 300,
    sql_mode = var.sql_mode // Needs to be inside a variable.tf file since it is a complex type (configuration is a dynamic object and the provider needs a precise type)
  }
}

```hcl
variable "sql_mode" {
  type    = list(string) // Imperative to work
  default = [
    "ERROR_FOR_DIVISION_BY_ZERO"
  ]
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
- `allowed_cidrs` (List of String) The list of allowed cidrs to access to the database.
- `configuration` (DynamicObject) DynamicObject to represent every possible configurations. For available params, please refer to [this documentation](https://developer.infomaniak.com/docs/api/put/1/public_clouds/%7Bpublic_cloud_id%7D/projects/%7Bpublic_cloud_project_id%7D/dbaas/%7Bdbaas_id%7D/configurations).
  - Dynamic Complex Attributes : if one of your attributes is a complex type (e.g: list(string)) please use a `variable.tf` file to specify the type, else the provider might see changes that don't exist because he may think the list is a set.

### Read-Only

- `id` (Integer) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `kube_identifier` (String) A computed value that gives the kubernetes identifier of the DbaaS
- `host` (String) The host to access the Database.
- `port` (String) The port to access the Database.
- `user` (String) The user to access the Database.
- `password` (String, Sensitive) The password to access the Database.
- `ca` (String) The database CA certificate.
- `effective_configuration` (DynamicObject) Specific engine params, including defaulted one (not set by `configuration`)
