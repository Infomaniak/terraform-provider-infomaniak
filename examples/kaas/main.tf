resource "infomaniak_kaas" "create_kluster" {
  public_cloud_id         = 42
  public_cloud_project_id = 54
  pack_name               = "test"

  region = "dc4"
}

resource "infomaniak_kaas_instance_pool" "create_instance_pool" {
  public_cloud_id         = 54
  public_cloud_project_id = 54
  kaas_id                 = infomaniak_kaas.create_kluster.id

  name          = "coucou"
  flavor_name   = "test"
  min_instances = 4
}

data "infomaniak_kaas" "get_kluster" {
  public_cloud_id         = infomaniak_kaas.create_kluster.public_cloud_id
  public_cloud_project_id = infomaniak_kaas.create_kluster.public_cloud_project_id
  id                      = infomaniak_kaas.create_kluster.id
}

output "kubeconfig" {
  sensitive = true
  value     = infomaniak_kaas.create_kluster.kubeconfig
}
