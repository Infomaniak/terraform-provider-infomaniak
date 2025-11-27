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

variable "project_name" {
  description = "Project name"
  type        = string
  nullable    = false
}

variable "project_user_description" {
  description = "Project user description"
  type        = string
  nullable    = true
}

variable "project_user_email" {
  description = "Project user email"
  type        = string
  nullable    = true
}

variable "project_user_password" {
  description = "Project user password"
  type        = string
  nullable    = false
}