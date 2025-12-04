package dbaasmigration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DBaasModelV0 struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64  `tfsdk:"id"`
	KubernetesIdentifier types.String `tfsdk:"kube_identifier"`

	Name     types.String `tfsdk:"name"`
	PackName types.String `tfsdk:"pack_name"`
	Region   types.String `tfsdk:"region"`
	Type     types.String `tfsdk:"type"`
	Version  types.String `tfsdk:"version"`

	Host     types.String `tfsdk:"host"`
	Port     types.String `tfsdk:"port"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Ca       types.String `tfsdk:"ca"`

	AllowedCIDRs types.List `tfsdk:"allowed_cidrs"`

	Configuration          types.Map `tfsdk:"configuration"`
	EffectiveConfiguration types.Map `tfsdk:"effective_configuration"`
}

type DBaasModelV1 struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64  `tfsdk:"id"`
	KubernetesIdentifier types.String `tfsdk:"kube_identifier"`

	Name     types.String `tfsdk:"name"`
	PackName types.String `tfsdk:"pack_name"`
	Region   types.String `tfsdk:"region"`
	Type     types.String `tfsdk:"type"`
	Version  types.String `tfsdk:"version"`

	Host     types.String `tfsdk:"host"`
	Port     types.String `tfsdk:"port"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Ca       types.String `tfsdk:"ca"`

	AllowedCIDRs types.List `tfsdk:"allowed_cidrs"`

	Configuration          types.Dynamic `tfsdk:"configuration"`
	EffectiveConfiguration types.Dynamic `tfsdk:"effective_configuration"`
}

func GetV0Schema() *schema.Schema {
	return &schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud project",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the DBaaS instance.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pack_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the pack associated to the DBaaS project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of database associated with the DBaaS being installed",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The version of database associated with the DBaaS being installed",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the DBaaS project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The region where the DBaaS will reside.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The host to access this database.",
			},
			"port": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The port to access this database.",
			},
			"user": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The username to access this database.",
			},
			"password": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The password to access this database.",
			},
			"ca": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Database CA Certificate",
			},
			"allowed_cidrs": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Allowed to query Database IP whitelist",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"kube_identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DbaaS kubernetes name",
			},
			"configuration": schema.MapAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"effective_configuration": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
		MarkdownDescription: "The dbaas resource allows the user to manage a dbaas project",
	}
}

func StateUpgrader(ctx context.Context, request resource.UpgradeStateRequest, response *resource.UpgradeStateResponse) {
	var state DBaasModelV0

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	var newState DBaasModelV1

	newState.PublicCloudId = state.PublicCloudId
	newState.PublicCloudProjectId = state.PublicCloudProjectId
	newState.Id = state.Id
	newState.KubernetesIdentifier = state.KubernetesIdentifier

	newState.Name = state.Name
	newState.PackName = state.PackName
	newState.Region = state.Region
	newState.Type = state.Type
	newState.Version = state.Version

	newState.Host = state.Host
	newState.Port = state.Port
	newState.User = state.User
	newState.Password = state.Password
	newState.Ca = state.Ca

	newState.AllowedCIDRs = state.AllowedCIDRs

	newState.EffectiveConfiguration = types.DynamicNull()
	if state.Configuration.IsUnknown() {
		newState.Configuration = types.DynamicUnknown()
		return
	}
	newState.Configuration = types.DynamicNull()
	response.Diagnostics.Append(response.State.Set(ctx, &newState)...)
}
