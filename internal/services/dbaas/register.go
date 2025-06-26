package dbaas

import "terraform-provider-infomaniak/internal/provider/registry"

func Register() {
	registry.RegisterResource(NewDBaasResource)
	registry.RegisterResource(NewDBaasBackupResource)
	registry.RegisterResource(NewDBaasRestoreResource)

	registry.RegisterDataSource(NewDBaasDataSource)
}
