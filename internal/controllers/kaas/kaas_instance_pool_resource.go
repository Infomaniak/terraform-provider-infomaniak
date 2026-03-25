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
	kaas_schemas "terraform-provider-infomaniak/internal/schemas/kaas"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

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

	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			err.Error(),
		)
		return
	}

	r.client = client
}

func (r *kaasInstancePoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = kaas_schemas.KaasInstancePoolResourceSchema
}

func (r *kaasInstancePoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data kaas_schemas.KaasInstancePoolModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := &kaas.InstancePool{
		KaasId:           data.KaasId.ValueInt64(),
		Name:             data.Name.ValueString(),
		AvailabilityZone: data.AvailabilityZone.ValueString(),
		FlavorName:       data.FlavorName.ValueString(),
		MinInstances:     data.MinInstances.ValueInt64(),
		MaxInstances:     data.MaxInstances.ValueInt64(),
		Labels:           r.getLabelsValues(data),
	}

	// CreateKaas API call logic
	instancePoolId, err := r.client.Kaas.CreateInstancePool(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		input,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating KaaS instance pool",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(instancePoolId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	isScalingDown := false
	instancePoolObject, err := r.waitUntilActive(ctx, data, instancePoolId, isScalingDown)
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

	data.Fill(instancePoolObject)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasInstancePoolResource) getLabelsValues(data kaas_schemas.KaasInstancePoolModel) map[string]string {
	labels := make(map[string]string)

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		for key, val := range data.Labels.Elements() {
			if strVal, ok := val.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				labels[key] = strVal.ValueString()
			}
		}
	}

	return labels
}

func (r *kaasInstancePoolResource) waitUntilActive(ctx context.Context, data kaas_schemas.KaasInstancePoolModel, id int64, scalingDown bool) (*kaas.InstancePool, error) {
	scaleDownFailedQuotaCount := 0
	scaleDownFailedQuotaAllowedRetrys := 5
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			found, err := r.client.Kaas.GetInstancePool(
				data.PublicCloudId.ValueInt64(),
				data.PublicCloudProjectId.ValueInt64(),
				data.KaasId.ValueInt64(),
				id,
			)
			if err != nil {
				return nil, err
			}

			if len(found.ErrorMessages) > 0 {
				// Special case when we hit quota failure but we are scaling down. OpenStack can take some time to update so we let him do his work
				if (found.Status == "ScalingDown" || scalingDown) && scaleDownFailedQuotaCount <= scaleDownFailedQuotaAllowedRetrys {
					scaleDownFailedQuotaCount++
					continue
				}
				return nil, errors.New(strings.Join(found.ErrorMessages, ","))
			}

			// We need the instance pool to be active, have the same state as us, be scaled properly and be in bound of the autoscaling
			isActive := found.Status == "Active"
			isEquivalent := found.MinInstances == data.MinInstances.ValueInt64()
			isScaledProperly := found.AvailableInstances == found.TargetInstances
			isInBound := found.MinInstances <= found.TargetInstances && found.TargetInstances <= found.MaxInstances
			if isActive && isEquivalent && isScaledProperly && isInBound {
				return found, nil
			}
		}
	}
}

func (r *kaasInstancePoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data kaas_schemas.KaasInstancePoolModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	obj, err := r.client.Kaas.GetInstancePool(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.KaasId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading KaaS Instance Pool",
			err.Error(),
		)
		return
	}

	if len(obj.ErrorMessages) > 0 {
		resp.Diagnostics.AddWarning(
			"KaaS was in error state:",
			strings.Join(obj.ErrorMessages, ","),
		)
	}

	data.Fill(obj)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasInstancePoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state kaas_schemas.KaasInstancePoolModel
	var data kaas_schemas.KaasInstancePoolModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	input := &kaas.InstancePool{
		KaasId: data.KaasId.ValueInt64(),
		Id:     state.Id.ValueInt64(),

		Name:         data.Name.ValueString(),
		FlavorName:   data.FlavorName.ValueString(),
		MinInstances: data.MinInstances.ValueInt64(),
		MaxInstances: data.MaxInstances.ValueInt64(),
		Labels:       r.getLabelsValues(data),
	}

	_, err := r.client.Kaas.UpdateInstancePool(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		input,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating KaaS Instance Pool",
			err.Error(),
		)
		return
	}

	scalingDown := data.MaxInstances.ValueInt64() < state.MaxInstances.ValueInt64()
	instancePoolObject, err := r.waitUntilActive(ctx, data, state.Id.ValueInt64(), scalingDown)
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

	data.Fill(instancePoolObject)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasInstancePoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data kaas_schemas.KaasInstancePoolModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteKaas API call logic
	_, err := r.client.Kaas.DeleteInstancePool(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.KaasId.ValueInt64(),
		data.Id.ValueInt64(),
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
