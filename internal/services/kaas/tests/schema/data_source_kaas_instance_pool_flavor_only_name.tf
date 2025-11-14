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

data "infomaniak_kaas_instance_pool_flavor" "create_instance_pool_flavor" {
  public_cloud_id         = 42
  public_cloud_project_id = 54
  region                  = "dc4-a"

  name = "a8-ram16-disk80-perf1"
}
