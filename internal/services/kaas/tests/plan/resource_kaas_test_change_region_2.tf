resource "infomaniak_kaas" "kluster" {
  public_cloud_id = 41
  public_cloud_project_id = 50

  kubernetes_version = "1.30"
  region = "dc2"
}
