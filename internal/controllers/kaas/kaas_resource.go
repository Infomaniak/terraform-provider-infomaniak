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
	kaas_schemas "terraform-provider-infomaniak/internal/schemas/kaas"
	"terraform-provider-infomaniak/internal/services"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

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
	client      *apis.Client
	kaasService *services.KaasService
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

	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			err.Error(),
		)
		return
	}

	r.client = client
	r.kaasService = services.NewKaasService(r.client)
}

func (r *kaasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = kaas_schemas.KaasResourceSchema
}

func (r *kaasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data kaas_schemas.KaasModel
	var state kaas_schemas.KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// On create the state does not exist, so we take the plan and we will add stuff to it
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	chosenPack, err := r.getPackId(data, &resp.Diagnostics)
	if err != nil {
		return
	}

	kaasId, diags := r.kaasService.CreateKaas(data, *chosenPack, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, diags = r.kaasService.GetKaasOnceActive(ctx, data, kaasId, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, diags = r.kaasService.GetKubeconfig(data, kaasId, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.kaasService.SetApiserverConfig(ctx, data, kaasId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *kaasResource) waitUntilActive(ctx context.Context, kaas *kaas.Kaas, id int64) (*kaas.Kaas, error) {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.C:
			found, err := r.client.Kaas.GetKaas(kaas.Project.PublicCloudId, kaas.Project.ProjectId, id)
			if err != nil {
				return nil, err
			}

			if found.Status == "Active" {
				return found, nil
			}
		}
	}
}

func (r *kaasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kaas_schemas.KaasModel
	var data kaas_schemas.KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	kaasId := data.Id.ValueInt64()

	_, diags := r.kaasService.GetKaasOnceActive(ctx, data, kaasId, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, diags = r.kaasService.GetKubeconfig(data, kaasId, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.kaasService.ReadApiserverConfig(ctx, data, kaasId, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *kaasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state kaas_schemas.KaasModel
	var data kaas_schemas.KaasModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	kaasId := data.Id.ValueInt64()

	chosenPackState, err := r.getPackId(state, &resp.Diagnostics)
	if err != nil {
		return
	}

	input := r.prepareUpdateInput(state, data, chosenPackState.Id)

	if _, err := r.client.Kaas.UpdateKaas(input); err != nil {
		resp.Diagnostics.AddError("Error when updating KaaS", err.Error())
		return
	}

	_, diags := r.kaasService.GetKaasOnceActive(ctx, data, kaasId, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, diags = r.kaasService.GetKubeconfig(data, kaasId, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.kaasService.SetApiserverConfig(ctx, data, input.Id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasResource) prepareUpdateInput(state, data kaas_schemas.KaasModel, packID int64) *kaas.Kaas {
	input := &kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: data.PublicCloudId.ValueInt64(),
			ProjectId:     data.PublicCloudProjectId.ValueInt64(),
		},
		Id:                state.Id.ValueInt64(),
		Name:              data.Name.ValueString(),
		PackId:            packID,
		Region:            state.Region.ValueString(),
		KubernetesVersion: data.KubernetesVersion.ValueString(),
	}

	if state.KubernetesVersion.ValueString() == data.KubernetesVersion.ValueString() {
		input.KubernetesVersion = ""
	}

	return input
}

func (r *kaasResource) fetchAndSetKubeconfig(data *kaas_schemas.KaasModel, input *kaas.Kaas) error {
	kubeconfig, err := r.client.Kaas.GetKubeconfig(
		input.Project.PublicCloudId,
		input.Project.ProjectId,
		input.Id,
	)
	if err != nil {
		return fmt.Errorf("could not get kubeconfig: %w", err)
	}
	data.Kubeconfig = types.StringValue(kubeconfig)
	return nil
}

func (r *kaasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data kaas_schemas.KaasModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteKaas API call logic
	_, err := r.client.Kaas.DeleteKaas(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.Id.ValueInt64(),
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

func (r *kaasResource) getPackId(data kaas_schemas.KaasModel, diagnostic *diag.Diagnostics) (*kaas.KaasPack, error) {
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
