package dbaas

import (
	"context"
	"encoding/json"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"terraform-provider-infomaniak/internal/dynamic"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dbaasDataSource{}
	_ datasource.DataSourceWithConfigure = &dbaasDataSource{}
)

type dbaasDataSource struct {
	client *apis.Client
}

// NewDBaasDataSource is a helper function to simplify the provider implementation.
func NewDBaasDataSource() datasource.DataSource {
	return &dbaasDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *dbaasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			err.Error(),
		)
		return
	}

	d.client = client
}

type DBaasDataModel struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64  `tfsdk:"id"`
	KubernetesIdentifier types.String `tfsdk:"kube_identifier"`

	Name     types.String `tfsdk:"name"`
	PackName types.String `tfsdk:"pack_name"`
	Region   types.String `tfsdk:"region"`
	Type     types.String `tfsdk:"type"`
	Version  types.String `tfsdk:"version"`

	Host types.String `tfsdk:"host"`
	Port types.String `tfsdk:"port"`
	User types.String `tfsdk:"user"`
	Ca   types.String `tfsdk:"ca"`

	AllowedCIDRs types.List `tfsdk:"allowed_cidrs"`

	EffectiveConfiguration types.Dynamic `tfsdk:"effective_configuration"`
}

func (data *DBaasDataModel) fill(obj *dbaas.DBaaS) {
	data.Region = types.StringValue(obj.Region)
	data.Name = types.StringValue(obj.Name)
	data.PackName = types.StringValue(obj.Pack.Name)
	data.Region = types.StringValue(obj.Region)
	data.Type = types.StringValue(obj.Type)
	data.Version = types.StringValue(obj.Version)
	if obj.Connection != nil {
		data.Host = types.StringValue(obj.Connection.Host)
		data.Port = types.StringValue(obj.Connection.Port)
		data.User = types.StringValue(obj.Connection.User)
		data.Ca = types.StringValue(obj.Connection.Ca)
	}
	data.KubernetesIdentifier = types.StringValue(obj.KubernetesIdentifier)
}

// Schema defines the schema for the data source.
func (d *dbaasDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getDbaasDataSourceSchema()
}

// Read refreshes the Terraform state with the latest data.
func (d *dbaasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DBaasDataModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	obj, err := d.client.DBaas.GetDBaaS(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find DBaaS",
			err.Error(),
		)
		return
	}

	newEffectiveConfig, diags := d.refreshEffectiveConfiguration(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.EffectiveConfiguration = newEffectiveConfig
	data.fill(obj)

	filteredIps, err := d.client.DBaas.GetIpFilters(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading DBaaS filtered IPs",
			err.Error(),
		)
		return
	}

	listFilteredIps, diags := types.ListValueFrom(ctx, types.StringType, filteredIps)
	data.AllowedCIDRs = listFilteredIps
	resp.Diagnostics.Append(diags...)

	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *dbaasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas"
}

func (r *dbaasDataSource) refreshEffectiveConfiguration(publicCloudId, publicCloudProjectId, id int64) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics
	effectiveSettings, err := r.client.DBaas.GetConfiguration(
		publicCloudId,
		publicCloudProjectId,
		id,
	)
	if err != nil {
		diags.AddError(
			"Error when reading DBaaS settings",
			err.Error(),
		)
		return types.DynamicNull(), diags
	}

	jsonEffectiveSettigs, err := json.Marshal(effectiveSettings)
	if err != nil {
		diags.AddError("could not marshall", "effective settings json marshall fail")
	}
	dynamicObj, err := dynamic.FromJSONImplied(jsonEffectiveSettigs)
	if err != nil {
		diags.AddError("could not create dynamic object", "effective settings dynamic object from json creation failure")
	}

	return dynamicObj, diags
}
