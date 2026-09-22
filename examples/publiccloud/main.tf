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
  token = var.infomaniak.token
}

# The top-level Public Cloud product cannot be ordered via the API. It must
# already exist (ordered from the Manager UI) and be brought under Terraform
# management with `terraform import`. Update of customer_name/description/
# bill_reference is supported via PATCH.
resource "infomaniak_public_cloud" "this" {
  customer_name  = var.cloud_customer_name
  description    = var.cloud_description
  bill_reference = var.cloud_bill_reference
}

resource "infomaniak_public_cloud_project" "this" {
  public_cloud_id = infomaniak_public_cloud.this.id
  name            = var.project_name
  user_password   = var.project_user_password
}

resource "infomaniak_public_cloud_user" "ops" {
  public_cloud_id         = infomaniak_public_cloud_project.this.public_cloud_id
  public_cloud_project_id = infomaniak_public_cloud_project.this.id

  password    = var.user_password
  description = "ops user managed by Terraform"
}

# An invited user receives an email and sets their own password. Useful for
# onboarding teammates.
resource "infomaniak_public_cloud_user" "invited" {
  public_cloud_id         = infomaniak_public_cloud_project.this.public_cloud_id
  public_cloud_project_id = infomaniak_public_cloud_project.this.id

  invite = true
  email  = var.invitee_email
}

data "infomaniak_public_cloud_openrc" "ops" {
  public_cloud_id         = infomaniak_public_cloud_user.ops.public_cloud_id
  public_cloud_project_id = infomaniak_public_cloud_user.ops.public_cloud_project_id
  public_cloud_user_id    = infomaniak_public_cloud_user.ops.id
  region                  = var.region
}
