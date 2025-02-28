resource "infomaniak_kaas" "kluster" {
  public_cloud_id = 42
  public_cloud_project_id = 54

  pack_name = "standard"
  name = "test"
  kubernetes_version = "1.30"
  region = "dc5"
}

resource "infomaniak_kaas_instance_pool" "instance_pool" {
  kaas_id = infomaniak_kaas.kluster.id

  availability_zone = "dc3-a-04"
  flavor_name = "test"
  min_instances   = 3
  #max_instances   = 6
}
