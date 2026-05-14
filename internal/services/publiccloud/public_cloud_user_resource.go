package publiccloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/publiccloud"
	"terraform-provider-infomaniak/internal/provider"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                   = &publicCloudUserResource{}
	_ resource.ResourceWithConfigure      = &publicCloudUserResource{}
	_ resource.ResourceWithImportState    = &publicCloudUserResource{}
	_ resource.ResourceWithValidateConfig = &publicCloudUserResource{}
)

const userCreateTimeout = 5 * time.Minute

type publicCloudUserResource struct {
	client *apis.Client
}

func NewPublicCloudUserResource() resource.Resource {
	return &publicCloudUserResource{}
}

// PublicCloudUserResourceModel represents the resource state. Email and
// password are sensitive and stored in state; the API never returns them so
// drift detection on those fields is impossible.
type PublicCloudUserResourceModel struct {
	PublicCloudId        types.Int64  `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64  `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64  `tfsdk:"id"`
	Description          types.String `tfsdk:"description"`

	Invite   types.Bool   `tfsdk:"invite"`
	Email    types.String `tfsdk:"email"`
	Password types.String `tfsdk:"password"`

	OpenStackName types.String `tfsdk:"open_stack_name"`
	Status        types.String `tfsdk:"status"`
	CreatedAt     types.Int64  `tfsdk:"created_at"`
	UpdatedAt     types.Int64  `tfsdk:"updated_at"`
}

func (r *publicCloudUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_user"
}

func (r *publicCloudUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", err.Error())
		return
	}
	r.client = client
}

func (r *publicCloudUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = getPublicCloudUserResourceSchema()
}

// ValidateConfig enforces the invite-vs-direct invariant *at create time only*.
// After creation, both fields are independently PATCH-able.
func (r *publicCloudUserResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg PublicCloudUserResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if cfg.Invite.IsUnknown() || cfg.Email.IsUnknown() || cfg.Password.IsUnknown() {
		return
	}

	invite := cfg.Invite.ValueBool()
	emailSet := !cfg.Email.IsNull() && cfg.Email.ValueString() != ""
	passwordSet := !cfg.Password.IsNull() && cfg.Password.ValueString() != ""

	if invite {
		if !emailSet {
			resp.Diagnostics.AddAttributeError(path.Root("email"),
				"Missing email",
				"`email` is required when `invite = true` so Infomaniak can send the invitation.")
		}
	} else {
		if !passwordSet {
			resp.Diagnostics.AddAttributeError(path.Root("password"),
				"Missing password",
				"`password` is required when `invite = false` (direct user creation).")
		}
	}
}

func (r *publicCloudUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PublicCloudUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &publiccloud.UserCreate{
		Description: plan.Description.ValueString(),
		Email:       plan.Email.ValueString(),
		Password:    plan.Password.ValueString(),
		Invite:      plan.Invite.ValueBool(),
	}

	id, err := r.client.PublicCloud.CreateUser(plan.PublicCloudId.ValueInt64(), plan.PublicCloudProjectId.ValueInt64(), input)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Public Cloud user", err.Error())
		return
	}

	if err := publiccloud.WaitForStatusOk(ctx, func() (string, error) {
		u, e := r.client.PublicCloud.GetUser(plan.PublicCloudId.ValueInt64(), plan.PublicCloudProjectId.ValueInt64(), id)
		if e != nil {
			return "", e
		}
		return u.Status, nil
	}, userCreateTimeout); err != nil {
		resp.Diagnostics.AddError("User did not reach ok status", err.Error())
		return
	}

	obj, err := r.client.PublicCloud.GetUser(plan.PublicCloudId.ValueInt64(), plan.PublicCloudProjectId.ValueInt64(), id)
	if err != nil {
		resp.Diagnostics.AddError("Unable to refresh user after create", err.Error())
		return
	}

	plan.Id = types.Int64Value(id)
	fillUserResourceModel(&plan, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *publicCloudUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PublicCloudUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := r.client.PublicCloud.GetUser(data.PublicCloudId.ValueInt64(), data.PublicCloudProjectId.ValueInt64(), data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Public Cloud user", err.Error())
		return
	}

	fillUserResourceModel(&data, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *publicCloudUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PublicCloudUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &publiccloud.UserUpdate{}
	if !plan.Description.Equal(state.Description) {
		input.Description = plan.Description.ValueString()
	}
	if !plan.Email.Equal(state.Email) {
		input.Email = plan.Email.ValueString()
	}
	if !plan.Password.Equal(state.Password) {
		input.Password = plan.Password.ValueString()
	}

	if _, err := r.client.PublicCloud.UpdateUser(plan.PublicCloudId.ValueInt64(), plan.PublicCloudProjectId.ValueInt64(), plan.Id.ValueInt64(), input); err != nil {
		resp.Diagnostics.AddError("Unable to update Public Cloud user", err.Error())
		return
	}

	obj, err := r.client.PublicCloud.GetUser(plan.PublicCloudId.ValueInt64(), plan.PublicCloudProjectId.ValueInt64(), plan.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to refresh user after update", err.Error())
		return
	}

	fillUserResourceModel(&plan, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *publicCloudUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PublicCloudUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.PublicCloud.DeleteUser(state.PublicCloudId.ValueInt64(), state.PublicCloudProjectId.ValueInt64(), state.Id.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Unable to delete Public Cloud user", err.Error())
		return
	}
}

// ImportState expects "<public_cloud_id>,<project_id>,<user_id>".
func (r *publicCloudUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected `<public_cloud_id>,<project_id>,<user_id>`; got: %q", req.ID),
		)
		return
	}
	cloudId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", "public_cloud_id must be numeric: "+parts[0])
		return
	}
	projectId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", "project_id must be numeric: "+parts[1])
		return
	}
	userId, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", "user_id must be numeric: "+parts[2])
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_id"), types.Int64Value(cloudId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_project_id"), types.Int64Value(projectId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(userId))...)
}

func fillUserResourceModel(m *PublicCloudUserResourceModel, obj *publiccloud.User) {
	m.Id = types.Int64Value(obj.Id)
	m.PublicCloudProjectId = types.Int64Value(obj.PublicCloudProjectId)
	m.OpenStackName = types.StringValue(obj.OpenStackName)
	m.Description = types.StringValue(obj.Description)
	m.Status = types.StringValue(obj.Status)
	m.CreatedAt = types.Int64Value(obj.CreatedAt)
	m.UpdatedAt = types.Int64Value(obj.UpdatedAt)
}
