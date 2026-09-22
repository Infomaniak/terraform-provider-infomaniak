package publiccloud

import (
	"terraform-provider-infomaniak/internal/apis/publiccloud"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func getPublicCloudUserAuthenticationDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Reads an authentication file (e.g. `clouds.yaml`) of a given type for a Public Cloud user.",
		Attributes: map[string]schema.Attribute{
			"public_cloud_id":         schema.Int64Attribute{Required: true, MarkdownDescription: "Identifier of the parent Public Cloud product."},
			"public_cloud_project_id": schema.Int64Attribute{Required: true, MarkdownDescription: "Identifier of the parent project."},
			"public_cloud_user_id":    schema.Int64Attribute{Required: true, MarkdownDescription: "Identifier of the user."},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Authentication file type as understood by the Infomaniak API (e.g. `clouds.yaml`).",
			},
			"region": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Target region (`pub1` or `pub2`).",
				Validators: []validator.String{
					stringvalidator.OneOf(publiccloud.RegionPub1, publiccloud.RegionPub2),
				},
			},
			"content": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Raw content of the authentication file.",
			},
		},
	}
}
