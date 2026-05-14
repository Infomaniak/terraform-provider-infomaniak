package publiccloud

import (
	"context"
	"strconv"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/publiccloud"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &publicCloudResource{}
	_ resource.ResourceWithConfigure   = &publicCloudResource{}
	_ resource.ResourceWithImportState = &publicCloudResource{}
)

type publicCloudResource struct {
	client *apis.Client
}

func NewPublicCloudResource() resource.Resource {
	return &publicCloudResource{}
}

func (r *publicCloudResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud"
}

func (r *publicCloudResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", err.Error())
		return
	}
	r.client = client
}

func (r *publicCloudResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = getPublicCloudResourceSchema()
}

// Create is intentionally a hard error: ordering a Public Cloud product is not
// exposed through the public Infomaniak API. Users must order the product via
// the Manager, then bring it under Terraform with `terraform import`.
func (r *publicCloudResource) Create(_ context.Context, _ resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Public Cloud cannot be created via Terraform",
		"Infomaniak's public API does not expose an endpoint to order a Public Cloud product. "+
			"Order the product through the Manager (https://manager.infomaniak.com), then run "+
			"`terraform import infomaniak_public_cloud.<name> <public_cloud_id>` to bring it under Terraform management.",
	)
}

func (r *publicCloudResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PublicCloudModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := r.client.PublicCloud.GetPublicCloud(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Public Cloud", err.Error())
		return
	}

	fillPublicCloudModel(&data, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *publicCloudResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PublicCloudModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &publiccloud.PublicCloud{
		Id:            plan.Id.ValueInt64(),
		CustomerName:  plan.CustomerName.ValueString(),
		Description:   plan.Description.ValueString(),
		BillReference: plan.BillReference.ValueString(),
	}
	if _, err := r.client.PublicCloud.UpdatePublicCloud(input); err != nil {
		resp.Diagnostics.AddError("Unable to update Public Cloud", err.Error())
		return
	}

	obj, err := r.client.PublicCloud.GetPublicCloud(plan.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to refresh Public Cloud after update", err.Error())
		return
	}

	fillPublicCloudModel(&plan, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op: the Public Cloud product cannot be removed via the public
// API. The resource is dropped from Terraform state with a warning explaining
// that the product remains active and billed.
func (r *publicCloudResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Public Cloud not deleted",
		"Removing this resource from Terraform does not delete the Public Cloud product at Infomaniak. "+
			"The product remains active and billed. Delete it via the Manager if needed.",
	)
}

func (r *publicCloudResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected a numeric Public Cloud id; got: "+req.ID,
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(id))...)
}
