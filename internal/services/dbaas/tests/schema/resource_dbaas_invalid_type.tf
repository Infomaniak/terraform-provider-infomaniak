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

resource "infomaniak_dbaas" "db" {
  public_cloud_id         = 42
  public_cloud_project_id = 54

  name    = "test"
  region  = "dc5-a"
  type    = "presgresql"
  version = "18.0"
  allowed_cidrs = [
    "0.0.0.0/0"
  ]
  pack_name = "essential-db-4"
}
