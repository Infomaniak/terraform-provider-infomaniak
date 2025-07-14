package dbaas

import "fmt"

type DBaaSPack struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type DBaaS struct {
	Id      int          `json:"id,omitempty"`
	Project DBaaSProject `json:"project,omitzero"`
	PackId  int          `json:"pack_id,omitempty"`
	Pack    *DBaaSPack   `json:"pack,omitempty"`

	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
	Name    string `json:"name,omitempty"`
	Region  string `json:"region,omitempty"`
	Status  string `json:"status,omitempty"`
}

type DBaaSConnectionInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Ca       string `json:"ca"`
}

type DBaaSBackup struct {
	Id     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

type DBaaSRestore struct {
	Id     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

func (dbaas *DBaaS) Key() string {
	return fmt.Sprintf("%d-%d-%d", dbaas.Project.PublicCloudId, dbaas.Project.ProjectId, dbaas.Id)
}

type DBaaSProject struct {
	PublicCloudId int `json:"public_cloud_id,omitempty"`
	ProjectId     int `json:"id,omitempty"`
}
