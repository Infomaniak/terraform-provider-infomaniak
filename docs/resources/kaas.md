---
page_title: "infomaniak_kaas"
subcategory: "KaaS"
description: |-
  The Kaas resource allows the user to manage a Kaas project
---

# infomaniak_kaas

The Kaas resource allows the user to manage a Kaas project.

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
resource "infomaniak_kaas" "kluster" {
  public_cloud_id = xxxxx
  public_cloud_project_id = yyyyy
  
  name = "kaastor"
  pack_name = "shared"
  kubernetes_version = "1.31"
  region = "zzzzz"
}
```

## Schema

### Required

- `public_cloud_id` (Integer) The id of the Public Cloud where KaaS is installed.
- `public_cloud_project_id` (Integer) The id of the public cloud project where KaaS is installed.
- `region` (String) Region where the instance live.
- `pack_name` (String) The name of the pack corresponding the KaaS project.
- `kubernetes_version` (String) The version of Kubernetes to use.
- `name` (String) The name of the KaaS shown on the manager.

### Read-Only

- `id` (Integer) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `kubeconfig` (String, Sensitive) The Kubeconfig to access the Kluster.
