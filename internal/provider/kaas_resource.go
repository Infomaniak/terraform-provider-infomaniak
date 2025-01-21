package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/kaas"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &kaasResource{}
	_ resource.ResourceWithConfigure   = &kaasResource{}
	_ resource.ResourceWithImportState = &kaasResource{}
)

func NewKaasResource() resource.Resource {
	return &kaasResource{}
}

type kaasResource struct {
	client *apis.Client
}

type KaasModel struct {
	PcpId types.String `tfsdk:"pcp_id"`
	Id    types.String `tfsdk:"id"`

	Region     types.String `tfsdk:"region"`
	Kubeconfig types.String `tfsdk:"kubeconfig"`
}

func (r *kaasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas"
}

// Configure adds the provider configured client to the data source.
func (r *kaasResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*IkProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *apis.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.Client
}

func (r *kaasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	var availablePacks []string

	packs, _ := r.client.Kaas.GetPacks()

	for _, pack := range packs {
		availablePacks = append(availablePacks, fmt.Sprint(pack.Name))
	}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pcp_id": schema.StringAttribute{
				Required:            true,
				Description:         "The id of the public cloud project where KaaS is installed",
				MarkdownDescription: "The id of the public cloud project where KaaS is installed",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pack": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the pack associated to the KaaS being installed",
				MarkdownDescription: "The name of the pack associated to the KaaS being installed",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(availablePacks...),
				},
			},
			"kubernetes_version": schema.StringAttribute{
				Required:            true,
				Description:         "The version of Kubernetes associated with the KaaS being installed",
				MarkdownDescription: "The version of Kubernetes associated with the KaaS being installed",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(availablePacks...),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				MarkdownDescription: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:            true,
				Description:         "The region where the KaaS will reside.",
				MarkdownDescription: "The region where the KaaS will reside.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kubeconfig": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "The kubeconfig generated to access to KaaS project",
				MarkdownDescription: "The kubeconfig generated to access to KaaS project",
			},
		},
		MarkdownDescription: "The kaas resource allows the user to manage a kaas project",
	}
}

func (r *kaasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := &kaas.Kaas{
		PcpId:  data.PcpId.ValueString(),
		Region: data.Region.ValueString(),
	}

	// CreateKaas API call logic
	obj, err := r.client.Kaas.CreateKaas(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating KaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(obj.Id)
	data.Kubeconfig = types.StringValue(obj.Kubeconfig)
	data.Region = types.StringValue(obj.Region)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KaasModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	obj, err := r.client.Kaas.GetKaas(data.PcpId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading KaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(obj.Id)
	data.Kubeconfig = types.StringValue(obj.Kubeconfig)
	data.Region = types.StringValue(obj.Region)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state KaasModel
	var data KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	input := &kaas.Kaas{
		PcpId: data.PcpId.ValueString(),
		Id:    state.Id.ValueString(),
	}

	obj, err := r.client.Kaas.UpdateKaas(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating KaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(obj.Id)
	data.Kubeconfig = types.StringValue(obj.Kubeconfig)
	data.Region = types.StringValue(obj.Region)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KaasModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteKaas API call logic
	err := r.client.Kaas.DeleteKaas(data.PcpId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting KaaS",
			err.Error(),
		)
		return
	}
}

func (r *kaasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: pcp_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pcp_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
