package kaas

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/kaas"
	"terraform-provider-infomaniak/internal/provider"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	PublicCloudId        types.Int64 `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64 `tfsdk:"public_cloud_project_id"`
	KaasId               types.Int64 `tfsdk:"kaas_id"`
	Id                   types.Int64 `tfsdk:"id"`

	Name             types.String `tfsdk:"name"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
	FlavorName       types.String `tfsdk:"flavor_name"`
	MinInstances     types.Int32  `tfsdk:"min_instances"`
	// MaxInstances types.Int32  `tfsdk:"max_instances"`
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

	data, ok := req.ProviderData.(*provider.IkProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *apis.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = apis.NewClient(data.Data.Host.ValueString(), data.Data.Token.ValueString())
	if data.Version.ValueString() == "test" {
		r.client = apis.NewMockClient()
	}
}

func (r *kaasInstancePoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud where KaaS is installed",
				MarkdownDescription: "The id of the public cloud where KaaS is installed",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud project where KaaS is installed",
				MarkdownDescription: "The id of the public cloud project where KaaS is installed",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"kaas_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the kaas project.",
				MarkdownDescription: "The id of the kaas project.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				MarkdownDescription: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the instance pool",
				MarkdownDescription: "The name of the instance pool",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"availability_zone": schema.StringAttribute{
				Required:            true,
				Description:         "The availability zone for the instances in the pool",
				MarkdownDescription: "The availability zone for the instances in the pool",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"flavor_name": schema.StringAttribute{
				Required:            true,
				Description:         "The flavor name for the instances in the pool",
				MarkdownDescription: "The flavor name for the instances in the pool",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"min_instances": schema.Int32Attribute{
				Required:            true,
				Description:         "The minimum instances in this instance pool (should be equal to max_instance until the AutoScaling feature is released)",
				MarkdownDescription: "The minimum instances in this instance pool (should be equal to max_instance until the AutoScaling feature is released)",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			// "max_instances": schema.Int32Attribute{
			// 	Required:            true,
			// 	Description:         "The maximum instances in this instance pool (should be equal to min_instance until the AutoScaling feature is released)",
			// 	MarkdownDescription: "The maximum instances in this instance pool (should be equal to min_instance until the AutoScaling feature is released)",
			// },
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
		KaasId:           int(data.KaasId.ValueInt64()),
		Name:             data.Name.ValueString(),
		AvailabilityZone: data.AvailabilityZone.ValueString(),
		FlavorName:       data.FlavorName.ValueString(),
		MinInstances:     data.MinInstances.ValueInt32(),
		// MaxInstances: data.MaxInstances.ValueInt32(),
	}

	// CreateKaas API call logic
	instancePoolId, err := r.client.Kaas.CreateInstancePool(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		input,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating KaaS instance pool",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(int64(instancePoolId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	instancePoolObject, err := r.waitUntilActive(ctx, data, instancePoolId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when waiting for KaaS Instance Pool to be Active",
			err.Error(),
		)
		return
	}

	if instancePoolObject == nil {
		return
	}

	data.fill(instancePoolObject)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasInstancePoolResource) waitUntilActive(ctx context.Context, data KaasInstancePoolModel, id int) (*kaas.InstancePool, error) {
	for {
		found, err := r.client.Kaas.GetInstancePool(
			int(data.PublicCloudId.ValueInt64()),
			int(data.PublicCloudProjectId.ValueInt64()),
			int(data.KaasId.ValueInt64()),
			id,
		)
		if err != nil {
			return nil, err
		}

		if ctx.Err() != nil {
			return nil, nil
		}

		// We need the instance pool to be active, have the same state as us, be scaled properly and be in bound of the autoscaling
		isActive := found.Status == "Active"
		isEquivalent := found.MinInstances == data.MinInstances.ValueInt32()
		isScaledProperly := found.AvailableInstances == found.TargetInstances
		isInBound := found.MinInstances <= found.TargetInstances && found.TargetInstances <= found.MaxInstances
		if isActive && isEquivalent && isScaledProperly && isInBound {
			return found, nil
		}

		time.Sleep(5 * time.Second)
	}
}

func (r *kaasInstancePoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KaasInstancePoolModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	obj, err := r.client.Kaas.GetInstancePool(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.KaasId.ValueInt64()),
		int(data.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading KaaS Instance Pool",
			err.Error(),
		)
		return
	}

	data.fill(obj)

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
		KaasId: int(data.KaasId.ValueInt64()),
		Id:     int(state.Id.ValueInt64()),

		Name:         data.Name.ValueString(),
		FlavorName:   data.FlavorName.ValueString(),
		MinInstances: data.MinInstances.ValueInt32(),
		// MaxInstances: data.MaxInstances.ValueInt32(),
	}

	_, err := r.client.Kaas.UpdateInstancePool(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		input,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating KaaS Instance Pool",
			err.Error(),
		)
		return
	}

	instancePoolObject, err := r.waitUntilActive(ctx, data, int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when waiting for KaaS Instance Pool to be Active",
			err.Error(),
		)
		return
	}

	if instancePoolObject == nil {
		return
	}

	data.fill(instancePoolObject)

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
	_, err := r.client.Kaas.DeleteInstancePool(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.KaasId.ValueInt64()),
		int(data.Id.ValueInt64()),
	)
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

	if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: public_cloud_id,public_cloud_project_id,kaas_id,id. Got: %q", req.ID),
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
	instancePoolId, err := strconv.ParseInt(idParts[3], 10, 64)
	errorList = errors.Join(errorList, err)

	if errorList != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: public_cloud_id,public_cloud_project_id,kaas_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_id"), publicCloudId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_project_id"), publicCloudProjectId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("kaas_id"), kaasId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), instancePoolId)...)
}

func (model *KaasInstancePoolModel) fill(instancePool *kaas.InstancePool) {
	model.Id = types.Int64Value(int64(instancePool.Id))
	model.Name = types.StringValue(instancePool.Name)
	model.FlavorName = types.StringValue(instancePool.FlavorName)
	model.MinInstances = types.Int32Value(instancePool.MinInstances)
	model.AvailabilityZone = types.StringValue(instancePool.AvailabilityZone)
	// data.MaxInstances = types.Int32Value(obj.MaxInstances)
}
