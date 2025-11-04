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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

	Id                 types.Int64  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	AddDefaultSchedule types.Bool   `tfsdk:"add_default_schedule"`
	Time               types.String `tfsdk:"time"`
	Keep               types.Int32  `tfsdk:"keep"`
	IsPitrEnabled      types.Bool   `tfsdk:"is_pitr_enabled"`
}

func (model *DBaasBackupScheduleModel) fill(backupSchedule *dbaas.DBaasBackupSchedule) {
	model.AddDefaultSchedule = types.BoolPointerValue(backupSchedule.AddDefaultSchedule)
	model.Time = types.StringPointerValue(backupSchedule.Time)
	model.Keep = types.Int32PointerValue(backupSchedule.Keep)
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
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud project",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"dbaas_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the dbaas",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "BackupSchedule identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the backup",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"add_default_schedule": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Add the default backup schedule to the Database Service. The default schedule is ran during a random hour:minute once a day keeping 7 days worth of backups",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"time": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Use the given time as the time to create the scheduled backup (24 hour, UTC)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"keep": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The number of backups to keep for the schedule",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"is_pitr_enabled": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Enable/Disable point in time recovery",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *dbaasBackupScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DBaasBackupScheduleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := &dbaas.DBaasBackupSchedule{
		AddDefaultSchedule: data.AddDefaultSchedule.ValueBoolPointer(),
		Time:               data.Time.ValueStringPointer(),
		Keep:               data.Keep.ValueInt32Pointer(),
		IsPitrEnabled:      data.IsPitrEnabled.ValueBoolPointer(),
	}

	created, err := r.client.DBaas.CreateDBaasScheduleBackup(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.DbaasId.ValueInt64()),
		input,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating Backup Schedule",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(created.Id)

	scheduleBackup, err := r.client.DBaas.GetDBaasScheduleBackup(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.DbaasId.ValueInt64()),
		int(data.Id.ValueInt64()),
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
		AddDefaultSchedule: data.AddDefaultSchedule.ValueBoolPointer(),
		Time:               data.Time.ValueStringPointer(),
		Keep:               data.Keep.ValueInt32Pointer(),
		IsPitrEnabled:      data.IsPitrEnabled.ValueBoolPointer(),
	}

	ok, err := r.client.DBaas.UpdateDBaasScheduleBackup(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.DbaasId.ValueInt64()),
		int(data.Id.ValueInt64()),
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
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.DbaasId.ValueInt64()),
		int(data.Id.ValueInt64()),
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
		int(state.PublicCloudId.ValueInt64()),
		int(state.PublicCloudProjectId.ValueInt64()),
		int(state.DbaasId.ValueInt64()),
		int(state.Id.ValueInt64()),
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
		int(state.PublicCloudId.ValueInt64()),
		int(state.PublicCloudProjectId.ValueInt64()),
		int(state.DbaasId.ValueInt64()),
		int(state.Id.ValueInt64()),
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
