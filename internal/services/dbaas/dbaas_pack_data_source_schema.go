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
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"group": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"instances": schema.Int64Attribute{
				Optional: true,
			},
			"cpu": schema.Int64Attribute{
				Optional: true,
			},
			"ram": schema.Int64Attribute{
				Optional: true,
			},
			"storage": schema.Int64Attribute{
				Optional: true,
			},
			"rates": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"chf": pricingObject,
					"eur": pricingObject,
				},
			},
		},
	}
}
