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
	_ resource.Resource                = &dbaasRestoreResource{}
	_ resource.ResourceWithConfigure   = &dbaasRestoreResource{}
	_ resource.ResourceWithImportState = &dbaasRestoreResource{}
)

func NewDBaasRestoreResource() resource.Resource {
	return &dbaasRestoreResource{}
}

type dbaasRestoreResource struct {
	client *apis.Client
}

type DBaasRestoreModel struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	DBaasId              types.Int64  `tfsdk:"dbaas_id"`
	BackupId             types.String `tfsdk:"backup_id"`
	Id                   types.String `tfsdk:"id"`
}

func (r *dbaasRestoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_restore"
}

// Configure adds the provider configured client to the data source.
func (r *dbaasRestoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dbaasRestoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"backup_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The id of the backup",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		MarkdownDescription: "The dbaas restore resource allows the user to restore a backup for a certain dbaas",
	}

	maps.Copy(resp.Schema.Attributes, scopes.DBaaS.Build())
}

func (r *dbaasRestoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DBaasRestoreModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// CreateBackup API call logic
	restore, err := r.client.DBaas.CreateRestore(
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.DBaasId.ValueInt64()),
		data.BackupId.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating Restore",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(restore.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	restore, err = r.waitUntilActive(ctx,
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.DBaasId.ValueInt64()),
		data.BackupId.ValueString(),
		restore.Id,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when waiting for Restore to be finished",
			err.Error(),
		)
		return
	}

	if restore == nil {
		return
	}

	data.fill(restore)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dbaasRestoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DBaasRestoreModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	restore, err := r.client.DBaas.GetRestore(
		int(state.PublicCloudId.ValueInt64()),
		int(state.PublicCloudProjectId.ValueInt64()),
		int(state.DBaasId.ValueInt64()),
		state.BackupId.ValueString(),
		state.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading Restore",
			err.Error(),
		)
		return
	}

	state.fill(restore)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dbaasRestoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// NOP
}

func (r *dbaasRestoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// NOP
}

func (r *dbaasRestoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *dbaasRestoreResource) waitUntilActive(ctx context.Context, publicCloudId int, publicCloudProjectId int, dbaasId int, backupId, id string) (*dbaas.DBaaSRestore, error) {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-t.C:
			found, err := r.client.DBaas.GetRestore(publicCloudId, publicCloudProjectId, dbaasId, backupId, id)
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

func (model *DBaasRestoreModel) fill(restore *dbaas.DBaaSRestore) {
	model.Id = types.StringValue(restore.Id)
}
