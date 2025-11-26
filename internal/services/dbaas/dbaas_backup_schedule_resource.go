package dbaas

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &dbaasBackupScheduleResource{}
	_ resource.ResourceWithConfigure   = &dbaasBackupScheduleResource{}
	_ resource.ResourceWithImportState = &dbaasBackupScheduleResource{}
)

func NewDBaasBackupScheduleResource() resource.Resource {
	return &dbaasBackupScheduleResource{}
}

type dbaasBackupScheduleResource struct {
	client *apis.Client
}

type DBaasBackupScheduleModel struct {
	PublicCloudId        types.Int64 `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64 `tfsdk:"public_cloud_project_id"`
	DbaasId              types.Int64 `tfsdk:"dbaas_id"`

	Id            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ScheduledAt   types.String `tfsdk:"scheduled_at"`
	Retention     types.Int64  `tfsdk:"retention"`
	IsPitrEnabled types.Bool   `tfsdk:"is_pitr_enabled"`
}

func (model *DBaasBackupScheduleModel) fill(backupSchedule *dbaas.DBaasBackupSchedule) {
	model.ScheduledAt = types.StringPointerValue(backupSchedule.ScheduledAt)
	model.Retention = types.Int64PointerValue(backupSchedule.Retention)
	model.IsPitrEnabled = types.BoolPointerValue(backupSchedule.IsPitrEnabled)
	model.Name = types.StringPointerValue(backupSchedule.Name)
	model.Id = types.Int64PointerValue(backupSchedule.Id)
}

func (r *dbaasBackupScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_backup_schedule"
}

// Configure adds the provider configured client to the data source.
func (r *dbaasBackupScheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dbaasBackupScheduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = getDbaasBackupScheduleResourceSchema()
}

func (r *dbaasBackupScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DBaasBackupScheduleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := &dbaas.DBaasBackupSchedule{
		ScheduledAt:   data.ScheduledAt.ValueStringPointer(),
		Retention:     data.Retention.ValueInt64Pointer(),
		IsPitrEnabled: data.IsPitrEnabled.ValueBoolPointer(),
	}

	scheduleId, err := r.client.DBaas.CreateDBaasScheduleBackup(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.DbaasId.ValueInt64(),
		input,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating Backup Schedule",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(scheduleId)

	scheduleBackup, err := r.client.DBaas.GetDBaasScheduleBackup(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.DbaasId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when getting Backup Schedule",
			err.Error(),
		)
		return
	}

	data.fill(scheduleBackup)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dbaasBackupScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DBaasBackupScheduleModel
	var state DBaasBackupScheduleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := &dbaas.DBaasBackupSchedule{
		ScheduledAt:   data.ScheduledAt.ValueStringPointer(),
		Retention:     data.Retention.ValueInt64Pointer(),
		IsPitrEnabled: data.IsPitrEnabled.ValueBoolPointer(),
	}

	ok, err := r.client.DBaas.UpdateDBaasScheduleBackup(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.DbaasId.ValueInt64(),
		data.Id.ValueInt64(),
		input,
	)
	if !ok && err == nil {
		resp.Diagnostics.AddError("Unknown Backup Schedule error", "")
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating Backup Schedule",
			err.Error(),
		)
		return
	}

	scheduleBackup, err := r.client.DBaas.GetDBaasScheduleBackup(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.DbaasId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when getting Backup Schedule",
			err.Error(),
		)
		return
	}

	state.fill(scheduleBackup)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dbaasBackupScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DBaasBackupScheduleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	scheduleBackup, err := r.client.DBaas.GetDBaasScheduleBackup(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.DbaasId.ValueInt64(),
		state.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when getting Backup Schedule",
			err.Error(),
		)
		return
	}

	state.fill(scheduleBackup)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dbaasBackupScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DBaasBackupScheduleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteDBaas API call logic
	_, err := r.client.DBaas.DeleteDBaasScheduleBackup(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.DbaasId.ValueInt64(),
		state.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting DBaaS backup schedule",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dbaasBackupScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
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
	id, err := strconv.ParseInt(idParts[3], 10, 64)
	errorList = errors.Join(errorList, err)

	if errorList != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: public_cloud_id,public_cloud_project_id,dbaas_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_id"), publicCloudId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_project_id"), publicCloudProjectId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("dbaas_id"), dbaasId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
