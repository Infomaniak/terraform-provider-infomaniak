package implementation

// Endpoint templates use resty path-param placeholders. Per the upstream
// OpenAPI spec, only top-level reads and the PATCH on /1/public_clouds/{id}
// are documented for the Public Cloud product itself; ordering and deletion
// are not exposed through the public API.
var (
	EndpointPublicClouds = "/1/public_clouds"
	EndpointPublicCloud  = "/1/public_clouds/{public_cloud_id}"
	EndpointConfig       = "/1/public_clouds/config"
	EndpointAccesses     = "/1/public_clouds/accesses"

	EndpointProjects       = "/1/public_clouds/{public_cloud_id}/projects"
	EndpointProjectsInvite = "/1/public_clouds/{public_cloud_id}/projects/invite"
	EndpointProject        = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}"

	EndpointUsers              = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/users"
	EndpointUsersInvite        = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/users/invite"
	EndpointUser               = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/users/{public_cloud_user_id}"
	EndpointUserOpenrc         = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/users/{public_cloud_user_id}/openrc"
	EndpointUserAuthentication = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/users/{public_cloud_user_id}/authentication/{type}"
)
