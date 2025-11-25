package dbaas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getDbaasResourceSchema() schema.Schema {
	return schema.Schema{
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
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the DBaaS instance.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
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
			"ca": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Database CA Certificate",
			},
			"allowed_cidrs": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Allowed to query Database IP whitelist",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"kube_identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DbaaS kubernetes name",
			},
			"configuration": schema.MapAttribute{
				Computed: true,
				Optional: true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"effective_configuration": schema.MapAttribute{
				Computed: true,
				ElementType: types.StringType,
			},
		},
		MarkdownDescription: "The dbaas resource allows the user to manage a dbaas project",
	}
}
