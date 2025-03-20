terraform {
  required_version = ">= 1.5"

  required_providers {
    infomaniak = {
      source  = "Infomaniak/infomaniak"
      version = "~> 1.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "infomaniak" {
  token = "changeme"
}

data "infomaniak_kaas" "get_kluster" {
  public_cloud_id         = 123
  public_cloud_project_id = 456
  id                      = 789
}

locals {
  kubeconfig = yamldecode(infomaniak_kaas.get_kluster.kubeconfig)
}

provider "kubernetes" {
  host                   = local.kubeconfig["clusters"][0]["cluster"]["server"]
  cluster_ca_certificate = base64decode(local.kubeconfig["clusters"][0]["cluster"]["certificate-authority-data"])
  client_certificate     = base64decode(local.kubeconfig["users"][0]["user"]["client-certificate-data"])
  client_key             = base64decode(local.kubeconfig["users"][0]["user"]["client-key-data"])
}

resource "kubernetes_namespace" "hello_world" {
  metadata {
    name = "hello-world"
  }
}
