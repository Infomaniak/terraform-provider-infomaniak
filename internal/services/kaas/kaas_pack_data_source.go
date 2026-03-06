package kaas

import (
	"context"
	"fmt"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/kaas"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &kaasPackDataSource{}
	_ datasource.DataSourceWithConfigure = &kaasPackDataSource{}
)

type kaasPackDataSource struct {
	client *apis.Client
}

type KaasPackModel struct {
	Id              types.Int64         `tfsdk:"id"`
	Name            types.String        `tfsdk:"name"`
	Description     types.String        `tfsdk:"description"`
	PricePerHour    *KaasPackRatesModel `tfsdk:"price_per_hour"`
	LimitPerProject types.Int64         `tfsdk:"limit_per_project"`
	IsActive        types.Bool          `tfsdk:"is_active"`
}

type KaasPackRatesModel struct {
	CHF types.Float64 `tfsdk:"chf"`
	EUR types.Float64 `tfsdk:"eur"`
}

// NewKaasPackDataSource is a helper function to simplify the provider implementation.
func NewKaasPackDataSource() datasource.DataSource {
	return &kaasPackDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *kaasPackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kaasPackDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getKaasPackDataSourceSchema()
}

// Read refreshes the Terraform state with the latest data.
func (d *kaasPackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KaasPackModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// error can be ignored as getPack already handles it for us
	obj, _ := d.getPack(data, &resp.Diagnostics)

	data.Id = types.Int64Value(int64(obj.Id))
	data.Name = types.StringValue((obj.Name))
	data.Description = types.StringValue(obj.Description)
	data.PricePerHour = &KaasPackRatesModel{
		CHF: types.Float64Value(obj.PricePerHour.CHF),
		EUR: types.Float64Value(obj.PricePerHour.EUR),
	}
	data.LimitPerProject = types.Int64Value(int64(obj.LimitPerProject))
	data.IsActive = types.BoolValue(obj.IsActive)

	// Set state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *kaasPackDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas_pack"
}

func (r *kaasPackDataSource) getPack(data KaasPackModel, diagnostic *diag.Diagnostics) (*kaas.KaasPack, error) {
	packs, err := r.client.Kaas.GetPacks()
	if err != nil {
		diagnostic.AddError(
			"Could not get KaaS Packs",
			err.Error(),
		)
		return nil, err
	}

	var chosenPack *kaas.KaasPack
	for _, pack := range packs {
		if pack.Name == data.Name.ValueString() {
			chosenPack = pack
			break
		}
	}

	if chosenPack == nil {
		var packNames []string
		for _, pack := range packs {
			packNames = append(packNames, pack.Name)
		}

		diagnostic.AddError(
			"Unknown KaaS Pack",
			fmt.Sprintf("pack_name must be one of : %v", packNames),
		)
		return nil, fmt.Errorf("pack name has not been found")
	}

	return chosenPack, nil
}
