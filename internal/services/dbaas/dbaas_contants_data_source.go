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
	_ datasource.DataSource              = &dbaasConstsDataSource{}
	_ datasource.DataSourceWithConfigure = &dbaasConstsDataSource{}
)

type dbaasConstsDataSource struct {
	client *apis.Client
}

// NewDBaasDataSource is a helper function to simplify the provider implementation.
func NewDBaasConstsDataSource() datasource.DataSource {
	return &dbaasConstsDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *dbaasConstsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type DBaasConstsDataModel struct {
	Regions types.List   `tfsdk:"regions"`
	Types   []DBaasTypes `tfsdk:"types"`
}

type DBaasTypes struct {
	Name     types.String `tfsdk:"name"`
	Versions types.List   `tfsdk:"versions"`
}

// Schema defines the schema for the data source.
func (d *dbaasConstsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"types": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"versions": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
				Computed: true,
			},
			"regions": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dbaasConstsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DBaasConstsDataModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	regions, err := d.client.DBaas.GetDbaasRegions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find DBaaS",
			err.Error(),
		)
		return
	}

	tfregions, diags := types.ListValueFrom(ctx, types.StringType, regions)
	resp.Diagnostics.Append(diags...)
	data.Regions = tfregions

	dbaasTypes, err := d.client.DBaas.GetDbaasTypes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find DBaaS",
			err.Error(),
		)
		return
	}

	var dbaasPacks []DBaasTypes
	for _, dbType := range dbaasTypes {
		versioned, diags := types.ListValueFrom(ctx, types.StringType, dbType.Versions)
		resp.Diagnostics.Append(diags...)
		dbaasPacks = append(dbaasPacks, DBaasTypes{
			Name:     types.StringValue(dbType.Name),
			Versions: versioned,
		})
	}
	data.Types = dbaasPacks

	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *dbaasConstsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_constants"
}
