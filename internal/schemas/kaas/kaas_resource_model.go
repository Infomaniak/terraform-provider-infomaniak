package kaas_schemas

import (
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

type ApiserverModel struct {
	IpFilters types.List   `tfsdk:"ip_filters"`
	Params    types.Map    `tfsdk:"params"`
	Oidc      types.Object `tfsdk:"oidc"`
	Audit     types.Object `tfsdk:"audit"`
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

type Audit struct {
	WebhookConfig types.String `tfsdk:"webhook_config"`
	Policy        types.String `tfsdk:"policy"`
}
