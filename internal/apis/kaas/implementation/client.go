package implementation

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-infomaniak/internal/apis/helpers"
	"terraform-provider-infomaniak/internal/apis/kaas"

	"resty.dev/v3"
)

// Ensure that our client implements Api
var (
	_ kaas.Api = (*Client)(nil)
)

type Client struct {
	resty *resty.Client
}

func New(baseUri, token, version string) *Client {
	return &Client{
		resty: resty.New().
			SetBaseURL(baseUri).
			SetAuthToken(token).
			SetHeader("User-Agent", helpers.GetUserAgent(version)),
	}
}

func (client *Client) GetPacks() ([]*kaas.KaasPack, error) {
	var result helpers.NormalizedApiResponse[[]*kaas.KaasPack]

	resp, err := client.resty.R().
		SetResult(&result).
		SetError(&result).
		Get(EndpointPacks)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetFlavor(publicCloudId int64, publicCloudProjectId int64, region string, params kaas.KaasFlavorLookupParameters) (*kaas.KaasFlavor, error) {
	var result helpers.NormalizedApiResponse[[]*kaas.KaasFlavor]

	builder := client.resty.R().
		SetResult(&result).
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetQueryParam("region", region).
		SetDebug(true).
		SetError(&result)

	if params.Name != nil {
		builder.SetQueryParam("filter[search]", *params.Name)
	}

	if params.Cpu != nil {
		builder.SetQueryParam("filter[cpu]", strconv.FormatInt(*params.Cpu, 10))
	}

	if params.Ram != nil {
		builder.SetQueryParam("filter[ram]", strconv.FormatInt(*params.Ram, 10))
	}

	if params.Storage != nil {
		builder.SetQueryParam("filter[storage]", strconv.FormatInt(*params.Storage, 10))
	}

	if params.IsMemoryOptimized != nil {
		builder.SetQueryParam("filter[memory_optimized]", strconv.FormatBool(*params.IsMemoryOptimized))
	}

	if params.IsIopsOptimized != nil {
		builder.SetQueryParam("filter[iops_optimized]", strconv.FormatBool(*params.IsIopsOptimized))
	}

	if params.IsGpuOptimized != nil {
		builder.SetQueryParam("filter[gpu_optimized]", strconv.FormatBool(*params.IsGpuOptimized))
	}

	resp, err := builder.Get(EndpointKaasFlavors)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("flavor not found")
	}

	if len(result.Data) != 1 {
		flavors := strings.Builder{}
		for _, flavor := range result.Data {
			flavors.WriteString(flavor.Name)
			flavors.WriteString(", ")
		}
		return nil, fmt.Errorf("multiple flavors found, please refine your search\nFound flavors: %s", flavors.String())
	}

	return result.Data[0], nil
}

func (client *Client) GetVersions() ([]string, error) {
	var result helpers.NormalizedApiResponse[[]string]

	resp, err := client.resty.R().
		SetResult(&result).
		SetError(&result).
		Get(EndpointVersions)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetKaas(publicCloudId int, publicCloudProjectId int, kaasId int) (*kaas.Kaas, error) {
	var result helpers.NormalizedApiResponse[*kaas.Kaas]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("kaas_id", fmt.Sprint(kaasId)).
		SetQueryParam("with", "packs,projects,instances,tags").
		SetResult(&result).
		SetError(&result).
		Get(EndpointKaas)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetKubeconfig(publicCloudId int, publicCloudProjectId int, kaasId int) (string, error) {
	var result helpers.NormalizedApiResponse[string]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("kaas_id", fmt.Sprint(kaasId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointKaasKubeconfig)
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", result.Error
	}

	return result.Data, nil
}

func (client *Client) CreateKaas(input *kaas.Kaas) (int, error) {
	var result helpers.NormalizedApiResponse[int]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(input.Project.PublicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(input.Project.ProjectId)).
		SetBody(input).
		SetResult(&result).
		SetError(&result).
		Post(EndpointKaases)
	if err != nil {
		return 0, err
	}

	if resp.IsError() {
		return 0, result.Error
	}

	return result.Data, nil
}

func (client *Client) UpdateKaas(input *kaas.Kaas) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(input.Project.PublicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(input.Project.ProjectId)).
		SetPathParam("kaas_id", fmt.Sprint(input.Id)).
		SetBody(input).
		SetResult(&result).
		SetError(&result).
		Patch(EndpointKaas)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) DeleteKaas(publicCloudId int, publicCloudProjectId int, kaasId int) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("kaas_id", fmt.Sprint(kaasId)).
		SetResult(&result).
		SetError(&result).
		Delete(EndpointKaas)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetInstancePool(publicCloudId int, publicCloudProjectId int, kaasId int, instancePoolId int) (*kaas.InstancePool, error) {
	var result helpers.NormalizedApiResponse[*kaas.InstancePool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("kaas_id", fmt.Sprint(kaasId)).
		SetPathParam("kaas_instance_pool_id", fmt.Sprint(instancePoolId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointInstancePool)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	// Default Max = Min
	if result.Data.MaxInstances == 0 {
		result.Data.MaxInstances = result.Data.MinInstances
	}

	return result.Data, nil
}

func (client *Client) CreateInstancePool(publicCloudId int, publicCloudProjectId int, input *kaas.InstancePool) (int, error) {
	var result helpers.NormalizedApiResponse[int]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("kaas_id", fmt.Sprint(input.KaasId)).
		SetBody(input).
		SetResult(&result).
		SetError(&result).
		Post(EndpointInstancePools)
	if err != nil {
		return 0, err
	}

	if resp.IsError() {
		return 0, result.Error
	}

	return result.Data, nil
}

func (client *Client) UpdateInstancePool(publicCloudId int, publicCloudProjectId int, input *kaas.InstancePool) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("kaas_id", fmt.Sprint(input.KaasId)).
		SetPathParam("kaas_instance_pool_id", fmt.Sprint(input.Id)).
		SetBody(input).
		SetResult(&result).
		SetError(&result).
		Patch(EndpointInstancePool)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) DeleteInstancePool(publicCloudId int, publicCloudProjectId int, kaasId int, instancePoolId int) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("kaas_id", fmt.Sprint(kaasId)).
		SetPathParam("kaas_instance_pool_id", fmt.Sprint(instancePoolId)).
		SetResult(&result).
		SetError(&result).
		Delete(EndpointInstancePool)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) PatchApiserverParams(input *kaas.Apiserver, publicCloudId int, projectId int, kaasId int) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]
	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(projectId)).
		SetPathParam("kaas_id", fmt.Sprint(kaasId)).
		SetBody(input).
		SetResult(&result).
		SetError(&result).
		Patch(EndpointApiserver)

	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetApiserverParams(publicCloudId int, projectId int, kaasId int) (*kaas.Apiserver, error) {

	var result helpers.NormalizedApiResponse[*kaas.Apiserver]
	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(projectId)).
		SetPathParam("kaas_id", fmt.Sprint(kaasId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointApiserver)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}
