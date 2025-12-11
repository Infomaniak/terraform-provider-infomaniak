package project

import "terraform-provider-infomaniak/internal/provider/registry"

func Register() {
	registry.RegisterResource(NewProjectResource)

	registry.RegisterDataSource(NewProjectDataSource)
}
