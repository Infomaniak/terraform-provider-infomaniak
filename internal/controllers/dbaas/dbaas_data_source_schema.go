package dbaas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getDbaasDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud where DBaaS is installed",
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of the public cloud project where DBaaS is installed",
			},
			"id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The id of this DBaaS",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the DBaaS project",
			},
			"pack_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the pack associated to the DBaaS project",
			},
			"region": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The region where the DBaaS project resides in.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the database associated with the DBaaS project",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The version of the database associated with the DBaaS project",
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
			"ca": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Database CA Certificate",
			},
			"allowed_cidrs": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Allowed to query Database IP whitelist",
			},
			"kube_identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DbaaS kubernetes name",
			},
			"effective_configuration": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: "The effective database configuration settings",
			},
		},
		MarkdownDescription: "The dbaas data source allows the user to manage a dbaas project",
	}
}
