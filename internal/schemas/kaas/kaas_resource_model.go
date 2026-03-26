package kaas_schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KaasModel struct {
	PublicCloudId        types.Int64 `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64 `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64 `tfsdk:"id"`

	Name              types.String `tfsdk:"name"`
	PackName          types.String `tfsdk:"pack_name"`
	Region            types.String `tfsdk:"region"`
	Kubeconfig        types.String `tfsdk:"kubeconfig"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
	Apiserver         types.Object `tfsdk:"apiserver"`
}

func (k KaasModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"public_cloud_id":         types.Int64Type,
		"public_cloud_project_id": types.Int64Type,
		"id":                      types.Int64Type,
		"name":                    types.StringType,
		"pack_name":               types.StringType,
		"region":                  types.StringType,
		"kubeconfig":              types.StringType,
		"kubernetes_version":      types.StringType,
		"apiserver":               types.ObjectType{AttrTypes: ApiserverModel{}.AttributeTypes()},
	}
}

type ApiserverModel struct {
	IpFilters types.List   `tfsdk:"ip_filters"`
	Params    types.Map    `tfsdk:"params"`
	Oidc      types.Object `tfsdk:"oidc"`
	Audit     types.Object `tfsdk:"audit"`
}

func (a ApiserverModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip_filters": types.ListType{ElemType: types.StringType},
		"params":     types.MapType{ElemType: types.StringType},
		"oidc":       types.ObjectType{AttrTypes: OidcModel{}.AttributeTypes()},
		"audit":      types.ObjectType{AttrTypes: Audit{}.AttributeTypes()},
	}
}

type OidcModel struct {
	IssuerUrl      types.String `tfsdk:"issuer_url"`
	ClientId       types.String `tfsdk:"client_id"`
	UsernameClaim  types.String `tfsdk:"username_claim"`
	UsernamePrefix types.String `tfsdk:"username_prefix"`
	SigningAlgs    types.String `tfsdk:"signing_algs"`
	GroupsClaim    types.String `tfsdk:"groups_claim"`
	GroupsPrefix   types.String `tfsdk:"groups_prefix"`
	RequiredClaim  types.String `tfsdk:"required_claim"`
	Ca             types.String `tfsdk:"ca"`
}

func (o OidcModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"issuer_url":      types.StringType,
		"client_id":       types.StringType,
		"username_claim":  types.StringType,
		"username_prefix": types.StringType,
		"signing_algs":    types.StringType,
		"groups_claim":    types.StringType,
		"groups_prefix":   types.StringType,
		"required_claim":  types.StringType,
		"ca":              types.StringType,
	}
}

type Audit struct {
	WebhookConfig types.String `tfsdk:"webhook_config"`
	Policy        types.String `tfsdk:"policy"`
}

func (a Audit) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"webhook_config": types.StringType,
		"policy":         types.StringType,
	}
}
