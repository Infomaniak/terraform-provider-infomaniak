---
page_title: "infomaniak_record"
subcategory: "Domain"
description: |-
  The infomaniak_record resource allows the user to manage a DNS Record
---

# infomaniak_record

The `infomaniak_record` resource allows the user to manage a DNS Record.  
You can either use the `data` field or the `target` field to specify the target of your Record, the `data` field contains helper for every supported type of DNS Record.

## Example

### Using raw record value

```hcl
resource "infomaniak_record" "recordC" {
  zone_fqdn = infomaniak_zone.zoneA.fqdn
  type = "SSHFP"
  source = var.records.sshfp2.source
  target = var.records.sshfp2.raw_record
}
```

### Using the `data` field

```hcl
resource "infomaniak_record" "recordB" {
  zone_fqdn = infomaniak_zone.zoneA.fqdn
  type = "SSHFP"
  source = var.records.sshfp.source

  data = {
    fingerprint = var.records.sshfp.fingerprint
    fingerprint_type = var.records.sshfp.fingerprint_type
    fingerprint_algorithm = var.records.sshfp.fingerprint_algorithm
  }
}
```

## Schema

### Required

- `zone_fqdn` (String) The FQDN of the zone where the record should be put in.
- `type` (String) Record Type. One of : "A", "AAAA", "CAA", "CNAME", "DNAME", "MX", "NS", "TXT", "DS", "HTTPS", "SMIMEA", "SRV", "SSHFP", "TLSA".
- `source` (String) The source of the record.

### Optional

- `description` (String) The description of the record.
- `ttl` (Integer) The TTL of the Record.
- `target` (String) The target of the Record (cannot be used with `data`).
- `data` ([Object](#nested-schema-for-data)) Components of a record.

### Read-Only

- `id` (Integer) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `computed_target` (String) The computed target for the record, either `target` or the representation of the `data` field.

### Nested Schema for `data`

Optional object containing the components of a DNS record.  
This field **conflicts with** the `target` field (i.e., cannot be set at the same time).

#### Optional

- `ip` (String) IP address of the record. Relevant for record types: A, AAAA.
- `priority` (Integer) Priority, usage, or weight depending on the record type. Relevant for: MX, SRV, TLSA, SMIMEA.
- `weight` (Integer) Weight of the service for load balancing. Relevant for: SRV.
- `port` (Integer) Port number used for the service. Relevant for: SRV.
- `flags` (Integer) Flags used in the record. Relevant for: CAA.
- `tag` (String) Tag name such as `issue`, `issuewild`, or `iodef`. Relevant for: CAA.
- `algorithm` (Integer) Cryptographic algorithm identifier. Relevant for: DS, SSHFP.
- `key_tag` (Integer) Key Tag of the DNSKEY. Relevant for: DS.
- `digest_type` (Integer) Digest algorithm type used to hash the DNSKEY. Relevant for: DS.
- `digest` (String) Digest value (usually a hex string). Relevant for: DS.
- `selector` (Integer) Specifies which part of the certificate is matched. Relevant for: TLSA, SMIMEA.
- `matching_type` (Integer) Specifies how the certificate data is matched (e.g., SHA256). Relevant for: TLSA, SMIMEA.
- `cert_assoc_data` (String) Certificate association data (usually a hash). Relevant for: TLSA, SMIMEA.
- `fingerprint_algorithm` (Integer) Algorithm used to create the SSH key fingerprint. Relevant for: SSHFP.
- `fingerprint_type` (Integer) Type of hash used for the fingerprint (e.g., SHA1, SHA256). Relevant for: SSHFP.
- `fingerprint` (String) Hex-encoded fingerprint of the SSH public key. Relevant for: SSHFP.
- `target` (String) Target FQDN of the record. Relevant for: MX, CNAME, DNAME, NS, PTR, etc.
- `value` (String) Generic string value for the record. Relevant for: TXT, CAA, and other textual records.
