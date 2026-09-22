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
	_ datasource.DataSource              = &publicCloudProjectDataSource{}
	_ datasource.DataSourceWithConfigure = &publicCloudProjectDataSource{}
)

type publicCloudProjectDataSource struct {
	client *apis.Client
}

func NewPublicCloudProjectDataSource() datasource.DataSource {
	return &publicCloudProjectDataSource{}
}

func (d *publicCloudProjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_project"
}

func (d *publicCloudProjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *publicCloudProjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getPublicCloudProjectDataSourceSchema()
}

func (d *publicCloudProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PublicCloudProjectModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := d.client.PublicCloud.GetProject(data.PublicCloudId.ValueInt64(), data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Public Cloud project", err.Error())
		return
	}

	fillPublicCloudProjectModel(&data, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func fillPublicCloudProjectModel(m *PublicCloudProjectModel, obj *publiccloud.Project) {
	m.PublicCloudId = types.Int64Value(obj.PublicCloudId)
	m.Id = types.Int64Value(obj.Id)
	m.Name = types.StringValue(obj.Name)
	m.OpenStackName = types.StringValue(obj.OpenStackName)
	m.Status = types.StringValue(obj.Status)
	m.Price = types.Float64Value(obj.Price)
	m.ResourceLevel = types.Int64Value(obj.ResourceLevel)
	m.UserCount = types.Int64Value(obj.UserCount)
	m.CreatedAt = types.Int64Value(obj.CreatedAt)
	m.UpdatedAt = types.Int64Value(obj.UpdatedAt)
	m.BillingStartAt = types.Int64Value(obj.BillingStartAt)
	m.BillingEndAt = types.Int64Value(obj.BillingEndAt)
	m.PriceUpdatedAt = types.Int64Value(obj.PriceUpdatedAt)
}
