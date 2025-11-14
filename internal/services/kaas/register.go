package kaas

import "terraform-provider-infomaniak/internal/provider/registry"

func Register() {
	registry.RegisterResource(NewKaasResource)
	registry.RegisterResource(NewKaasInstancePoolResource)

	registry.RegisterDataSource(NewKaasDataSource)
	registry.RegisterDataSource(NewKaasInstancePoolDataSource)
	registry.RegisterDataSource(NewKaasInstancePoolFlavorDataSource)
	registry.RegisterDataSource(NewKaasRegionsDataSource)
}
