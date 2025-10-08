---
page_title: "infomaniak_zone"
subcategory: "Domain"
description: |-
  The infomaniak_zone resource allows the user to manage a DNS Zone
---

## infomaniak_zone

The `infomaniak_zone` resource allows the user to manage a zone for a domain project.

### Example Usage

```hcl
resource "infomaniak_zone" "example" {
  fqdn = "example.com."
}
```

### Attributes

#### Required

- `fqdn` (String) The FQDN of the zone.

#### Computed

- `id` (Number) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
