package dbaas

type Api interface {
	FindPack(dbType string, name string) (*DBaaSPack, error)

	GetDBaaS(publicCloudId int, publicCloudProjectId int, DBaaSId int) (*DBaaS, error)
	GetPassword(publicCloudId int, publicCloudProjectId int, DBaaSId int) (*DBaaSConnectionInfo, error)
	CreateDBaaS(input *DBaaS) (int, error)
	UpdateDBaaS(input *DBaaS) (bool, error)
	DeleteDBaaS(publicCloudId int, publicCloudProjectId int, DBaaSId int) (bool, error)

	CreateBackup(publicCloudId int, publicCloudProjectId int, dbaasId int) (*DBaaSBackup, error)
	GetBackup(publicCloudId int, publicCloudProjectId int, dbaasId int, backupId string) (*DBaaSBackup, error)
	DeleteBackup(publicCloudId int, publicCloudProjectId int, dbaasId int, backupId string) (bool, error)

	CreateRestore(publicCloudId int, publicCloudProjectId int, dbaasId int, backupId string) (*DBaaSRestore, error)
	GetRestore(publicCloudId int, publicCloudProjectId int, dbaasId int, restoreId string) (*DBaaSRestore, error)
}
