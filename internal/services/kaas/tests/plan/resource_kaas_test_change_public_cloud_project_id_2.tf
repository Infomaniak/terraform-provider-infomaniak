resource "infomaniak_kaas" "kluster" {
  public_cloud_id = 41
  public_cloud_project_id = 51

  kubernetes_version = "1.30"
  region = "dc1"
}
