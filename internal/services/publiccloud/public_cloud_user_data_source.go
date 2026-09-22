package publiccloud

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/publiccloud"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &publicCloudUserDataSource{}
	_ datasource.DataSourceWithConfigure = &publicCloudUserDataSource{}
)

type publicCloudUserDataSource struct {
	client *apis.Client
}

func NewPublicCloudUserDataSource() datasource.DataSource {
	return &publicCloudUserDataSource{}
}

func (d *publicCloudUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_user"
}

func (d *publicCloudUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *publicCloudUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getPublicCloudUserDataSourceSchema()
}

func (d *publicCloudUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PublicCloudUserModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := d.client.PublicCloud.GetUser(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Public Cloud user", err.Error())
		return
	}

	fillPublicCloudUserModel(&data, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func fillPublicCloudUserModel(m *PublicCloudUserModel, obj *publiccloud.User) {
	m.Id = types.Int64Value(obj.Id)
	m.PublicCloudProjectId = types.Int64Value(obj.PublicCloudProjectId)
	m.OpenStackName = types.StringValue(obj.OpenStackName)
	m.Description = types.StringValue(obj.Description)
	m.Status = types.StringValue(obj.Status)
	m.CreatedAt = types.Int64Value(obj.CreatedAt)
	m.UpdatedAt = types.Int64Value(obj.UpdatedAt)
}
