package implementation

var (
	EndpointRegions = "/1/public_clouds/dbaas/regions"
	EndpointFlavors = "/1/public_clouds/dbaas/flavors"
	EndpointTypes   = "/1/public_clouds/dbaas/types"
	EndpointPacks   = "/1/public_clouds/dbaas/packs"
	EndpointPack    = "/1/public_clouds/dbaas/packs/1"

	EndpointDatabases        = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas"
	EndpointDatabase         = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}"
	EndpointDatabaseBackups  = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}/backups"
	EndpointDatabaseBackup   = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}/backups/{backup_id}"
	EndpointDatabaseRestores = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}/restores"
	EndpointDatabaseRestore  = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}/restores/{restore_id}"

	EndpointDatabaseIpFilter = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}/ip_filters"
	EndpointDatabaseBackupSchedules = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}/backup_schedules"
)
