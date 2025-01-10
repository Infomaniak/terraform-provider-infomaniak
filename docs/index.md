---
page_title: "Provider: Infomaniak"
sidebar_current: "docs-infomaniak-index"
description: |-
  The Infomaniak Provider is used to interact with Infomaniak's resources. The provider needs to be configured with the proper credentials before it can be used.
---

# Infomaniak Provider

The Infomaniak Provider is used to interact with Infomaniak's resources.

-> __NOTE__ According on your needs, you may need to use additional providers. 

Use the navigation on the left to read about the available resources.

## Configuration

To configure the project you will need an API token, instructions on how to generate one may be found [https://www.infomaniak.com/en/support/faq/2582/generate-and-manage-infomaniak-api-tokens](here).

```hcl
terraform {
  required_providers {
    infomaniak = {
      source = "infomaniak/infomaniak"
    }
  }
}

provider "infomaniak" {
  host          = "https://api.infomaniak.com"
  token         = "xxxxxxxxxxx"
}
```

Alternatively, you can configure these variables using these environment variables :

- `INFOMANIAK_HOST`
- `INFOMANIAK_TOKEN`

## Schema

### Optional

- `host` (String) The base endpoint for Infomaniak's API (including scheme).
- `token` (String, Sensitive) The token used for authenticating against Infomaniak's API.
