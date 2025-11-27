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

type StringMap map[string]string

func (sm *StringMap) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*sm = make(StringMap)
	for key, value := range raw {
		switch v := value.(type) {
		case string:
			(*sm)[key] = v
		default:
			(*sm)[key] = fmt.Sprintf("%v", v)
		}
	}
	return nil
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

	Settings StringMap
}

type AllowedCIDRs struct {
	IpFilters []string `json:"ip_filters"`
}

type MySqlConfig struct {
	AutoIncrementIncrement          *int64   `json:"auto_increment_increment,omitempty"`
	AutoIncrementOffset             *int64   `json:"auto_increment_offset,omitempty"`
	CharacterSetServer              *string  `json:"character_set_server,omitempty"`
	ConnectTimeout                  *int64   `json:"connect_timeout,omitempty"`
	GroupConcatMaxLen               *int64   `json:"group_concat_max_len,omitempty"`
	InformationSchemaStatsExpiry    *int64   `json:"information_schema_stats_expiry,omitempty"`
	InnodbChangeBufferMaxSize       *int64   `json:"innodb_change_buffer_max_size,omitempty"`
	InnodbFlushNeighbors            *int64   `json:"innodb_flush_neighbors,omitempty"`
	InnodbFtMaxTokenSize            *int64   `json:"innodb_ft_max_token_size,omitempty"`
	InnodbFtMinTokenSize            *int64   `json:"innodb_ft_min_token_size,omitempty"`
	InnodbFtServerStopwordTable     *string  `json:"innodb_ft_server_stopword_table,omitempty"`
	InnodbLockWaitTimeout           *int64   `json:"innodb_lock_wait_timeout,omitempty"`
	InnodbLogBufferSize             *int64   `json:"innodb_log_buffer_size,omitempty"`
	InnodbOnlineAlterLogMaxSize     *int64   `json:"innodb_online_alter_log_max_size,omitempty"`
	InnodbPrintAllDeadlocks         *string  `json:"innodb_print_all_deadlocks,omitempty"`
	InnodbReadIoThreads             *int64   `json:"innodb_read_io_threads,omitempty"`
	InnodbRollbackOnTimeout         *string  `json:"innodb_rollback_on_timeout,omitempty"`
	InnodbStatsPersistentSamplePages *int64  `json:"innodb_stats_persistent_sample_pages,omitempty"`
	InnodbThreadConcurrency         *int64   `json:"innodb_thread_concurrency,omitempty"`
	InnodbWriteIoThreads            *int64   `json:"innodb_write_io_threads,omitempty"`
	InteractiveTimeout              *int64   `json:"interactive_timeout,omitempty"`
	LockWaitTimeout                 *int64   `json:"lock_wait_timeout,omitempty"`
	LogBinTrustFunctionCreators     *string  `json:"log_bin_trust_function_creators,omitempty"`
	LongQueryTime                   *float64 `json:"long_query_time,omitempty"`
	MaxAllowedPacket                *int64   `json:"max_allowed_packet,omitempty"`
	MaxConnections                  *int64   `json:"max_connections,omitempty"`
	MaxDigestLength                 *int64   `json:"max_digest_length,omitempty"`
	MaxHeapTableSize                *int64   `json:"max_heap_table_size,omitempty"`
	MaxPreparedStmtCount            *int64   `json:"max_prepared_stmt_count,omitempty"`
	MinExaminedRowLimit             *int64   `json:"min_examined_row_limit,omitempty"`
	NetBufferLength                 *int64   `json:"net_buffer_length,omitempty"`
	NetReadTimeout                  *int64   `json:"net_read_timeout,omitempty"`
	NetWriteTimeout                 *int64   `json:"net_write_timeout,omitempty"`
	PerformanceSchemaMaxDigestLength *int64  `json:"performance_schema_max_digest_length,omitempty"`
	RequireSecureTransport          *string  `json:"require_secure_transport,omitempty"`
	SortBufferSize                  *int64   `json:"sort_buffer_size,omitempty"`
	SqlMode                         []string `json:"sql_mode,omitempty"`
	TableDefinitionCache            *int64   `json:"table_definition_cache,omitempty"`
	TableOpenCache                  *int64   `json:"table_open_cache,omitempty"`
	TableOpenCacheInstances         *int64   `json:"table_open_cache_instances,omitempty"`
	ThreadStack                     *int64   `json:"thread_stack,omitempty"`
	TransactionIsolation            *string  `json:"transaction_isolation,omitempty"`
	WaitTimeout                     *int64   `json:"wait_timeout,omitempty"`
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
