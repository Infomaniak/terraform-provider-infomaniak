resource "infomaniak_kaas" "kluster" {
  public_cloud_id = 41
  public_cloud_project_id = 50

  pack_name = "standard"
  name = "test"
  kubernetes_version = "1.30"
  region = "dc4"
}

resource "infomaniak_kaas_instance_pool" "instance_pool" {
  public_cloud_id = infomaniak_kaas.kluster.public_cloud_id
  public_cloud_project_id  = infomaniak_kaas.kluster.public_cloud_project_id
  kaas_id = infomaniak_kaas.kluster.id

  name        = "coucou"
  availability_zone = "dc3-a-04"
  flavor_name = "test"
  min_instances   = 3
  #max_instances   = 6
}
