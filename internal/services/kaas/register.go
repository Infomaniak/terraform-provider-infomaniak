package kaas

import "terraform-provider-infomaniak/internal/provider/registry"

func init() {
	registry.RegisterResource(NewKaasResource)
	registry.RegisterResource(NewKaasInstancePoolResource)

	registry.RegisterDataSource(NewKaasDataSource)
	registry.RegisterDataSource(NewKaasInstancePoolDataSource)
}
