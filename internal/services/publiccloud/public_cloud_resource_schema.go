package publiccloud

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func getPublicCloudResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages an Infomaniak Public Cloud product. " +
			"Public Cloud products cannot be ordered through the Infomaniak public API; this resource " +
			"can only manage products that already exist. Order the product via the Manager UI, then run " +
			"`terraform import infomaniak_public_cloud.<name> <public_cloud_id>` to bring it under Terraform. " +
			"Removing the resource only clears Terraform state — the product remains billed at Infomaniak.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the Public Cloud product. Set during import.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"customer_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Customer-defined name for the Public Cloud product. Minimum 2 characters; the API rejects empty values.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:          []validator.String{stringvalidator.LengthAtLeast(2)},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Free-form description. Minimum 2 characters; the API rejects empty values.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:          []validator.String{stringvalidator.LengthAtLeast(2)},
			},
			"bill_reference": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Billing reference set on the product. Minimum 2 characters; the API rejects empty values and once set the value cannot be cleared via PATCH.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:          []validator.String{stringvalidator.LengthAtLeast(2)},
			},
			"account_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the Infomaniak account owning the product.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"service_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the Infomaniak service for Public Cloud.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"service_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service name.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"internal_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal Infomaniak identifier. May be empty.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"created_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp (Unix seconds).",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"expired_at":                schema.Int64Attribute{Computed: true, MarkdownDescription: "Expiration timestamp (Unix seconds)."},
			"is_free":                   schema.BoolAttribute{Computed: true, MarkdownDescription: "Is the product billed at zero (free tier).", PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}},
			"is_zero_price":             schema.BoolAttribute{Computed: true, MarkdownDescription: "Is the product priced at zero.", PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}},
			"is_trial":                  schema.BoolAttribute{Computed: true, MarkdownDescription: "Is the product currently in trial mode.", PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}},
			"is_locked":                 schema.BoolAttribute{Computed: true, MarkdownDescription: "Is the product locked."},
			"has_maintenance":           schema.BoolAttribute{Computed: true, MarkdownDescription: "Is a maintenance currently active."},
			"has_operation_in_progress": schema.BoolAttribute{Computed: true, MarkdownDescription: "Is an asynchronous operation currently running."},
		},
	}
}
