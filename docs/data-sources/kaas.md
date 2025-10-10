---
page_title: "infomaniak_kaas"
subcategory: "KaaS"
description: |-
  The Kaas Data Source allows the user to read information about a Kaas project
---

# infomaniak_kaas (Data Source)

The Kaas Data Source allows the user to read information about a Kaas project.

## Example

```hcl
data "infomaniak_kaas" "kluster" {
  public_cloud_id = xxxxx
  public_cloud_project_id = yyyyy
  id     = zzzzz
}
```

## Schema

### Required

- `id` (Integer) The id of the KaaS project.
- `public_cloud_project_id` (Integer) The id of the Public Cloud Project where KaaS is installed.
- `public_cloud_id` (Integer) The id of the Public Cloud where KaaS is installed.

### Read-Only

- `kubeconfig` (String, Sensitive) The Kubeconfig to access the Kluster.
- `region` (String) Region where the instance live.
- `pack_name` (String) The name of the pack corresponding the KaaS project.
- `kubernetes_version` (String) The version of Kubernetes to use.
- `name` (String) The name of the KaaS shown on the manager.
- `apiserver` (Object): The object to configure Kubernetes Apiserver settings. This configuration allows you to customize the behavior of the Apiserver, including audit logging and authentication settings.
  - `audit` (Object): The object to configure Kubernetes audit logs using [Kubernetes YAML resources](https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/). Audit logs provide a record of all requests made to the Apiserver, and can be used for security and compliance purposes.
    - `webhook_config` (File): The YAML file specifying the Webhook Config for audit logs. This file defines the endpoint where audit logs will be sent, and can be used to integrate with external logging and monitoring systems.
    - `policy` (File): The YAML file defining the [Audit Policy](https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/#audit-policy) for the cluster. This file specifies the types of events that will be audited, and the level of logging that will be performed.
  - `oidc` (Object): The object to configure OpenID Connect (OIDC) for authentication in the Kubernetes Cluster using [Apiserver flags](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#configuring-the-api-server). OIDC provides a standardized way to authenticate users and services, and can be used to integrate with external identity providers.
    - `issuer_url` (String): The OIDC issuer URL. This is the URL of the OIDC issuer, and is used to verify the authenticity of OIDC tokens.
    - `client_id` (String): The OIDC client ID. This is the client ID of the OIDC application, and is used to identify the application to the OIDC issuer.
    - `username_claim` (String): The claim in the OIDC token that contains the username. This claim is used to extract the username from the OIDC token, and can be used to authenticate users.
    - `username_prefix` (String): The prefix to be added to the username. This prefix can be used to distinguish between different types of users, or to integrate with existing user management systems.
    - `signing_algs` (String): The signing algorithms supported by the OIDC issuer. This specifies the algorithms that can be used to sign OIDC tokens, and can be used to ensure that tokens are properly verified.
    - `required_claim` (String): A key=value pair that describes a required claim in the ID Token.
    - `ca` (File): The OIDC CA Certificate file. This file contains the CA certificate used to verify the authenticity of OIDC tokens, and is used to establish trust with the OIDC issuer.
