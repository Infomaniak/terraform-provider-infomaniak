output "openrc" {
  description = "openrc.sh content for the ops user."
  value       = data.infomaniak_public_cloud_openrc.ops.content
  sensitive   = true
}

output "openstack_project_name" {
  description = "Underlying OpenStack tenant name of the project."
  value       = infomaniak_public_cloud_project.this.open_stack_name
}
