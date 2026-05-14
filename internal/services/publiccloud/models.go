package publiccloud

import "github.com/hashicorp/terraform-plugin-framework/types"

// PublicCloudModel mirrors apis/publiccloud.PublicCloud for Terraform state.
type PublicCloudModel struct {
	Id                     types.Int64  `tfsdk:"id"`
	AccountId              types.Int64  `tfsdk:"account_id"`
	ServiceId              types.Int64  `tfsdk:"service_id"`
	ServiceName            types.String `tfsdk:"service_name"`
	CustomerName           types.String `tfsdk:"customer_name"`
	InternalName           types.String `tfsdk:"internal_name"`
	Description            types.String `tfsdk:"description"`
	BillReference          types.String `tfsdk:"bill_reference"`
	CreatedAt              types.Int64  `tfsdk:"created_at"`
	ExpiredAt              types.Int64  `tfsdk:"expired_at"`
	IsFree                 types.Bool   `tfsdk:"is_free"`
	IsZeroPrice            types.Bool   `tfsdk:"is_zero_price"`
	IsTrial                types.Bool   `tfsdk:"is_trial"`
	IsLocked               types.Bool   `tfsdk:"is_locked"`
	HasMaintenance         types.Bool   `tfsdk:"has_maintenance"`
	HasOperationInProgress types.Bool   `tfsdk:"has_operation_in_progress"`
}

// PublicCloudsModel is the data-source state for the list endpoint.
type PublicCloudsModel struct {
	AccountId    types.Int64        `tfsdk:"account_id"`
	PublicClouds []PublicCloudModel `tfsdk:"public_clouds"`
}

// PublicCloudConfigModel mirrors apis/publiccloud.Config.
type PublicCloudConfigModel struct {
	AccountId            types.Int64   `tfsdk:"account_id"`
	FreeTier             types.Float64 `tfsdk:"free_tier"`
	FreeTierUsed         types.Float64 `tfsdk:"free_tier_used"`
	AccountResourceLevel types.Int64   `tfsdk:"account_resource_level"`
	ProjectCount         types.Int64   `tfsdk:"project_count"`
	ValidFrom            types.Int64   `tfsdk:"valid_from"`
	ValidTo              types.Int64   `tfsdk:"valid_to"`
}

// PublicCloudAccessesModel exposes the maintenance status of the Public Cloud
// API surface for the given account.
type PublicCloudAccessesModel struct {
	AccountId            types.Int64 `tfsdk:"account_id"`
	IsMaintenanceOngoing types.Bool  `tfsdk:"is_maintenance_ongoing"`
}

// PublicCloudProjectModel mirrors apis/publiccloud.Project.
type PublicCloudProjectModel struct {
	PublicCloudId  types.Int64   `tfsdk:"public_cloud_id"`
	Id             types.Int64   `tfsdk:"id"`
	Name           types.String  `tfsdk:"name"`
	OpenStackName  types.String  `tfsdk:"open_stack_name"`
	Status         types.String  `tfsdk:"status"`
	Price          types.Float64 `tfsdk:"price"`
	ResourceLevel  types.Int64   `tfsdk:"resource_level"`
	UserCount      types.Int64   `tfsdk:"user_count"`
	CreatedAt      types.Int64   `tfsdk:"created_at"`
	UpdatedAt      types.Int64   `tfsdk:"updated_at"`
	BillingStartAt types.Int64   `tfsdk:"billing_start_at"`
	BillingEndAt   types.Int64   `tfsdk:"billing_end_at"`
	PriceUpdatedAt types.Int64   `tfsdk:"price_updated_at"`
}

// PublicCloudUserModel mirrors apis/publiccloud.User. The API never returns
// the password or email so they are absent on read-only flows.
type PublicCloudUserModel struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64  `tfsdk:"id"`
	OpenStackName        types.String `tfsdk:"open_stack_name"`
	Description          types.String `tfsdk:"description"`
	Status               types.String `tfsdk:"status"`
	CreatedAt            types.Int64  `tfsdk:"created_at"`
	UpdatedAt            types.Int64  `tfsdk:"updated_at"`
}

// PublicCloudOpenrcModel exposes the openrc.sh file content for a user.
type PublicCloudOpenrcModel struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	PublicCloudUserId    types.Int64  `tfsdk:"public_cloud_user_id"`
	Region               types.String `tfsdk:"region"`
	Content              types.String `tfsdk:"content"`
}

// PublicCloudUserAuthenticationModel exposes an authentication file (e.g.
// clouds.yaml) for a user.
type PublicCloudUserAuthenticationModel struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	PublicCloudUserId    types.Int64  `tfsdk:"public_cloud_user_id"`
	Type                 types.String `tfsdk:"type"`
	Region               types.String `tfsdk:"region"`
	Content              types.String `tfsdk:"content"`
}
