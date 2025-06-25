package dbaas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var publicCloudAttributes = map[string]schema.Attribute{
	"public_cloud_id": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The id of the public cloud where DBaaS is installed",
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
	},
	"public_cloud_project_id": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The id of the public cloud project where DBaaS is installed",
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
	},
}

var publicCloudAttributes = map[string]schema.Attribute{
	"public_cloud_id": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The id of the public cloud where DBaaS is installed",
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
	},
	"public_cloud_project_id": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The id of the public cloud project where DBaaS is installed",
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
	},
}
