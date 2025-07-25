package domain

import (
	"net"
	"strings"
	"terraform-provider-infomaniak/internal/apis/domain"

	"github.com/miekg/dns"
)

func (model *RecordModel) ComputeRawTarget() string {
	// don't do anything if it's already set
	if !model.Target.IsUnknown() && !model.Target.IsNull() {
		return model.Target.ValueString()
	}

	var record dns.RR

	switch model.Type.ValueString() {
	case domain.RecordA:
		record = &dns.A{
			A: net.ParseIP(model.Data.IP.ValueString()),
		}
	case domain.RecordAAAA:
		record = &dns.AAAA{
			AAAA: net.ParseIP(model.Data.IP.ValueString()),
		}
	case domain.RecordCAA:
		record = &dns.CAA{
			Flag:  uint8(model.Data.Flags.ValueInt64()),
			Tag:   model.Data.Tag.ValueString(),
			Value: model.Data.Value.ValueString(),
		}
	case domain.RecordCNAME:
		record = &dns.CNAME{
			Target: dns.Fqdn(model.Data.Target.ValueString()),
		}
	case domain.RecordDNAME:
		record = &dns.DNAME{
			Target: dns.Fqdn(model.Data.Target.ValueString()),
		}
	case domain.RecordDNSKEY:
		record = &dns.DNSKEY{
			Flags:     uint16(model.Data.Flags.ValueInt64()),
			Protocol:  3,
			Algorithm: uint8(model.Data.Algorithm.ValueInt64()),
			PublicKey: model.Data.PublicKey.ValueString(),
		}
	case domain.RecordDS:
		record = &dns.DS{
			KeyTag:     uint16(model.Data.KeyTag.ValueInt64()),
			Algorithm:  uint8(model.Data.Algorithm.ValueInt64()),
			DigestType: uint8(model.Data.DigestType.ValueInt64()),
			Digest:     model.Data.Digest.ValueString(),
		}
	case domain.RecordMX:
		record = &dns.MX{
			Preference: uint16(model.Data.Priority.ValueInt64()),
			Mx:         dns.Fqdn(model.Data.Target.ValueString()),
		}
	case domain.RecordNS:
		record = &dns.NS{
			Ns: dns.Fqdn(model.Data.Target.ValueString()),
		}
	case domain.RecordPTR:
		record = &dns.PTR{
			Ptr: dns.Fqdn(model.Data.Target.ValueString()),
		}
	case domain.RecordSMIMEA:
		record = &dns.SMIMEA{
			Usage:        uint8(model.Data.Priority.ValueInt64()),
			Selector:     uint8(model.Data.Selector.ValueInt64()),
			MatchingType: uint8(model.Data.MatchingType.ValueInt64()),
			Certificate:  model.Data.CertAssocData.ValueString(),
		}
	case domain.RecordSOA:
		record = &dns.SOA{
			Ns:      dns.Fqdn(model.Data.MName.ValueString()),
			Mbox:    dns.Fqdn(model.Data.RName.ValueString()),
			Serial:  uint32(model.Data.Serial.ValueInt64()),
			Refresh: uint32(model.Data.Refresh.ValueInt64()),
			Retry:   uint32(model.Data.Retry.ValueInt64()),
			Expire:  uint32(model.Data.Expire.ValueInt64()),
			Minttl:  uint32(model.Data.Minimum.ValueInt64()),
		}
	case domain.RecordSRV:
		record = &dns.SRV{
			Priority: uint16(model.Data.Priority.ValueInt64()),
			Weight:   uint16(model.Data.Weight.ValueInt64()),
			Port:     uint16(model.Data.Port.ValueInt64()),
			Target:   dns.Fqdn(model.Data.Target.ValueString()),
		}
	case domain.RecordSSHFP:
		record = &dns.SSHFP{
			Algorithm:   uint8(model.Data.FingerprintAlgorithm.ValueInt64()),
			Type:        uint8(model.Data.FingerprintType.ValueInt64()),
			FingerPrint: model.Data.Fingerprint.ValueString(),
		}
	case domain.RecordTLSA:
		record = &dns.TLSA{
			Usage:        uint8(model.Data.Priority.ValueInt64()),
			Selector:     uint8(model.Data.Selector.ValueInt64()),
			MatchingType: uint8(model.Data.MatchingType.ValueInt64()),
			Certificate:  model.Data.CertAssocData.ValueString(),
		}
	case domain.RecordTXT:
		record = &dns.TXT{
			Txt: []string{model.Data.Value.ValueString()},
		}
	}

	return strings.TrimPrefix(record.String(), record.Header().String())
}
