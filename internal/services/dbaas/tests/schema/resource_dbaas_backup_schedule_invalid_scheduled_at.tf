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
  type    = "mysql"
  version = "8.0"
  allowed_cidrs = [
    "0.0.0.0/0"
  ]
  pack_name = "essential-db-4"
}

resource "infomaniak_dbaas_backup_schedule" "my_schedule" {
  public_cloud_id         = infomaniak_dbaas.db.public_cloud_id
  public_cloud_project_id = infomaniak_dbaas.db.public_cloud_project_id
  dbaas_id                = infomaniak_dbaas.db.id
  scheduled_at            = "berries"
  is_pitr_enabled         = true
  retention               = 30
}
