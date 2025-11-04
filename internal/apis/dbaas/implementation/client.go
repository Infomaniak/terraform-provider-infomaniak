package implementation

import (
	"fmt"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"terraform-provider-infomaniak/internal/apis/helpers"

	"resty.dev/v3"
)

// Ensure that our client implements Api
var (
	_ dbaas.Api = (*Client)(nil)
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

func (client *Client) FindPack(dbType string, name string) (*dbaas.DBaaSPack, error) {
	var result helpers.NormalizedApiResponse[[]*dbaas.DBaaSPack]

	resp, err := client.resty.R().
		SetResult(&result).
		SetError(&result).
		SetQueryParam("filter[type]", dbType).
		SetQueryParam("filter[names][]", name).
		Get(EndpointPacks)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	data := result.Data
	if len(data) != 1 || data[0].Name != name {
		return nil, fmt.Errorf("pack not found")
	}

	return data[0], nil
}

func (client *Client) GetDBaaS(publicCloudId int, publicCloudProjectId int, dbaasId int) (*dbaas.DBaaS, error) {
	var result helpers.NormalizedApiResponse[*dbaas.DBaaS]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetQueryParam("with", "packs,projects,tags,connection").
		SetResult(&result).
		SetError(&result).
		Get(EndpointDatabase)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) CreateDBaaS(input *dbaas.DBaaS) (*dbaas.DBaaSCreateInfo, error) {
	var result helpers.NormalizedApiResponse[*dbaas.DBaaSCreateInfo]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(input.Project.PublicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(input.Project.ProjectId)).
		SetBody(input).
		SetResult(&result).
		SetError(&result).
		Post(EndpointDatabases)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) UpdateDBaaS(input *dbaas.DBaaS) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(input.Project.PublicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(input.Project.ProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(input.Id)).
		SetBody(input).
		SetResult(&result).
		SetError(&result).
		Patch(EndpointDatabase)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) DeleteDBaaS(publicCloudId int, publicCloudProjectId int, dbaasId int) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetResult(&result).
		SetError(&result).
		Delete(EndpointDatabase)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) PatchIpFilters(publicCloudId int, publicCloudProjectId int, dbaasId int, filters []string) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetBody(map[string][]string{
			"ip_filters": filters,
		}).
		SetResult(&result).
		SetError(&result).
		Put(EndpointDatabaseIpFilter)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetIpFilters(publicCloudId int, publicCloudProjectId int, dbaasId int) ([]string, error) {
	var result helpers.NormalizedApiResponse[[]string]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointDatabaseIpFilter)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) CreateDBaasScheduleBackup(publicCloudId int, publicCloudProjectId int, dbaasId int, backupSchedules *dbaas.DBaasBackupSchedule) (*dbaas.DBaasBackupScheduleCreateInfo, error) {
	var result helpers.NormalizedApiResponse[*dbaas.DBaasBackupScheduleCreateInfo]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetBody(backupSchedules).
		SetResult(&result).
		SetError(&result).
		Post(EndpointDatabaseBackupSchedules)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) UpdateDBaasScheduleBackup(publicCloudId int, publicCloudProjectId int, dbaasId int, id int, backupSchedules *dbaas.DBaasBackupSchedule) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetPathParam("schedule_id", fmt.Sprint(id)).
		SetBody(backupSchedules).
		SetResult(&result).
		SetError(&result).
		Patch(EndpointDatabaseBackupSchedule)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetDBaasScheduleBackup(publicCloudId int, publicCloudProjectId int, dbaasId int, id int) (*dbaas.DBaasBackupSchedule, error) {
	var result helpers.NormalizedApiResponse[*dbaas.DBaasBackupSchedule]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetPathParam("schedule_id", fmt.Sprint(id)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointDatabaseBackupSchedule)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) DeleteDBaasScheduleBackup(publicCloudId int, publicCloudProjectId int, dbaasId int, id int) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetPathParam("schedule_id", fmt.Sprint(id)).
		SetResult(&result).
		SetError(&result).
		Delete(EndpointDatabaseBackupSchedule)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetDbaasRegions() ([]string, error) {
	var result helpers.NormalizedApiResponse[[]string]

	resp, err := client.resty.R().
		SetResult(&result).
		SetError(&result).
		Get(EndpointDbaasDataRegion)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetDbaasTypes() ([]*dbaas.DbaasType, error) {
	var result helpers.NormalizedApiResponse[[]*dbaas.DbaasType]

	resp, err := client.resty.R().
		SetResult(&result).
		SetError(&result).
		Get(EndpointDbaasDataTypes)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetDbaasPacks(dbType string) ([]*dbaas.Pack, error) {
	var result helpers.NormalizedApiResponse[[]*dbaas.Pack]

	resp, err := client.resty.R().
		SetQueryParams(map[string]string{
			"type":     dbType,
			"per_page": "1000",
		}).
		SetResult(&result).
		SetError(&result).
		Get(EndpointDbaasDataPacks)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}
