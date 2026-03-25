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


data "infomaniak_kaas" "kluster" {
  public_cloud_id = 42
  public_cloud_project_id = 54

  id = infomaniak_kaas.kluster.id

  kubeconfig = "ds;lkhjf;lsdhflsk"
}
