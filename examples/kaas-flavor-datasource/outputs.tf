output "kubeconfig" {
  description = "Cluster kubeconfig content"
  value       = infomaniak_kaas.create_kluster.kubeconfig
  sensitive   = true
}

