package services

import (
	"context"
	"fmt"
	"net/netip"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/kaas"
	kaas_schemas "terraform-provider-infomaniak/internal/schemas/kaas"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type KaasService struct {
	client *apis.Client
}

func NewKaasService(client *apis.Client) *KaasService {
	return &KaasService{
		client: client,
	}
}

func (s *KaasService) CreateKaas(model kaas_schemas.KaasModel, pack kaas.KaasPack, state *kaas_schemas.KaasModel) (int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	input := &kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: model.PublicCloudId.ValueInt64(),
			ProjectId:     model.PublicCloudProjectId.ValueInt64(),
		},
		Region:            model.Region.ValueString(),
		KubernetesVersion: model.KubernetesVersion.ValueString(),
		Name:              model.Name.ValueString(),
		PackId:            pack.Id,
	}

	kaasId, err := s.client.Kaas.CreateKaas(input)
	if err != nil {
		diags.AddError("error when creating KaaS", err.Error())
		return 0, diags
	}

	state.Id = types.Int64Value(kaasId)

	return kaasId, diags
}

func (s *KaasService) UpdateKaas(model kaas_schemas.KaasModel, pack kaas.KaasPack, state *kaas_schemas.KaasModel) diag.Diagnostics {
	var diags diag.Diagnostics

	input := &kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: model.PublicCloudId.ValueInt64(),
			ProjectId:     model.PublicCloudProjectId.ValueInt64(),
		},
		Region:            model.Region.ValueString(),
		KubernetesVersion: model.KubernetesVersion.ValueString(),
		Name:              model.Name.ValueString(),
		PackId:            pack.Id,
		Id:                model.Id.ValueInt64(),
	}

	if state.KubernetesVersion.ValueString() == model.KubernetesVersion.ValueString() {
		input.KubernetesVersion = ""
	}

	_, err := s.client.Kaas.UpdateKaas(input)
	if err != nil {
		diags.AddError("error when updating KaaS", err.Error())
		return diags
	}

	state.PublicCloudProjectId = model.PublicCloudProjectId
	state.PublicCloudId = model.PublicCloudId
	state.Region = model.Region
	state.Name = model.Name
	state.PackName = model.PackName
	state.Id = model.Id

	return diags
}

func (s *KaasService) GetKaasOnceActive(ctx context.Context, model kaas_schemas.KaasModel, kaasId int64, state *kaas_schemas.KaasModel) (*kaas.Kaas, diag.Diagnostics) {
	var diags diag.Diagnostics

	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			diags.AddError(fmt.Sprintf("context expired while waiting for KaaS %d to be active", kaasId), ctx.Err().Error())
			return nil, diags
		case <-t.C:
			found, err := s.client.Kaas.GetKaas(model.PublicCloudId.ValueInt64(), model.PublicCloudProjectId.ValueInt64(), kaasId)
			if err != nil {
				diags.AddError(fmt.Sprintf("error while finding KaaS %d", kaasId), err.Error())
				return nil, diags
			}

			if found.Status == "Active" {
				state.Id = types.Int64Value(found.Id)
				state.Region = types.StringValue(found.Region)
				state.KubernetesVersion = types.StringValue(found.KubernetesVersion)
				state.Name = types.StringValue(found.Name)
				state.PackName = types.StringValue(found.Pack.Name)
				return found, diags
			}
		}
	}
}

func (s *KaasService) GetKubeconfig(model kaas_schemas.KaasModel, kaasId int64, state *kaas_schemas.KaasModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	kubeconfig, err := s.client.Kaas.GetKubeconfig(
		model.PublicCloudId.ValueInt64(),
		model.PublicCloudProjectId.ValueInt64(),
		kaasId,
	)
	if err != nil {
		diags.AddError(fmt.Sprintf("could not retrieve KaaS %d kubeconfig file", kaasId), err.Error())
		return "", diags
	}

	state.Kubeconfig = types.StringValue(kubeconfig)

	return kubeconfig, diags
}

