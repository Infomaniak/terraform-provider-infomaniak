package implementation

var (
	EndpointRegions = "/1/public_clouds/dbaas/regions"
	EndpointFlavors = "/1/public_clouds/dbaas/flavors"
	EndpointTypes   = "/1/public_clouds/dbaas/types"
	EndpointPacks   = "/1/public_clouds/dbaas/packs"
	EndpointPack    = "/1/public_clouds/dbaas/packs/1"

	EndpointDatabases = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas"
	EndpointDatabase  = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}"

	EndpointDatabaseIpFilter = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}/ip_filters"

	EndpointDatabaseBackupSchedules = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}/backup_schedules"
	EndpointDatabaseBackupSchedule  = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/dbaas/{dbaas_id}/backup_schedules/{schedule_id}"

	EndpointDbaasDataRegion = "/1/public_clouds/dbaas/regions"
	EndpointDbaasDataPacks  = "/1/public_clouds/dbaas/packs"
	EndpointDbaasDataTypes  = "/1/public_clouds/dbaas/types"
)
