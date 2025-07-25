variable "infomaniak" {
  description = "Infomaniak credentials"
  type = object({
    token      = string
  })
  sensitive = true
  nullable  = false
}

variable "zone_fqdn" {
  description = "FQDN for the specific zone"
  type        = string
  nullable    = false
}

variable "records" {
  description = "Infomaniak credentials"
  type = object({
    a = object({
      source = string
      ip = string
    })

    sshfp = object({
      source = string
      fingerprint = string
      fingerprint_type = number
      fingerprint_algorithm = number
    })

    sshfp2 = object({
      source = string
      raw_record = string
    })
  })
  nullable  = false

  default = {
    a = {
      source = "server1"
      ip     = "192.0.2.1"
    }

    sshfp = {
      source               = "server1"
      fingerprint          = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
      fingerprint_type     = 1
      fingerprint_algorithm = 1
    }

    sshfp2 = {
      source     = "server1"
      raw_record = "1 1 bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
    }
  }
}