func (s *KaasService) SetApiserverConfig(ctx context.Context, model kaas_schemas.KaasModel, kaasId int64, state *kaas_schemas.KaasModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.Apiserver.IsNull() || model.Apiserver.IsUnknown() {
		state.Apiserver = model.Apiserver
		return diags
	}

	var apiserverModel kaas_schemas.ApiserverModel
	diags.Append(model.Apiserver.As(ctx, &apiserverModel, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	diags.Append(s.setIpFilters(ctx, apiserverModel, model.PublicCloudId.ValueInt64(), model.PublicCloudProjectId.ValueInt64(), kaasId)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(s.setApiServerParams(ctx, apiserverModel, model.PublicCloudId.ValueInt64(), model.PublicCloudProjectId.ValueInt64(), kaasId)...)
	if diags.HasError() {
		return diags
	}

	state.Apiserver = model.Apiserver

	return diags
}

func (s *KaasService) ReadApiserverConfig(ctx context.Context, model kaas_schemas.KaasModel, kaasId int64, state *kaas_schemas.KaasModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.Apiserver.IsNull() || model.Apiserver.IsUnknown() {
		state.Apiserver = model.Apiserver
		return diags
	}

	var apiserverState kaas_schemas.ApiserverModel
	diags.Append(state.Apiserver.As(ctx, &apiserverState, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	var apiserverModel kaas_schemas.ApiserverModel
	diags.Append(model.Apiserver.As(ctx, &apiserverModel, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	diags.Append(s.readIpFilters(ctx, apiserverModel, &apiserverState, model.PublicCloudId.ValueInt64(), model.PublicCloudProjectId.ValueInt64(), kaasId)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(s.readApiserverParams(ctx, &apiserverState, model.PublicCloudId.ValueInt64(), model.PublicCloudProjectId.ValueInt64(), kaasId)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(s.readApiserverOidc(ctx, apiserverModel, &apiserverState, model.PublicCloudId.ValueInt64(), model.PublicCloudProjectId.ValueInt64(), kaasId)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(s.readApiserverAudit(ctx, apiserverModel, &apiserverState, model.PublicCloudId.ValueInt64(), model.PublicCloudProjectId.ValueInt64(), kaasId)...)
	if diags.HasError() {
		return diags
	}

	apiserverTfObject, diags := types.ObjectValueFrom(ctx, apiserverState.AttributeTypes(), apiserverState)
	if diags.HasError() {
		return diags
	}
	state.Apiserver = apiserverTfObject

	return diags
}

func (s *KaasService) readApiserverOidc(ctx context.Context, model kaas_schemas.ApiserverModel, state *kaas_schemas.ApiserverModel, publicCloudId, projectId, kaasId int64) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.Oidc.IsNull() {
		state.Oidc = model.Oidc
		return diags
	}

	var oidcState kaas_schemas.OidcModel
	diags.Append(state.Oidc.As(ctx, &oidcState, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	apiserverParams, err := s.client.Kaas.GetApiserverParams(publicCloudId, projectId, kaasId)
	if err != nil {
		diags.AddError("could not read apiserver params", err.Error())
		return diags
	}

	oidcState.Ca = types.StringPointerValue(apiserverParams.OidcCa)
	if apiserverParams.Params != nil {
		oidcState.ClientId = types.StringPointerValue(apiserverParams.Params.ClientId)
		oidcState.IssuerUrl = types.StringPointerValue(apiserverParams.Params.IssuerUrl)
		oidcState.UsernameClaim = types.StringPointerValue(apiserverParams.Params.UsernameClaim)
		oidcState.UsernamePrefix = types.StringPointerValue(apiserverParams.Params.UsernamePrefix)
		oidcState.SigningAlgs = types.StringPointerValue(apiserverParams.Params.SigningAlgs)
		oidcState.GroupsClaim = types.StringPointerValue(apiserverParams.Params.GroupsClaim)
		oidcState.GroupsPrefix = types.StringPointerValue(apiserverParams.Params.GroupsPrefix)
		oidcState.RequiredClaim = types.StringPointerValue(apiserverParams.Params.RequiredClaim)
	}
	oidcTfObject, diags := types.ObjectValueFrom(ctx, oidcState.AttributeTypes(), oidcState)
	if diags.HasError() {
		return diags
	}
	state.Oidc = oidcTfObject

	return diags
}

func (s *KaasService) readApiserverAudit(ctx context.Context, model kaas_schemas.ApiserverModel, state *kaas_schemas.ApiserverModel, publicCloudId, projectId, kaasId int64) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.Audit.IsNull() {
		state.Audit = model.Audit
		return diags
	}

	var auditState kaas_schemas.Audit
	diags.Append(state.Audit.As(ctx, &auditState, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	apiserverParams, err := s.client.Kaas.GetApiserverParams(publicCloudId, projectId, kaasId)
	if err != nil {
		diags.AddError("could not read apiserver params", err.Error())
		return diags
	}

	auditState.Policy = types.StringPointerValue(apiserverParams.AuditLogPolicy)
	auditState.WebhookConfig = types.StringPointerValue(apiserverParams.AuditLogWebhook)
	auditTfObject, diags := types.ObjectValueFrom(ctx, auditState.AttributeTypes(), auditState)
	if diags.HasError() {
		return diags
	}
	state.Audit = auditTfObject

	return diags
}

func (s *KaasService) readApiserverParams(ctx context.Context, state *kaas_schemas.ApiserverModel, publicCloudId, projectId, kaasId int64) diag.Diagnostics {
	var diags diag.Diagnostics

	state.Params = basetypes.NewMapNull(types.StringType)

	return diags
}

func (s *KaasService) readIpFilters(ctx context.Context, model kaas_schemas.ApiserverModel, state *kaas_schemas.ApiserverModel, publicCloudId, projectId, kaasId int64) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.IpFilters.IsNull() {
		state.IpFilters = model.IpFilters
		return diags
	}

	ipFilters, err := s.client.Kaas.GetIPFilters(publicCloudId, projectId, kaasId)
	if err != nil {
		diags.AddError("could not read Ip Filters", err.Error())
		return diags
	}

	stringIpFilter := make([]string, len(ipFilters))
	for i, ipFilter := range ipFilters {
		stringIpFilter[i] = ipFilter.String()
	}

	terraformIpFilters, diagnostics := types.ListValueFrom(ctx, types.StringType, stringIpFilter)
	diags.Append(diagnostics...)
	if diags.HasError() {
		return diags
	}

	state.IpFilters = terraformIpFilters

	return diags
}

func (s *KaasService) setIpFilters(ctx context.Context, model kaas_schemas.ApiserverModel, publicCloudId, projectId, kaasId int64) diag.Diagnostics {
	var diags diag.Diagnostics

	ipFilters := make([]string, 0, len(model.IpFilters.Elements()))
	diags.Append(model.IpFilters.ElementsAs(ctx, &ipFilters, true)...)
	if diags.HasError() {
		return diags
	}

	convertedIpFilters := make([]netip.Prefix, len(ipFilters))
	for i, cidr := range ipFilters {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			diags.AddError("invalid cidr format", err.Error())
			return diags
		}

		convertedIpFilters[i] = prefix
	}

	ok, err := s.client.Kaas.PutIPFilters(convertedIpFilters, publicCloudId, projectId, kaasId)
	if !ok || err != nil {
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		} else {
			errMsg = "PutIPFilters returned false but no error was provided"
		}
		diags.AddError("Error when applying ip filters", errMsg)
	}

	return diags
}

func (s *KaasService) setApiServerParams(ctx context.Context, model kaas_schemas.ApiserverModel, publicCloudId, projectId, kaasId int64) diag.Diagnostics {
	var diags diag.Diagnostics

	unmanagedParams := make(map[string]string)
	if !model.Params.IsNull() && !model.Params.IsUnknown() {
		for key, val := range model.Params.Elements() {
			if strVal, ok := val.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				unmanagedParams[key] = strVal.ValueString()
			}
		}
	}

	apiserverParamsInput := &kaas.Apiserver{
		NonSpecificApiServerParams: unmanagedParams,
	}

	if !model.Audit.IsNull() && !model.Audit.IsUnknown() {
		var auditModel kaas_schemas.Audit
		diags.Append(model.Audit.As(ctx, &auditModel, basetypes.ObjectAsOptions{})...)

		apiserverParamsInput.AuditLogPolicy = auditModel.Policy.ValueStringPointer()
		apiserverParamsInput.AuditLogWebhook = auditModel.WebhookConfig.ValueStringPointer()
	}

	if !model.Oidc.IsNull() && !model.Oidc.IsUnknown() {
		var oidcModel kaas_schemas.OidcModel
		diags.Append(model.Oidc.As(ctx, &oidcModel, basetypes.ObjectAsOptions{})...)

		apiserverParamsInput.OidcCa = oidcModel.Ca.ValueStringPointer()
		apiserverParamsInput.Params = &kaas.ApiServerParams{
			IssuerUrl:      oidcModel.IssuerUrl.ValueStringPointer(),
			ClientId:       oidcModel.ClientId.ValueStringPointer(),
			UsernameClaim:  oidcModel.UsernameClaim.ValueStringPointer(),
			UsernamePrefix: oidcModel.UsernamePrefix.ValueStringPointer(),
			SigningAlgs:    oidcModel.SigningAlgs.ValueStringPointer(),
			GroupsClaim:    oidcModel.GroupsClaim.ValueStringPointer(),
			GroupsPrefix:   oidcModel.GroupsPrefix.ValueStringPointer(),
			RequiredClaim:  oidcModel.RequiredClaim.ValueStringPointer(),
		}
	}

	created, err := s.client.Kaas.PatchApiserverParams(apiserverParamsInput, publicCloudId, projectId, kaasId)
	if !created || err != nil {
		diags.AddError("error while setting apiserver params", err.Error())
		return diags
	}

	return diags
}
