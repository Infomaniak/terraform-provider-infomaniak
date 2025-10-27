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
  token = var.infomaniak.token
}

resource "infomaniak_dbaas" "db-0" {
  public_cloud_id         = var.infomaniak.cloud_id
  public_cloud_project_id = var.infomaniak.project_id

  name      = var.name
  pack_name = var.pack_name
  type      = var.db_type
  version   = var.db_version
  region    = var.db_region

  allowedCIDRs = var.allowed_cidrs
}
