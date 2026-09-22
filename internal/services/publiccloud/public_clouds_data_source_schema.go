package publiccloud

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func getPublicCloudsDataSourceSchema() schema.Schema {
	itemAttrs := map[string]schema.Attribute{
		"id":                        schema.Int64Attribute{Computed: true, MarkdownDescription: "Public Cloud product identifier."},
		"account_id":                schema.Int64Attribute{Computed: true, MarkdownDescription: "Account identifier."},
		"service_id":                schema.Int64Attribute{Computed: true, MarkdownDescription: "Service identifier."},
		"service_name":              schema.StringAttribute{Computed: true, MarkdownDescription: "Service name."},
		"customer_name":             schema.StringAttribute{Computed: true, MarkdownDescription: "Customer-defined name."},
		"internal_name":             schema.StringAttribute{Computed: true, MarkdownDescription: "Internal Infomaniak identifier."},
		"description":               schema.StringAttribute{Computed: true, MarkdownDescription: "Free-form description."},
		"bill_reference":            schema.StringAttribute{Computed: true, MarkdownDescription: "Billing reference."},
		"created_at":                schema.Int64Attribute{Computed: true, MarkdownDescription: "Creation timestamp (Unix seconds)."},
		"expired_at":                schema.Int64Attribute{Computed: true, MarkdownDescription: "Expiration timestamp (Unix seconds)."},
		"is_free":                   schema.BoolAttribute{Computed: true, MarkdownDescription: "Free tier flag."},
		"is_zero_price":             schema.BoolAttribute{Computed: true, MarkdownDescription: "Zero price flag."},
		"is_trial":                  schema.BoolAttribute{Computed: true, MarkdownDescription: "Trial flag."},
		"is_locked":                 schema.BoolAttribute{Computed: true, MarkdownDescription: "Locked flag."},
		"has_maintenance":           schema.BoolAttribute{Computed: true, MarkdownDescription: "Maintenance flag."},
		"has_operation_in_progress": schema.BoolAttribute{Computed: true, MarkdownDescription: "Operation-in-progress flag."},
	}

	return schema.Schema{
		MarkdownDescription: "Lists all Public Cloud products owned by the given account.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Infomaniak account identifier. Discover it with `curl -s -H \"Authorization: Bearer $INFOMANIAK_TOKEN\" https://api.infomaniak.com/2/profile | jq '.data.preferences.account.current_account_id'`.",
			},
			"public_clouds": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of Public Cloud products.",
				NestedObject:        schema.NestedAttributeObject{Attributes: itemAttrs},
			},
		},
	}
}
