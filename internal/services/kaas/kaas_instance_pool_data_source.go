package kaas

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &kaasInstancePoolDataSource{}
	_ datasource.DataSourceWithConfigure = &kaasInstancePoolDataSource{}
)

type kaasInstancePoolDataSource struct {
	client *apis.Client
}

// NewKaasInstancePoolDataSource is a helper function to simplify the provider implementation.
func NewKaasInstancePoolDataSource() datasource.DataSource {
	return &kaasInstancePoolDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *kaasInstancePoolDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Schema defines the schema for the data source.
func (d *kaasInstancePoolDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getKaasInstancePoolDataSourceSchema()
}

// Read refreshes the Terraform state with the latest data.
func (d *kaasInstancePoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KaasInstancePoolModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	obj, err := d.client.Kaas.GetInstancePool(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.KaasId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find KaaS instance pool",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(obj.Id)
	data.Name = types.StringValue(obj.Name)
	data.FlavorName = types.StringValue(obj.FlavorName)
	data.MinInstances = types.Int64Value(obj.MinInstances)
	data.MaxInstances = types.Int64Value(obj.MaxInstances)
	labels, diags := types.MapValueFrom(ctx, types.StringType, obj.Labels)
	resp.Diagnostics.Append(diags...)
	data.Labels = labels

	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *kaasInstancePoolDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas_instance_pool"
}
