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
	_ datasource.DataSource              = &publicCloudDataSource{}
	_ datasource.DataSourceWithConfigure = &publicCloudDataSource{}
)

type publicCloudDataSource struct {
	client *apis.Client
}

func NewPublicCloudDataSource() datasource.DataSource {
	return &publicCloudDataSource{}
}

func (d *publicCloudDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud"
}

func (d *publicCloudDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *publicCloudDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getPublicCloudDataSourceSchema()
}

func (d *publicCloudDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PublicCloudModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := d.client.PublicCloud.GetPublicCloud(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Public Cloud", err.Error())
		return
	}

	fillPublicCloudModel(&data, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func fillPublicCloudModel(m *PublicCloudModel, obj *publiccloud.PublicCloud) {
	m.Id = types.Int64Value(obj.Id)
	m.AccountId = types.Int64Value(obj.AccountId)
	m.ServiceId = types.Int64Value(obj.ServiceId)
	m.ServiceName = types.StringValue(obj.ServiceName)
	m.CustomerName = types.StringValue(obj.CustomerName)
	m.InternalName = types.StringValue(obj.InternalName)
	m.Description = types.StringValue(obj.Description)
	m.BillReference = types.StringValue(obj.BillReference)
	m.CreatedAt = types.Int64Value(obj.CreatedAt)
	m.ExpiredAt = types.Int64Value(obj.ExpiredAt)
	m.IsFree = types.BoolValue(obj.IsFree)
	m.IsZeroPrice = types.BoolValue(obj.IsZeroPrice)
	m.IsTrial = types.BoolValue(obj.IsTrial)
	m.IsLocked = types.BoolValue(obj.IsLocked)
	m.HasMaintenance = types.BoolValue(obj.HasMaintenance)
	m.HasOperationInProgress = types.BoolValue(obj.HasOperationInProgress)
}
