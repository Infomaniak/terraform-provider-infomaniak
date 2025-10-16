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

resource "infomaniak_kaas" "kluster" {
  public_cloud_id = 42
  public_cloud_project_id = 54

  pack_name = "standard"
  name = "test"
  kubernetes_version = "1.30"
  region = "dc5"
}

resource "infomaniak_kaas_instance_pool" "instance_pool" {
  public_cloud_id  = infomaniak_kaas.kluster.public_cloud_id
  public_cloud_project_id  = infomaniak_kaas.kluster.public_cloud_project_id
  kaas_id = infomaniak_kaas.kluster.id

  name        = "coucou"
  availability_zone = "dc3-a-04"
  flavor_name = "test"
  min_instances   = 3
  max_instances   = 6
}

data "infomaniak_kaas_instance_pool" "instance_pool" {
  public_cloud_id  = infomaniak_kaas.kluster.public_cloud_id
  public_cloud_project_id = infomaniak_kaas_instance_pool.instance_pool.public_cloud_project_id
  kaas_id = infomaniak_kaas_instance_pool.instance_pool.kaas_id
  id = infomaniak_kaas_instance_pool.instance_pool.id

  max_instances = 6
}
