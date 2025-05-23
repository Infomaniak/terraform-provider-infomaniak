package domain

import (
	"net"
	"strings"
	"terraform-provider-infomaniak/internal/apis/domain"

	"github.com/miekg/dns"
)

func (model *RecordModel) ComputeRawTarget() string {
	// don't do anything if it's already set
	if !model.RawTarget.IsUnknown() {
		return model.RawTarget.ValueString()
	}

	var record dns.RR

	switch model.Type.ValueString() {
	case domain.RecordA:
		record = &dns.A{
			A: net.ParseIP(model.IP.ValueString()),
		}
	case domain.RecordAAAA:
		record = &dns.AAAA{
			AAAA: net.ParseIP(model.IP.ValueString()),
		}
	case domain.RecordCAA:
		record = &dns.CAA{
			Flag:  uint8(model.Flags.ValueInt64()),
			Tag:   model.Tag.ValueString(),
			Value: model.Value.ValueString(),
		}
	case domain.RecordCNAME:
		record = &dns.CNAME{
			Target: dns.Fqdn(model.Target.ValueString()),
		}
	case domain.RecordDNAME:
		record = &dns.DNAME{
			Target: dns.Fqdn(model.Target.ValueString()),
		}
	case domain.RecordDNSKEY:
		record = &dns.DNSKEY{
			Flags:     uint16(model.Flags.ValueInt64()),
			Protocol:  3,
			Algorithm: uint8(model.Algorithm.ValueInt64()),
			PublicKey: model.PublicKey.ValueString(),
		}
	case domain.RecordDS:
		record = &dns.DS{
			KeyTag:     uint16(model.Id.ValueInt64()), // or model.KeyTag.ValueInt64() if separate
			Algorithm:  uint8(model.Algorithm.ValueInt64()),
			DigestType: uint8(model.DigestType.ValueInt64()),
			Digest:     model.Digest.ValueString(),
		}
	case domain.RecordMX:
		record = &dns.MX{
			Preference: uint16(model.Priority.ValueInt64()),
			Mx:         dns.Fqdn(model.Target.ValueString()),
		}
	case domain.RecordNS:
		record = &dns.NS{
			Ns: dns.Fqdn(model.Target.ValueString()),
		}
	case domain.RecordPTR:
		record = &dns.PTR{
			Ptr: dns.Fqdn(model.Target.ValueString()),
		}
	case domain.RecordSMIMEA:
		record = &dns.SMIMEA{
			Usage:        uint8(model.Priority.ValueInt64()),
			Selector:     uint8(model.Selector.ValueInt64()),
			MatchingType: uint8(model.MatchingType.ValueInt64()),
			Certificate:  model.CertAssocData.ValueString(),
		}
	case domain.RecordSOA:
		record = &dns.SOA{
			Ns:      dns.Fqdn(model.MName.ValueString()),
			Mbox:    dns.Fqdn(model.RName.ValueString()),
			Serial:  uint32(model.Serial.ValueInt64()),
			Refresh: uint32(model.Refresh.ValueInt64()),
			Retry:   uint32(model.Retry.ValueInt64()),
			Expire:  uint32(model.Expire.ValueInt64()),
			Minttl:  uint32(model.Minimum.ValueInt64()),
		}
	case domain.RecordSRV:
		record = &dns.SRV{
			Priority: uint16(model.Priority.ValueInt64()),
			Weight:   uint16(model.Weight.ValueInt64()),
			Port:     uint16(model.Port.ValueInt64()),
			Target:   dns.Fqdn(model.Target.ValueString()),
		}
	case domain.RecordSSHFP:
		record = &dns.SSHFP{
			Algorithm:   uint8(model.Algorithm.ValueInt64()),
			Type:        uint8(model.Fptype.ValueInt64()),
			FingerPrint: model.Fingerprint.ValueString(),
		}
	case domain.RecordTLSA:
		record = &dns.TLSA{
			Usage:        uint8(model.Priority.ValueInt64()),
			Selector:     uint8(model.Selector.ValueInt64()),
			MatchingType: uint8(model.MatchingType.ValueInt64()),
			Certificate:  model.CertAssocData.ValueString(),
		}
	case domain.RecordTXT:
		record = &dns.TXT{
			Txt: []string{model.Value.ValueString()},
		}
	}

	return strings.TrimPrefix(record.String(), record.Header().String())
}
