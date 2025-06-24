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

variable "oidc_params" {
  description = "oidc kubernetes params"
  type = map(string)
  default = {
    "oidc-issuer-url" = "https://your-oidc-url.ch",
    "oidc-client-id" = "kube-login",
    "oidc-username-claim" = "email",
    "oidc-username-prefix" = "oidc:",
  }
}

variable "oidc_ca" {
  description = "oidc CA certificate local file path"
  type        = string
  default     = "./oidc_ca.crt"
}
