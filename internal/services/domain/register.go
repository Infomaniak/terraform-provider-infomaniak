package domain

import "terraform-provider-infomaniak/internal/provider/registry"

func Register() {
	registry.RegisterResource(NewZoneResource)
	registry.RegisterResource(NewRecordResource)
}
