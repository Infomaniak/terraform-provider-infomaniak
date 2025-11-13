package project

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/project"
	"terraform-provider-infomaniak/internal/provider"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &projectResource{}
	_ resource.ResourceWithConfigure = &projectResource{}
)

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

type projectResource struct {
	client *apis.Client
}

type ProjectModel struct {
	PublicCloudId types.Int64 `tfsdk:"public_cloud_id"`
	Id            types.Int64 `tfsdk:"id"`

	Name            types.String `tfsdk:"name"`
	UserDescription types.String `tfsdk:"user_description"`
	UserEmail       types.String `tfsdk:"user_email"`
	UserPassword    types.String `tfsdk:"user_password"`

	OpenStackName  types.String `tfsdk:"open_stack_name"`
	Price          types.Int64  `tfsdk:"price"`
	ResourceLevel  types.Int64  `tfsdk:"resource_level"`
	Status         types.String `tfsdk:"status"`
	PriceUpdatedAt types.Int64  `tfsdk:"price_updated_at"`
	BillingStartAt types.Int64  `tfsdk:"billing_start_at"`
	BillingEndAt   types.Int64  `tfsdk:"billing_end_at"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
	DeletedAt      types.Int64  `tfsdk:"deleted_at"`
	UserCount      types.Int64  `tfsdk:"user_count"`
	Tags           types.List   `tfsdk:"tags"`
}

type ProjectTagModel struct {
	Id           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Color        types.Int64  `tfsdk:"color"`
	ProductCount types.Int64  `tfsdk:"product_count"`
}

func (r *projectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Project Resource Configure Type",
			err.Error(),
		)
		return
	}

	r.client = client
}

func (r *projectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud where the project is installed",
				MarkdownDescription: "The id of the public cloud where the project is installed",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"user_description": schema.StringAttribute{
				Required:            false,
				Optional:            true,
				MarkdownDescription: "The description of the project",
			},
			"user_email": schema.StringAttribute{
				Required: false,
				Optional: true,

				MarkdownDescription: "The email of the owner creating the project",
			},
			"user_password": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The password of the project?",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"open_stack_name": schema.StringAttribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The name of the open stack project",
			},
			"price": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The current cost of the project",
			},
			"resource_level": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The project resource level limit",
			},
			"status": schema.StringAttribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The project status",
			},
			"price_updated_at": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The date of the last time the price was updated",
			},
			"billing_start_at": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The date of when the billing started",
			},
			"billing_end_at": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The date of when the billing will stop",
			},
			"created_at": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The date of when the project was created.",
			},
			"updated_at": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The date of the last time the project was updated.",
			},
			"deleted_at": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The date of when the project was deleted.",
			},
			"user_count": schema.Int64Attribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The total number of users in the project",
			},
			"tags": schema.ListNestedAttribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "List of resource tags",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Required:            false,
							Computed:            true,
							MarkdownDescription: "The id of the resource tag",
						},
						"name": schema.StringAttribute{
							Required:            false,
							Computed:            true,
							MarkdownDescription: "The name of the resource tag",
						},
						"color": schema.Int64Attribute{
							Required:            false,
							Computed:            true,
							MarkdownDescription: "The color of the resource tag",
						},
						"product_count": schema.Int64Attribute{
							Required:            false,
							Computed:            true,
							MarkdownDescription: "The product count related to the resource tag",
						},
					},
				},
			},
		},
		MarkdownDescription: "The Project resource allows the user to manage a public cloud project for a product",
	}
}

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &project.CreateProject{
		PublicCloudId:   int(data.PublicCloudId.ValueInt64()),
		Name:            data.Name.ValueString(),
		UserDescription: data.UserDescription.ValueString(),
		UserEmail:       data.UserEmail.ValueString(),
		UserPassword:    data.UserPassword.ValueString(),
	}

	projectId, err := r.client.Project.CreateProject(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating Project",
			err.Error(),
		)
		return
	}

	projectObject, err := r.waitUntilActive(ctx, input, projectId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when waiting for Project to be Active",
			err.Error(),
		)
		return
	}

	if projectObject == nil {
		return
	}

	data.fill(projectObject)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	project, err := r.client.Project.GetProject(int(state.PublicCloudId.ValueInt64()), int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading Project",
			err.Error(),
		)
		return
	}

	state.fill(project)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state ProjectModel
	var data ProjectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := &project.UpdateProject{
		Name: state.Name.String(),
	}

	_, err := r.client.Project.UpdateProject(int(state.PublicCloudId.ValueInt64()), int(state.Id.ValueInt64()), input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when updating Project",
			err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteZone API call logic
	_, err := r.client.Project.DeleteProject(int(data.PublicCloudId.ValueInt64()), int(data.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting Project",
			err.Error(),
		)
		return
	}
}

func (r *projectResource) waitUntilActive(ctx context.Context, project *project.CreateProject, projectId int) (*project.Project, error) {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-t.C:
			found, err := r.client.Project.GetProject(project.PublicCloudId, projectId)
			if err != nil {
				return nil, err
			}

			if ctx.Err() != nil {
				return nil, nil
			}

			if found.Status == "ok" {
				return found, nil
			}
		}
	}
}

func (model *ProjectModel) fill(project *project.Project) {
	model.Id = types.Int64Value(int64(project.ProjectId))
	model.Name = types.StringValue(project.Name)
	model.OpenStackName = types.StringValue(project.OpenStackName)
	model.Price = types.Int64Value(int64(project.Price))
	model.ResourceLevel = types.Int64Value(int64(project.ResourceLevel))
	model.Status = types.StringValue(project.Status)
	model.PriceUpdatedAt = types.Int64Value(int64(project.PriceUpdatedAt))
	model.BillingStartAt = types.Int64Value(int64(project.BillingStartAt))
	model.BillingEndAt = types.Int64Value(int64(project.BillingEndAt))
	model.CreatedAt = types.Int64Value(int64(project.CreatedAt))
	model.UpdatedAt = types.Int64Value(int64(project.UpdatedAt))
	model.DeletedAt = types.Int64Value(int64(project.DeletedAt))
	model.UserCount = types.Int64Value(int64(project.UserCount))

	tags := []attr.Value{}
	for _, tag := range project.Tags {
		tagObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"id":            types.Int64Type,
				"name":          types.StringType,
				"color":         types.Int64Type,
				"product_count": types.Int64Type,
			},
			map[string]attr.Value{"id": types.Int64Value(int64(tag.Id)),
				"name":          types.StringValue(tag.Name),
				"color":         types.Int64Value(int64(tag.Color)),
				"product_count": types.Int64Value(int64(tag.ProductCount)),
			},
		)

		tags = append(tags, tagObj)
	}

	tagsList, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":            types.Int64Type,
				"name":          types.StringType,
				"color":         types.Int64Type,
				"product_count": types.Int64Type,
			},
		},
		tags,
	)

	model.Tags = tagsList
}
