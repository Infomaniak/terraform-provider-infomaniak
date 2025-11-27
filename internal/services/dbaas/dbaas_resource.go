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
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &dbaasResource{}
	_ resource.ResourceWithConfigure   = &dbaasResource{}
	_ resource.ResourceWithImportState = &dbaasResource{}
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

	MySqlConfiguration *MySqlConfig `tfsdk:"mysqlconfig"`
}

type MySqlConfig struct {
	AutoIncrementIncrement           types.Int64   `tfsdk:"auto_increment_increment"`
	AutoIncrementOffset              types.Int64   `tfsdk:"auto_increment_offset"`
	CharacterSetServer               types.String  `tfsdk:"character_set_server"`
	ConnectTimeout                   types.Int64   `tfsdk:"connect_timeout"`
	GroupConcatMaxLen                types.Int64   `tfsdk:"group_concat_max_len"`
	InformationSchemaStatsExpiry     types.Int64   `tfsdk:"information_schema_stats_expiry"`
	InnodbChangeBufferMaxSize        types.Int64   `tfsdk:"innodb_change_buffer_max_size"`
	InnodbFlushNeighbors             types.Int64   `tfsdk:"innodb_flush_neighbors"`
	InnodbFtMaxTokenSize             types.Int64   `tfsdk:"innodb_ft_max_token_size"`
	InnodbFtMinTokenSize             types.Int64   `tfsdk:"innodb_ft_min_token_size"`
	InnodbFtServerStopwordTable      types.String  `tfsdk:"innodb_ft_server_stopword_table"`
	InnodbLockWaitTimeout            types.Int64   `tfsdk:"innodb_lock_wait_timeout"`
	InnodbLogBufferSize              types.Int64   `tfsdk:"innodb_log_buffer_size"`
	InnodbOnlineAlterLogMaxSize      types.Int64   `tfsdk:"innodb_online_alter_log_max_size"`
	InnodbPrintAllDeadlocks          types.String  `tfsdk:"innodb_print_all_deadlocks"`
	InnodbReadIoThreads              types.Int64   `tfsdk:"innodb_read_io_threads"`
	InnodbRollbackOnTimeout          types.String  `tfsdk:"innodb_rollback_on_timeout"`
	InnodbStatsPersistentSamplePages types.Int64   `tfsdk:"innodb_stats_persistent_sample_pages"`
	InnodbThreadConcurrency          types.Int64   `tfsdk:"innodb_thread_concurrency"`
	InnodbWriteIoThreads             types.Int64   `tfsdk:"innodb_write_io_threads"`
	InteractiveTimeout               types.Int64   `tfsdk:"interactive_timeout"`
	LockWaitTimeout                  types.Int64   `tfsdk:"lock_wait_timeout"`
	LogBinTrustFunctionCreators      types.String  `tfsdk:"log_bin_trust_function_creators"`
	LongQueryTime                    types.Float64 `tfsdk:"long_query_time"`
	MaxAllowedPacket                 types.Int64   `tfsdk:"max_allowed_packet"`
	MaxConnections                   types.Int64   `tfsdk:"max_connections"`
	MaxDigestLength                  types.Int64   `tfsdk:"max_digest_length"`
	MaxHeapTableSize                 types.Int64   `tfsdk:"max_heap_table_size"`
	MaxPreparedStmtCount             types.Int64   `tfsdk:"max_prepared_stmt_count"`
	MinExaminedRowLimit              types.Int64   `tfsdk:"min_examined_row_limit"`
	NetBufferLength                  types.Int64   `tfsdk:"net_buffer_length"`
	NetReadTimeout                   types.Int64   `tfsdk:"net_read_timeout"`
	NetWriteTimeout                  types.Int64   `tfsdk:"net_write_timeout"`
	PerformanceSchemaMaxDigestLength types.Int64   `tfsdk:"performance_schema_max_digest_length"`
	RequireSecureTransport           types.String  `tfsdk:"require_secure_transport"`
	SortBufferSize                   types.Int64   `tfsdk:"sort_buffer_size"`
	SqlMode                          types.List    `tfsdk:"sql_mode"`
	TableDefinitionCache             types.Int64   `tfsdk:"table_definition_cache"`
	TableOpenCache                   types.Int64   `tfsdk:"table_open_cache"`
	TableOpenCacheInstances          types.Int64   `tfsdk:"table_open_cache_instances"`
	ThreadStack                      types.Int64   `tfsdk:"thread_stack"`
	TransactionIsolation             types.String  `tfsdk:"transaction_isolation"`
	WaitTimeout                      types.Int64   `tfsdk:"wait_timeout"`
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

	if data.MySqlConfiguration != nil {
		mysqlConfig, diags := data.MySqlConfiguration.toAPIModel(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		ok, err := r.client.DBaas.PutConfiguration(
			data.PublicCloudId.ValueInt64(),
			data.PublicCloudProjectId.ValueInt64(),
			data.Id.ValueInt64(),
			*mysqlConfig,
		)
		if !ok && err == nil {
			resp.Diagnostics.AddError("Unknown MySQL Settings error", "")
			return
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Error when updating DBaaS MySQL Settings",
				err.Error(),
			)
			return
		}
	}

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

	state.fill(dbaasObject)

	listFilteredIps, diags := types.ListValueFrom(ctx, types.StringType, filteredIps)
	state.AllowedCIDRs = listFilteredIps
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current configuration
	apiConfig, err := r.client.DBaas.GetConfiguration(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading DBaaS configuration",
			err.Error(),
		)
		return
	}

	if apiConfig != nil {
		if state.MySqlConfiguration == nil {
			state.MySqlConfiguration = &MySqlConfig{}
		}
		diags = state.MySqlConfiguration.fromAPIModel(ctx, apiConfig)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

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
	if resp.Diagnostics.HasError() {
		return
	}
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

	if data.MySqlConfiguration != nil {
		mysqlConfig, diags := data.MySqlConfiguration.toAPIModel(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		ok, err = r.client.DBaas.PutConfiguration(
			state.PublicCloudId.ValueInt64(),
			state.PublicCloudProjectId.ValueInt64(),
			state.Id.ValueInt64(),
			*mysqlConfig,
		)
		if !ok && err == nil {
			resp.Diagnostics.AddError("Unknown MySQL Settings error", "")
			return
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Error when updating DBaaS MySQL Settings",
				err.Error(),
			)
			return
		}

		state.MySqlConfiguration = data.MySqlConfiguration
	}

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
		model.Password = types.StringValue(dbaas.Connection.Password)
		model.Ca = types.StringValue(dbaas.Connection.Ca)
	}
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

