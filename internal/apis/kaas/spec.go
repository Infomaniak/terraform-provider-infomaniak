package kaas

type Api interface {
	GetPacks() ([]*KaasPack, error)
	GetVersions() ([]string, error)

	GetKaas(publicCloudId int, publicCloudProjectId int, kaasId int) (*Kaas, error)
	CreateKaas(input *Kaas) (*Kaas, error)
	UpdateKaas(input *Kaas) (*Kaas, error)
	DeleteKaas(publicCloudId int, publicCloudProjectId int, kaasId int) error

	GetInstancePool(publicCloudId int, publicCloudProjectId int, kaasId int, instancePoolId int) (*InstancePool, error)
	CreateInstancePool(publicCloudId int, publicCloudProjectId int, input *InstancePool) (*InstancePool, error)
	UpdateInstancePool(publicCloudId int, publicCloudProjectId int, input *InstancePool) (*InstancePool, error)
	DeleteInstancePool(publicCloudId int, publicCloudProjectId int, kaasId int, instancePoolId int) error
}
