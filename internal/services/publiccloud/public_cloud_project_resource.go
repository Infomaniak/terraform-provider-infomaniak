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
	_ resource.Resource                   = &publicCloudProjectResource{}
	_ resource.ResourceWithConfigure      = &publicCloudProjectResource{}
	_ resource.ResourceWithImportState    = &publicCloudProjectResource{}
	_ resource.ResourceWithValidateConfig = &publicCloudProjectResource{}
)

const projectCreateTimeout = 10 * time.Minute

type publicCloudProjectResource struct {
	client *apis.Client
}

func NewPublicCloudProjectResource() resource.Resource {
	return &publicCloudProjectResource{}
}

type PublicCloudProjectResourceModel struct {
	PublicCloudId types.Int64  `tfsdk:"public_cloud_id"`
	Id            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`

	Invite          types.Bool   `tfsdk:"invite"`
	UserEmail       types.String `tfsdk:"user_email"`
	UserPassword    types.String `tfsdk:"user_password"`
	UserDescription types.String `tfsdk:"user_description"`

	OpenStackName  types.String  `tfsdk:"open_stack_name"`
	Status         types.String  `tfsdk:"status"`
	Price          types.Float64 `tfsdk:"price"`
	ResourceLevel  types.Int64   `tfsdk:"resource_level"`
	UserCount      types.Int64   `tfsdk:"user_count"`
	CreatedAt      types.Int64   `tfsdk:"created_at"`
	UpdatedAt      types.Int64   `tfsdk:"updated_at"`
	BillingStartAt types.Int64   `tfsdk:"billing_start_at"`
	BillingEndAt   types.Int64   `tfsdk:"billing_end_at"`
	PriceUpdatedAt types.Int64   `tfsdk:"price_updated_at"`
}

func (r *publicCloudProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_project"
}

func (r *publicCloudProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *publicCloudProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = getPublicCloudProjectResourceSchema()
}

// ValidateConfig enforces the invite-vs-direct invariant: invite=true requires
// user_email and forbids user_password; invite=false requires user_password.
func (r *publicCloudProjectResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg PublicCloudProjectResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Unknown values mean a referenced input — defer to apply.
	if cfg.Invite.IsUnknown() || cfg.UserEmail.IsUnknown() || cfg.UserPassword.IsUnknown() {
		return
	}

	invite := cfg.Invite.ValueBool()
	emailSet := !cfg.UserEmail.IsNull() && cfg.UserEmail.ValueString() != ""
	passwordSet := !cfg.UserPassword.IsNull() && cfg.UserPassword.ValueString() != ""

	if invite {
		if !emailSet {
			resp.Diagnostics.AddAttributeError(path.Root("user_email"),
				"Missing user_email",
				"`user_email` is required when `invite = true` so Infomaniak can send the invitation.")
		}
		if passwordSet {
			resp.Diagnostics.AddAttributeError(path.Root("user_password"),
				"Unexpected user_password",
				"`user_password` must not be set when `invite = true`; the invited user picks their own password.")
		}
	} else {
		if !passwordSet {
			resp.Diagnostics.AddAttributeError(path.Root("user_password"),
				"Missing user_password",
				"`user_password` is required when `invite = false` (direct creation of the bootstrap user).")
		}
	}
}

func (r *publicCloudProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PublicCloudProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &publiccloud.ProjectCreate{
		Name:            plan.Name.ValueString(),
		UserEmail:       plan.UserEmail.ValueString(),
		UserPassword:    plan.UserPassword.ValueString(),
		UserDescription: plan.UserDescription.ValueString(),
		Invite:          plan.Invite.ValueBool(),
	}

	id, err := r.client.PublicCloud.CreateProject(plan.PublicCloudId.ValueInt64(), input)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Public Cloud project", err.Error())
		return
	}

	if err := publiccloud.WaitForStatusOk(ctx, func() (string, error) {
		p, e := r.client.PublicCloud.GetProject(plan.PublicCloudId.ValueInt64(), id)
		if e != nil {
			return "", e
		}
		return p.Status, nil
	}, projectCreateTimeout); err != nil {
		resp.Diagnostics.AddError("Project did not reach ok status", err.Error())
		return
	}

	obj, err := r.client.PublicCloud.GetProject(plan.PublicCloudId.ValueInt64(), id)
	if err != nil {
		resp.Diagnostics.AddError("Unable to refresh project after create", err.Error())
		return
	}

	plan.Id = types.Int64Value(id)
	fillProjectResourceModel(&plan, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *publicCloudProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PublicCloudProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := r.client.PublicCloud.GetProject(data.PublicCloudId.ValueInt64(), data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Public Cloud project", err.Error())
		return
	}

	fillProjectResourceModel(&data, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *publicCloudProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PublicCloudProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &publiccloud.Project{
		Id:            plan.Id.ValueInt64(),
		PublicCloudId: plan.PublicCloudId.ValueInt64(),
		Name:          plan.Name.ValueString(),
	}
	if err := r.client.PublicCloud.UpdateProject(input); err != nil {
		resp.Diagnostics.AddError("Unable to update Public Cloud project", err.Error())
		return
	}

	obj, err := r.client.PublicCloud.GetProject(plan.PublicCloudId.ValueInt64(), plan.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to refresh project after update", err.Error())
		return
	}

	fillProjectResourceModel(&plan, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *publicCloudProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PublicCloudProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.PublicCloud.DeleteProject(state.PublicCloudId.ValueInt64(), state.Id.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Unable to delete Public Cloud project", err.Error())
		return
	}
}

// ImportState expects "<public_cloud_id>,<project_id>".
func (r *publicCloudProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			errInvalidImportID,
			fmt.Sprintf("Expected `<public_cloud_id>,<project_id>`; got: %q", req.ID),
		)
		return
	}
	publicCloudId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(errInvalidImportID, "public_cloud_id must be numeric: "+parts[0])
		return
	}
	projectId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(errInvalidImportID, "project_id must be numeric: "+parts[1])
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_id"), types.Int64Value(publicCloudId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(projectId))...)
}

func fillProjectResourceModel(m *PublicCloudProjectResourceModel, obj *publiccloud.Project) {
	m.PublicCloudId = types.Int64Value(obj.PublicCloudId)
	m.Id = types.Int64Value(obj.Id)
	m.Name = types.StringValue(obj.Name)
	m.OpenStackName = types.StringValue(obj.OpenStackName)
	m.Status = types.StringValue(obj.Status)
	m.Price = types.Float64Value(obj.Price)
	m.ResourceLevel = types.Int64Value(obj.ResourceLevel)
	m.UserCount = types.Int64Value(obj.UserCount)
	m.CreatedAt = types.Int64Value(obj.CreatedAt)
	m.UpdatedAt = types.Int64Value(obj.UpdatedAt)
	m.BillingStartAt = types.Int64Value(obj.BillingStartAt)
	m.BillingEndAt = types.Int64Value(obj.BillingEndAt)
	m.PriceUpdatedAt = types.Int64Value(obj.PriceUpdatedAt)
}
