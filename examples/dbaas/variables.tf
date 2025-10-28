variable "infomaniak" {
  description = "Infomaniak credentials"
  type = object({
    token      = string
    cloud_id   = number
    project_id = number
  })
  sensitive = true
  nullable  = false
}

variable "name" {
  description = "Database name"
  type        = string
  nullable    = false
}

variable "db_region" {
  description = "Database region"
  type        = string
  default     = "dc4-a"
}

variable "pack_name" {
  description = "Pack name"
  type        = string
  default     = "pro-4"
}

variable "db_type" {
  description = "Database type"
  type        = string
  default     = "mysql"
}

variable "db_version" {
  description = "Database version"
  type        = string
  default     = "8.0.42"
}

variable "allowed_cidrs" {
  description = "CIDR whitelist"
  type        = list(string)
  default     = [
    "162.1.15.122/32",
    "1.1.1.1",
    "2345:425:2CA1:0000:0000:567:5673:23b5/64",
  ]
}

variable "time" {
  description = "Backup daily hours in UTC"
  type = string
  default = "12:00"
}

variable "keep" {
  description = "Number of backups to keep"
  type = number
  default = 5
}

variable "is_pitr_enabled" {
  description = "activates / deactivate point in time recovery"
  type = bool
  default = true
}
