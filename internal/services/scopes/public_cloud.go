package scopes

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var PublicCloud = New(map[string]schema.Attribute{
	"public_cloud_id": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The id of the public cloud",
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
	},
	"public_cloud_project_id": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The id of the public cloud project",
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
	},
})
