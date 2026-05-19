package publiccloud

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// descriptionPattern matches the live-API constraint observed on user
// description fields: letters, digits, spaces, hyphens, underscores, periods.
// The OpenAPI spec doesn't document this; we mirror the live error message.
var descriptionPattern = regexp.MustCompile(`^[A-Za-z0-9 \-_.]+$`)

func getPublicCloudUserResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages an OpenStack user inside a Public Cloud project. Two creation modes are " +
			"supported: direct (`invite = false`, `password` required) and invitation (`invite = true`, " +
			"`email` required). After creation, `email`, `password` and `description` are independently " +
			"PATCH-able. Sensitive fields are stored in state (the API never returns them, so drift detection " +
			"on those values is not possible).",
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Identifier of the parent Public Cloud product. Forces replacement on change.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Identifier of the parent project. Forces replacement on change.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the user.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Free-form description. The API rejects values containing characters " +
					"outside `[A-Za-z0-9 -_.]` and longer than 250 characters.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
					stringvalidator.RegexMatches(descriptionPattern,
						"must contain only letters, numbers, spaces, hyphens, underscores and periods"),
				},
			},
			"invite": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "If `true`, the user is created via invitation email instead of direct password " +
					"creation. Forces replacement on change.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"email": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Email address. Required when `invite = true`. PATCH-able after creation.",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password. Required when `invite = false`. PATCH-able after creation. Stored in state.",
			},
			"open_stack_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Underlying OpenStack user name.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"status":     schema.StringAttribute{Computed: true, MarkdownDescription: "Lifecycle status of the user."},
			"created_at": schema.Int64Attribute{Computed: true, MarkdownDescription: "Creation timestamp (Unix seconds).", PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"updated_at": schema.Int64Attribute{Computed: true, MarkdownDescription: "Last update timestamp (Unix seconds)."},
		},
	}
}
