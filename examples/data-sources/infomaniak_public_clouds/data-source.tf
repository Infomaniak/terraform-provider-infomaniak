data "infomaniak_public_clouds" "all" {}

output "cloud_ids" {
  value = [for c in data.infomaniak_public_clouds.all.public_clouds : c.id]
}
