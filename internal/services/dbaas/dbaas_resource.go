package dbaas

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/services/scopes"
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
	PublicCloudId        types.Int64 `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64 `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64 `tfsdk:"id"`

	Name     types.String `tfsdk:"name"`
	PackName types.String `tfsdk:"pack_name"`
	Region   types.String `tfsdk:"region"`
	Type     types.String `tfsdk:"type"`
	Version  types.String `tfsdk:"version"`

	Host     types.String `tfsdk:"host"`
	Port     types.String `tfsdk:"port"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
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
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pack_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the pack associated to the DBaaS project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of database associated with the DBaaS being installed",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The version of database associated with the DBaaS being installed",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the DBaaS project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The region where the DBaaS will reside.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The host to access this database.",
			},
			"port": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The port to access this database.",
			},
			"user": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The username to access this database.",
			},
			"password": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The password to access this database.",
			},
		},
		MarkdownDescription: "The dbaas resource allows the user to manage a dbaas project",
	}

	maps.Copy(resp.Schema.Attributes, scopes.PublicCloud.Build())
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
			PublicCloudId: int(data.PublicCloudId.ValueInt64()),
			ProjectId:     int(data.PublicCloudProjectId.ValueInt64()),
		},
		Region:  data.Region.ValueString(),
		Version: data.Version.ValueString(),
		Type:    data.Type.ValueString(),
		Name:    data.Name.ValueString(),
		PackId:  chosenPack.Id,
	}

	// CreateDBaas API call logic
	dbaasId, err := r.client.DBaas.CreateDBaaS(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating DBaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(int64(dbaasId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	dbaasObject, err := r.waitUntilActive(ctx, input, dbaasId)
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

	connectionInfos, err := r.client.DBaas.GetPassword(
		input.Project.PublicCloudId,
		input.Project.ProjectId,
		dbaasId,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading DBaaS connection infos",
			err.Error(),
		)
		return
	}

	data.fill(dbaasObject, connectionInfos)

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
		int(state.PublicCloudId.ValueInt64()),
		int(state.PublicCloudProjectId.ValueInt64()),
		int(state.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading DBaaS",
			err.Error(),
		)
		return
	}

	connectionInfos, err := r.client.DBaas.GetPassword(
		int(state.PublicCloudId.ValueInt64()),
		int(state.PublicCloudProjectId.ValueInt64()),
		int(state.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading DBaaS connection infos",
			err.Error(),
		)
		return
	}

	state.fill(dbaasObject, connectionInfos)

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
			PublicCloudId: int(data.PublicCloudId.ValueInt64()),
			ProjectId:     int(data.PublicCloudProjectId.ValueInt64()),
		},
		Id:      int(state.Id.ValueInt64()),
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

	connectionInfos, err := r.client.DBaas.GetPassword(
		int(state.PublicCloudId.ValueInt64()),
		int(state.PublicCloudProjectId.ValueInt64()),
		int(state.Id.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading DBaaS connection infos",
			err.Error(),
		)
		return
	}

	state.fill(dbaasObject, connectionInfos)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
		int(data.PublicCloudId.ValueInt64()),
		int(data.PublicCloudProjectId.ValueInt64()),
		int(data.Id.ValueInt64()),
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

func (model *DBaasModel) fill(dbaas *dbaas.DBaaS, connectionInfos *dbaas.DBaaSConnectionInfo) {
	model.Id = types.Int64Value(int64(dbaas.Id))
	model.Region = types.StringValue(dbaas.Region)
	model.Type = types.StringValue(dbaas.Type)
	model.Version = types.StringValue(dbaas.Version)
	model.Name = types.StringValue(dbaas.Name)
	model.PackName = types.StringValue(dbaas.Pack.Name)

	model.Host = types.StringValue(connectionInfos.Host)
	model.Port = types.StringValue(connectionInfos.Port)
	model.User = types.StringValue(connectionInfos.User)
	model.Password = types.StringValue(connectionInfos.Password)
}
func (r *dbaasResource) waitUntilActive(ctx context.Context, dbaas *dbaas.DBaaS, id int) (*dbaas.DBaaS, error) {
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
