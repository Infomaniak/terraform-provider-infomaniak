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
    oidc = {
      issuer_url      = var.issuer_url
      client_id       = var.client_id
      username_claim  = var.username_claim
      username_prefix = var.username_prefix
      signing_algs    = var.signing_algs
      required_claim  = var.required_claim
      ca              = file(var.oidc_ca_filename)
    }
  }
}
