package kaas

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"
	kaas_schemas "terraform-provider-infomaniak/internal/schemas/kaas"
	"terraform-provider-infomaniak/internal/services"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &kaasDataSource{}
	_ datasource.DataSourceWithConfigure = &kaasDataSource{}
)

type kaasDataSource struct {
	client      *apis.Client
	kaasService *services.KaasService
}

// NewKaasDataSource is a helper function to simplify the provider implementation.
func NewKaasDataSource() datasource.DataSource {
	return &kaasDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *kaasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.kaasService = services.NewKaasService(d.client)
}

// Schema defines the schema for the data source.
func (d *kaasDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = kaas_schemas.KaasDataSourceSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *kaasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state kaas_schemas.KaasModel
	var data kaas_schemas.KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	kaasId := data.Id.ValueInt64()

	_, diags := d.kaasService.GetKaasOnceActive(ctx, data, kaasId, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, diags = d.kaasService.GetKubeconfig(data, kaasId, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.kaasService.ReadApiserverConfig(ctx, data, kaasId, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Metadata returns the data source type name.
func (d *kaasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas"
}
