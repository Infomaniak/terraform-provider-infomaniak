package dbaas

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dbaasPackDataSource{}
	_ datasource.DataSourceWithConfigure = &dbaasPackDataSource{}
)

type dbaasPackDataSource struct {
	client *apis.Client
}

// NewDBaasDataSource is a helper function to simplify the provider implementation.
func NewDBaasPackDataSource() datasource.DataSource {
	return &dbaasPackDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *dbaasPackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type DBaasPackDataModel struct {
	Type      types.String `tfsdk:"type"`
	ID        types.Int64  `tfsdk:"id"`
	Group     types.String `tfsdk:"group"`
	Name      types.String `tfsdk:"name"`
	Instances types.Int64  `tfsdk:"instances"`
	CPU       types.Int64  `tfsdk:"cpu"`
	RAM       types.Int64  `tfsdk:"ram"`
	Storage   types.Int64  `tfsdk:"storage"`
	Rates     *RatesModel  `tfsdk:"rates"`
}

type RatesModel struct {
	CHF *PricingModel `tfsdk:"chf"`
	EUR *PricingModel `tfsdk:"eur"`
}

type PricingModel struct {
	HourlyExcludingTaxes types.Float64 `tfsdk:"hour_excl_tax"`
	HourlyIncludingTaxes types.Float64 `tfsdk:"hour_incl_tax"`
}

// Schema defines the schema for the data source.
func (d *dbaasPackDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getDbaasPacksDataSourceSchema()
}

// Read refreshes the Terraform state with the latest data.
func (d *dbaasPackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DBaasPackDataModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	pack, err := d.client.DBaas.GetDbaasPack(dbaas.PackFilter{
		DbType:    data.Type.ValueString(),
		Group:     data.Group.ValueStringPointer(),
		Name:      data.Name.ValueStringPointer(),
		Instances: data.Instances.ValueInt64Pointer(),
		Cpu:       data.CPU.ValueInt64Pointer(),
		Ram:       data.RAM.ValueInt64Pointer(),
		Storage:   data.Storage.ValueInt64Pointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find DBaaS packs",
			err.Error(),
		)
		return
	}

	data.Type = types.StringValue(pack.Type)
	data.Group = types.StringValue(pack.Group)
	data.ID = types.Int64Value(pack.ID)
	data.Name = types.StringValue(pack.Name)
	data.Instances = types.Int64Value(pack.Instances)
	data.CPU = types.Int64Value(pack.CPU)
	data.RAM = types.Int64Value(pack.RAM)
	data.Storage = types.Int64Value(pack.Storage)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *dbaasPackDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_pack"
}
