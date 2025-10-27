package dbaas

import "terraform-provider-infomaniak/internal/provider/registry"

func Register() {
	registry.RegisterResource(NewDBaasResource)

	registry.RegisterDataSource(NewDBaasDataSource)
}
