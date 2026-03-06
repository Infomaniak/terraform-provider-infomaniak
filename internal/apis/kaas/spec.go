package kaas

type Api interface {
	GetPacks() ([]*KaasPack, error)
	GetRegions() ([]string, error)
	GetFlavor(publicCloudId int64, publicCloudProjectId int64, region string, params KaasFlavorLookupParameters) (*KaasFlavor, error)
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

	GetApiserverParams(publicCloudId int64, projectId int64, kaasId int64) (*Apiserver, error)
	PatchApiserverParams(input *Apiserver, publicCloudId int64, projectId int64, kaasId int64) (bool, error)
}
