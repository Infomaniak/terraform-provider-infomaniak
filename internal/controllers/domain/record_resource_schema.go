package domain

import (
	"terraform-provider-infomaniak/internal/apis/domain"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func getRecordResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"zone_fqdn": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The FQDN of the zone where the record should be put in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The source of the Record.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of the Record.",
				Validators: []validator.String{
					stringvalidator.OneOf(domain.RecordTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the Record.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ttl": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The TTL of the Record.",
				Default:             int64default.StaticInt64(3600),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"computed_target": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The computed target of the Record.",
			},
			"target": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The target of the Record.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"data": schema.SingleNestedAttribute{
				Description: "Components of a DNS record.",
				Optional:    true,
				Validators: []validator.Object{
					objectvalidator.All(
						objectvalidator.ConflictsWith(path.MatchRoot("target")),
					),
				},
				Attributes: map[string]schema.Attribute{
					"ip": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "IP for the record",
					},
					// For MX, SRV, TLSA, SMIMEA, SSHFP
					"priority": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The priority/usage/weight of the Record (MX, SRV, TLSA, SMIMEA).",
					},
					// For SRV
					"weight": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The weight of the Record (SRV).",
					},
					"port": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The port of the Record (SRV).",
					},

					// For CAA
					"flags": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The flags of the Record (CAA).",
					},
					"tag": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The tag of the Record (CAA).",
					},

					// For DNSKEY
					"algorithm": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The algorithm of the Record (DNSKEY, DS, SSHFP).",
					},

					// For DS
					"key_tag": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The Key Tag of the Record (DS).",
					},
					"digest_type": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The digest type of the Record (DS).",
					},
					"digest": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The digest of the Record (DS).",
					},

					// For TLSA / SMIMEA
					"selector": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The selector of the Record (TLSA, SMIMEA).",
					},
					"matching_type": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The matching type of the Record (TLSA, SMIMEA).",
					},
					"cert_assoc_data": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The certificate association data (TLSA, SMIMEA).",
					},

					// For SSHFP
					"fingerprint_algorithm": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The algorithm of the Record (DNSKEY, DS, SSHFP).",
					},
					"fingerprint_type": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The fingerprint type of the Record (SSHFP).",
					},
					"fingerprint": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The fingerprint of the Record (SSHFP).",
					},

					"target": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The target of the Record (MX, CNAME, DNAME, NS, PTR, etc).",
					},
					// For generic text value (e.g. TXT, CAA value, etc)
					"value": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The value of the Record (TXT, CAA, etc).",
					},
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The id of the Record.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
		MarkdownDescription: "The record resource allows the user to manage a record inside a zone of a domain",
	}
}
