package kaas

type Api interface {
	GetPacks() ([]*KaasPack, error)
	GetVersions() ([]string, error)

	GetKaas(publicCloudId int64, publicCloudProjectId int64, kaasId int64) (*Kaas, error)
	CreateKaas(input *Kaas) (int64, error)
	UpdateKaas(input *Kaas) (bool, error)
	DeleteKaas(publicCloudId int64, publicCloudProjectId int64, kaasId int64) (bool, error)

	GetKubeconfig(publicCloudId int64, publicCloudProjectId int64, kaasId int64) (string, error)

	GetInstancePool(publicCloudId int64, publicCloudProjectId int64, kaasId int64, instancePoolId int64) (*InstancePool, error)
	CreateInstancePool(publicCloudId int64, publicCloudProjectId int64, input *InstancePool) (int64, error)
	UpdateInstancePool(publicCloudId int64, publicCloudProjectId int64, input *InstancePool) (bool, error)
	DeleteInstancePool(publicCloudId int64, publicCloudProjectId int64, kaasId int64, instancePoolId int64) (bool, error)

	GetApiserverParams(publicCloudId int, projectId int, kaasId int) (*Apiserver, error)
	PatchApiserverParams(input *Apiserver, publicCloudId int, projectId int, kaasId int) (bool, error)
	PatchIPFilters(cidrs []string, publicCloudId int, projectId int, kaasId int) (bool, error)
	GetIPFilters(publicCloudId int, projectId int, kaasId int) ([]string, error)
}
