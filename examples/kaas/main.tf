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
    params = {
      "--oidc-groups-claim" = "custom-set-param"
    }
    oidc = {
      issuer_url      = var.issuer_url
      client_id       = var.client_id
      username_claim  = var.username_claim
      username_prefix = var.username_prefix
      signing_algs    = var.signing_algs
      ca              = file(var.oidc_ca_filename)
    }
  }
}

resource "infomaniak_kaas_instance_pool" "create_instance_pool_1" {
  public_cloud_id         = infomaniak_kaas.create_kluster.public_cloud_id
  public_cloud_project_id = infomaniak_kaas.create_kluster.public_cloud_project_id
  kaas_id                 = infomaniak_kaas.create_kluster.id

  name              = "${infomaniak_kaas.create_kluster.name}-pool-1"
  flavor_name       = var.pool_type
  min_instances     = var.pool_min
  availability_zone = var.pool_az

  labels            = var.pool_labels
}
