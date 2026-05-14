package publiccloud

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func getPublicCloudConfigDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Reads the Public Cloud configuration of the given account (free tier, resource level, validity window).",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Infomaniak account identifier. Discover it with `curl -s -H \"Authorization: Bearer $INFOMANIAK_TOKEN\" https://api.infomaniak.com/2/profile | jq '.data.preferences.account.current_account_id'`.",
			},
			"free_tier":              schema.Float64Attribute{Computed: true, MarkdownDescription: "Free-tier amount limit, in the account currency."},
			"free_tier_used":         schema.Float64Attribute{Computed: true, MarkdownDescription: "Free-tier amount already consumed."},
			"account_resource_level": schema.Int64Attribute{Computed: true, MarkdownDescription: "Resource level of the account (1–4)."},
			"project_count":          schema.Int64Attribute{Computed: true, MarkdownDescription: "Number of projects currently owned by the account."},
			"valid_from":             schema.Int64Attribute{Computed: true, MarkdownDescription: "Timestamp this configuration became valid. May be 0 when the API returns null."},
			"valid_to":               schema.Int64Attribute{Computed: true, MarkdownDescription: "Timestamp this configuration expires. May be 0 when the API returns null."},
		},
	}
}
