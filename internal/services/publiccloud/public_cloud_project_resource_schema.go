package publiccloud

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func getPublicCloudProjectResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Public Cloud project (OpenStack tenant). Creating a project also " +
			"bootstraps an admin user: either directly with `user_password` (and optional `user_email`), or " +
			"by setting `invite = true` and providing `user_email` so Infomaniak sends an invitation email.",
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Identifier of the parent Public Cloud product. Forces replacement on change.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the project.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project name. Max 250 characters.",
				Validators:          []validator.String{stringvalidator.LengthAtMost(250)},
			},
			"invite": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "If `true`, the bootstrap user is created via invitation email instead of " +
					"direct password creation. Forces replacement on change.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"user_email": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Email address of the bootstrap admin user. Required when `invite = true`. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				MarkdownDescription: "Initial password for the bootstrap admin user. Required when `invite = false`. " +
					"Stored in Terraform state (sensitive). Forces replacement.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_description": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Description applied to the bootstrap admin user. The API rejects values containing " +
					"characters outside `[A-Za-z0-9 -_.]` and longer than 250 characters. Forces replacement.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
					stringvalidator.RegexMatches(descriptionPattern,
						"must contain only letters, numbers, spaces, hyphens, underscores and periods"),
				},
			},
			"open_stack_name":  schema.StringAttribute{Computed: true, MarkdownDescription: "Underlying OpenStack tenant name.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"status":           schema.StringAttribute{Computed: true, MarkdownDescription: "Lifecycle status of the project."},
			"price":            schema.Float64Attribute{Computed: true, MarkdownDescription: "Current cumulative cost."},
			"resource_level":   schema.Int64Attribute{Computed: true, MarkdownDescription: "Resource-level limit."},
			"user_count":       schema.Int64Attribute{Computed: true, MarkdownDescription: "Total number of users."},
			"created_at":       schema.Int64Attribute{Computed: true, MarkdownDescription: "Creation timestamp (Unix seconds).", PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"updated_at":       schema.Int64Attribute{Computed: true, MarkdownDescription: "Last update timestamp (Unix seconds)."},
			"billing_start_at": schema.Int64Attribute{Computed: true, MarkdownDescription: "Start of the current billing window."},
			"billing_end_at":   schema.Int64Attribute{Computed: true, MarkdownDescription: "End of the current billing window."},
			"price_updated_at": schema.Int64Attribute{Computed: true, MarkdownDescription: "Last time the price was recomputed."},
		},
	}
}
