package publiccloud

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &publicCloudAccessesDataSource{}
	_ datasource.DataSourceWithConfigure = &publicCloudAccessesDataSource{}
)

type publicCloudAccessesDataSource struct {
	client *apis.Client
}

func NewPublicCloudAccessesDataSource() datasource.DataSource {
	return &publicCloudAccessesDataSource{}
}

func (d *publicCloudAccessesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_accesses"
}

func (d *publicCloudAccessesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *publicCloudAccessesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getPublicCloudAccessesDataSourceSchema()
}

func (d *publicCloudAccessesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PublicCloudAccessesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := d.client.PublicCloud.GetAccesses(data.AccountId.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Public Cloud accesses", err.Error())
		return
	}

	data.IsMaintenanceOngoing = types.BoolValue(obj.IsMaintenanceOngoing)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
