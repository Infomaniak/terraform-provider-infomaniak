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

resource "infomaniak_public_cloud_project" "this" {
  public_cloud_id = 99
  name            = "plan-test"
  user_password   = "Secret!Password1"
}
