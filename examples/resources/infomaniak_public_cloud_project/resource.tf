# Direct creation: bootstrap an admin user with a password.
resource "infomaniak_public_cloud_project" "direct" {
  public_cloud_id = 1234
  name            = "dev-project"
  user_password   = "Sup3r-Secret!"
}

# Invitation: Infomaniak emails the user; they pick their own password.
resource "infomaniak_public_cloud_project" "invited" {
  public_cloud_id = 1234
  name            = "team-project"
  invite          = true
  user_email      = "lead@example.test"
}
