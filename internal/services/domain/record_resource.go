package domain

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/domain"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &recordResource{}
	_ resource.ResourceWithConfigure   = &recordResource{}
	_ resource.ResourceWithImportState = &recordResource{}
	_ resource.ResourceWithModifyPlan  = &recordResource{}
)

func NewRecordResource() resource.Resource {
	return &recordResource{}
}

type recordResource struct {
	client *apis.Client
}

type RecordModel struct {
	ZoneFqdn types.String `tfsdk:"zone_fqdn"`
	Id       types.Int64  `tfsdk:"id"`

	Type           types.String     `tfsdk:"type"`
	Source         types.String     `tfsdk:"source"`
	ComputedTarget types.String     `tfsdk:"computed_target"`
	Target         types.String     `tfsdk:"target"`
	TTL            types.Int64      `tfsdk:"ttl"`
	Description    types.String     `tfsdk:"description"`
	Data           *RecordDataModel `tfsdk:"data"`
}

type RecordDataModel struct {
	IP                   types.String `tfsdk:"ip"`                    // A, AAAA
	Priority             types.Int64  `tfsdk:"priority"`              // MX, SRV
	Target               types.String `tfsdk:"target"`                // MX, SRV, CNAME, NS, PTR
	Weight               types.Int64  `tfsdk:"weight"`                // SRV
	Port                 types.Int64  `tfsdk:"port"`                  // SRV
	Algorithm            types.Int64  `tfsdk:"algorithm"`             // DNSKEY, DS, SSHFP, TLSA
	DigestType           types.Int64  `tfsdk:"digest_type"`           // DS, TLSA
	Digest               types.String `tfsdk:"digest"`                // DS, TLSA, SSHFP
	Selector             types.Int64  `tfsdk:"selector"`              // SMIMEA, TLSA
	MatchingType         types.Int64  `tfsdk:"matching_type"`         // SMIMEA, TLSA
	CertAssocData        types.String `tfsdk:"cert_assoc_data"`       // SMIMEA, TLSA
	Flags                types.Int64  `tfsdk:"flags"`                 // CAA, DNSKEY
	Tag                  types.String `tfsdk:"tag"`                   // CAA
	Value                types.String `tfsdk:"value"`                 // CAA, TXT
	Serial               types.Int64  `tfsdk:"serial"`                // SOA
	Refresh              types.Int64  `tfsdk:"refresh"`               // SOA
	Retry                types.Int64  `tfsdk:"retry"`                 // SOA
	Expire               types.Int64  `tfsdk:"expire"`                // SOA
	Minimum              types.Int64  `tfsdk:"minimum"`               // SOA
	MName                types.String `tfsdk:"mname"`                 // SOA
	RName                types.String `tfsdk:"rname"`                 // SOA
	PublicKey            types.String `tfsdk:"public_key"`            // DNSKEY
	Fingerprint          types.String `tfsdk:"fingerprint"`           // SSHFP
	FingerprintType      types.Int64  `tfsdk:"fingerprint_type"`      // SSHFP
	FingerprintAlgorithm types.Int64  `tfsdk:"fingerprint_algorithm"` // SSHFP
}

func (r *recordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_record"
}

// Configure adds the provider configured client to the data source.
func (r *recordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Record Resource Configure Type",
			err.Error(),
		)
		return
	}

	r.client = client
}

func (r *recordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
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
				Default:             int64default.StaticInt64(2500),
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
					"public_key": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The public key of the Record (DNSKEY).",
					},

					// For DS
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

					// For SOA
					"mname": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The primary nameserver for the SOA record.",
					},
					"rname": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The responsible party email for the SOA record.",
					},
					"serial": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The serial number for the SOA record.",
					},
					"refresh": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The refresh interval for the SOA record.",
					},
					"retry": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The retry interval for the SOA record.",
					},
					"expire": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The expire time for the SOA record.",
					},
					"minimum": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The minimum TTL for the SOA record.",
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

func (r *recordResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		// Handle destroy plan (optional)
		return
	}

	var plan RecordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	computedTarget := plan.ComputeRawTarget()
	plan.ComputedTarget = types.StringValue(computedTarget)

	// Set the modified plan back
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (r *recordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RecordModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rawTarget := data.ComputeRawTarget()

	// CreateKaas API call logic
	record, err := r.client.Domain.CreateRecord(
		data.ZoneFqdn.ValueString(),
		data.Type.ValueString(),
		data.Source.ValueString(),
		rawTarget,
		data.TTL.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating KaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(int64(record.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *recordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RecordModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	record, err := r.client.Domain.GetRecord(state.ZoneFqdn.ValueString(), state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading KaaS",
			err.Error(),
		)
		return
	}

	state.Id = types.Int64Value(int64(record.ID))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *recordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state RecordModel
	var data RecordModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *recordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RecordModel
	var data RecordModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteKaas API call logic
	_, err := r.client.Domain.DeleteRecord(
		data.ZoneFqdn.ValueString(),
		state.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting KaaS",
			err.Error(),
		)
		return
	}
}

func (r *recordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: public_cloud_id,public_cloud_project_id,id. Got: %q", req.ID),
		)
		return
	}

	var errorList error

	publicCloudId, err := strconv.ParseInt(idParts[0], 10, 64)
	errorList = errors.Join(errorList, err)
	publicCloudProjectId, err := strconv.ParseInt(idParts[1], 10, 64)
	errorList = errors.Join(errorList, err)
	kaasId, err := strconv.ParseInt(idParts[2], 10, 64)
	errorList = errors.Join(errorList, err)

	if errorList != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: public_cloud_id,public_cloud_project_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_id"), publicCloudId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_project_id"), publicCloudProjectId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), kaasId)...)
}
