package publiccloud

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &publicCloudsDataSource{}
	_ datasource.DataSourceWithConfigure = &publicCloudsDataSource{}
)

type publicCloudsDataSource struct {
	client *apis.Client
}

func NewPublicCloudsDataSource() datasource.DataSource {
	return &publicCloudsDataSource{}
}

func (d *publicCloudsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_clouds"
}

func (d *publicCloudsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *publicCloudsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getPublicCloudsDataSourceSchema()
}

func (d *publicCloudsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PublicCloudsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	objs, err := d.client.PublicCloud.ListPublicClouds(data.AccountId.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Public Clouds", err.Error())
		return
	}

	data.PublicClouds = make([]PublicCloudModel, 0, len(objs))
	for _, obj := range objs {
		var item PublicCloudModel
		fillPublicCloudModel(&item, obj)
		data.PublicClouds = append(data.PublicClouds, item)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
