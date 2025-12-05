package implementation

import (
	"fmt"
	"terraform-provider-infomaniak/internal/apis/project"

	"terraform-provider-infomaniak/internal/apis/helpers"

	"resty.dev/v3"
)

// Ensure that our client implements Api
var (
	_ project.Api = (*Client)(nil)
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

func (client *Client) GetProject(publicCloudId int, publicCloudProjectId int) (*project.Project, error) {
	var result helpers.NormalizedApiResponse[*project.Project]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointProject)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) CreateProject(input *project.CreateProject) (int, error) {
	var result helpers.NormalizedApiResponse[int]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(input.PublicCloudId)).
		SetBody(input).
		SetResult(&result).
		SetError(&result).
		Post(EndpointProjects)
	if err != nil {
		return 0, err
	}

	if resp.IsError() {
		return 0, result.Error
	}

	return result.Data, nil
}

func (client *Client) UpdateProject(publicCloudId int, publicCloudProjectId int, input *project.UpdateProject) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetBody(input).
		SetResult(&result).
		SetError(&result).
		Patch(EndpointProject)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) DeleteProject(publicCloudId int, publicCloudProjectId int) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(publicCloudProjectId)).
		SetResult(&result).
		SetError(&result).
		Delete(EndpointProject)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}
