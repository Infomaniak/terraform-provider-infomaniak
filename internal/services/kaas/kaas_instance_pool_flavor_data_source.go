package kaas

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/kaas"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &kaasInstancePoolFlavorDataSource{}
	_ datasource.DataSourceWithConfigure = &kaasInstancePoolFlavorDataSource{}
)

type kaasInstancePoolFlavorDataSource struct {
	client *apis.Client
}

type KaasInstancePoolFlavorModel struct {
	PublicCloudId        types.Int64   `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64   `tfsdk:"public_cloud_project_id"`
	Region               types.String  `tfsdk:"region"`
	Name                 types.String  `tfsdk:"name"`
	Cpu                  types.Int64   `tfsdk:"cpu"`
	Ram                  types.Int64   `tfsdk:"ram"`
	Storage              types.Int64   `tfsdk:"storage"`
	IsAvailable          types.Bool    `tfsdk:"is_available"`
	IsMemoryOptimized    types.Bool    `tfsdk:"is_memory_optimized"`
	IsIopsOptimized      types.Bool    `tfsdk:"is_iops_optimized"`
	IsGpuOptimized       types.Bool    `tfsdk:"is_gpu_optimized"`
	Rates                *PricingModel `tfsdk:"rates"`
}

type PricingModel struct {
	HourExclTax types.Float64 `tfsdk:"hour_excl_tax"`
	HourInclTax types.Float64 `tfsdk:"hour_incl_tax"`
}

// NewKaasInstancePoolFlavorDataSource is a helper function to simplify the provider implementation.
func NewKaasInstancePoolFlavorDataSource() datasource.DataSource {
	return &kaasInstancePoolFlavorDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *kaasInstancePoolFlavorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kaasInstancePoolFlavorDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getKaasInstancePoolFlavorDataSourceSchema()
}

// Read refreshes the Terraform state with the latest data.
func (d *kaasInstancePoolFlavorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KaasInstancePoolFlavorModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	obj, err := d.client.Kaas.GetFlavor(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.Region.ValueString(),
		kaas.KaasFlavorLookupParameters{
			Name:              data.Name.ValueStringPointer(),
			Cpu:               data.Cpu.ValueInt64Pointer(),
			Ram:               data.Ram.ValueInt64Pointer(),
			Storage:           data.Storage.ValueInt64Pointer(),
			IsMemoryOptimized: data.IsMemoryOptimized.ValueBoolPointer(),
			IsIopsOptimized:   data.IsIopsOptimized.ValueBoolPointer(),
			IsGpuOptimized:    data.IsGpuOptimized.ValueBoolPointer(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find KaaS instance pool flavor",
			err.Error(),
		)
		return
	}

	data.Name = types.StringValue((obj.Name))
	data.Cpu = types.Int64Value(obj.Cpu)
	data.Ram = types.Int64Value(obj.Ram)
	data.Storage = types.Int64Value(obj.Storage)
	data.IsAvailable = types.BoolValue(obj.IsAvailable)
	data.IsMemoryOptimized = types.BoolValue(obj.IsMemoryOptimized)
	data.IsIopsOptimized = types.BoolValue(obj.IsIopsOptimized)
	data.IsGpuOptimized = types.BoolValue(obj.IsGpuOptimized)
	data.Rates = &PricingModel{
		HourExclTax: types.Float64Value(obj.Rates.HourlyExcludingTaxes),
		HourInclTax: types.Float64Value(obj.Rates.HourlyIncludingTaxes),
	}

	// Set state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *kaasInstancePoolFlavorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas_instance_pool_flavor"
}
