package dbaas

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dbaasPacksDataSource{}
	_ datasource.DataSourceWithConfigure = &dbaasPacksDataSource{}
)

type dbaasPacksDataSource struct {
	client *apis.Client
}

// NewDBaasDataSource is a helper function to simplify the provider implementation.
func NewDBaasPacksDataSource() datasource.DataSource {
	return &dbaasPacksDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *dbaasPacksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type DBaasConstantsDataModel struct {
	Type  types.String      `tfsdk:"type"`
	Packs []DBaasPacksModel `tfsdk:"packs"`
}

type DBaasPacksModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Type      types.String `tfsdk:"type"`
	Group     types.String `tfsdk:"group"`
	Name      types.String `tfsdk:"name"`
	Instances types.Int64  `tfsdk:"instances"`
	CPU       types.Int64  `tfsdk:"cpu"`
	RAM       types.Int64  `tfsdk:"ram"`
	Storage   types.Int64  `tfsdk:"storage"`
	Rates     RatesModel   `tfsdk:"rates"`
}

type RatesModel struct {
	CHF PricingModel `tfsdk:"chf"`
	EUR PricingModel `tfsdk:"eur"`
}

type PricingModel struct {
	HourlyExcludingTaxes types.Float64 `tfsdk:"hour_excl_tax"`
	HourlyIncludingTaxes types.Float64 `tfsdk:"hour_incl_tax"`
}

// Schema defines the schema for the data source.
func (d *dbaasPacksDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	pricingObject := schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"hour_excl_tax": schema.Float64Attribute{
				Computed: true,
			},
			"hour_incl_tax": schema.Float64Attribute{
				Computed: true,
			},
		},
	}
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Required: true,
			},
			"packs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"group": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"instances": schema.Int64Attribute{
							Computed: true,
						},
						"cpu": schema.Int64Attribute{
							Computed: true,
						},
						"ram": schema.Int64Attribute{
							Computed: true,
						},
						"storage": schema.Int64Attribute{
							Computed: true,
						},
						"rates": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"chf": pricingObject,
								"eur": pricingObject,
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dbaasPacksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DBaasConstantsDataModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	packs, err := d.client.DBaas.GetDbaasPacks(data.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find DBaaS packs",
			err.Error(),
		)
		return
	}

	var tfPacks []DBaasPacksModel
	for _, pack := range packs {
		tfPacks = append(tfPacks, DBaasPacksModel{
			ID:        types.Int64Value(pack.ID),
			Type:      types.StringValue(pack.Type),
			Group:     types.StringValue(pack.Group),
			Name:      types.StringValue(pack.Name),
			Instances: types.Int64Value(pack.Instances),
			CPU:       types.Int64Value(pack.CPU),
			RAM:       types.Int64Value(pack.RAM),
			Storage:   types.Int64Value(pack.Storage),
			Rates: RatesModel{
				CHF: PricingModel{
					HourlyExcludingTaxes: types.Float64Value(pack.Rates.CHF.HourExclTax),
					HourlyIncludingTaxes: types.Float64Value(pack.Rates.CHF.HourInclTax),
				},
				EUR: PricingModel{
					HourlyExcludingTaxes: types.Float64Value(pack.Rates.EUR.HourExclTax),
					HourlyIncludingTaxes: types.Float64Value(pack.Rates.EUR.HourInclTax),
				},
			},
		})
	}

	data.Packs = tfPacks

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *dbaasPacksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_packs"
}
