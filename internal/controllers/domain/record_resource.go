package domain

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

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
	KeyTag               types.Int64  `tfsdk:"key_tag"`               // DS
	Algorithm            types.Int64  `tfsdk:"algorithm"`             // DNSKEY, DS, SSHFP, TLSA
	DigestType           types.Int64  `tfsdk:"digest_type"`           // DS, TLSA
	Digest               types.String `tfsdk:"digest"`                // DS, TLSA, SSHFP
	Selector             types.Int64  `tfsdk:"selector"`              // SMIMEA, TLSA
	MatchingType         types.Int64  `tfsdk:"matching_type"`         // SMIMEA, TLSA
	CertAssocData        types.String `tfsdk:"cert_assoc_data"`       // SMIMEA, TLSA
	Flags                types.Int64  `tfsdk:"flags"`                 // CAA, DNSKEY
	Tag                  types.String `tfsdk:"tag"`                   // CAA
	Value                types.String `tfsdk:"value"`                 // CAA, TXT
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
	resp.Schema = getRecordResourceSchema()
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

	computedTarget, isKnown := plan.ComputeRawTarget()
	plan.ComputedTarget = types.StringUnknown()
	if isKnown {
		plan.ComputedTarget = types.StringValue(computedTarget)
	}

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

	rawTarget, _ := data.ComputeRawTarget()

	record, err := r.client.Domain.CreateRecord(
		data.ZoneFqdn.ValueString(),
		data.Type.ValueString(),
		data.Source.ValueString(),
		rawTarget,
		data.TTL.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating Record",
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
			"Error when reading Record",
			err.Error(),
		)
		return
	}

	state.Id = types.Int64Value(int64(record.ID))
	state.TTL = types.Int64Value(int64(record.TTL))
	state.Source = types.StringValue(record.Source)
	state.Type = types.StringValue(record.Type)

	// If we have neither of them, we fill them with the API
	// However in this state (import), we can't know which field is planned by the user
	if state.Target.IsNull() && state.Data == nil {
		state.Target = types.StringValue(record.Target)
		state.Data = &RecordDataModel{}
		state.ParseRawTarget(record.Target)
	}

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

	rawTarget, _ := data.ComputeRawTarget()

	record, err := r.client.Domain.UpdateRecord(
		data.ZoneFqdn.ValueString(),
		state.Id.ValueInt64(),
		data.Type.ValueString(),
		data.Source.ValueString(),
		rawTarget,
		data.TTL.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating Record",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(int64(record.ID))
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

	_, err := r.client.Domain.DeleteRecord(
		data.ZoneFqdn.ValueString(),
		state.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting Record",
			err.Error(),
		)
		return
	}
}

func (r *recordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: zone_fqdn,id. Got: %q", req.ID),
		)
		return
	}

	var errorList error

	zoneFQDN := idParts[0]
	recordId, err := strconv.ParseInt(idParts[1], 10, 64)
	errorList = errors.Join(errorList, err)

	if errorList != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: zone_fqdn,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_fqdn"), zoneFQDN)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), recordId)...)
}
