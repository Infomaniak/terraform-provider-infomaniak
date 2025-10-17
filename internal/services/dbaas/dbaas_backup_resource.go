package dbaas

import (
	"context"
	"maps"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/services/scopes"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &dbaasBackupResource{}
	_ resource.ResourceWithConfigure   = &dbaasBackupResource{}
	_ resource.ResourceWithImportState = &dbaasBackupResource{}
)

func NewDBaasBackupResource() resource.Resource {
	return &dbaasBackupResource{}
}

type dbaasBackupResource struct {
	client *apis.Client
}

type DBaasBackupModel struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	DBaasId              types.Int64  `tfsdk:"dbaas_id"`
	Id                   types.String `tfsdk:"id"`
}

func (r *dbaasBackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_backup"
}

// Configure adds the provider configured client to the data source.
func (r *dbaasBackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dbaasBackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		MarkdownDescription: "The dbaas backup resource allows the user to manage a backup for a certain dbaas",
	}

	maps.Copy(resp.Schema.Attributes, scopes.DBaaS.Build())
}

func (r *dbaasBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DBaasBackupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// CreateBackup API call logic
	backupId, err := r.client.DBaas.CreateBackup(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.DBaasId.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating Backup",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(backupId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	backup, err := r.waitUntilActive(ctx,
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.DBaasId.ValueInt64()),
		backupId,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when waiting for Backup to be finished",
			err.Error(),
		)
		return
	}

	if backup == nil {
		return
	}

	data.fill(backup)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dbaasBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DBaasBackupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	backup, err := r.client.DBaas.GetBackup(
		int(state.PublicCloudId.ValueInt64()),
		int(state.PublicCloudProjectId.ValueInt64()),
		int(state.DBaasId.ValueInt64()),
		state.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading Backup",
			err.Error(),
		)
		return
	}

	state.fill(backup)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dbaasBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Cannot update",
		"This resource cannot be updated.",
	)
}

func (r *dbaasBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DBaasBackupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteBackup API call logic
	_, err := r.client.DBaas.DeleteBackup(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.DBaasId.ValueInt64()),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting Backup",
			err.Error(),
		)
		return
	}
}

func (r *dbaasBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	imports, err := parseBackupRestoreImport(req)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_id"), imports.PublicCloudId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_project_id"), imports.PublicCloudProjectId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("dbaas_id"), imports.DbaasId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), imports.Id)...)
}

func (r *dbaasBackupResource) waitUntilActive(ctx context.Context, publicCloudId int, publicCloudProjectId int, dbaasId int, id string) (*dbaas.DBaaSBackup, error) {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-t.C:
			found, err := r.client.DBaas.GetBackup(publicCloudId, publicCloudProjectId, dbaasId, id)
			if err != nil {
				return nil, err
			}

			if ctx.Err() != nil {
				return nil, nil
			}

			if found.Status != "Unknown" {
				return found, nil
			}
		}
	}
}

func (model *DBaasBackupModel) fill(backup *dbaas.DBaaSBackup) {
	model.Id = types.StringValue(backup.Id)
}
