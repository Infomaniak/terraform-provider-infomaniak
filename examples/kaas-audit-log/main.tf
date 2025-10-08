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

resource "infomaniak_kaas" "create_kluster" {
  public_cloud_id         = var.infomaniak.cloud_id
  public_cloud_project_id = var.infomaniak.project_id

  name               = var.cluster_name
  pack_name          = var.cluster_type
  kubernetes_version = var.cluster_version
  region             = var.cluster_region

  apiserver = {
    audit = {
      webhook_config = file(var.audit_logs_webhook_filename)
      policy = file(var.audit_logs_policy_filename)
    }
  }
}
