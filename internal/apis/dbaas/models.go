package dbaas

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DBaaSPack struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type DBaaS struct {
	Id         int                  `json:"id,omitempty"`
	Project    DBaaSProject         `json:"project,omitzero"`
	PackId     int                  `json:"pack_id,omitempty"`
	Pack       *DBaaSPack           `json:"pack,omitempty"`
	Connection *DBaaSConnectionInfo `json:"connection,omitempty"`

	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
	Name    string `json:"name,omitempty"`
	Region  string `json:"region,omitempty"`
	Status  string `json:"status,omitempty"`
}

// temporary fix until the backend returns null instead of []
// TODO: Remove me
func (d *DBaaS) UnmarshalJSON(data []byte) error {
    // Create an alias to avoid infinite recursion
    type Alias DBaaS
    aux := &struct {
        Connection json.RawMessage `json:"connection,omitempty"`
        *Alias
    }{
        Alias: (*Alias)(d),
    }

    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    // Handle the 'connection' field specifically
    if len(aux.Connection) > 0 {
        // Check if it's an empty array
        if strings.TrimSpace(string(aux.Connection)) == "[]" {
            d.Connection = nil // Explicitly set to nil
        } else {
            // Assume it's a JSON object
            d.Connection = &DBaaSConnectionInfo{}
            if err := json.Unmarshal(aux.Connection, d.Connection); err != nil {
                return err
            }
        }
    }
    return nil
}

type DBaaSCreateInfo struct {
	Id 							int 		`json:"id"`
	RootPassword	 	string 	`json:"root_password"`
	KubeIdentifier 	string 	`json:"kube_identifier"`
}

type DBaaSConnectionInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Ca       string `json:"ca"`
}

type DBaaSBackup struct {
	Id string `json:"id,omitempty"`
	Location   string `json:"location,omitempty"`
	CreatedAt  uint64  `json:"created_at,omitempty"`
	CompletedAt  uint64  `json:"completed_at,omitempty"`
	Status   string `json:"status,omitempty"`
}

type DBaaSRestore struct {
	Id     string `json:"id,omitempty"`
	BackupSource string `json:"backup_source,omitempty"`
	CreatedAt uint64 `json:"created_at,omitempty"`
	Status string `json:"status,omitempty"`
	NewService *DBaaSCreateInfo `json:"new_service,omitempty"`
}

func (dbaas *DBaaS) Key() string {
	return fmt.Sprintf("%d-%d-%d", dbaas.Project.PublicCloudId, dbaas.Project.ProjectId, dbaas.Id)
}

type DBaaSProject struct {
	PublicCloudId int `json:"public_cloud_id,omitempty"`
	ProjectId     int `json:"id,omitempty"`
}
