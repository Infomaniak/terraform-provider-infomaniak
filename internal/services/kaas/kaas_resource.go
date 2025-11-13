package kaas

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/apis/kaas"
	"terraform-provider-infomaniak/internal/provider"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &kaasResource{}
	_ resource.ResourceWithConfigure   = &kaasResource{}
	_ resource.ResourceWithImportState = &kaasResource{}
)

func NewKaasResource() resource.Resource {
	return &kaasResource{}
}

type kaasResource struct {
	client *apis.Client
}

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

type ApiserverModel struct {
	Params types.Map  `tfsdk:"params"`
	Oidc   *OidcModel `tfsdk:"oidc"`
	Audit  *Audit     `tfsdk:"audit"`
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

func (r *kaasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kaas"
}

// Configure adds the provider configured client to the data source.
func (r *kaasResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, err := provider.GetApiClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			err.Error(),
		)
		return
	}

	r.client = client
}

func (r *kaasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
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

func (r *kaasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	chosenPack, err := r.getPackId(data, &resp.Diagnostics)
	if err != nil {
		return
	}

	input := &kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: data.PublicCloudId.ValueInt64(),
			ProjectId:     data.PublicCloudProjectId.ValueInt64(),
		},
		Region:            data.Region.ValueString(),
		KubernetesVersion: data.KubernetesVersion.ValueString(),
		Name:              data.Name.ValueString(),
		PackId:            chosenPack.Id,
	}

	// CreateKaas API call logic
	kaasId, err := r.client.Kaas.CreateKaas(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when creating KaaS",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(kaasId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	kaasObject, err := r.waitUntilActive(ctx, input, kaasId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when waiting for KaaS to be Active",
			err.Error(),
		)
		return
	}

	if kaasObject == nil {
		return
	}

	err = r.fetchAndSetKubeconfig(&data, kaasObject)
	if err != nil {
		resp.Diagnostics.AddWarning("could not fetch and set kubeconfig", err.Error())
	}

	data.fill(kaasObject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apiserver != nil {
		apiserverParamsInput := r.buildApiserverParamsInput(data)
		created, err := r.client.Kaas.PatchApiserverParams(apiserverParamsInput, input.Project.PublicCloudId, input.Project.ProjectId, kaasId)
		if !created || err != nil {
			resp.Diagnostics.AddError(
				"Error when creating Oidc",
				err.Error(),
			)
			return
		}

		data.fillApiserverState(ctx, apiserverParamsInput)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (state *KaasModel) fillApiserverState(ctx context.Context, apiserverParams *kaas.Apiserver) {
	if state.shouldUpdateApiserver() {
		state.SetDefaultValues(ctx)
		state.updateAuditConfig(apiserverParams)
		state.updateOIDCConfig(apiserverParams)
		if state.canSetApiserverToNil() {
			state.Apiserver = nil
		}
	}
}

func (state *KaasModel) shouldUpdateApiserver() bool {
	apiserver := state.Apiserver
	return apiserver != nil && (apiserver.Audit != nil || apiserver.Oidc != nil || !apiserver.Params.IsNull())
}

func (state *KaasModel) updateAuditConfig(apiserverParams *kaas.Apiserver) {
	if apiserverParams.AuditLogPolicy == nil && apiserverParams.AuditLogWebhook == nil {
		state.Apiserver.Audit = nil
	} else {
		state.Apiserver.Audit.Policy = types.StringPointerValue(apiserverParams.AuditLogPolicy)
		state.Apiserver.Audit.WebhookConfig = types.StringPointerValue(apiserverParams.AuditLogWebhook)
	}
}

func (state *KaasModel) updateOIDCConfig(apiserverParams *kaas.Apiserver) {
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

func (state *KaasModel) canSetApiserverToNil() bool {
	apiserver := state.Apiserver
	return apiserver.Audit == nil && apiserver.Oidc == nil && apiserver.Params.IsNull()
}

func (r *kaasResource) waitUntilActive(ctx context.Context, kaas *kaas.Kaas, id int64) (*kaas.Kaas, error) {
	for {
		found, err := r.client.Kaas.GetKaas(kaas.Project.PublicCloudId, kaas.Project.ProjectId, id)
		if err != nil {
			return nil, err
		}

		if ctx.Err() != nil {
			return nil, nil
		}

		if found.Status == "Active" {
			return found, nil
		}

		time.Sleep(5 * time.Second)
	}
}

func (r *kaasResource) getApiserverParamsValues(data KaasModel) map[string]string {
	params := make(map[string]string)
	if !data.Apiserver.Params.IsNull() && !data.Apiserver.Params.IsUnknown() {
		for key, val := range data.Apiserver.Params.Elements() {
			if strVal, ok := val.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				params[key] = strVal.ValueString()
			}
		}
	}

	return params
}

func (r *kaasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KaasModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	kaasObject, err := r.client.Kaas.GetKaas(
		state.PublicCloudId.ValueInt64(),
		state.PublicCloudProjectId.ValueInt64(),
		state.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading KaaS",
			err.Error(),
		)
		return
	}

	state.fill(kaasObject)

	err = r.fetchAndSetKubeconfig(&state, kaasObject)
	if err != nil {
		resp.Diagnostics.AddWarning("could not fetch and set kubeconfig", err.Error())
	}

	apiserverParams, err := r.client.Kaas.GetApiserverParams(state.PublicCloudId.ValueInt64(), state.PublicCloudProjectId.ValueInt64(), kaasObject.Id)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Could not get Oidc",
			err.Error(),
		)
	}

	if apiserverParams != nil {
		state.fillApiserverState(ctx, apiserverParams)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *kaasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state KaasModel
	var data KaasModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	chosenPackState, err := r.getPackId(state, &resp.Diagnostics)
	if err != nil {
		return
	}

	input := r.prepareUpdateInput(state, data, chosenPackState.Id)

	if _, err := r.client.Kaas.UpdateKaas(input); err != nil {
		resp.Diagnostics.AddError("Error when updating KaaS", err.Error())
		return
	}

	kaasObject, err := r.waitUntilActive(ctx, input, input.Id)
	if err != nil || kaasObject == nil {
		resp.Diagnostics.AddError("Error waiting for KaaS activation", err.Error())
		return
	}

	err = r.fetchAndSetKubeconfig(&data, kaasObject)
	if err != nil {
		resp.Diagnostics.AddWarning("could not fetch and set kubeconfig", err.Error())
	}

	data.fill(kaasObject)

	if data.Apiserver != nil {
		r.handleApiserverConfig(ctx, &data, input, resp)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kaasResource) prepareUpdateInput(state, data KaasModel, packID int64) *kaas.Kaas {
	input := &kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: data.PublicCloudId.ValueInt64(),
			ProjectId:     data.PublicCloudProjectId.ValueInt64(),
		},
		Id:                state.Id.ValueInt64(),
		Name:              data.Name.ValueString(),
		PackId:            packID,
		Region:            state.Region.ValueString(),
		KubernetesVersion: data.KubernetesVersion.ValueString(),
	}

	if state.KubernetesVersion.ValueString() == data.KubernetesVersion.ValueString() {
		input.KubernetesVersion = ""
	}

	return input
}

