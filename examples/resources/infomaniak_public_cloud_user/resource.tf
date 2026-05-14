resource "infomaniak_public_cloud_user" "ops" {
  public_cloud_id         = 1234
  public_cloud_project_id = 5678
  password                = "Sup3r-Secret!"
  description             = "operator account"
}

resource "infomaniak_public_cloud_user" "invited" {
  public_cloud_id         = 1234
  public_cloud_project_id = 5678
  invite                  = true
  email                   = "ops@example.test"
}
