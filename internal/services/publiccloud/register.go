package publiccloud

import "terraform-provider-infomaniak/internal/provider/registry"

// Register wires this service's resources and data sources into the shared
// provider registry. Entries are added per phase as resources land.
func Register() {
	registry.RegisterResource(NewPublicCloudResource)
	registry.RegisterResource(NewPublicCloudProjectResource)
	registry.RegisterResource(NewPublicCloudUserResource)

	registry.RegisterDataSource(NewPublicCloudDataSource)
	registry.RegisterDataSource(NewPublicCloudsDataSource)
	registry.RegisterDataSource(NewPublicCloudConfigDataSource)
	registry.RegisterDataSource(NewPublicCloudAccessesDataSource)
	registry.RegisterDataSource(NewPublicCloudProjectDataSource)
	registry.RegisterDataSource(NewPublicCloudUserDataSource)
	registry.RegisterDataSource(NewPublicCloudOpenrcDataSource)
	registry.RegisterDataSource(NewPublicCloudUserAuthenticationDataSource)
}
