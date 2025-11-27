output "project_id" {
  description = "Public cloud project id"
  value       = infomaniak_project.create_project.id
  sensitive   = true
}

output "open_stack_name" {
  description = "Open stack project name"
  value       = infomaniak_project.create_project.open_stack_name
  sensitive   = true

}

