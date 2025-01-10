package kaas

type Api interface {
	GetKaas(pcpId, kaasId string) (*Kaas, error)
	CreateKaas(input *Kaas) (*Kaas, error)
	UpdateKaas(input *Kaas) (*Kaas, error)
	DeleteKaas(pcpId, kaasId string) error

	GetInstancePool(pcpId, kaasId, instancePoolId string) (*InstancePool, error)
	CreateInstancePool(input *InstancePool) (*InstancePool, error)
	UpdateInstancePool(input *InstancePool) (*InstancePool, error)
	DeleteInstancePool(pcpId, kaasId, instancePoolId string) error
}
