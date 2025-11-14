package kaas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getKaasRegionDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"items": schema.ListAttribute{
				Computed:    true,
				Description: "A list of every available regions",
				ElementType: types.StringType,
			},
		},
		MarkdownDescription: "The kaas region data source allows the user to get all available regions.",
	}
}
