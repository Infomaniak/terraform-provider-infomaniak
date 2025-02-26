package kaas

import "fmt"

type KaasPack struct {
	Id          int    `json:"kaas_pack_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Kaas struct {
	Id      int         `json:"kaas_id,omitempty"`
	Project KaasProject `json:"project,omitempty"`

	Region            string `json:"region,omitempty"`
	Kubeconfig        string `json:"kubeconfig,omitempty"`
	KubernetesVersion string `json:"kubernetes_version,omitempty"`
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
	Id     int `json:"id,omitempty"`

	Name             string `json:"name,omitempty"`
	FlavorName       string `json:"flavor_name,omitempty"`
	AvailabilityZone string `json:"availability_zone,omitempty"`
	MinInstances     int32  `json:"minimum_instances,omitempty"`
	// MaxInstances     int32  `json:"maximum_instances,omitempty"`
}

func (instancePool *InstancePool) Key() string {
	return fmt.Sprintf("%d-%d", instancePool.KaasId, instancePool.Id)
}
