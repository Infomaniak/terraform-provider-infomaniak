package dbaas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"terraform-provider-infomaniak/internal/dynamic"
	"terraform-provider-infomaniak/internal/provider"
	dbaasmigration "terraform-provider-infomaniak/internal/services/dbaas/dbaas_migration"
	"terraform-provider-infomaniak/internal/utils"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                 = &dbaasResource{}
	_ resource.ResourceWithConfigure    = &dbaasResource{}
	_ resource.ResourceWithImportState  = &dbaasResource{}
	_ resource.ResourceWithUpgradeState = &dbaasResource{}
)

func NewDBaasResource() resource.Resource {
	return &dbaasResource{}
}

type dbaasResource struct {
	client *apis.Client
}

type DBaasModel struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64  `tfsdk:"id"`
	KubernetesIdentifier types.String `tfsdk:"kube_identifier"`

	Name     types.String `tfsdk:"name"`
	PackName types.String `tfsdk:"pack_name"`
	Region   types.String `tfsdk:"region"`
	Type     types.String `tfsdk:"type"`
	Version  types.String `tfsdk:"version"`

	Host     types.String `tfsdk:"host"`
	Port     types.String `tfsdk:"port"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Ca       types.String `tfsdk:"ca"`

	AllowedCIDRs types.List `tfsdk:"allowed_cidrs"`

	Configuration          types.Dynamic `tfsdk:"configuration"`
	EffectiveConfiguration types.Dynamic `tfsdk:"effective_configuration"`
}

func (r *dbaasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas"
}

// Configure adds the provider configured client to the data source.
func (r *dbaasResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dbaasResource) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   dbaasmigration.GetV0Schema(),
			StateUpgrader: dbaasmigration.StateUpgrader,
		},
	}
}

func (r *dbaasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = getDbaasResourceSchema()
}

func (r *dbaasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DBaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	chosenPack, err := r.getPackId(data, &resp.Diagnostics)
	if err != nil {
		return
	}

	input := &dbaas.DBaaS{
		Project: dbaas.DBaaSProject{
			PublicCloudId: data.PublicCloudId.ValueInt64(),
			ProjectId:     data.PublicCloudProjectId.ValueInt64(),
		},
		Region:  data.Region.ValueString(),
		Version: data.Version.ValueString(),
		Type:    data.Type.ValueString(),
		Name:    data.Name.ValueString(),
		PackId:  chosenPack.Id,
	}

	// CreateDBaas API call logic
	createInfos, err := r.client.DBaas.CreateDBaaS(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating DBaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(createInfos.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	dbaasObject, err := r.waitUntilActive(ctx, input, createInfos.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when waiting for DBaaS to be Active",
			err.Error(),
		)
		return
	}

	if dbaasObject == nil {
		return
	}

	cidrs := make([]string, 0, len(data.AllowedCIDRs.Elements()))
	resp.Diagnostics.Append(data.AllowedCIDRs.ElementsAs(ctx, &cidrs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	allowedCIDRs := dbaas.AllowedCIDRs{
		IpFilters: cidrs,
	}
	ok, err := r.client.DBaas.PatchIpFilters(
		input.Project.PublicCloudId,
		input.Project.ProjectId,
		dbaasObject.Id,
		allowedCIDRs,
	)
	if !ok {
		resp.Diagnostics.AddError("Unknown IP filter error", "")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating IP Filters",
			err.Error(),
		)
		return
	}

	if !data.Configuration.IsNull() && !data.Configuration.IsUnknown() {
		configuration, diags := utils.ConvertDynamicObjectToMapAny(data.Configuration)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		ok, err = r.client.DBaas.PutConfiguration(
			data.PublicCloudId.ValueInt64(),
			data.PublicCloudProjectId.ValueInt64(),
			data.Id.ValueInt64(),
			configuration,
		)
		if !ok && err == nil {
			resp.Diagnostics.AddError("Unknown Settings error", "")
			return
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Error when updating DBaaS Settings",
				err.Error(),
			)
			return
		}
	} else {
		data.Configuration = types.DynamicNull()
	}

	newEffectiveConfig, diags := r.refreshEffectiveConfiguration(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.EffectiveConfiguration = newEffectiveConfig
	data.fill(dbaasObject)
	data.Password = types.StringValue(createInfos.RootPassword)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dbaasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DBaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	dbaasObject, err := r.client.DBaas.GetDBaaS(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading DBaaS",
			err.Error(),
		)
		return
	}

	filteredIps, err := r.client.DBaas.GetIpFilters(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading DBaaS filtered IPs",
			err.Error(),
		)
		return
	}

	listFilteredIps, diags := types.ListValueFrom(ctx, types.StringType, filteredIps)
	state.AllowedCIDRs = listFilteredIps
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newEffectiveConfig, diags := r.refreshEffectiveConfiguration(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.Id.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update values of configuration and effective_configuration
	newEffectiveConfig, newConfig, diags := utils.ObjectStateManager(
		ctx,
		newEffectiveConfig,
		state.EffectiveConfiguration,
		state.Configuration,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	shouldUpdateConfiguration, diags := hasObjectChanged(state.Configuration, newConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if shouldUpdateConfiguration {
		state.Configuration = newConfig
	}

	state.EffectiveConfiguration = newEffectiveConfig

	state.fill(dbaasObject)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dbaasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state DBaasModel
	var data DBaasModel

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
	input := &dbaas.DBaaS{
		Project: dbaas.DBaaSProject{
			PublicCloudId: data.PublicCloudId.ValueInt64(),
			ProjectId:     data.PublicCloudProjectId.ValueInt64(),
		},
		Id:      state.Id.ValueInt64(),
		Name:    data.Name.ValueString(),
		PackId:  chosenPackState.Id,
		Region:  state.Region.ValueString(),
		Version: state.Version.ValueString(),
		Type:    state.Type.ValueString(),
	}

	_, err = r.client.DBaas.UpdateDBaaS(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating DBaaS",
			err.Error(),
		)
		return
	}

	dbaasObject, err := r.waitUntilActive(ctx, input, input.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when getting DBaaS",
			err.Error(),
		)
		return
	}

	if dbaasObject == nil {
		return
	}

	cidrs := make([]string, 0, len(data.AllowedCIDRs.Elements()))
	resp.Diagnostics.Append(data.AllowedCIDRs.ElementsAs(ctx, &cidrs, false)...)
	allowedCIDRs := dbaas.AllowedCIDRs{
		IpFilters: cidrs,
	}
	ok, err := r.client.DBaas.PatchIpFilters(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.Id.ValueInt64(),
		allowedCIDRs,
	)
	if !ok && err == nil {
		resp.Diagnostics.AddError("Unknown IP filter error", "")
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating IP Filters",
			err.Error(),
		)
		return
	}

	state.AllowedCIDRs = data.AllowedCIDRs

	if !data.Configuration.IsNull() && !data.Configuration.IsUnknown() {
		configuration, diags := utils.ConvertDynamicObjectToMapAny(data.Configuration)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		ok, err = r.client.DBaas.PutConfiguration(
			state.PublicCloudId.ValueInt64(),
			state.PublicCloudProjectId.ValueInt64(),
			state.Id.ValueInt64(),
			configuration,
		)
		if !ok && err == nil {
			resp.Diagnostics.AddError("Unknown Settings error", "")
			return
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Error when updating DBaaS Settings",
				err.Error(),
			)
			return
		}

		state.Configuration = data.Configuration
	} else {
		state.Configuration = types.DynamicNull()
	}

	newEffectiveConfig, diags := r.refreshEffectiveConfiguration(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.Id.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.EffectiveConfiguration = newEffectiveConfig
	state.fill(dbaasObject)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dbaasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DBaasModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteDBaas API call logic
	_, err := r.client.DBaas.DeleteDBaaS(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting DBaaS",
			err.Error(),
		)
		return
	}
}

func (r *dbaasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
	dbaasId, err := strconv.ParseInt(idParts[2], 10, 64)
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), dbaasId)...)
}

func (r *dbaasResource) getPackId(data DBaasModel, diagnostic *diag.Diagnostics) (*dbaas.DBaaSPack, error) {
	pack, err := r.client.DBaas.FindPack(data.Type.ValueString(), data.PackName.ValueString())
	if err != nil {
		diagnostic.AddError(
			"Could not find DBaaS Pack",
			err.Error(),
		)
		return nil, err
	}

	return pack, nil
}

func (model *DBaasModel) fill(dbaas *dbaas.DBaaS) {
	model.Id = types.Int64Value(dbaas.Id)
	model.KubernetesIdentifier = types.StringValue(dbaas.KubernetesIdentifier)
	model.Region = types.StringValue(dbaas.Region)
	model.Type = types.StringValue(dbaas.Type)
	model.Version = types.StringValue(dbaas.Version)
	model.Name = types.StringValue(dbaas.Name)
	model.PackName = types.StringValue(dbaas.Pack.Name)

	if dbaas.Connection != nil {
		model.Host = types.StringValue(dbaas.Connection.Host)
		model.Port = types.StringValue(dbaas.Connection.Port)
		model.User = types.StringValue(dbaas.Connection.User)
		model.Ca = types.StringValue(dbaas.Connection.Ca)

		if model.Password == types.StringNull() || model.Password == types.StringUnknown() {
			model.Password = types.StringValue(dbaas.Connection.Password)
		}
	}
}

func (r *dbaasResource) refreshEffectiveConfiguration(publicCloudId, publicCloudProjectId, id int64) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics
	effectiveSettings, err := r.client.DBaas.GetConfiguration(
		publicCloudId,
		publicCloudProjectId,
		id,
	)
	if err != nil {
		diags.AddError(
			"Error when reading DBaaS settings",
			err.Error(),
		)
		return types.DynamicNull(), diags
	}

	jsonEffectiveSettigs, err := json.Marshal(effectiveSettings)
	if err != nil {
		diags.AddError("could not marshall", "effective settings json marshall fail")
	}
	dynamicObj, err := dynamic.FromJSONImplied(jsonEffectiveSettigs)
	if err != nil {
		diags.AddError("could not create dynamic object", "effective settings dynamic object from json creation failure")
	}

	return dynamicObj, diags
}

// hasObjectChanged convert the state configuration and newly generated configuration
// It then compares equivalent values ("100" is equal to 100) and tells if we need to update the state
func hasObjectChanged(stateConfig, newConfig types.Dynamic) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newConfigMap, d := utils.ConvertDynamicObjectToMapAny(newConfig)
	diags.Append(d...)

	newConfigConverted := utils.ConvertIntsToStrings(newConfigMap)

	stateConfigMap, d := utils.ConvertDynamicObjectToMapAny(stateConfig)
	diags.Append(d...)

	stateConfigConverted := utils.ConvertIntsToStrings(stateConfigMap)

	return !reflect.DeepEqual(newConfigConverted, stateConfigConverted), diags
}

func (r *dbaasResource) waitUntilActive(ctx context.Context, dbaas *dbaas.DBaaS, id int64) (*dbaas.DBaaS, error) {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-t.C:
			found, err := r.client.DBaas.GetDBaaS(dbaas.Project.PublicCloudId, dbaas.Project.ProjectId, id)
			if err != nil {
				return nil, err
			}

			if ctx.Err() != nil {
				return nil, nil
			}

			if found.Status == "ready" {
				return found, nil
			}
		}
	}
}
