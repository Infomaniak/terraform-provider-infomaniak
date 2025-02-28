resource "infomaniak_kaas" "kluster" {
  public_cloud_id = 42
  public_cloud_project_id = 54

  pack_name = "standard"
  name = "test"
  kubernetes_version = "1.30"
  region = "dc5"
}
