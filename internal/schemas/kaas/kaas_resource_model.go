package kaas_schemas

import (
	"context"
	"net/netip"
	"terraform-provider-infomaniak/internal/apis/kaas"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KaasModel struct {
	PublicCloudId        types.Int64 `tfsdk:"public_cloud_id"`
	PublicCloudProjectId types.Int64 `tfsdk:"public_cloud_project_id"`
	Id                   types.Int64 `tfsdk:"id"`

	Name              types.String    `tfsdk:"name"`
	PackName          types.String    `tfsdk:"pack_name"`
	Region            types.String    `tfsdk:"region"`
	Kubeconfig        types.String    `tfsdk:"kubeconfig"`
	KubernetesVersion types.String    `tfsdk:"kubernetes_version"`
	Apiserver         *ApiserverModel `tfsdk:"apiserver"`
}

type ApiserverModel struct {
	IpFilters types.List `tfsdk:"ip_filters"`
	Params    types.Map  `tfsdk:"params"`
	Oidc      *OidcModel `tfsdk:"oidc"`
	Audit     *Audit     `tfsdk:"audit"`
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

func (m *KaasModel) SetDefaultValues(ctx context.Context) {
	if m.Apiserver == nil {
		defaultParams, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		m.Apiserver = &ApiserverModel{
			Params: defaultParams,
		}
	}
	if m.Apiserver.Audit == nil {
		m.Apiserver.Audit = &Audit{}
	}
	if m.Apiserver.Oidc == nil {
		m.Apiserver.Oidc = &OidcModel{}
	}
}

func (state *KaasModel) FillApiserverState(ctx context.Context, apiserverParams *kaas.Apiserver) {
	if state.ShouldUpdateApiserver() {
		state.SetDefaultValues(ctx)
		state.UpdateAuditConfig(apiserverParams)
		state.UpdateOIDCConfig(apiserverParams)
		if state.CanSetApiserverToNil() {
			state.Apiserver = nil
		}
	}
}

func (state *KaasModel) FillIpFilters(ctx context.Context, cidrs []netip.Prefix) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	if len(cidrs) == 0 {
		return diagnostics
	}

	convertedCirds := make([]string, len(cidrs))
	for i, cidr := range cidrs {
		convertedCirds[i] = cidr.String()
	}

	listValue, diags := types.ListValueFrom(ctx, types.StringType, convertedCirds)
	state.Apiserver.IpFilters = listValue
	diagnostics = diags

	return diagnostics
}

func (state *KaasModel) ShouldUpdateApiserver() bool {
	apiserver := state.Apiserver
	return apiserver != nil && (apiserver.Audit != nil || apiserver.Oidc != nil || !apiserver.Params.IsNull())
}

func (state *KaasModel) UpdateAuditConfig(apiserverParams *kaas.Apiserver) {
	if apiserverParams.AuditLogPolicy == nil && apiserverParams.AuditLogWebhook == nil {
		state.Apiserver.Audit = nil
	} else {
		state.Apiserver.Audit.Policy = types.StringPointerValue(apiserverParams.AuditLogPolicy)
		state.Apiserver.Audit.WebhookConfig = types.StringPointerValue(apiserverParams.AuditLogWebhook)
	}
}

func (state *KaasModel) UpdateOIDCConfig(apiserverParams *kaas.Apiserver) {
	if apiserverParams.Params != nil {
		params := apiserverParams.Params
		state.Apiserver.Oidc = &OidcModel{
			ClientId:       types.StringPointerValue(params.ClientId),
			IssuerUrl:      types.StringPointerValue(params.IssuerUrl),
			UsernameClaim:  types.StringPointerValue(params.UsernameClaim),
			UsernamePrefix: types.StringPointerValue(params.UsernamePrefix),
			SigningAlgs:    types.StringPointerValue(params.SigningAlgs),
			GroupsClaim:    types.StringPointerValue(params.GroupsClaim),
			GroupsPrefix:   types.StringPointerValue(params.GroupsPrefix),
			RequiredClaim:  types.StringPointerValue(params.RequiredClaim),
			Ca:             types.StringPointerValue(apiserverParams.OidcCa),
		}
	} else {
		state.Apiserver.Oidc = nil
		state.Apiserver.Params = types.MapNull(types.StringType)
	}
}

func (state *KaasModel) CanSetApiserverToNil() bool {
	apiserver := state.Apiserver
	return apiserver.Audit == nil && apiserver.Oidc == nil && apiserver.Params.IsNull()
}

func (model *KaasModel) Fill(kaas *kaas.Kaas) {
	model.Id = types.Int64Value(kaas.Id)
	model.Region = types.StringValue(kaas.Region)
	model.KubernetesVersion = types.StringValue(kaas.KubernetesVersion)
	model.Name = types.StringValue(kaas.Name)
	model.PackName = types.StringValue(kaas.Pack.Name)
}
