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

data "infomaniak_dbaas_pack" "pack" {
  type      = "mysql"
  instances = 1
  cpu       = 2
  ram       = 8
  storage   = 160
}
