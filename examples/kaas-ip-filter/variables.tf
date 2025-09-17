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

variable "cluster_name" {
  description = "Cluster name"
  type        = string
  nullable    = false
}

variable "cluster_region" {
  description = "Cluster region"
  type        = string
  default     = "dc4-a"
}

variable "cluster_type" {
  description = "Cluster type"
  type        = string
  default     = "shared"
}

variable "cluster_version" {
  description = "Cluster version"
  type        = string
  default     = "1.31"
}

variable "allow_requests_from_cidr" {
  description = "Allow only specific cidrs / IPs to access to control plane"
  type = list(string)
  default = [
    "1.1.1.1",
    "2.2.2.2/32",
    "192.168.0.0/24"
  ]
}
