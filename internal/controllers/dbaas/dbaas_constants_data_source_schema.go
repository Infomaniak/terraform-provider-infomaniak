package dbaas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getDbaasConstantsDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"types": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"versions": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
				Computed: true,
			},
			"regions": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}
