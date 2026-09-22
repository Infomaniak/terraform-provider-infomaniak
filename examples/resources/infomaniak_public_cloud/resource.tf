# The Public Cloud product cannot be ordered via the API. Order it from the
# Manager UI first, then bring it under Terraform with `terraform import`.
resource "infomaniak_public_cloud" "this" {
  customer_name  = "my-cloud"
  description    = "Production Public Cloud"
  bill_reference = "PO-2026-001"
}
