package publiccloud

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func getPublicCloudProjectDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Reads a Public Cloud project by id.",
		Attributes: map[string]schema.Attribute{
			"public_cloud_id":  schema.Int64Attribute{Required: true, MarkdownDescription: "Identifier of the parent Public Cloud product."},
			"id":               schema.Int64Attribute{Required: true, MarkdownDescription: "Identifier of the project."},
			"name":             schema.StringAttribute{Computed: true, MarkdownDescription: "Project name."},
			"open_stack_name":  schema.StringAttribute{Computed: true, MarkdownDescription: "Underlying OpenStack tenant name."},
			"status":           schema.StringAttribute{Computed: true, MarkdownDescription: "Lifecycle status (`creating`, `ok`, `error`, …)."},
			"price":            schema.Float64Attribute{Computed: true, MarkdownDescription: "Current cumulative cost of the project."},
			"resource_level":   schema.Int64Attribute{Computed: true, MarkdownDescription: "Resource-level limit of the project."},
			"user_count":       schema.Int64Attribute{Computed: true, MarkdownDescription: "Total number of users in the project."},
			"created_at":       schema.Int64Attribute{Computed: true, MarkdownDescription: "Creation timestamp (Unix seconds)."},
			"updated_at":       schema.Int64Attribute{Computed: true, MarkdownDescription: "Last update timestamp (Unix seconds)."},
			"billing_start_at": schema.Int64Attribute{Computed: true, MarkdownDescription: "Start of the current billing window."},
			"billing_end_at":   schema.Int64Attribute{Computed: true, MarkdownDescription: "End of the current billing window."},
			"price_updated_at": schema.Int64Attribute{Computed: true, MarkdownDescription: "Last time the price was recomputed."},
		},
	}
}
