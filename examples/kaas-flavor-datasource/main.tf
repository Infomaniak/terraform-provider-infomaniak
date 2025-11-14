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

data "infomaniak_kaas_instance_pool_flavor" "create_instance_pool_flavor" {
  public_cloud_id         = var.infomaniak.cloud_id
  public_cloud_project_id = var.infomaniak.project_id
  region                  = var.cluster_region

  cpu               = 4
  ram               = 8
  storage           = 20
  is_iops_optimized = false
}

resource "infomaniak_kaas" "create_kluster" {
  public_cloud_id         = var.infomaniak.cloud_id
  public_cloud_project_id = var.infomaniak.project_id

  name               = var.cluster_name
  pack_name          = var.cluster_type
  kubernetes_version = var.cluster_version
  region             = var.cluster_region
}

resource "infomaniak_kaas_instance_pool" "create_instance_pool_1" {
  public_cloud_id         = infomaniak_kaas.create_kluster.public_cloud_id
  public_cloud_project_id = infomaniak_kaas.create_kluster.public_cloud_project_id
  kaas_id                 = infomaniak_kaas.create_kluster.id

  name              = "${infomaniak_kaas.create_kluster.name}-pool-1"
  flavor_name       = data.infomaniak_kaas_instance_pool_flavor.create_instance_pool_flavor.name
  min_instances     = var.pool_min
  max_instances     = var.pool_max
  availability_zone = var.pool_az

  labels = var.pool_labels
}
