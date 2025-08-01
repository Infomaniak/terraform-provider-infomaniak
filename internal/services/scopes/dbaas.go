package scopes

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var DBaaS = PublicCloud.Subscope(map[string]schema.Attribute{
	"dbaas_id": schema.Int64Attribute{
		Required:    true,
		Description: "The id of the dbaas project.",
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
	},
})
