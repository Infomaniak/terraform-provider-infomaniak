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


# See https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/#audit-policy for more details on Kubernetes Auditing configuration
variable "audit_logs_webhook_filename" {
  description = "Audit logs webhook yaml filename"
  type        = string
}


variable "audit_logs_policy_filename" {
  description = "Audit logs policy yaml filename"
  type        = string
}
