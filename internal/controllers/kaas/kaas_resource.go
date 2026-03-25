package kaas

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/kaas"
	"terraform-provider-infomaniak/internal/provider"
	kaas_schemas "terraform-provider-infomaniak/internal/schemas/kaas"
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
	client *apis.Client
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
}

func (r *kaasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = kaas_schemas.GetKaasResourceSchema()
}

func (r *kaasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data kaas_schemas.KaasModel

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
			PublicCloudId: data.PublicCloudId.ValueInt64(),
			ProjectId:     data.PublicCloudProjectId.ValueInt64(),
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

	data.Id = types.Int64Value(kaasId)
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

	err = r.fetchAndSetKubeconfig(&data, kaasObject)
	if err != nil {
		resp.Diagnostics.AddWarning("could not fetch and set kubeconfig", err.Error())
	}

	data.Fill(kaasObject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apiserver != nil {
		apiserverParamsInput := r.buildApiserverParamsInput(data)
		created, err := r.client.Kaas.PatchApiserverParams(apiserverParamsInput, input.Project.PublicCloudId, input.Project.ProjectId, kaasId)
		if !created || err != nil {
			resp.Diagnostics.AddError(
				"Error when creating Oidc",
				err.Error(),
			)
			return
		}

		applyFiltersDiags := r.applyIPFilters(ctx, data.Apiserver.IpFilters, input.Project.PublicCloudId, input.Project.ProjectId, kaasId)
		resp.Diagnostics.Append(applyFiltersDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		data.FillApiserverState(ctx, apiserverParamsInput)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

func (r *kaasResource) getApiserverParamsValues(data kaas_schemas.KaasModel) map[string]string {
	params := make(map[string]string)
	if !data.Apiserver.Params.IsNull() && !data.Apiserver.Params.IsUnknown() {
		for key, val := range data.Apiserver.Params.Elements() {
			if strVal, ok := val.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				params[key] = strVal.ValueString()
			}
		}
	}

	return params
}

func (r *kaasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kaas_schemas.KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	kaasObject, err := r.client.Kaas.GetKaas(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading KaaS",
			err.Error(),
		)
		return
	}

	state.Fill(kaasObject)

	err = r.fetchAndSetKubeconfig(&state, kaasObject)
	if err != nil {
		resp.Diagnostics.AddWarning("could not fetch and set kubeconfig", err.Error())
	}

	apiserverParams, err := r.client.Kaas.GetApiserverParams(state.PublicCloudId.ValueInt64(), state.PublicCloudProjectId.ValueInt64(), kaasObject.Id)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Could not get Oidc",
			err.Error(),
		)
	}

	if apiserverParams != nil {
		state.FillApiserverState(ctx, apiserverParams)
	}

	if state.Apiserver != nil {
		ipFilters, err := r.client.Kaas.GetIPFilters(state.PublicCloudId.ValueInt64(), state.PublicCloudProjectId.ValueInt64(), kaasObject.Id)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not get IP filter",
				err.Error(),
			)
			return
		}
		resp.Diagnostics.Append(state.FillIpFilters(ctx, ipFilters)...)
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

	chosenPackState, err := r.getPackId(state, &resp.Diagnostics)
	if err != nil {
		return
	}

	input := r.prepareUpdateInput(state, data, chosenPackState.Id)

	if _, err := r.client.Kaas.UpdateKaas(input); err != nil {
		resp.Diagnostics.AddError("Error when updating KaaS", err.Error())
		return
	}

	kaasObject, err := r.waitUntilActive(ctx, input, input.Id)
	if err != nil || kaasObject == nil {
		resp.Diagnostics.AddError("Error waiting for KaaS activation", err.Error())
		return
	}

	err = r.fetchAndSetKubeconfig(&data, kaasObject)
	if err != nil {
		resp.Diagnostics.AddWarning("could not fetch and set kubeconfig", err.Error())
	}

	data.Fill(kaasObject)

	if data.Apiserver != nil {
		r.handleApiserverConfig(ctx, &data, input, resp)

		applyFiltersDiags := r.applyIPFilters(ctx, data.Apiserver.IpFilters, input.Project.PublicCloudId, input.Project.ProjectId, input.Id)
		resp.Diagnostics.Append(applyFiltersDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
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

func (r *kaasResource) handleApiserverConfig(ctx context.Context, data *kaas_schemas.KaasModel, input *kaas.Kaas, resp *resource.UpdateResponse) {
	apiserverParamsInput := r.buildApiserverParamsInput(*data)
	patched, err := r.client.Kaas.PatchApiserverParams(apiserverParamsInput, input.Project.PublicCloudId, input.Project.ProjectId, input.Id)
	if !patched || err != nil {
		resp.Diagnostics.AddError("Error when patching Apiserver params", err.Error())
		return
	}
	data.FillApiserverState(ctx, apiserverParamsInput)
}

func (r *kaasResource) applyIPFilters(ctx context.Context, terraformIpFilters types.List, publicCloudId, projectId, kaasId int64) diag.Diagnostics {
	var diags diag.Diagnostics
	ipFilters := make([]string, 0, len(terraformIpFilters.Elements()))
	diags.Append(terraformIpFilters.ElementsAs(ctx, &ipFilters, true)...)
	if diags.HasError() {
		return diags
	}

	convertedIpFilters := make([]netip.Prefix, len(ipFilters))
	for i, cidr := range ipFilters {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			diags.AddError("invalid cidr format", err.Error())
			return diags
		}

		convertedIpFilters[i] = prefix
	}

	ok, err := r.client.Kaas.PutIPFilters(convertedIpFilters, publicCloudId, projectId, kaasId)
	if !ok || err != nil {
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		} else {
			errMsg = "PutIPFilters returned false but no error was provided"
		}
		diags.AddError("Error when applying ip filters", errMsg)
	}

	return diags
}

func (r *kaasResource) buildApiserverParamsInput(data kaas_schemas.KaasModel) *kaas.Apiserver {
	apiserverParamsInput := &kaas.Apiserver{
		NonSpecificApiServerParams: r.getApiserverParamsValues(data),
	}
	if data.Apiserver.Audit != nil {
		apiserverParamsInput.AuditLogPolicy = data.Apiserver.Audit.Policy.ValueStringPointer()
		apiserverParamsInput.AuditLogWebhook = data.Apiserver.Audit.WebhookConfig.ValueStringPointer()
	}
	if data.Apiserver.Oidc != nil {
		apiserverParamsInput.OidcCa = data.Apiserver.Oidc.Ca.ValueStringPointer()
		apiserverParamsInput.Params = &kaas.ApiServerParams{
			IssuerUrl:      data.Apiserver.Oidc.IssuerUrl.ValueStringPointer(),
			ClientId:       data.Apiserver.Oidc.ClientId.ValueStringPointer(),
			UsernameClaim:  data.Apiserver.Oidc.UsernameClaim.ValueStringPointer(),
			UsernamePrefix: data.Apiserver.Oidc.UsernamePrefix.ValueStringPointer(),
			SigningAlgs:    data.Apiserver.Oidc.SigningAlgs.ValueStringPointer(),
			GroupsClaim:    data.Apiserver.Oidc.GroupsClaim.ValueStringPointer(),
			GroupsPrefix:   data.Apiserver.Oidc.GroupsPrefix.ValueStringPointer(),
			RequiredClaim:  data.Apiserver.Oidc.RequiredClaim.ValueStringPointer(),
		}
	}
	return apiserverParamsInput
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
