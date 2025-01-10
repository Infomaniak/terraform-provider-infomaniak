package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/kaas"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &kaasInstancePoolResource{}
	_ resource.ResourceWithConfigure = &kaasInstancePoolResource{}
)

func NewKaasInstancePoolResource() resource.Resource {
	return &kaasInstancePoolResource{}
}

type kaasInstancePoolResource struct {
	client *apis.Client
}

type KaasInstancePoolModel struct {
	PcpId  types.String `tfsdk:"pcp_id"`
	KaasId types.String `tfsdk:"kaas_id"`
	Id     types.String `tfsdk:"id"`

	Name         types.String `tfsdk:"name"`
	FlavorName   types.String `tfsdk:"flavor_name"`
	MinInstances types.Int32  `tfsdk:"min_instances"`
	MaxInstances types.Int32  `tfsdk:"max_instances"`
}

func (r *kaasInstancePoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas_instance_pool"
}

// Configure adds the provider configured client to the data source.
func (r *kaasInstancePoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *kaasInstancePoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"kaas_id": schema.StringAttribute{
				Required:            true,
				Description:         "The id of the kaas project.",
				MarkdownDescription: "The id of the kaas project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name",
				MarkdownDescription: "The name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"flavor_name": schema.StringAttribute{
				Required:            true,
				Description:         "The flavor name",
				MarkdownDescription: "The flavor name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"min_instances": schema.Int32Attribute{
				Required:            true,
				Description:         "The kubeconfig",
				MarkdownDescription: "The kubeconfig",
			},
			"max_instances": schema.Int32Attribute{
				Required:            true,
				Description:         "The kubeconfig",
				MarkdownDescription: "The kubeconfig",
			},
		},
		MarkdownDescription: "The kaas instance pool resource is used to manage instance pools inside a kaas project",
	}
}

func (r *kaasInstancePoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KaasInstancePoolModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := &kaas.InstancePool{
		PcpId:        data.PcpId.ValueString(),
		KaasId:       data.KaasId.ValueString(),
		Name:         data.Name.ValueString(),
		FlavorName:   data.FlavorName.ValueString(),
		MinInstances: data.MinInstances.ValueInt32(),
		MaxInstances: data.MaxInstances.ValueInt32(),
	}

	// CreateKaas API call logic
	obj, err := r.client.Kaas.CreateInstancePool(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating KaaS instance pool",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(obj.Id)
	data.Name = types.StringValue(obj.Name)
	data.FlavorName = types.StringValue(obj.FlavorName)
	data.MinInstances = types.Int32Value(obj.MinInstances)
	data.MaxInstances = types.Int32Value(obj.MaxInstances)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasInstancePoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KaasInstancePoolModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	obj, err := r.client.Kaas.GetInstancePool(data.PcpId.ValueString(), data.KaasId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading KaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(obj.Id)
	data.Name = types.StringValue(obj.Name)
	data.FlavorName = types.StringValue(obj.FlavorName)
	data.MinInstances = types.Int32Value(obj.MinInstances)
	data.MaxInstances = types.Int32Value(obj.MaxInstances)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasInstancePoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state KaasInstancePoolModel
	var data KaasInstancePoolModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	input := &kaas.InstancePool{
		PcpId:  data.PcpId.ValueString(),
		KaasId: data.KaasId.ValueString(),
		Id:     state.Id.ValueString(),

		Name:         data.Name.ValueString(),
		FlavorName:   data.FlavorName.ValueString(),
		MinInstances: data.MinInstances.ValueInt32(),
		MaxInstances: data.MaxInstances.ValueInt32(),
	}

	obj, err := r.client.Kaas.UpdateInstancePool(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating KaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(obj.Id)
	data.Name = types.StringValue(obj.Name)
	data.FlavorName = types.StringValue(obj.FlavorName)
	data.MinInstances = types.Int32Value(obj.MinInstances)
	data.MaxInstances = types.Int32Value(obj.MaxInstances)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasInstancePoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KaasInstancePoolModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteKaas API call logic
	err := r.client.Kaas.DeleteInstancePool(data.PcpId.ValueString(), data.KaasId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting KaaS",
			err.Error(),
		)
		return
	}
}

func (r *kaasInstancePoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: pcp_id,kaas_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pcp_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("kaas_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[2])...)
}
