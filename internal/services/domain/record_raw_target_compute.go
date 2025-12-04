package domain

import (
	"fmt"
	"net"
	"strings"
	"terraform-provider-infomaniak/internal/apis/domain"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/miekg/dns"
)

type Unknowable interface {
	IsUnknown() bool
}

func (model *RecordModel) ComputeRawTarget() (value string, isKnown bool) {
	// If user explicitly set target, return it directly
	if !model.Target.IsNull() {
		return model.Target.ValueString(), !model.Target.IsUnknown()
	}

	isAnyFieldUnknown := func(vals ...Unknowable) bool {
		for _, v := range vals {
			if v.IsUnknown() {
				return true
			}
		}
		return false
	}

	var record dns.RR
	isFullyKnown := true

	switch model.Type.ValueString() {

	case domain.RecordA:
		isFullyKnown = !model.Data.IP.IsUnknown()
		if isFullyKnown {
			record = &dns.A{A: net.ParseIP(model.Data.IP.ValueString())}
		}

	case domain.RecordAAAA:
		isFullyKnown = !model.Data.IP.IsUnknown()
		if isFullyKnown {
			record = &dns.AAAA{AAAA: net.ParseIP(model.Data.IP.ValueString())}
		}

	case domain.RecordCAA:
		isFullyKnown = !isAnyFieldUnknown(model.Data.Flags, model.Data.Tag, model.Data.Value)
		if isFullyKnown {
			record = &dns.CAA{
				Flag:  uint8(model.Data.Flags.ValueInt64()),
				Tag:   model.Data.Tag.ValueString(),
				Value: model.Data.Value.ValueString(),
			}
		}

	case domain.RecordCNAME:
		isFullyKnown = !model.Data.Target.IsUnknown()
		if isFullyKnown {
			record = &dns.CNAME{Target: dns.Fqdn(model.Data.Target.ValueString())}
		}

	case domain.RecordDNAME:
		isFullyKnown = !model.Data.Target.IsUnknown()
		if isFullyKnown {
			record = &dns.DNAME{Target: dns.Fqdn(model.Data.Target.ValueString())}
		}

	case domain.RecordDS:
		isFullyKnown = !isAnyFieldUnknown(model.Data.KeyTag, model.Data.Algorithm, model.Data.DigestType, model.Data.Digest)
		if isFullyKnown {
			record = &dns.DS{
				KeyTag:     uint16(model.Data.KeyTag.ValueInt64()),
				Algorithm:  uint8(model.Data.Algorithm.ValueInt64()),
				DigestType: uint8(model.Data.DigestType.ValueInt64()),
				Digest:     model.Data.Digest.ValueString(),
			}
		}

	case domain.RecordMX:
		isFullyKnown = !isAnyFieldUnknown(model.Data.Priority, model.Data.Target)
		if isFullyKnown {
			record = &dns.MX{
				Preference: uint16(model.Data.Priority.ValueInt64()),
				Mx:         dns.Fqdn(model.Data.Target.ValueString()),
			}
		}

	case domain.RecordNS:
		isFullyKnown = !model.Data.Target.IsUnknown()
		if isFullyKnown {
			record = &dns.NS{
				Ns: dns.Fqdn(model.Data.Target.ValueString()),
			}
		}

	case domain.RecordSMIMEA:
		isFullyKnown = !isAnyFieldUnknown(model.Data.Priority, model.Data.Selector, model.Data.MatchingType, model.Data.CertAssocData)
		if isFullyKnown {
			record = &dns.SMIMEA{
				Usage:        uint8(model.Data.Priority.ValueInt64()),
				Selector:     uint8(model.Data.Selector.ValueInt64()),
				MatchingType: uint8(model.Data.MatchingType.ValueInt64()),
				Certificate:  model.Data.CertAssocData.ValueString(),
			}
		}

	case domain.RecordSRV:
		isFullyKnown = !isAnyFieldUnknown(model.Data.Priority, model.Data.Weight, model.Data.Port, model.Data.Target)
		if isFullyKnown {
			record = &dns.SRV{
				Priority: uint16(model.Data.Priority.ValueInt64()),
				Weight:   uint16(model.Data.Weight.ValueInt64()),
				Port:     uint16(model.Data.Port.ValueInt64()),
				Target:   dns.Fqdn(model.Data.Target.ValueString()),
			}
		}

	case domain.RecordSSHFP:
		isFullyKnown = !isAnyFieldUnknown(model.Data.FingerprintAlgorithm, model.Data.FingerprintType, model.Data.Fingerprint)
		if isFullyKnown {
			record = &dns.SSHFP{
				Algorithm:   uint8(model.Data.FingerprintAlgorithm.ValueInt64()),
				Type:        uint8(model.Data.FingerprintType.ValueInt64()),
				FingerPrint: model.Data.Fingerprint.ValueString(),
			}
		}

	case domain.RecordTLSA:
		isFullyKnown = !isAnyFieldUnknown(model.Data.Priority, model.Data.Selector, model.Data.MatchingType, model.Data.CertAssocData)
		if isFullyKnown {
			record = &dns.TLSA{
				Usage:        uint8(model.Data.Priority.ValueInt64()),
				Selector:     uint8(model.Data.Selector.ValueInt64()),
				MatchingType: uint8(model.Data.MatchingType.ValueInt64()),
				Certificate:  model.Data.CertAssocData.ValueString(),
			}
		}

	case domain.RecordTXT:
		isFullyKnown = !model.Data.Value.IsUnknown()
		if isFullyKnown {
			record = &dns.TXT{
				Txt: []string{model.Data.Value.ValueString()},
			}
		}
	}

	// If any referenced value is still unknown, return unknown
	if !isFullyKnown || record == nil {
		return "", false
	}

	raw := strings.TrimPrefix(record.String(), record.Header().String())
	return raw, true
}

