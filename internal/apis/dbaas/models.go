package dbaas

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DBaaSPack struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type DbaasType struct {
	Name     string   `json:"name,omitempty"`
	Versions []string `json:"versions,omitempty"`
}

type PackFilter struct {
	DbType    string
	Group     *string
	Name      *string
	Instances *int64
	Cpu       *int64
	Ram       *int64
	Storage   *int64
}

type Pack struct {
	ID        int64  `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Group     string `json:"group,omitempty"`
	Name      string `json:"name,omitempty"`
	Instances int64  `json:"instances,omitempty"`
	CPU       int64  `json:"cpu,omitempty"`
	RAM       int64  `json:"ram,omitempty"`
	Storage   int64  `json:"storage,omitempty"`
	Rates     Rates  `json:"rates"`
}

type Rates struct {
	CHF Pricing `json:"CHF"`
	EUR Pricing `json:"EUR"`
}

type Pricing struct {
	HourExclTax float64 `json:"hour_excl_tax,omitempty"`
	HourInclTax float64 `json:"hour_incl_tax,omitempty"`
}

type DBaaS struct {
	Id         int64                `json:"id,omitempty"`
	Project    DBaaSProject         `json:"project,omitzero"`
	PackId     int64                `json:"pack_id,omitempty"`
	Pack       *DBaaSPack           `json:"pack,omitempty"`
	Connection *DBaaSConnectionInfo `json:"connection,omitempty"`

	Type                 string `json:"type,omitempty"`
	Version              string `json:"version,omitempty"`
	Name                 string `json:"name,omitempty"`
	KubernetesIdentifier string `json:"kube_identifier,omitempty"`
	Region               string `json:"region,omitempty"`
	Status               string `json:"status,omitempty"`
}

type AllowedCIDRs struct {
	IpFilters []string `json:"ip_filters,omitempty"`
}

// avoid crashes when the backend returns [] instead of null when connection is not yet avaialble
func (d *DBaaS) UnmarshalJSON(data []byte) error {
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

	if len(aux.Connection) > 0 {
		if strings.TrimSpace(string(aux.Connection)) == "[]" {
			d.Connection = nil
		} else {
			d.Connection = &DBaaSConnectionInfo{}
			if err := json.Unmarshal(aux.Connection, d.Connection); err != nil {
				return err
			}
		}
	}
	return nil
}

type DBaasBackupSchedule struct {
	Id            *int64  `json:"id,omitempty"`
	Name          *string `json:"name,omitempty"`
	ScheduledAt   *string `json:"scheduled_at,omitempty"`
	Retention     *int64  `json:"retention,omitempty"`
	IsPitrEnabled *bool   `json:"is_pitr_enabled,omitempty"`
}

type DBaaSCreateInfo struct {
	Id             int64  `json:"id"`
	RootPassword   string `json:"admin_password"`
	KubeIdentifier string `json:"kube_identifier"`
}

type DBaaSConnectionInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Ca       string `json:"ca"`
}

type DBaaSBackup struct {
	Id          string `json:"id,omitempty"`
	Location    string `json:"location,omitempty"`
	CreatedAt   uint64 `json:"created_at,omitempty"`
	CompletedAt uint64 `json:"completed_at,omitempty"`
	Status      string `json:"status,omitempty"`
}

type DBaaSRestore struct {
	Id           string           `json:"id,omitempty"`
	BackupSource string           `json:"backup_source,omitempty"`
	CreatedAt    uint64           `json:"created_at,omitempty"`
	Status       string           `json:"status,omitempty"`
	NewService   *DBaaSCreateInfo `json:"new_service,omitempty"`
}

func (dbaas *DBaaS) Key() string {
	return fmt.Sprintf("%d-%d-%d", dbaas.Project.PublicCloudId, dbaas.Project.ProjectId, dbaas.Id)
}

type DBaaSProject struct {
	PublicCloudId int64 `json:"public_cloud_id,omitempty"`
	ProjectId     int64 `json:"id,omitempty"`
}
