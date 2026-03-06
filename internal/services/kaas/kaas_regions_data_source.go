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
	_ datasource.DataSource              = &kaasRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &kaasRegionsDataSource{}
)

type kaasRegionsDataSource struct {
	client *apis.Client
}

type KaasRegionModel struct {
	Items types.List `tfsdk:"items"`
}

// NewKaasRegionsDataSource is a helper function to simplify the provider implementation.
func NewKaasRegionsDataSource() datasource.DataSource {
	return &kaasRegionsDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *kaasRegionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kaasRegionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getKaasRegionDataSourceSchema()
}

// Read refreshes the Terraform state with the latest data.
func (d *kaasRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KaasRegionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	obj, err := d.client.Kaas.GetRegions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find KaaS regions",
			err.Error(),
		)
		return
	}

	items, diags := types.ListValueFrom(ctx, types.StringType, obj)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = items

	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *kaasRegionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas_regions"
}
