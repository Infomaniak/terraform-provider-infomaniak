package implementation

import (
	"fmt"
	"strconv"
	"strings"
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

func (client *Client) GetDBaaS(publicCloudId int64, publicCloudProjectId int64, dbaasId int64) (*dbaas.DBaaS, error) {
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

func (client *Client) DeleteDBaaS(publicCloudId int64, publicCloudProjectId int64, dbaasId int64) (bool, error) {
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

func (client *Client) PatchIpFilters(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, filters dbaas.AllowedCIDRs) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetBody(filters).
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

func (client *Client) PutConfiguration(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, configuration dbaas.MySqlConfig) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetBody(configuration).
		SetResult(&result).
		SetError(&result).
		SetDebug(true).
		Put(EndpointDatabaseConfiguration)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetConfiguration(publicCloudId int64, publicCloudProjectId int64, dbaasId int64) (*dbaas.MySqlConfig, error) {
	var result helpers.NormalizedApiResponse[*dbaas.MySqlConfig]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetDebug(true).
		SetResult(&result).
		SetError(&result).
		Get(EndpointDatabaseConfiguration)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	fmt.Printf("icilastruct %+#v \n", result.Data)

	return result.Data, nil
}

func (client *Client) GetIpFilters(publicCloudId int64, publicCloudProjectId int64, dbaasId int64) ([]string, error) {
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

func (client *Client) CreateDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, backupSchedules *dbaas.DBaasBackupSchedule) (int64, error) {
	var result helpers.NormalizedApiResponse[int64]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetPathParam("dbaas_id", fmt.Sprint(dbaasId)).
		SetBody(backupSchedules).
		SetResult(&result).
		SetError(&result).
		Post(EndpointDatabaseBackupSchedules)
	if err != nil {
		return 0, err
	}

	if resp.IsError() {
		return 0, result.Error
	}

	return result.Data, nil
}

func (client *Client) UpdateDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64, backupSchedules *dbaas.DBaasBackupSchedule) (bool, error) {
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

func (client *Client) GetDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64) (*dbaas.DBaasBackupSchedule, error) {
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

func (client *Client) DeleteDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64) (bool, error) {
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

func (client *Client) GetDbaasPack(params dbaas.PackFilter) (*dbaas.Pack, error) {
	var result helpers.NormalizedApiResponse[[]*dbaas.Pack]

	builder := client.resty.R().
		SetResult(&result).
		SetError(&result).
		SetQueryParam("filter[type]", params.DbType)

	if params.Name != nil {
		builder = builder.SetQueryParam("filter[names][]", *params.Name)
	}

	if params.Group != nil {
		builder = builder.SetQueryParam("filter[groups][]", *params.Group)
	}

	if params.Instances != nil {
		builder = builder.SetQueryParam("filter[instances]", strconv.FormatInt(*params.Instances, 10))
	}

	if params.Cpu != nil {
		builder = builder.SetQueryParam("filter[cpu]", strconv.FormatInt(*params.Cpu, 10))
	}

	if params.Ram != nil {
		builder = builder.SetQueryParam("filter[ram]", strconv.FormatInt(*params.Ram, 10))
	}

	if params.Storage != nil {
		builder = builder.SetQueryParam("filter[storage]", strconv.FormatInt(*params.Storage, 10))
	}

	resp, err := builder.Get(EndpointPacks)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	data := result.Data
	if len(data) == 0 {
		return nil, fmt.Errorf("pack not found")
	}

	if len(data) != 1 {
		packs := strings.Builder{}
		for _, pack := range data {
			packs.WriteString(pack.Name)
			packs.WriteString(", ")
		}
		return nil, fmt.Errorf("multiple packs found, please refine your search\nfound packs: %s", packs.String())
	}

	return data[0], nil
}
