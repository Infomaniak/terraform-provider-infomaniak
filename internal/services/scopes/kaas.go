package scopes

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

var KaaS = PublicCloud.Subscope(map[string]schema.Attribute{
	"kaas_id": schema.Int64Attribute{
		Required:    true,
		Description: "The id of the kaas project.",
	},
})
