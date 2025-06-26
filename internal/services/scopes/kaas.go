package scopes

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var KaaS = PublicCloud.Subscope(map[string]schema.Attribute{
	"kaas_id": schema.Int64Attribute{
		Required:    true,
		Description: "The id of the kaas project.",
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
	},
})
