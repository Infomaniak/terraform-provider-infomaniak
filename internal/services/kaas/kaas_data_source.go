package kaas

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
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "The name of the KaaS project",
				MarkdownDescription: "The name of the KaaS project",
			},
			"pack_name": schema.StringAttribute{
				Computed:            true,
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
			"apiserver": schema.SingleNestedAttribute{
				Description:         "Kubernetes Apiserver editable params",
				MarkdownDescription: "Kubernetes Apiserver editable params",
				Attributes: map[string]schema.Attribute{
					"params": schema.MapAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "Map of Kubernetes Apiserver params in case the terraform provider does not already abstracts them",
						MarkdownDescription: "Map of Kubernetes Apiserver params in case the terraform provider does not already abstracts them",
					},
					"audit": schema.SingleNestedAttribute{
						MarkdownDescription: "Kubernetes audit logs specification files",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"webhook_config": schema.StringAttribute{
								MarkdownDescription: "YAML manifest for audit webhook config",
								Computed:            true,
							},
							"policy": schema.StringAttribute{
								MarkdownDescription: "YAML manifest for audit policy",
								Computed:            true,
							},
						},
					},
					"oidc": schema.SingleNestedAttribute{
						Description:         "OIDC specific Apiserver params",
						MarkdownDescription: "OIDC specific Apiserver params",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"ca": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC Ca Certificate",
								MarkdownDescription: "OIDC Ca Certificate",
							},
							"groups_claim": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "OIDC groups claim",
							},
							"groups_prefix": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "OIDC groups prefix",
							},
							"issuer_url": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC issuer URL",
								MarkdownDescription: "OIDC issuer URL",
							},
							"client_id": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC client identifier",
								MarkdownDescription: "OIDC client identifier",
							},
							"username_claim": schema.StringAttribute{
								Optional:            true,
								Description:         "OIDC username claim",
								MarkdownDescription: "OIDC username claim",
							},
							"username_prefix": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC username prefix",
								MarkdownDescription: "OIDC username prefix",
							},
							"required_claim": schema.StringAttribute{
								Computed: true,
								MarkdownDescription: "A key=value pair that describes a required claim in the ID Token.",
							},
							"signing_algs": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC signing algorithm. Kubernetes will default it to RS256",
								MarkdownDescription: "OIDC signing algorithm. Kubernetes will default it to RS256",
							},
						},
					},
				},
				Optional: true,
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

	kubeconfig, err := d.client.Kaas.GetKubeconfig(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get kubeconfig from KaaS",
			err.Error(),
		)
		return
	}

	data.Kubeconfig = types.StringValue(kubeconfig)
	data.Region = types.StringValue(obj.Region)
	data.KubernetesVersion = types.StringValue(obj.KubernetesVersion)

	apiserverParams, err := d.client.Kaas.GetApiserverParams(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Oidc from KaaS",
			err.Error(),
		)
		return
	}

	if apiserverParams != nil {
		data.fillApiserverState(ctx, apiserverParams)
	}

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
