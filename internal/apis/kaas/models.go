package kaas

import (
	"encoding/json"
	"fmt"
	"maps"
)

type KaasPack struct {
	Id          int    `json:"kaas_pack_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Apiserver struct {
	Params                     *ApiServerParams   `json:"apiserver_params"`
	NonSpecificApiServerParams map[string]string `json:"-"`

	OidcCa          *string `json:"oidc_ca"`
	AuditLogWebhook *string `json:"audit-webhook-config"`
	AuditLogPolicy  *string `json:"audit-policy"`
}

var _ json.Marshaler = (*Apiserver)(nil)

// We can delete this once json v2 is out, so we can flatten everything without having to do this
func (a *Apiserver) MarshalJSON() ([]byte, error) {
	paramBytes, err := json.Marshal(a.Params)
	if err != nil {
		paramBytes = []byte("{}")
	}
	paramsMap := make(map[string]string)
	json.Unmarshal(paramBytes, &paramsMap)
	nonSpecificMap := a.NonSpecificApiServerParams
	res := make(map[string]string)
	maps.Copy(res, paramsMap)
	maps.Copy(res, nonSpecificMap)
	result, err := json.Marshal(map[string]any{
		"apiserver_params":     res,
		"oidc_ca":              a.OidcCa,
		"audit-policy":         a.AuditLogPolicy,
		"audit-webhook-config": a.AuditLogWebhook,
	})
	return result, err
}

type ApiServerParams struct {
	IssuerUrl      string `json:"--oidc-issuer-url,omitempty"`
	ClientId       string `json:"--oidc-client-id,omitempty"`
	UsernameClaim  string `json:"--oidc-username-claim,omitempty"`
	UsernamePrefix string `json:"--oidc-username-prefix,omitempty"`
	SigningAlgs    string `json:"--oidc-signing-algs,omitempty"`
	GroupsClaim    string `json:"--oidc-groups-claim,omitempty"`
	GroupsPrefix   string `json:"--oidc-groups-prefix,omitempty"`
}

type Kaas struct {
	Name    string      `json:"name,omitempty"`
	Id      int         `json:"kaas_id,omitempty"`
	Project KaasProject `json:"project,omitzero"`
	PackId  int         `json:"kaas_pack_id,omitempty"`
	Pack    *KaasPack   `json:"pack,omitempty"`

	Region            string `json:"region,omitempty"`
	KubernetesVersion string `json:"kubernetes_version,omitempty"`
	Status            string `json:"status,omitempty"`
}

func (kaas *Kaas) Key() string {
	return fmt.Sprintf("%d-%d-%d", kaas.Project.PublicCloudId, kaas.Project.ProjectId, kaas.Id)
}

type KaasProject struct {
	PublicCloudId int `json:"public_cloud_id,omitempty"`
	ProjectId     int `json:"id,omitempty"`
}

type InstancePool struct {
	KaasId int `json:"kaas_id,omitempty"`
	Id     int `json:"instance_pool_id,omitempty"`

	Name             string            `json:"name,omitempty"`
	FlavorName       string            `json:"flavor,omitempty"`
	AvailabilityZone string            `json:"availability_zone,omitempty"`
	MinInstances     int32             `json:"minimum_instances,omitempty"`
	MaxInstances     int32             `json:"maximum_instances,omitempty"`
	Status           string            `json:"status,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`

	TargetInstances    int32 `json:"target_instances,omitempty"`
	AvailableInstances int32 `json:"available_instances,omitempty"`
}

func (instancePool *InstancePool) Key() string {
	return fmt.Sprintf("%d-%d", instancePool.KaasId, instancePool.Id)
}
