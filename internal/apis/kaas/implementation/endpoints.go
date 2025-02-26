package implementation

import (
	"net/http"
	"terraform-provider-infomaniak/internal/apis/endpoints"
)

var (
	GetPacks    = endpoints.NewEndpoint(http.MethodGet, "/1/public_clouds/kaas/packs")
	GetVersions = endpoints.NewEndpoint(http.MethodGet, "/1/public_clouds/kaas/versions")

	CreateKaas = endpoints.NewEndpoint(http.MethodPost, "/1/public_clouds/{public_cloud_id}/projects/{project_id}/kaas/")
	GetKaas    = endpoints.NewEndpoint(http.MethodGet, "/1/public_clouds/{public_cloud_id}/projects/{project_id}/kaas/{kaas_id}/")
	UpdateKaas = endpoints.NewEndpoint(http.MethodPatch, "/1/public_clouds/{public_cloud_id}/projects/{project_id}/kaas/{kaas_id}/")
	DeleteKaas = endpoints.NewEndpoint(http.MethodDelete, "/1/public_clouds/{public_cloud_id}/projects/{project_id}/kaas/{kaas_id}/")

	CreateInstancePool = endpoints.NewEndpoint(http.MethodPost, "/1/public_clouds/{public_cloud_id}/projects/{project_id}/kaas/{kaas_id}/")
	GetInstancePool    = endpoints.NewEndpoint(http.MethodGet, "/1/public_clouds/{public_cloud_id}/projects/{project_id}/kaas/{kaas_id}/instance_pools/{kaas_instance_pool_id}/")
	UpdateInstancePool = endpoints.NewEndpoint(http.MethodPatch, "/1/public_clouds/{public_cloud_id}/projects/{project_id}/kaas/{kaas_id}/instance_pools/{kaas_instance_pool_id}/")
	DeleteInstancePool = endpoints.NewEndpoint(http.MethodDelete, "/1/public_clouds/{public_cloud_id}/projects/{project_id}/kaas/{kaas_id}/instance_pools/{kaas_instance_pool_id}/")
)
