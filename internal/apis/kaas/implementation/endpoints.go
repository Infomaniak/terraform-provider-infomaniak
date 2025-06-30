package implementation

var (
	EndpointPacks    = "/1/public_clouds/kaas/packs"
	EndpointVersions = "/1/public_clouds/kaas/versions"

	EndpointKaases         = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/kaas"
	EndpointKaas           = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/kaas/{kaas_id}"
	EndpointKaasKubeconfig = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/kaas/{kaas_id}/kube_config"

	EndpointInstancePools = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/kaas/{kaas_id}/instance_pools"
	EndpointInstancePool  = "/1/public_clouds/{public_cloud_id}/projects/{public_cloud_project_id}/kaas/{kaas_id}/instance_pools/{kaas_instance_pool_id}"
)
