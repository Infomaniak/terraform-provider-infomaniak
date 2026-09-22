package publiccloud

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &publicCloudConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &publicCloudConfigDataSource{}
)

type publicCloudConfigDataSource struct {
	client *apis.Client
}

func NewPublicCloudConfigDataSource() datasource.DataSource {
	return &publicCloudConfigDataSource{}
}

func (d *publicCloudConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_config"
}

func (d *publicCloudConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *publicCloudConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getPublicCloudConfigDataSourceSchema()
}

func (d *publicCloudConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PublicCloudConfigModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := d.client.PublicCloud.GetConfig(data.AccountId.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Public Cloud config", err.Error())
		return
	}

	data.FreeTier = types.Float64Value(obj.FreeTier)
	data.FreeTierUsed = types.Float64Value(obj.FreeTierUsed)
	data.AccountResourceLevel = types.Int64Value(obj.AccountResourceLevel)
	data.ProjectCount = types.Int64Value(obj.ProjectCount)
	data.ValidFrom = types.Int64Value(obj.ValidFrom)
	data.ValidTo = types.Int64Value(obj.ValidTo)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
