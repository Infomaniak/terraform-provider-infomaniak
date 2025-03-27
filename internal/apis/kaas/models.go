package kaas

import "fmt"

type KaasPack struct {
	Id          int    `json:"kaas_pack_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Kaas struct {
	Name    string      `json:"name,omitempty"`
	Id      int         `json:"kaas_id,omitempty"`
	Project KaasProject `json:"project,omitempty"`
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
	ProjectId     int `json:"project_cloud_project_id,omitempty"`
}

type InstancePool struct {
	KaasId int `json:"kaas_id,omitempty"`
	Id     int `json:"instance_pool_id,omitempty"`

	Name             string `json:"name,omitempty"`
	FlavorName       string `json:"flavor,omitempty"`
	AvailabilityZone string `json:"availability_zone,omitempty"`
	MinInstances     int32  `json:"minimum_instances,omitempty"`
	MaxInstances     int32  `json:"maximum_instances,omitempty"`
	Status           string `json:"status,omitempty"`

	TargetInstances    int32 `json:"target_instances,omitempty"`
	AvailableInstances int32 `json:"available_instances,omitempty"`
}

func (instancePool *InstancePool) Key() string {
	return fmt.Sprintf("%d-%d", instancePool.KaasId, instancePool.Id)
}
