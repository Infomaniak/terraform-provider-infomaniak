package kaas

import "fmt"

type KaasPack struct {
	Id          int    `json:"kaas_pack_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Kaas struct {
	PcpId string `json:"pcp_id,omitempty"`
	Id    string `json:"id,omitempty"`

	Region     string `json:"region,omitempty"`
	Kubeconfig string `json:"kubeconfig,omitempty"`
}

func (kaas *Kaas) Key() string {
	return fmt.Sprintf("%s-%s", kaas.PcpId, kaas.Id)
}

type InstancePool struct {
	PcpId  string `json:"pcp_id,omitempty"`
	KaasId string `json:"kaas_id,omitempty"`
	Id     string `json:"id,omitempty"`

	Name         string `json:"name,omitempty"`
	FlavorName   string `json:"flavor_name,omitempty"`
	MinInstances int32  `json:"min_instances,omitempty"`
	MaxInstances int32  `json:"max_instances,omitempty"`
}

func (instancePool *InstancePool) Key() string {
	return fmt.Sprintf("%s-%s-%s", instancePool.PcpId, instancePool.KaasId, instancePool.Id)
}
