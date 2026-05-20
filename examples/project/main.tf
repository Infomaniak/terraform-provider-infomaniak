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

resource "infomaniak_project" "create_project" {
  public_cloud_id = var.infomaniak.cloud_id

  name             = var.project_name
  user_description = var.project_user_description
  user_email       = var.project_user_email
  user_password    = var.project_user_password
}


/*
data "infomaniak_project" "create_project" {
  public_cloud_id = var.infomaniak.cloud_id
  id              = var.infomaniak.project_id
}
*/
