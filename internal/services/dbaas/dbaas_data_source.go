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

// Schema defines the schema for the data source.
func (d *dbaasDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud where DBaaS is installed",
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud project where DBaaS is installed",
			},
			"id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of this DBaaS",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the DBaaS project",
			},
			"pack_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the pack associated to the DBaaS project",
			},
			"region": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The region where the DBaaS project resides in.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the database associated with the DBaaS project",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The version of the database associated with the DBaaS project",
			},
			"allowedCIDRs": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Allowed to query Database IP whitelist",
			},
		},
		MarkdownDescription: "The dbaas data source allows the user to manage a dbaas project",
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dbaasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DBaasModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	obj, err := d.client.DBaas.GetDBaaS(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find DBaaS",
			err.Error(),
		)
		return
	}

	data.Region = types.StringValue(obj.Region)
	data.Version = types.StringValue(obj.Version)
	data.Type = types.StringValue(obj.Type)

	// Set state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *dbaasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas"
}
