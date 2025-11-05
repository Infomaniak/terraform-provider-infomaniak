package kaas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getKaasInstancePoolDataSourceSchema() schema.Schema {
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
			"kaas_id": schema.Int64Attribute{
				Required:    true,
				Description: "The id of the kaas project.",
			},
			"id": schema.Int64Attribute{
				Required:    true,
				Description: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of this instance pool",
			},
			"availability_zone": schema.StringAttribute{
				Computed:            true,
				Description:         "The availability zone for the instances in the pool",
				MarkdownDescription: "The availability zone for the instances in the pool",
			},
			"flavor_name": schema.StringAttribute{
				Computed:    true,
				Description: "The flavor name of the instance in this instance pool",
			},
			"min_instances": schema.Int64Attribute{
				Computed:    true,
				Description: "The minimum amount of instances in the instance pool",
			},
			"max_instances": schema.Int64Attribute{
				Computed:    true,
				Description: "The maximum amount of instances in the instance pool",
			},
			"labels": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Kubernetes node labels",
			},
		},
		MarkdownDescription: "The KaaS Instance Pool data source retrieves information about a KaaS instance pool.",
	}
}
