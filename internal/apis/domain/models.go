package domain

import "slices"

type Zone struct {
	ID             int        `json:"id,omitempty"`
	FQDN           string     `json:"fqdn,omitempty"`
	DNSSEC         ZoneDNSSEC `json:"dnssec,omitempty"`
	Nameservers    []string   `json:"nameservers,omitempty"`
	Records        []Record   `json:"records,omitempty"`
	ClusterRecords []Record   `json:"cluster_records,omitempty"`
}

type ZoneDNSSEC struct {
	IsEnabled bool `json:"is_enabled,omitempty"`
}

type RecordType = string

var (
	RecordA      RecordType = "A"
	RecordAAAA   RecordType = "AAAA"
	RecordCAA    RecordType = "CAA"
	RecordCNAME  RecordType = "CNAME"
	RecordDNAME  RecordType = "DNAME"
	RecordDNSKEY RecordType = "DNSKEY"
	RecordDS     RecordType = "DS"
	RecordMX     RecordType = "MX"
	RecordNS     RecordType = "NS"
	RecordPTR    RecordType = "PTR"
	RecordSMIMEA RecordType = "SMIMEA"
	RecordSOA    RecordType = "SOA"
	RecordSRV    RecordType = "SRV"
	RecordSSHFP  RecordType = "SSHFP"
	RecordTLSA   RecordType = "TLSA"
	RecordTXT    RecordType = "TXT"
)

var RecordTypes = []RecordType{RecordA, RecordAAAA, RecordCAA, RecordCNAME, RecordDNAME, RecordDNSKEY, RecordDS, RecordMX, RecordPTR, RecordSMIMEA, RecordSOA, RecordSRV, RecordSSHFP, RecordTLSA, RecordTXT}

func IsValidRecordType(t RecordType) bool {
	return slices.Contains(RecordTypes, t)
}

type Record struct {
	ID        int        `json:"id,omitempty"`
	Source    string     `json:"source,omitempty"`
	SourceIDN *string    `json:"source_idn,omitempty"`
	Type      RecordType `json:"type,omitempty"`
	TTL       int        `json:"ttl,omitempty"`
	Target    string     `json:"target,omitempty"`
	DynDNSID  int        `json:"dyndns_id,omitempty"`
	// Description string     `json:"description,omitempty"`
}

type (
	recordTypeA      struct{ string }
	recordTypeAAAA   struct{ string }
	recordTypeCAA    struct{ string }
	recordTypeCNAME  struct{ string }
	recordTypeDNAME  struct{ string }
	recordTypeDNSKEY struct{ string }
	recordTypeDS     struct{ string }
	recordTypeMX     struct{ string }
	recordTypeNS     struct{ string }
	recordTypePTR    struct{ string }
	recordTypeSMIMEA struct{ string }
	recordTypeSOA    struct{ string }
	recordTypeSRV    struct{ string }
	recordTypeSSHFP  struct{ string }
	recordTypeTLSA   struct{ string }
	recordTypeTXT    struct{ string }
)

var (
	RecordTypeA      = recordTypeA{"A"}
	RecordTypeAAAA   = recordTypeAAAA{"AAAA"}
	RecordTypeCAA    = recordTypeCAA{"CAA"}
	RecordTypeCNAME  = recordTypeCNAME{"CNAME"}
	RecordTypeDNAME  = recordTypeDNAME{"DNAME"}
	RecordTypeDNSKEY = recordTypeDNSKEY{"DNSKEY"}
	RecordTypeDS     = recordTypeDS{"DS"}
	RecordTypeMX     = recordTypeMX{"MX"}
	RecordTypeNS     = recordTypeNS{"NS"}
	RecordTypePTR    = recordTypePTR{"PTR"}
	RecordTypeSMIMEA = recordTypeSMIMEA{"SMIMEA"}
	RecordTypeSOA    = recordTypeSOA{"SOA"}
	RecordTypeSRV    = recordTypeSRV{"SRV"}
	RecordTypeSSHFP  = recordTypeSSHFP{"SSHFP"}
	RecordTypeTLSA   = recordTypeTLSA{"TLSA"}
	RecordTypeTXT    = recordTypeTXT{"TXT"}
)

type RecordConstraint interface {
	recordTypeA | recordTypeAAAA | recordTypeCAA | recordTypeCNAME | recordTypeDNAME | recordTypeDNSKEY | recordTypeDS | recordTypeMX | recordTypePTR | recordTypeSMIMEA | recordTypeSOA | recordTypeSRV | recordTypeSSHFP | recordTypeTLSA | recordTypeTXT | recordTypeNS
}
