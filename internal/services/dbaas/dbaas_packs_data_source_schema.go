package dbaas

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func getDbaasPacksDataSourceSchema() schema.Schema {
	pricingObject := schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"hour_excl_tax": schema.Float64Attribute{
				Computed: true,
			},
			"hour_incl_tax": schema.Float64Attribute{
				Computed: true,
			},
		},
	}
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Required: true,
			},
			"packs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"group": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"instances": schema.Int64Attribute{
							Computed: true,
						},
						"cpu": schema.Int64Attribute{
							Computed: true,
						},
						"ram": schema.Int64Attribute{
							Computed: true,
						},
						"storage": schema.Int64Attribute{
							Computed: true,
						},
						"rates": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"chf": pricingObject,
								"eur": pricingObject,
							},
						},
					},
				},
			},
		},
	}
}
