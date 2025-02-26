resource "infomaniak_kaas" "kluster" {
  public_cloud_id = 42
  public_cloud_project_id = 54

  kubernetes_version = "1.30"
  region = "dc5"
}


data "infomaniak_kaas" "kluster" {
  public_cloud_id = 42
  public_cloud_project_id = 54

  id = infomaniak_kaas.kluster.id
}
