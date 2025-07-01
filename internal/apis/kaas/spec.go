package kaas

type Api interface {
	GetPacks() ([]*KaasPack, error)
	GetVersions() ([]string, error)

	GetKaas(publicCloudId int, publicCloudProjectId int, kaasId int) (*Kaas, error)
	CreateKaas(input *Kaas) (int, error)
	UpdateKaas(input *Kaas) (bool, error)
	DeleteKaas(publicCloudId int, publicCloudProjectId int, kaasId int) (bool, error)

	GetKubeconfig(publicCloudId int, publicCloudProjectId int, kaasId int) (string, error)

	GetInstancePool(publicCloudId int, publicCloudProjectId int, kaasId int, instancePoolId int) (*InstancePool, error)
	CreateInstancePool(publicCloudId int, publicCloudProjectId int, input *InstancePool) (int, error)
	UpdateInstancePool(publicCloudId int, publicCloudProjectId int, input *InstancePool) (bool, error)
	DeleteInstancePool(publicCloudId int, publicCloudProjectId int, kaasId int, instancePoolId int) (bool, error)

	GetApiserverParams(publicCloudId int, projectId int, kaasId int) (*Apiserver, error)
	CreateApiserverParams(input *Apiserver, publicCloudId int, projectId int, kaasId int) (bool, error)
	PatchApiserverParams(input *Apiserver, publicCloudId int, projectId int, kaasId int) (bool, error)
}