func (r *kaasResource) fetchAndSetKubeconfig(data *KaasModel, input *kaas.Kaas) error {
	kubeconfig, err := r.client.Kaas.GetKubeconfig(
		input.Project.PublicCloudId,
		input.Project.ProjectId,
		input.Id,
	)
	if err != nil {
		return fmt.Errorf("could not get kubeconfig: %w", err)
	}
	data.Kubeconfig = types.StringValue(kubeconfig)
	return nil
}

func (r *kaasResource) handleApiserverConfig(ctx context.Context, data *KaasModel, input *kaas.Kaas, resp *resource.UpdateResponse) {
	apiserverParamsInput := r.buildApiserverParamsInput(*data)
	patched, err := r.client.Kaas.PatchApiserverParams(apiserverParamsInput, input.Project.PublicCloudId, input.Project.ProjectId, input.Id)
	if !patched || err != nil {
		resp.Diagnostics.AddError("Error when patching Apiserver params", err.Error())
		return
	}
	data.fillApiserverState(ctx, apiserverParamsInput)
}

func (r *kaasResource) buildApiserverParamsInput(data KaasModel) *kaas.Apiserver {
	apiserverParamsInput := &kaas.Apiserver{
		NonSpecificApiServerParams: r.getApiserverParamsValues(data),
	}
	if data.Apiserver.Audit != nil {
		apiserverParamsInput.AuditLogPolicy = data.Apiserver.Audit.Policy.ValueStringPointer()
		apiserverParamsInput.AuditLogWebhook = data.Apiserver.Audit.WebhookConfig.ValueStringPointer()
	}
	if data.Apiserver.Oidc != nil {
		apiserverParamsInput.OidcCa = data.Apiserver.Oidc.Ca.ValueStringPointer()
		apiserverParamsInput.Params = &kaas.ApiServerParams{
			IssuerUrl:      data.Apiserver.Oidc.IssuerUrl.ValueStringPointer(),
			ClientId:       data.Apiserver.Oidc.ClientId.ValueStringPointer(),
			UsernameClaim:  data.Apiserver.Oidc.UsernameClaim.ValueStringPointer(),
			UsernamePrefix: data.Apiserver.Oidc.UsernamePrefix.ValueStringPointer(),
			SigningAlgs:    data.Apiserver.Oidc.SigningAlgs.ValueStringPointer(),
			GroupsClaim:    data.Apiserver.Oidc.GroupsClaim.ValueStringPointer(),
			GroupsPrefix:   data.Apiserver.Oidc.GroupsPrefix.ValueStringPointer(),
			RequiredClaim:  data.Apiserver.Oidc.RequiredClaim.ValueStringPointer(),
		}
	}
	return apiserverParamsInput
}

