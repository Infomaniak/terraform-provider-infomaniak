package project

import (
	"context"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *apis.Client
}

func (r *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *projectDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud where the project is installed",
				MarkdownDescription: "The id of the public cloud where the project is installed",
			},
			"id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud project.",
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
				Required:            false,
				Optional:            true,
				MarkdownDescription: "The password of the project?",
			},
			"name": schema.StringAttribute{
				Required:            false,
				Computed:            true,
				MarkdownDescription: "The name of the project",
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
		MarkdownDescription: "The Project data source allows the user to manage a public cloud project",
	}
}

func (r *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ProjectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

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