func (model *RecordModel) ParseRawTarget(raw string) error {
	// We need to prepend a fake name to make dns.NewRR happy
	full := fmt.Sprintf("example.com. 3600 IN %s %s", model.Type.ValueString(), raw)

	rr, err := dns.NewRR(full)
	if err != nil {
		return fmt.Errorf("failed to parse DNS record: %w", err)
	}

	switch v := rr.(type) {
	case *dns.A:
		model.Data.IP = types.StringValue(v.A.String())

	case *dns.AAAA:
		model.Data.IP = types.StringValue(v.AAAA.String())

	case *dns.CAA:
		model.Data.Flags = types.Int64Value(int64(v.Flag))
		model.Data.Tag = types.StringValue(v.Tag)
		model.Data.Value = types.StringValue(v.Value)

	case *dns.CNAME:
		model.Data.Target = types.StringValue(strings.TrimSuffix(v.Target, "."))

	case *dns.DNAME:
		model.Data.Target = types.StringValue(strings.TrimSuffix(v.Target, "."))

	case *dns.DS:
		model.Data.KeyTag = types.Int64Value(int64(v.KeyTag))
		model.Data.Algorithm = types.Int64Value(int64(v.Algorithm))
		model.Data.DigestType = types.Int64Value(int64(v.DigestType))
		model.Data.Digest = types.StringValue(v.Digest)

	case *dns.MX:
		model.Data.Priority = types.Int64Value(int64(v.Preference))
		model.Data.Target = types.StringValue(strings.TrimSuffix(v.Mx, "."))

	case *dns.NS:
		model.Data.Target = types.StringValue(strings.TrimSuffix(v.Ns, "."))

	case *dns.PTR:
		model.Data.Target = types.StringValue(strings.TrimSuffix(v.Ptr, "."))

	case *dns.SMIMEA:
		model.Data.Priority = types.Int64Value(int64(v.Usage))
		model.Data.Selector = types.Int64Value(int64(v.Selector))
		model.Data.MatchingType = types.Int64Value(int64(v.MatchingType))
		model.Data.CertAssocData = types.StringValue(v.Certificate)

	case *dns.SRV:
		model.Data.Priority = types.Int64Value(int64(v.Priority))
		model.Data.Weight = types.Int64Value(int64(v.Weight))
		model.Data.Port = types.Int64Value(int64(v.Port))
		model.Data.Target = types.StringValue(strings.TrimSuffix(v.Target, "."))

	case *dns.SSHFP:
		model.Data.FingerprintAlgorithm = types.Int64Value(int64(v.Algorithm))
		model.Data.FingerprintType = types.Int64Value(int64(v.Type))
		model.Data.Fingerprint = types.StringValue(v.FingerPrint)

	case *dns.TLSA:
		model.Data.Priority = types.Int64Value(int64(v.Usage))
		model.Data.Selector = types.Int64Value(int64(v.Selector))
		model.Data.MatchingType = types.Int64Value(int64(v.MatchingType))
		model.Data.CertAssocData = types.StringValue(v.Certificate)

	case *dns.TXT:
		if len(v.Txt) > 0 {
			model.Data.Value = types.StringValue(v.Txt[0])
		}

	default:
		return fmt.Errorf("unsupported record type: %T", rr)
	}

	return nil
}
