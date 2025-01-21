package provider

import (
	"context"
	"fmt"
	"terraform-provider-infomaniak/internal/apis"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

	data, ok := req.ProviderData.(*IkProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *apis.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = data.Client
}

// Schema defines the schema for the data source.
func (d *kaasInstancePoolDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pcp_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the public cloud project where KaaS is installed",
			},
			"kaas_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the kaas project.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of this instance pool",
			},
			"flavor_name": schema.StringAttribute{
				Computed:    true,
				Description: "The flavor name of the instance in this instance pool",
			},
			"min_instances": schema.Int32Attribute{
				Computed:    true,
				Description: "The minimum amount of instances in the instance pool",
			},
			"max_instances": schema.Int32Attribute{
				Computed:    true,
				Description: "The maximum amount of instances in the instance pool",
			},
		},
		MarkdownDescription: "The kaas data source allows the user to manage a kaas project",
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kaasInstancePoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KaasInstancePoolModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	obj, err := d.client.Kaas.GetInstancePool(data.PcpId.ValueString(), data.KaasId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find KaaS instance pool",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(obj.Id)
	data.Name = types.StringValue(obj.Name)
	data.FlavorName = types.StringValue(obj.FlavorName)
	data.MinInstances = types.Int32Value(obj.MinInstances)
	data.MaxInstances = types.Int32Value(obj.MaxInstances)

	// Set state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *kaasInstancePoolDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas_instance_pool"
}
