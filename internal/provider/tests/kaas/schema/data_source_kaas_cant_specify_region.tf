resource "infomaniak_kaas" "kluster" {
  pcp_id = "54"

  region = "dc5"
}


data "infomaniak_kaas" "kluster" {
  pcp_id = "54"

  id = infomaniak_kaas.kluster.id

  region = "ds;lkhjf;lsdhflsk"
}
