package kaas

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/kaas"
	"terraform-provider-infomaniak/internal/provider"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &kaasResource{}
	_ resource.ResourceWithConfigure   = &kaasResource{}
	_ resource.ResourceWithImportState = &kaasResource{}
)

func NewKaasResource() resource.Resource {
	return &kaasResource{}
}

type kaasResource struct {
	client *apis.Client
}

type KaasModel struct {
	PublicCloudId        types.Int64 `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64 `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64 `tfsdk:"id"`

	Name              types.String `tfsdk:"name"`
	PackName          types.String `tfsdk:"pack_name"`
	Region            types.String `tfsdk:"region"`
	Kubeconfig        types.String `tfsdk:"kubeconfig"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
}

func (r *kaasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas"
}

// Configure adds the provider configured client to the data source.
func (r *kaasResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*provider.IkProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *apis.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = apis.NewClient(data.Data.Host.ValueString(), data.Data.Token.ValueString())
	if data.Version.ValueString() == "test" {
		r.client = apis.NewMockClient()
	}
}

func (r *kaasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud where KaaS is installed",
				MarkdownDescription: "The id of the public cloud where KaaS is installed",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud project where KaaS is installed",
				MarkdownDescription: "The id of the public cloud project where KaaS is installed",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"pack_name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the pack associated to the KaaS project",
				MarkdownDescription: "The name of the pack associated to the KaaS project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kubernetes_version": schema.StringAttribute{
				Required:            true,
				Description:         "The version of Kubernetes associated with the KaaS being installed",
				MarkdownDescription: "The version of Kubernetes associated with the KaaS being installed",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the KaaS project",
				MarkdownDescription: "The name of the KaaS project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				MarkdownDescription: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:            true,
				Description:         "The region where the KaaS will reside.",
				MarkdownDescription: "The region where the KaaS will reside.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kubeconfig": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "The kubeconfig generated to access to KaaS project",
				MarkdownDescription: "The kubeconfig generated to access to KaaS project",
			},
		},
		MarkdownDescription: "The kaas resource allows the user to manage a kaas project",
	}
}

func (r *kaasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	chosenPack, err := r.getPackId(data, &resp.Diagnostics)
	if err != nil {
		return
	}

	input := &kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: int(data.PublicCloudId.ValueInt64()),
			ProjectId:     int(data.PublicCloudProjectId.ValueInt64()),
		},
		Region:            data.Region.ValueString(),
		KubernetesVersion: data.KubernetesVersion.ValueString(),
		Name:              data.Name.ValueString(),
		PackId:            chosenPack.Id,
	}

	// CreateKaas API call logic
	kaasId, err := r.client.Kaas.CreateKaas(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating KaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(int64(kaasId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	kaasObject, err := r.waitUntilActive(ctx, input, kaasId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when waiting for KaaS to be Active",
			err.Error(),
		)
		return
	}

	if kaasObject == nil {
		return
	}

	// Wait for kubeconfig to be available
	kubeconfig, err := r.client.Kaas.GetKubeconfig(input.Project.PublicCloudId, input.Project.ProjectId, kaasId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Could not get kubeconfig for KaaS",
			err.Error(),
		)
	}

	data.Kubeconfig = types.StringValue(kubeconfig)
	data.Region = types.StringValue(kaasObject.Region)
	data.KubernetesVersion = types.StringValue(kaasObject.KubernetesVersion)
	data.Name = types.StringValue(kaasObject.Name)
	data.PackName = types.StringValue(kaasObject.Pack.Name)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasResource) waitUntilActive(ctx context.Context, kaas *kaas.Kaas, id int) (*kaas.Kaas, error) {
	for {
		found, err := r.client.Kaas.GetKaas(kaas.Project.PublicCloudId, kaas.Project.ProjectId, id)
		if err != nil {
			return nil, err
		}

		if ctx.Err() != nil {
			return nil, nil
		}

		if found.Status == "Active" {
			return found, nil
		}

		time.Sleep(5 * time.Second)
	}
}

func (r *kaasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	kaasObject, err := r.client.Kaas.GetKaas(
		int(state.PublicCloudId.ValueInt64()),
		int(state.PublicCloudProjectId.ValueInt64()),
		int(state.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading KaaS",
			err.Error(),
		)
		return
	}

	// Wait for kubeconfig to be available
	kubeconfig, err := r.client.Kaas.GetKubeconfig(kaasObject.Project.PublicCloudId, kaasObject.Project.ProjectId, kaasObject.Id)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Could not get kubeconfig",
			err.Error(),
		)
	}

	state.Id = types.Int64Value(int64(kaasObject.Id))
	state.Kubeconfig = types.StringValue(kubeconfig)
	state.Region = types.StringValue(kaasObject.Region)
	state.KubernetesVersion = types.StringValue(kaasObject.KubernetesVersion)
	state.Name = types.StringValue(kaasObject.Name)
	state.PackName = types.StringValue(kaasObject.Pack.Name)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *kaasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state KaasModel
	var data KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	chosenPackState, err := r.getPackId(state, &resp.Diagnostics)
	if err != nil {
		return
	}

	// Update API call logic
	input := &kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: int(data.PublicCloudId.ValueInt64()),
			ProjectId:     int(data.PublicCloudProjectId.ValueInt64()),
		},
		Id:                int(state.Id.ValueInt64()),
		Name:              data.Name.ValueString(),
		PackId:            chosenPackState.Id,
		Region:            state.Region.ValueString(),
		KubernetesVersion: state.KubernetesVersion.ValueString(),
	}

	_, err = r.client.Kaas.UpdateKaas(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating KaaS",
			err.Error(),
		)
		return
	}

	kaasObject, err := r.waitUntilActive(ctx, input, input.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when getting KaaS",
			err.Error(),
		)
		return
	}

	if kaasObject == nil {
		return
	}

	// Wait for kubeconfig to be available
	kubeconfig, err := r.client.Kaas.GetKubeconfig(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(state.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Could not get kubeconfig",
			err.Error(),
		)
	}

	data.Id = types.Int64Value(int64(kaasObject.Id))
	data.Kubeconfig = types.StringValue(kubeconfig)
	data.Region = types.StringValue(kaasObject.Region)
	data.PackName = types.StringValue(chosenPackState.Name)
	data.KubernetesVersion = types.StringValue(kaasObject.KubernetesVersion)
	data.Name = types.StringValue(kaasObject.Name)
	state.PackName = types.StringValue(kaasObject.Pack.Name)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KaasModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteKaas API call logic
	_, err := r.client.Kaas.DeleteKaas(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting KaaS",
			err.Error(),
		)
		return
	}
}

func (r *kaasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: public_cloud_id,public_cloud_project_id,id. Got: %q", req.ID),
		)
		return
	}

	var errorList error

	publicCloudId, err := strconv.ParseInt(idParts[0], 10, 64)
	errorList = errors.Join(errorList, err)
	publicCloudProjectId, err := strconv.ParseInt(idParts[1], 10, 64)
	errorList = errors.Join(errorList, err)
	kaasId, err := strconv.ParseInt(idParts[2], 10, 64)
	errorList = errors.Join(errorList, err)

	if errorList != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: public_cloud_id,public_cloud_project_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_id"), publicCloudId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_project_id"), publicCloudProjectId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), kaasId)...)
}

func (r *kaasResource) getPackId(data KaasModel, diagnostic *diag.Diagnostics) (*kaas.KaasPack, error) {
	packs, err := r.client.Kaas.GetPacks()
	if err != nil {
		diagnostic.AddError(
			"Could not get KaaS Packs",
			err.Error(),
		)
		return nil, err
	}

	var chosenPack *kaas.KaasPack
	for _, pack := range packs {
		if pack.Name == data.PackName.ValueString() {
			chosenPack = pack
			break
		}
	}

	if chosenPack == nil {
		var packNames []string
		for _, pack := range packs {
			packNames = append(packNames, pack.Name)
		}

		diagnostic.AddError(
			"Unknown KaaS Pack",
			fmt.Sprintf("pack_name must be one of : %v", packNames),
		)
		return nil, fmt.Errorf("pack name has not been found")
	}

	return chosenPack, nil
}
