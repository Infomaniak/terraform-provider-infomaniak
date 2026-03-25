package kaas_schemas

import (
	"terraform-provider-infomaniak/internal/apis/kaas"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KaasInstancePoolModel struct {
	PublicCloudId        types.Int64 `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64 `tfsdk:"public_cloud_project_id"`
	KaasId               types.Int64 `tfsdk:"kaas_id"`
	Id                   types.Int64 `tfsdk:"id"`

	Name             types.String `tfsdk:"name"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
	FlavorName       types.String `tfsdk:"flavor_name"`
	MinInstances     types.Int64  `tfsdk:"min_instances"`
	MaxInstances     types.Int64  `tfsdk:"max_instances"`
	Labels           types.Map    `tfsdk:"labels"`
}

func (model *KaasInstancePoolModel) Fill(instancePool *kaas.InstancePool) {
	model.Id = types.Int64Value(instancePool.Id)
	model.Name = types.StringValue(instancePool.Name)
	model.FlavorName = types.StringValue(instancePool.FlavorName)
	model.MinInstances = types.Int64Value(instancePool.MinInstances)
	model.MaxInstances = types.Int64Value(instancePool.MaxInstances)
	model.AvailabilityZone = types.StringValue(instancePool.AvailabilityZone)
}
