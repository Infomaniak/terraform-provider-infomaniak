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

  apiserver = {
    acl_rules = [
      "1.2.3.4/5",
      "127.126.125.124/32",
      "127.0.0.1"
    ]

    audit = {
      webhook_config = file("some/file/path/webhook.yaml")
      policy = file("some/file/path/policy.yaml")
    }

    oidc = {
      issuer_url      = "https://issuer.oidc.com"
      client_id       = "oidc-id"
      username_claim  = "email"
      username_prefix = "oidc-refix"
      signing_algs    = "RS256"
      ca              = file("some/file/path/ca.crt")
    }
  }
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

### Optional Configuration

- `apiserver` (Object): The object to configure Kubernetes Apiserver settings. This configuration allows you to customize the behavior of the Apiserver, including audit logging and authentication settings.
  - `acl_rules` (List): The whitelisted CIDRs/IPs allowed to access the Kubernetes API Server.
  - `audit` (Object): The object to configure Kubernetes audit logs using [Kubernetes YAML resources](https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/). Audit logs provide a record of all requests made to the Apiserver, and can be used for security and compliance purposes.
    - `webhook_config` (File): The YAML file specifying the Webhook Config for audit logs. This file defines the endpoint where audit logs will be sent, and can be used to integrate with external logging and monitoring systems.
    - `policy` (File): The YAML file defining the [Audit Policy](https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/#audit-policy) for the cluster. This file specifies the types of events that will be audited, and the level of logging that will be performed.
  - `oidc` (Object): The object to configure OpenID Connect (OIDC) for authentication in the Kubernetes Cluster using [Apiserver flags](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#configuring-the-api-server). OIDC provides a standardized way to authenticate users and services, and can be used to integrate with external identity providers.
    - `issuer_url` (String): The OIDC issuer URL. This is the URL of the OIDC issuer, and is used to verify the authenticity of OIDC tokens.
    - `client_id` (String): The OIDC client ID. This is the client ID of the OIDC application, and is used to identify the application to the OIDC issuer.
    - `username_claim` (String): The claim in the OIDC token that contains the username. This claim is used to extract the username from the OIDC token, and can be used to authenticate users.
    - `username_prefix` (String): The prefix to be added to the username. This prefix can be used to distinguish between different types of users, or to integrate with existing user management systems.
    - `signing_algs` (String): The signing algorithms supported by the OIDC issuer. This specifies the algorithms that can be used to sign OIDC tokens, and can be used to ensure that tokens are properly verified.
    - `ca` (File): The OIDC CA Certificate file. This file contains the CA certificate used to verify the authenticity of OIDC tokens, and is used to establish trust with the OIDC issuer.

### Read-Only

- `id` (Integer) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `kubeconfig` (String, Sensitive) The Kubeconfig to access the Kluster.
