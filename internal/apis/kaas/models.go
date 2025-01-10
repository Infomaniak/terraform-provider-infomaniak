package kaas

import "fmt"

type Kaas struct {
	PcpId string `json:"pcp_id"`
	Id    string `json:"id"`

	Region     string `json:"region"`
	Kubeconfig string `json:"kubeconfig"`
}

func (kaas *Kaas) Key() string {
	return fmt.Sprintf("%s-%s", kaas.PcpId, kaas.Id)
}

type InstancePool struct {
	PcpId  string `json:"pcp_id"`
	KaasId string `json:"kaas_id"`
	Id     string `json:"id"`

	Name         string `json:"name"`
	FlavorName   string `json:"flavor_name"`
	MinInstances int32  `json:"min_instances"`
	MaxInstances int32  `json:"max_instances"`
}

func (instancePool *InstancePool) Key() string {
	return fmt.Sprintf("%s-%s-%s", instancePool.PcpId, instancePool.KaasId, instancePool.Id)
}
