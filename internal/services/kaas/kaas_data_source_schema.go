package kaas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getKaasDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_cloud_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud where KaaS is installed",
				MarkdownDescription: "The id of the public cloud where KaaS is installed",
			},
			"public_cloud_project_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the public cloud project where KaaS is installed",
				MarkdownDescription: "The id of the public cloud project where KaaS is installed",
			},
			"id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of this KaaS",
				MarkdownDescription: "The id of this KaaS",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "The name of the KaaS project",
				MarkdownDescription: "The name of the KaaS project",
			},
			"pack_name": schema.StringAttribute{
				Computed:            true,
				Description:         "The name of the pack associated to the KaaS project",
				MarkdownDescription: "The name of the pack associated to the KaaS project",
			},
			"region": schema.StringAttribute{
				Computed:            true,
				Description:         "The region where the KaaS project resides in.",
				MarkdownDescription: "The region where the KaaS project resides in.",
			},
			"kubeconfig": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "The kubeconfig generated to access to KaaS project",
				MarkdownDescription: "The kubeconfig generated to access to KaaS project",
			},
			"kubernetes_version": schema.StringAttribute{
				Computed:            true,
				Description:         "The version of Kubernetes associated with the KaaS project",
				MarkdownDescription: "The version of Kubernetes associated with the KaaS project",
			},
			"apiserver": schema.SingleNestedAttribute{
				Description:         "Kubernetes Apiserver editable params",
				MarkdownDescription: "Kubernetes Apiserver editable params",
				Attributes: map[string]schema.Attribute{
					"params": schema.MapAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "Map of Kubernetes Apiserver params in case the terraform provider does not already abstracts them",
						MarkdownDescription: "Map of Kubernetes Apiserver params in case the terraform provider does not already abstracts them",
					},
					"audit": schema.SingleNestedAttribute{
						MarkdownDescription: "Kubernetes audit logs specification files",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"webhook_config": schema.StringAttribute{
								MarkdownDescription: "YAML manifest for audit webhook config",
								Computed:            true,
							},
							"policy": schema.StringAttribute{
								MarkdownDescription: "YAML manifest for audit policy",
								Computed:            true,
							},
						},
					},
					"oidc": schema.SingleNestedAttribute{
						Description:         "OIDC specific Apiserver params",
						MarkdownDescription: "OIDC specific Apiserver params",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"ca": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC Ca Certificate",
								MarkdownDescription: "OIDC Ca Certificate",
							},
							"groups_claim": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "OIDC groups claim",
							},
							"groups_prefix": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "OIDC groups prefix",
							},
							"issuer_url": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC issuer URL",
								MarkdownDescription: "OIDC issuer URL",
							},
							"client_id": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC client identifier",
								MarkdownDescription: "OIDC client identifier",
							},
							"username_claim": schema.StringAttribute{
								Optional:            true,
								Description:         "OIDC username claim",
								MarkdownDescription: "OIDC username claim",
							},
							"username_prefix": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC username prefix",
								MarkdownDescription: "OIDC username prefix",
							},
							"required_claim": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "A key=value pair that describes a required claim in the ID Token.",
							},
							"signing_algs": schema.StringAttribute{
								Computed:            true,
								Description:         "OIDC signing algorithm. Kubernetes will default it to RS256",
								MarkdownDescription: "OIDC signing algorithm. Kubernetes will default it to RS256",
							},
						},
					},
				},
				Optional: true,
			},
		},
		MarkdownDescription: "The kaas data source allows the user to manage a kaas project",
	}
}
