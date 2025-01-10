resource "infomaniak_kaas" "create_kluster" {
  pcp_id = "54"

  region = "dc4"
}

resource "infomaniak_kaas_instance_pool" "create_instance_pool" {
  pcp_id  = "54"
  kaas_id = infomaniak_kaas.create_kluster.id

  name        = "coucou"
  flavor_name = "test"
  min_instances   = 4
  max_instances   = 6
}

data "infomaniak_kaas" "get_kluster" {
  pcp_id = infomaniak_kaas.create_kluster.pcp_id
  id = infomaniak_kaas.create_kluster.id
}

output "kubeconfig" {
  sensitive = true
  value     = infomaniak_kaas.create_kluster.kubeconfig
}
