package domain

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &zoneResource{}
	_ resource.ResourceWithConfigure   = &zoneResource{}
	_ resource.ResourceWithImportState = &zoneResource{}
)

func NewZoneResource() resource.Resource {
	return &zoneResource{}
}

type zoneResource struct {
	client *apis.Client
}

type ZoneModel struct {
	Fqdn types.String `tfsdk:"fqdn"`
	Id   types.Int64  `tfsdk:"id"`
}

func (r *zoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

// Configure adds the provider configured client to the data source.
func (r *zoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Zone Resource Configure Type",
			err.Error(),
		)
		return
	}

	r.client = client
}

func (r *zoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"fqdn": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The fqdn of the zone",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
		MarkdownDescription: "The Zone resource allows the user to manage a zone for a domain project",
	}
}

func (r *zoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ZoneModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// CreateZone API call logic
	zone, err := r.client.Domain.CreateZone(data.Fqdn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating Zone",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(int64(zone.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *zoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ZoneModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	zone, err := r.client.Domain.GetZone(state.Fqdn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading Zone",
			err.Error(),
		)
		return
	}

	state.Id = types.Int64Value(int64(zone.ID))
	state.Fqdn = types.StringValue(zone.FQDN)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *zoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state ZoneModel
	var data ZoneModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *zoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ZoneModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteZone API call logic
	_, err := r.client.Domain.DeleteZone(data.Fqdn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting Zone",
			err.Error(),
		)
		return
	}
}

func (r *zoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("fqdn"), req.ID)...)
}
