package dbaas

type Api interface {
	FindPack(dbType string, name string) (*DBaaSPack, error)

	GetDBaaS(publicCloudId int, publicCloudProjectId int, DBaaSId int) (*DBaaS, error)
	CreateDBaaS(input *DBaaS) (*DBaaSCreateInfo, error)
	UpdateDBaaS(input *DBaaS) (bool, error)
	DeleteDBaaS(publicCloudId int, publicCloudProjectId int, DBaaSId int) (bool, error)

	PatchIpFilters(publicCloudId int, publicCloudProjectId int, dbaasId int, filters []string) (bool, error)
	GetIpFilters(publicCloudId int, publicCloudProjectId int, dbaasId int) ([]string, error)

	GetDBaasScheduleBackup(publicCloudId int, publicCloudProjectId int, dbaasId int) ([]DBaasBackupSchedule, error)

	UpdateDBaasScheduleBackup(publicCloudId int, publicCloudProjectId int, dbaasId int, backupSchedules *DBaasBackupSchedule) (bool, error)
	DeleteDBaasScheduleBackup(publicCloudId int, publicCloudProjectId int, dbaasId int) (bool, error)

	GetDbaasRegions() ([]string, error)
	GetDbaasTypes() ([]DbaasType, error)
	GetDbaasPacks(dbType string) ([]Pack, error)
}
