resource "infomaniak_kaas" "kluster" {
  public_cloud_id = 45
  pack_name = "standard"
  name = "test"
  kubernetes_version = "1.30"
  region = "dc5"
}
