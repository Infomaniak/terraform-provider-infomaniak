variable "infomaniak" {
  description = "Infomaniak API credentials."
  type = object({
    token = string
  })
  sensitive = true
  nullable  = false
}

variable "cloud_customer_name" {
  description = "Customer-defined name for the imported Public Cloud product."
  type        = string
}

variable "cloud_description" {
  description = "Free-form description for the imported Public Cloud product."
  type        = string
  default     = ""
}

variable "cloud_bill_reference" {
  description = "Billing reference for the imported Public Cloud product."
  type        = string
  default     = ""
}

variable "project_name" {
  description = "Name of the project to create."
  type        = string
}

variable "project_user_password" {
  description = "Initial password for the project's bootstrap admin user. Forces replacement on change."
  type        = string
  sensitive   = true
}

variable "user_password" {
  description = "Password for the additional ops user. PATCH-able after creation."
  type        = string
  sensitive   = true
}

variable "invitee_email" {
  description = "Email address of the teammate invited into the project."
  type        = string
}

variable "region" {
  description = "Target region for the openrc.sh file."
  type        = string
  default     = "pub1"
}