func (config *MySqlConfig) toAPIModel(ctx context.Context) (*dbaas.MySqlConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	apiConfig := &dbaas.MySqlConfig{}

	if !config.AutoIncrementIncrement.IsUnknown() {
		apiConfig.AutoIncrementIncrement = config.AutoIncrementIncrement.ValueInt64Pointer()
	} else {
		config.AutoIncrementIncrement = types.Int64Null()
	}
	if !config.AutoIncrementOffset.IsUnknown() {
		apiConfig.AutoIncrementOffset = config.AutoIncrementOffset.ValueInt64Pointer()
	} else {
		config.AutoIncrementOffset = types.Int64Null()
	}
	if !config.CharacterSetServer.IsUnknown() {
		apiConfig.CharacterSetServer = config.CharacterSetServer.ValueStringPointer()
	} else {
		config.CharacterSetServer = types.StringNull()
	}
	if !config.ConnectTimeout.IsUnknown() {
		apiConfig.ConnectTimeout = config.ConnectTimeout.ValueInt64Pointer()
	} else {
		config.ConnectTimeout = types.Int64Null()
	}
	if !config.GroupConcatMaxLen.IsUnknown() {
		apiConfig.GroupConcatMaxLen = config.GroupConcatMaxLen.ValueInt64Pointer()
	} else {
		config.GroupConcatMaxLen = types.Int64Null()
	}
	if !config.InformationSchemaStatsExpiry.IsUnknown() {
		apiConfig.InformationSchemaStatsExpiry = config.InformationSchemaStatsExpiry.ValueInt64Pointer()
	} else {
		config.InformationSchemaStatsExpiry = types.Int64Null()
	}
	if !config.InnodbChangeBufferMaxSize.IsUnknown() {
		apiConfig.InnodbChangeBufferMaxSize = config.InnodbChangeBufferMaxSize.ValueInt64Pointer()
	} else {
		config.InnodbChangeBufferMaxSize = types.Int64Null()
	}
	if !config.InnodbFlushNeighbors.IsUnknown() {
		apiConfig.InnodbFlushNeighbors = config.InnodbFlushNeighbors.ValueInt64Pointer()
	} else {
		config.InnodbFlushNeighbors = types.Int64Null()
	}
	if !config.InnodbFtMaxTokenSize.IsUnknown() {
		apiConfig.InnodbFtMaxTokenSize = config.InnodbFtMaxTokenSize.ValueInt64Pointer()
	} else {
		config.InnodbFtMaxTokenSize = types.Int64Null()
	}
	if !config.InnodbFtMinTokenSize.IsUnknown() {
		apiConfig.InnodbFtMinTokenSize = config.InnodbFtMinTokenSize.ValueInt64Pointer()
	} else {
		config.InnodbFtMinTokenSize = types.Int64Null()
	}
	if !config.InnodbFtServerStopwordTable.IsUnknown() {
		apiConfig.InnodbFtServerStopwordTable = config.InnodbFtServerStopwordTable.ValueStringPointer()
	} else {
		config.InnodbFtServerStopwordTable = types.StringNull()
	}
	if !config.InnodbLockWaitTimeout.IsUnknown() {
		apiConfig.InnodbLockWaitTimeout = config.InnodbLockWaitTimeout.ValueInt64Pointer()
	} else {
		config.InnodbLockWaitTimeout = types.Int64Null()
	}
	if !config.InnodbLogBufferSize.IsUnknown() {
		apiConfig.InnodbLogBufferSize = config.InnodbLogBufferSize.ValueInt64Pointer()
	} else {
		config.InnodbLogBufferSize = types.Int64Null()
	}
	if !config.InnodbOnlineAlterLogMaxSize.IsUnknown() {
		apiConfig.InnodbOnlineAlterLogMaxSize = config.InnodbOnlineAlterLogMaxSize.ValueInt64Pointer()
	} else {
		config.InnodbOnlineAlterLogMaxSize = types.Int64Null()
	}
	if !config.InnodbPrintAllDeadlocks.IsUnknown() {
		apiConfig.InnodbPrintAllDeadlocks = config.InnodbPrintAllDeadlocks.ValueStringPointer()
	} else {
		config.InnodbPrintAllDeadlocks = types.StringNull()
	}
	if !config.InnodbReadIoThreads.IsUnknown() {
		apiConfig.InnodbReadIoThreads = config.InnodbReadIoThreads.ValueInt64Pointer()
	} else {
		config.InnodbReadIoThreads = types.Int64Null()
	}
	if !config.InnodbRollbackOnTimeout.IsUnknown() {
		apiConfig.InnodbRollbackOnTimeout = config.InnodbRollbackOnTimeout.ValueStringPointer()
	} else {
		config.InnodbRollbackOnTimeout = types.StringNull()
	}
	if !config.InnodbStatsPersistentSamplePages.IsUnknown() {
		apiConfig.InnodbStatsPersistentSamplePages = config.InnodbStatsPersistentSamplePages.ValueInt64Pointer()
	} else {
		config.InnodbStatsPersistentSamplePages = types.Int64Null()
	}
	if !config.InnodbThreadConcurrency.IsUnknown() {
		apiConfig.InnodbThreadConcurrency = config.InnodbThreadConcurrency.ValueInt64Pointer()
	} else {
		config.InnodbThreadConcurrency = types.Int64Null()
	}
	if !config.InnodbWriteIoThreads.IsUnknown() {
		apiConfig.InnodbWriteIoThreads = config.InnodbWriteIoThreads.ValueInt64Pointer()
	} else {
		config.InnodbWriteIoThreads = types.Int64Null()
	}
	if !config.InteractiveTimeout.IsUnknown() {
		apiConfig.InteractiveTimeout = config.InteractiveTimeout.ValueInt64Pointer()
	} else {
		config.InteractiveTimeout = types.Int64Null()
	}
	if !config.LockWaitTimeout.IsUnknown() {
		apiConfig.LockWaitTimeout = config.LockWaitTimeout.ValueInt64Pointer()
	} else {
		config.LockWaitTimeout = types.Int64Null()
	}
	if !config.LogBinTrustFunctionCreators.IsUnknown() {
		apiConfig.LogBinTrustFunctionCreators = config.LogBinTrustFunctionCreators.ValueStringPointer()
	} else {
		config.LogBinTrustFunctionCreators = types.StringNull()
	}
	if !config.LongQueryTime.IsUnknown() {
		apiConfig.LongQueryTime = config.LongQueryTime.ValueFloat64Pointer()
	} else {
		config.LongQueryTime = types.Float64Null()
	}
	if !config.MaxAllowedPacket.IsUnknown() {
		apiConfig.MaxAllowedPacket = config.MaxAllowedPacket.ValueInt64Pointer()
	} else {
		config.MaxAllowedPacket = types.Int64Null()
	}
	if !config.MaxConnections.IsUnknown() {
		apiConfig.MaxConnections = config.MaxConnections.ValueInt64Pointer()
	} else {
		config.MaxConnections = types.Int64Null()
	}
	if !config.MaxDigestLength.IsUnknown() {
		apiConfig.MaxDigestLength = config.MaxDigestLength.ValueInt64Pointer()
	} else {
		config.MaxDigestLength = types.Int64Null()
	}
	if !config.MaxHeapTableSize.IsUnknown() {
		apiConfig.MaxHeapTableSize = config.MaxHeapTableSize.ValueInt64Pointer()
	} else {
		config.MaxHeapTableSize = types.Int64Null()
	}
	if !config.MaxPreparedStmtCount.IsUnknown() {
		apiConfig.MaxPreparedStmtCount = config.MaxPreparedStmtCount.ValueInt64Pointer()
	} else {
		config.MaxPreparedStmtCount = types.Int64Null()
	}
	if !config.MinExaminedRowLimit.IsUnknown() {
		apiConfig.MinExaminedRowLimit = config.MinExaminedRowLimit.ValueInt64Pointer()
	} else {
		config.MinExaminedRowLimit = types.Int64Null()
	}
	if !config.NetBufferLength.IsUnknown() {
		apiConfig.NetBufferLength = config.NetBufferLength.ValueInt64Pointer()
	} else {
		config.NetBufferLength = types.Int64Null()
	}
	if !config.NetReadTimeout.IsUnknown() {
		apiConfig.NetReadTimeout = config.NetReadTimeout.ValueInt64Pointer()
	} else {
		config.NetReadTimeout = types.Int64Null()
	}
	if !config.NetWriteTimeout.IsUnknown() {
		apiConfig.NetWriteTimeout = config.NetWriteTimeout.ValueInt64Pointer()
	} else {
		config.NetWriteTimeout = types.Int64Null()
	}
	if !config.PerformanceSchemaMaxDigestLength.IsUnknown() {
		apiConfig.PerformanceSchemaMaxDigestLength = config.PerformanceSchemaMaxDigestLength.ValueInt64Pointer()
	} else {
		config.PerformanceSchemaMaxDigestLength = types.Int64Null()
	}
	if !config.RequireSecureTransport.IsUnknown() {
		apiConfig.RequireSecureTransport = config.RequireSecureTransport.ValueStringPointer()
	} else {
		config.RequireSecureTransport = types.StringNull()
	}
	if !config.SortBufferSize.IsUnknown() {
		apiConfig.SortBufferSize = config.SortBufferSize.ValueInt64Pointer()
	} else {
		config.SortBufferSize = types.Int64Null()
	}

	if !config.SqlMode.IsNull() && !config.SqlMode.IsUnknown() {
		var sqlModes []string
		diags.Append(config.SqlMode.ElementsAs(ctx, &sqlModes, false)...)
		apiConfig.SqlMode = sqlModes
	} else {
		config.SqlMode = types.ListNull(types.StringType)
	}

	if !config.TableDefinitionCache.IsUnknown() {
		apiConfig.TableDefinitionCache = config.TableDefinitionCache.ValueInt64Pointer()
	} else {
		config.TableDefinitionCache = types.Int64Null()
	}
	if !config.TableOpenCache.IsUnknown() {
		apiConfig.TableOpenCache = config.TableOpenCache.ValueInt64Pointer()
	} else {
		config.TableOpenCache = types.Int64Null()
	}
	if !config.TableOpenCacheInstances.IsUnknown() {
		apiConfig.TableOpenCacheInstances = config.TableOpenCacheInstances.ValueInt64Pointer()
	} else {
		config.TableOpenCacheInstances = types.Int64Null()
	}
	if !config.ThreadStack.IsUnknown() {
		apiConfig.ThreadStack = config.ThreadStack.ValueInt64Pointer()
	} else {
		config.ThreadStack = types.Int64Null()
	}
	if !config.TransactionIsolation.IsUnknown() {
		apiConfig.TransactionIsolation = config.TransactionIsolation.ValueStringPointer()
	} else {
		config.TransactionIsolation = types.StringNull()
	}
	if !config.WaitTimeout.IsUnknown() {
		apiConfig.WaitTimeout = config.WaitTimeout.ValueInt64Pointer()
	} else {
		config.WaitTimeout = types.Int64Null()
	}

	return apiConfig, diags
}