func (r *kaasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KaasModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DeleteKaas API call logic
	_, err := r.client.Kaas.DeleteKaas(
		data.PublicCloudId.ValueInt64(),
		data.PublicCloudProjectId.ValueInt64(),
		data.Id.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when deleting KaaS",
			err.Error(),
		)
		return
	}
}

func (r *kaasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: public_cloud_id,public_cloud_project_id,id. Got: %q", req.ID),
		)
		return
	}

	var errorList error

	publicCloudId, err := strconv.ParseInt(idParts[0], 10, 64)
	errorList = errors.Join(errorList, err)
	publicCloudProjectId, err := strconv.ParseInt(idParts[1], 10, 64)
	errorList = errors.Join(errorList, err)
	kaasId, err := strconv.ParseInt(idParts[2], 10, 64)
	errorList = errors.Join(errorList, err)

	if errorList != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: public_cloud_id,public_cloud_project_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_id"), publicCloudId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_cloud_project_id"), publicCloudProjectId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), kaasId)...)
}

func (r *kaasResource) getPackId(data KaasModel, diagnostic *diag.Diagnostics) (*kaas.KaasPack, error) {
	packs, err := r.client.Kaas.GetPacks()
	if err != nil {
		diagnostic.AddError(
			"Could not get KaaS Packs",
			err.Error(),
		)
		return nil, err
	}

	var chosenPack *kaas.KaasPack
	for _, pack := range packs {
		if pack.Name == data.PackName.ValueString() {
			chosenPack = pack
			break
		}
	}

	if chosenPack == nil {
		var packNames []string
		for _, pack := range packs {
			packNames = append(packNames, pack.Name)
		}

		diagnostic.AddError(
			"Unknown KaaS Pack",
			fmt.Sprintf("pack_name must be one of : %v", packNames),
		)
		return nil, fmt.Errorf("pack name has not been found")
	}

	return chosenPack, nil
}

func (model *KaasModel) fill(kaas *kaas.Kaas) {
	model.Id = types.Int64Value(kaas.Id)
	model.Region = types.StringValue(kaas.Region)
	model.KubernetesVersion = types.StringValue(kaas.KubernetesVersion)
	model.Name = types.StringValue(kaas.Name)
	model.PackName = types.StringValue(kaas.Pack.Name)
}
