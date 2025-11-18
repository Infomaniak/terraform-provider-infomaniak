package kaas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getKaasResourceSchema() schema.Schema {
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
			"pack_name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the pack associated to the KaaS project",
				MarkdownDescription: "The name of the pack associated to the KaaS project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kubernetes_version": schema.StringAttribute{
				Required:            true,
				Description:         "The version of Kubernetes associated with the KaaS being installed",
				MarkdownDescription: "The version of Kubernetes associated with the KaaS being installed",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the KaaS project",
				MarkdownDescription: "The name of the KaaS project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				MarkdownDescription: "A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:            true,
				Description:         "The region where the KaaS will reside.",
				MarkdownDescription: "The region where the KaaS will reside.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kubeconfig": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "The kubeconfig generated to access to KaaS project",
				MarkdownDescription: "The kubeconfig generated to access to KaaS project",
			},
			"apiserver": schema.SingleNestedAttribute{
				MarkdownDescription: "Kubernetes Apiserver editable params",
				Attributes: map[string]schema.Attribute{
					"params": schema.MapAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "Map of Kubernetes Apiserver params in case the terraform provider does not already abstracts them",
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.UseStateForUnknown(),
						},
					},
					"audit": schema.SingleNestedAttribute{
						MarkdownDescription: "Kubernetes audit logs specification files",
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"webhook_config": schema.StringAttribute{
								MarkdownDescription: "YAML manifest for audit webhook config",
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"policy": schema.StringAttribute{
								MarkdownDescription: "YAML manifest for audit policy",
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"oidc": schema.SingleNestedAttribute{
						MarkdownDescription: "OIDC specific Apiserver params",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ca": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "OIDC Ca Certificate",
							},
							"groups_claim": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "OIDC groups claim",
							},
							"groups_prefix": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "OIDC groups prefix",
							},
							"issuer_url": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "OIDC issuer URL",
							},
							"client_id": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "OIDC client identifier",
							},
							"username_claim": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "OIDC username claim",
							},
							"username_prefix": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "OIDC username prefix",
							},
							"required_claim": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "A key=value pair that describes a required claim in the ID Token. If set, the claim is verified to be present in the ID Token with a matching value. Repeat this flag to specify multiple claims.",
							},
							"signing_algs": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "OIDC signing algorithm. Kubernetes will default it to RS256",
							},
						},
					},
				},
				Optional: true,
			},
		},
		MarkdownDescription: "The kaas resource allows the user to manage a kaas project",
	}
}
