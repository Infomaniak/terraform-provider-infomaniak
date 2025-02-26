resource "infomaniak_kaas" "kluster" {
  pcp_id = "54"

  region = "dc5"
}

resource "infomaniak_kaas_instance_pool" "instance_pool" {
  pcp_id  = infomaniak_kaas.kluster.pcp_id
  kaas_id = infomaniak_kaas.kluster.id

  name        = "coucou"
  flavor_name = "test"
  max_instances   = 6
}
