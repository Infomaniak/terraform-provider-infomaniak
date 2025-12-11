package kaas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func getKaasPackDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The id of the KaaS pack.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the KaaS pack",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the KaaS pack.",
			},
			"price_per_hour": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"chf": schema.Float64Attribute{
						Computed: true,
					},
					"eur": schema.Float64Attribute{
						Computed: true,
					},
				},
			},
			"limit_per_project": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The maximum number of Kubernetes services of the given pack per project.",
			},
			"is_active": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Is the pack currently active for use in Kubernetes service creations.",
			},
		},
		MarkdownDescription: "The kaas pack data source allows the user to retrieve information about a kaas pack.",
	}
}
