package dbaas

type Api interface {
	FindPack(dbType string, name string) (*DBaaSPack, error)

	GetDBaaS(publicCloudId int64, publicCloudProjectId int64, DBaaSId int64) (*DBaaS, error)
	CreateDBaaS(input *DBaaS) (*DBaaSCreateInfo, error)
	UpdateDBaaS(input *DBaaS) (bool, error)
	DeleteDBaaS(publicCloudId int64, publicCloudProjectId int64, DBaaSId int64) (bool, error)

	PatchIpFilters(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, filters AllowedCIDRs) (bool, error)
	GetIpFilters(publicCloudId int64, publicCloudProjectId int64, dbaasId int64) ([]string, error)

	GetDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64) (*DBaasBackupSchedule, error)
	CreateDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, backupSchedules *DBaasBackupSchedule) (*DBaasBackupScheduleCreateInfo, error)
	UpdateDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64, backupSchedules *DBaasBackupSchedule) (bool, error)
	DeleteDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64) (bool, error)

	GetDbaasRegions() ([]string, error)
	GetDbaasTypes() ([]*DbaasType, error)
	GetDbaasPacks(dbType string) ([]*Pack, error)
}
