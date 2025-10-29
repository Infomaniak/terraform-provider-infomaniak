package dbaas

import "terraform-provider-infomaniak/internal/provider/registry"

func Register() {
	registry.RegisterResource(NewDBaasResource)
	registry.RegisterResource(NewDBaasBackupScheduleResource)

	registry.RegisterDataSource(NewDBaasDataSource)
	registry.RegisterDataSource(NewDBaasPacksDataSource)
	registry.RegisterDataSource(NewDBaasConstsDataSource)
}
