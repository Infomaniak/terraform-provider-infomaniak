package publiccloud

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func getPublicCloudDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Reads an Infomaniak Public Cloud product by id.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Unique identifier of the Public Cloud product.",
			},
			"account_id":                schema.Int64Attribute{Computed: true, MarkdownDescription: "Identifier of the Infomaniak account owning the product."},
			"service_id":                schema.Int64Attribute{Computed: true, MarkdownDescription: "Identifier of the service (140 for Public Cloud)."},
			"service_name":              schema.StringAttribute{Computed: true, MarkdownDescription: "Service name."},
			"customer_name":             schema.StringAttribute{Computed: true, MarkdownDescription: "Customer-defined name for the Public Cloud product."},
			"internal_name":             schema.StringAttribute{Computed: true, MarkdownDescription: "Internal Infomaniak identifier (PC-...)."},
			"description":               schema.StringAttribute{Computed: true, MarkdownDescription: "Free-form description."},
			"bill_reference":            schema.StringAttribute{Computed: true, MarkdownDescription: "Billing reference set on the product."},
			"created_at":                schema.Int64Attribute{Computed: true, MarkdownDescription: "Creation timestamp (Unix seconds)."},
			"expired_at":                schema.Int64Attribute{Computed: true, MarkdownDescription: "Expiration timestamp (Unix seconds)."},
			"is_free":                   schema.BoolAttribute{Computed: true, MarkdownDescription: "Is the product billed at zero (free tier)."},
			"is_zero_price":             schema.BoolAttribute{Computed: true, MarkdownDescription: "Is the product priced at zero."},
			"is_trial":                  schema.BoolAttribute{Computed: true, MarkdownDescription: "Is the product currently in trial mode."},
			"is_locked":                 schema.BoolAttribute{Computed: true, MarkdownDescription: "Is the product locked."},
			"has_maintenance":           schema.BoolAttribute{Computed: true, MarkdownDescription: "Is a maintenance currently active."},
			"has_operation_in_progress": schema.BoolAttribute{Computed: true, MarkdownDescription: "Is an asynchronous operation currently running."},
		},
	}
}
