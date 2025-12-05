package kaas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func getKaasInstancePoolFlavorDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:    true,
				Description: "The id of the public cloud project where KaaS is installed",
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:    true,
				Description: "The id of the public cloud project where KaaS is installed",
			},
			"region": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The region where the kaas will be installed.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name of the KaaS instance pool flavor",
				Validators: []validator.String{
					NameOrResourcesValidator{},
				},
			},
			"cpu": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The numbers of CPU cores that will be allocated to each instances with this flavor.",
			},
			"ram": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The numbers of GB of ram that will be allocated to each instances with this flavor.",
			},
			"storage": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The numbers of GB of storage that will be allocated to each instances with this flavor.",
			},
			"is_available": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Is the flavor available.",
			},
			"is_memory_optimized": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Is the flavor optimized for memory intensive operations.",
			},
			"is_iops_optimized": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Is the flavor optimized for disk intensive operations.",
			},
			"is_gpu_optimized": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Is the flavor optimized for GPU intensive operations.",
			},
			"rates": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"hour_excl_tax": schema.Float64Attribute{
						Computed: true,
					},
					"hour_incl_tax": schema.Float64Attribute{
						Computed: true,
					},
				},
			},
		},
		MarkdownDescription: "The kaas instance pool flavor data source allows the user to manage a kaas instance pool flavor.",
	}
}
