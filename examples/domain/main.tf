terraform {
  required_version = ">= 1.5"

  required_providers {
    infomaniak = {
      source  = "Infomaniak/infomaniak"
      version = "~> 1.1.9"
    }
  }
}

provider "infomaniak" {
  token = var.infomaniak.token
}

resource "infomaniak_zone" "zoneA" {
  fqdn = var.zone_fqdn
}

resource "infomaniak_record" "recordA" {
  zone_fqdn = infomaniak_zone.zoneA.fqdn
  type = "A"
  source = var.records.a.source

  data = {
    ip = var.records.a.ip
  }
}

resource "infomaniak_record" "recordB" {
  zone_fqdn = infomaniak_zone.zoneA.fqdn
  type = "SSHFP"
  source = var.records.sshfp.source

  data = {
    fingerprint = var.records.sshfp.fingerprint
    fingerprint_type = var.records.sshfp.fingerprint_type
    fingerprint_algorithm = var.records.sshfp.fingerprint_algorithm
  }
}

resource "infomaniak_record" "recordC" {
  zone_fqdn = infomaniak_zone.zoneA.fqdn
  type = "SSHFP"
  source = var.records.sshfp2.source
  target = var.records.sshfp2.raw_record
}
