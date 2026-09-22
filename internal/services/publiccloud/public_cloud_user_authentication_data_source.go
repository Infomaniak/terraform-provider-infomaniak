package publiccloud

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &publicCloudUserAuthenticationDataSource{}
	_ datasource.DataSourceWithConfigure = &publicCloudUserAuthenticationDataSource{}
)

type publicCloudUserAuthenticationDataSource struct {
	client *apis.Client
}

func NewPublicCloudUserAuthenticationDataSource() datasource.DataSource {
	return &publicCloudUserAuthenticationDataSource{}
}

func (d *publicCloudUserAuthenticationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_user_authentication"
}

func (d *publicCloudUserAuthenticationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", err.Error())
		return
	}
	d.client = client
}

func (d *publicCloudUserAuthenticationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getPublicCloudUserAuthenticationDataSourceSchema()
}

func (d *publicCloudUserAuthenticationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PublicCloudUserAuthenticationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	content, err := d.client.PublicCloud.GetAuthentication(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.PublicCloudUserId.ValueInt64(),
		data.Type.ValueString(),
		data.Region.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Unable to fetch authentication file", err.Error())
		return
	}

	data.Content = types.StringValue(content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
