package publiccloud

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func getPublicCloudAccessesDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Reports the maintenance status of the Public Cloud API surface for the given account. " +
			"Despite the endpoint name, it does not list region accesses.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Infomaniak account identifier. Discover it with `curl -s -H \"Authorization: Bearer $INFOMANIAK_TOKEN\" https://api.infomaniak.com/2/profile | jq '.data.preferences.account.current_account_id'`.",
			},
			"is_maintenance_ongoing": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True when a Public Cloud maintenance window is currently active for this account.",
			},
		},
	}
}
