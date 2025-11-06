package kaas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getKaasInstancePoolResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud where KaaS is installed",
				MarkdownDescription: "The id of the public cloud where KaaS is installed",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud project where KaaS is installed",
				MarkdownDescription: "The id of the public cloud project where KaaS is installed",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"kaas_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the kaas project.",
				MarkdownDescription: "The id of the kaas project.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "A computed value representing the unique identifier for the instance pool. Mandatory for acceptance testing.",
				MarkdownDescription: "A computed value representing the unique identifier for the instance pool. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the instance pool",
				MarkdownDescription: "The name of the instance pool",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"availability_zone": schema.StringAttribute{
				Required:            true,
				Description:         "The availability zone for the instances in the pool",
				MarkdownDescription: "The availability zone for the instances in the pool",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"flavor_name": schema.StringAttribute{
				Required:            true,
				Description:         "The flavor name for the instances in the pool",
				MarkdownDescription: "The flavor name for the instances in the pool",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"min_instances": schema.Int64Attribute{
				Required:            true,
				Description:         "The minimum amount of instances in this instance pool",
				MarkdownDescription: "The minimum amount of instances in this instance pool",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_instances": schema.Int64Attribute{
				Required:            true,
				Description:         "The maximum amount of instances in this instance pool",
				MarkdownDescription: "The maximum amount of instances in this instance pool",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"labels": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
					mapplanmodifier.RequiresReplace(),
				},
				Description:         "Kubernetes labels to apply to the instances. The label must have a prefix of node-role.kubernetes.io or belong to the domains node-restriction.kubernetes.io or custom.kaas.infomaniak.cloud.",
				MarkdownDescription: "Kubernetes labels to apply to the instances. The label must have a prefix of node-role.kubernetes.io or belong to the domains node-restriction.kubernetes.io or custom.kaas.infomaniak.cloud.",
			},
		},
		MarkdownDescription: "The kaas instance pool resource is used to manage instance pools inside a kaas project",
	}
}
