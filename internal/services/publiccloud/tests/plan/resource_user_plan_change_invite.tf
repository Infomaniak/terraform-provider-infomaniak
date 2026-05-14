terraform {
  required_version = ">= 1.5"

  required_providers {
    infomaniak = {
      source  = "Infomaniak/infomaniak"
      version = "~> 1.0"
    }
  }
}

provider "infomaniak" {
  token = "fake-token"
}

resource "infomaniak_public_cloud_user" "this" {
  public_cloud_id         = 42
  public_cloud_project_id = 54
  invite                  = true
  email                   = "ops@example.test"
  description             = "v1"
}