func (config *MySqlConfig) fromAPIModel(ctx context.Context, apiConfig *dbaas.MySqlConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	config.AutoIncrementIncrement = types.Int64PointerValue(apiConfig.AutoIncrementIncrement)
	config.AutoIncrementOffset = types.Int64PointerValue(apiConfig.AutoIncrementOffset)
	config.CharacterSetServer = types.StringPointerValue(apiConfig.CharacterSetServer)
	config.ConnectTimeout = types.Int64PointerValue(apiConfig.ConnectTimeout)
	config.GroupConcatMaxLen = types.Int64PointerValue(apiConfig.GroupConcatMaxLen)
	config.InformationSchemaStatsExpiry = types.Int64PointerValue(apiConfig.InformationSchemaStatsExpiry)
	config.InnodbChangeBufferMaxSize = types.Int64PointerValue(apiConfig.InnodbChangeBufferMaxSize)
	config.InnodbFlushNeighbors = types.Int64PointerValue(apiConfig.InnodbFlushNeighbors)
	config.InnodbFtMaxTokenSize = types.Int64PointerValue(apiConfig.InnodbFtMaxTokenSize)
	config.InnodbFtMinTokenSize = types.Int64PointerValue(apiConfig.InnodbFtMinTokenSize)
	config.InnodbFtServerStopwordTable = types.StringPointerValue(apiConfig.InnodbFtServerStopwordTable)
	config.InnodbLockWaitTimeout = types.Int64PointerValue(apiConfig.InnodbLockWaitTimeout)
	config.InnodbLogBufferSize = types.Int64PointerValue(apiConfig.InnodbLogBufferSize)
	config.InnodbOnlineAlterLogMaxSize = types.Int64PointerValue(apiConfig.InnodbOnlineAlterLogMaxSize)
	config.InnodbPrintAllDeadlocks = types.StringPointerValue(apiConfig.InnodbPrintAllDeadlocks)
	config.InnodbReadIoThreads = types.Int64PointerValue(apiConfig.InnodbReadIoThreads)
	config.InnodbRollbackOnTimeout = types.StringPointerValue(apiConfig.InnodbRollbackOnTimeout)
	config.InnodbStatsPersistentSamplePages = types.Int64PointerValue(apiConfig.InnodbStatsPersistentSamplePages)
	config.InnodbThreadConcurrency = types.Int64PointerValue(apiConfig.InnodbThreadConcurrency)
	config.InnodbWriteIoThreads = types.Int64PointerValue(apiConfig.InnodbWriteIoThreads)
	config.InteractiveTimeout = types.Int64PointerValue(apiConfig.InteractiveTimeout)
	config.LockWaitTimeout = types.Int64PointerValue(apiConfig.LockWaitTimeout)
	config.LogBinTrustFunctionCreators = types.StringPointerValue(apiConfig.LogBinTrustFunctionCreators)
	config.LongQueryTime = types.Float64PointerValue(apiConfig.LongQueryTime)
	config.MaxAllowedPacket = types.Int64PointerValue(apiConfig.MaxAllowedPacket)
	config.MaxConnections = types.Int64PointerValue(apiConfig.MaxConnections)
	config.MaxDigestLength = types.Int64PointerValue(apiConfig.MaxDigestLength)
	config.MaxHeapTableSize = types.Int64PointerValue(apiConfig.MaxHeapTableSize)
	config.MaxPreparedStmtCount = types.Int64PointerValue(apiConfig.MaxPreparedStmtCount)
	config.MinExaminedRowLimit = types.Int64PointerValue(apiConfig.MinExaminedRowLimit)
	config.NetBufferLength = types.Int64PointerValue(apiConfig.NetBufferLength)
	config.NetReadTimeout = types.Int64PointerValue(apiConfig.NetReadTimeout)
	config.NetWriteTimeout = types.Int64PointerValue(apiConfig.NetWriteTimeout)
	config.PerformanceSchemaMaxDigestLength = types.Int64PointerValue(apiConfig.PerformanceSchemaMaxDigestLength)
	config.RequireSecureTransport = types.StringPointerValue(apiConfig.RequireSecureTransport)
	config.SortBufferSize = types.Int64PointerValue(apiConfig.SortBufferSize)

	if len(apiConfig.SqlMode) > 0 {
		sqlModeList, diags := types.ListValueFrom(ctx, types.StringType, apiConfig.SqlMode)
		config.SqlMode = sqlModeList
		diags.Append(diags...)
	} else {
		config.SqlMode = types.ListNull(types.StringType)
	}

	config.TableDefinitionCache = types.Int64PointerValue(apiConfig.TableDefinitionCache)
	config.TableOpenCache = types.Int64PointerValue(apiConfig.TableOpenCache)
	config.TableOpenCacheInstances = types.Int64PointerValue(apiConfig.TableOpenCacheInstances)
	config.ThreadStack = types.Int64PointerValue(apiConfig.ThreadStack)
	config.TransactionIsolation = types.StringPointerValue(apiConfig.TransactionIsolation)
	config.WaitTimeout = types.Int64PointerValue(apiConfig.WaitTimeout)

	return diags
}
