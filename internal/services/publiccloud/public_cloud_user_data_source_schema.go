package publiccloud

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func getPublicCloudUserDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Reads a Public Cloud user (OpenStack identity) by id.",
		Attributes: map[string]schema.Attribute{
			"public_cloud_id":         schema.Int64Attribute{Required: true, MarkdownDescription: "Identifier of the parent Public Cloud product."},
			"public_cloud_project_id": schema.Int64Attribute{Required: true, MarkdownDescription: "Identifier of the parent project."},
			"id":                      schema.Int64Attribute{Required: true, MarkdownDescription: "Identifier of the user."},
			"open_stack_name":         schema.StringAttribute{Computed: true, MarkdownDescription: "Underlying OpenStack user name."},
			"description":             schema.StringAttribute{Computed: true, MarkdownDescription: "Free-form description."},
			"status":                  schema.StringAttribute{Computed: true, MarkdownDescription: "Lifecycle status."},
			"created_at":              schema.Int64Attribute{Computed: true, MarkdownDescription: "Creation timestamp (Unix seconds)."},
			"updated_at":              schema.Int64Attribute{Computed: true, MarkdownDescription: "Last update timestamp (Unix seconds)."},
		},
	}
}
