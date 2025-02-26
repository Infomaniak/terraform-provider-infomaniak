package kaas

import (
	"context"
	"fmt"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &kaasDataSource{}
	_ datasource.DataSourceWithConfigure = &kaasDataSource{}
)

type kaasDataSource struct {
	client *apis.Client
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

	data, ok := req.ProviderData.(*provider.IkProviderData)
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
func (d *kaasDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud where KaaS is installed",
				MarkdownDescription: "The id of the public cloud where KaaS is installed",
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud project where KaaS is installed",
				MarkdownDescription: "The id of the public cloud project where KaaS is installed",
			},
			"id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of this KaaS",
				MarkdownDescription: "The id of this KaaS",
			},
			"pack_name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the pack associated to the KaaS project",
				MarkdownDescription: "The name of the pack associated to the KaaS project",
			},
			"region": schema.StringAttribute{
				Computed:            true,
				Description:         "The region where the KaaS project resides in.",
				MarkdownDescription: "The region where the KaaS project resides in.",
			},
			"kubeconfig": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "The kubeconfig generated to access to KaaS project",
				MarkdownDescription: "The kubeconfig generated to access to KaaS project",
			},
			"kubernetes_version": schema.StringAttribute{
				Computed:            true,
				Description:         "The version of Kubernetes associated with the KaaS project",
				MarkdownDescription: "The version of Kubernetes associated with the KaaS project",
			},
		},
		MarkdownDescription: "The kaas data source allows the user to manage a kaas project",
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kaasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KaasModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	obj, err := d.client.Kaas.GetKaas(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find KaaS",
			err.Error(),
		)
		return
	}

	data.Kubeconfig = types.StringValue(obj.Kubeconfig)
	data.Region = types.StringValue(obj.Region)
	data.KubernetesVersion = types.StringValue(obj.KubernetesVersion)

	// Set state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata returns the data source type name.
func (d *kaasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas"
}
