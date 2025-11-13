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

resource "infomaniak_dbaas_backup_schedule" "db-0-backup-0" {
  public_cloud_id         = infomaniak_dbaas.db-0.public_cloud_id
  public_cloud_project_id = infomaniak_dbaas.db-0.public_cloud_project_id
  dbaas_id = infomaniak_dbaas.db-0.id

  scheduled_at = var.scheduled_at
  retention = var.retention
  is_pitr_enabled = var.is_pitr_enabled
}
