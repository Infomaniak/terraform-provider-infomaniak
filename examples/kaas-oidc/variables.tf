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

variable "pool_type" {
  description = "Pool instance type"
  type        = string
  default     = "a1-ram2-disk20-perf1"
}

variable "pool_min" {
  description = "Minimum pool instance number"
  type        = number
  default     = 3
}

variable "pool_az" {
  description = "Pool instance availability zone"
  type        = string
  default     = "az-2"
}

variable "pool_labels" {
  description = "Pool instance custom labels"
  type = map(string)
  default = {
    "node-role.kubernetes.io/worker" = "high"
  }
}

# See https://kubernetes.io/docs/reference/access-authn-authz/authentication/#configuring-the-api-server for more details on Kubernetes OIDC Apiserver parameters configuration
variable "issuer_url" {
  description = "OIDC issuer url"
  type        = string
}

variable "client_id" {
  description = "OIDC client id"
  type        = string
}

variable "username_claim" {
  description = "OIDC username claim"
  type        = string
}

variable "username_prefix" {
  description = "OIDC username prefix"
  type        = string
}

variable "signing_algs" {
  description = "OIDC effective signing algorithm"
  type        = string
}


variable "audit_logs_webhook_filename" {
  description = "Audit logs webhook yaml filename"
  type        = string
}

variable "oidc_ca_filename" {
  description = "OIDC CA certificate local file path"
  type        = string
}
