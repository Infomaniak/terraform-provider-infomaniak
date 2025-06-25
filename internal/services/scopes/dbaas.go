package scopes

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

var DBaaS = PublicCloud.Subscope(map[string]schema.Attribute{
	"dbaas_id": schema.Int64Attribute{
		Required:    true,
		Description: "The id of the dbaas project.",
	},
})
