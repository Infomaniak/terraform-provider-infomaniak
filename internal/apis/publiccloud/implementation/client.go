package implementation

import (
	"fmt"
	"terraform-provider-infomaniak/internal/apis/helpers"
	"terraform-provider-infomaniak/internal/apis/publiccloud"
	"time"

	"resty.dev/v3"
)

// asyncTimeout is the upper bound for resolving delayed Public Cloud
// operations (project/user create-delete). Projects normally finish under a
// minute; we leave generous headroom for slow zones.
const asyncTimeout = 10 * time.Minute

var _ publiccloud.Api = (*Client)(nil)

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

func (c *Client) ListPublicClouds(accountId int64) ([]*publiccloud.PublicCloud, error) {
	var result helpers.NormalizedApiResponse[[]*publiccloud.PublicCloud]
	resp, err := c.resty.R().
		SetQueryParam("account_id", fmt.Sprint(accountId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointPublicClouds)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, result.Error
	}
	return result.Data, nil
}

func (c *Client) GetPublicCloud(publicCloudId int64) (*publiccloud.PublicCloud, error) {
	var result helpers.NormalizedApiResponse[*publiccloud.PublicCloud]
	resp, err := c.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointPublicCloud)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, result.Error
	}
	return result.Data, nil
}

// UpdatePublicCloud PATCHes the writable fields (customer_name, description,
// bill_reference) of an existing Public Cloud product.
func (c *Client) UpdatePublicCloud(input *publiccloud.PublicCloud) error {
	body := map[string]string{}
	if input.CustomerName != "" {
		body["customer_name"] = input.CustomerName
	}
	if input.Description != "" {
		body["description"] = input.Description
	}
	if input.BillReference != "" {
		body["bill_reference"] = input.BillReference
	}

	return c.patchAsync(EndpointPublicCloud, map[string]string{
		"public_cloud_id": fmt.Sprint(input.Id),
	}, body)
}

func (c *Client) GetConfig(accountId int64) (*publiccloud.Config, error) {
	var result helpers.NormalizedApiResponse[*publiccloud.Config]
	resp, err := c.resty.R().
		SetQueryParam("account_id", fmt.Sprint(accountId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointConfig)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, result.Error
	}
	return result.Data, nil
}

func (c *Client) GetAccesses(accountId int64) (*publiccloud.Accesses, error) {
	var result helpers.NormalizedApiResponse[*publiccloud.Accesses]
	resp, err := c.resty.R().
		SetQueryParam("account_id", fmt.Sprint(accountId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointAccesses)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, result.Error
	}
	return result.Data, nil
}

// CreateProject calls POST /projects or POST /projects/invite depending on
// input.Invite. Returns the new project id.
func (c *Client) CreateProject(publicCloudId int64, input *publiccloud.ProjectCreate) (int64, error) {
	endpoint := EndpointProjects
	if input.Invite {
		endpoint = EndpointProjectsInvite
	}
	return c.postAsyncInt64(endpoint, map[string]string{
		"public_cloud_id": fmt.Sprint(publicCloudId),
	}, input)
}

func (c *Client) UpdateProject(input *publiccloud.Project) error {
	body := map[string]string{}
	if input.Name != "" {
		body["name"] = input.Name
	}
	return c.patchAsync(EndpointProject, map[string]string{
		"public_cloud_id":         fmt.Sprint(input.PublicCloudId),
		"public_cloud_project_id": fmt.Sprint(input.Id),
	}, body)
}

func (c *Client) DeleteProject(publicCloudId, projectId int64) error {
	return c.deleteAsync(EndpointProject, map[string]string{
		"public_cloud_id":         fmt.Sprint(publicCloudId),
		"public_cloud_project_id": fmt.Sprint(projectId),
	})
}

func (c *Client) GetProject(publicCloudId, projectId int64) (*publiccloud.Project, error) {
	var result helpers.NormalizedApiResponse[*publiccloud.Project]
	resp, err := c.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(projectId)).
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

func (c *Client) GetUser(publicCloudId, projectId, userId int64) (*publiccloud.User, error) {
	var result helpers.NormalizedApiResponse[*publiccloud.User]
	resp, err := c.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(projectId)).
		SetPathParam("public_cloud_user_id", fmt.Sprint(userId)).
		SetResult(&result).
		SetError(&result).
		Get(EndpointUser)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, result.Error
	}
	return result.Data, nil
}

func (c *Client) CreateUser(publicCloudId, projectId int64, input *publiccloud.UserCreate) (int64, error) {
	endpoint := EndpointUsers
	if input.Invite {
		endpoint = EndpointUsersInvite
	}
	return c.postAsyncInt64(endpoint, map[string]string{
		"public_cloud_id":         fmt.Sprint(publicCloudId),
		"public_cloud_project_id": fmt.Sprint(projectId),
	}, input)
}

func (c *Client) UpdateUser(publicCloudId, projectId, userId int64, input *publiccloud.UserUpdate) error {
	return c.patchAsync(EndpointUser, map[string]string{
		"public_cloud_id":         fmt.Sprint(publicCloudId),
		"public_cloud_project_id": fmt.Sprint(projectId),
		"public_cloud_user_id":    fmt.Sprint(userId),
	}, input)
}

func (c *Client) DeleteUser(publicCloudId, projectId, userId int64) error {
	return c.deleteAsync(EndpointUser, map[string]string{
		"public_cloud_id":         fmt.Sprint(publicCloudId),
		"public_cloud_project_id": fmt.Sprint(projectId),
		"public_cloud_user_id":    fmt.Sprint(userId),
	})
}

// postAsyncInt64 issues a POST on an endpoint known to be asynchronous and
// returns the final integer id resolved from the async task.
func (c *Client) postAsyncInt64(endpoint string, pathParams map[string]string, body any) (int64, error) {
	var raw helpers.AsyncResponse
	req := c.resty.R().SetBody(body).SetResult(&raw).SetError(&raw)
	for k, v := range pathParams {
		req = req.SetPathParam(k, v)
	}
	resp, err := req.Post(endpoint)
	if err != nil {
		return 0, err
	}
	return helpers.ResolveAsync[int64](c.resty, resp, raw, asyncTimeout)
}

// patchAsync issues a PATCH on an endpoint known to be asynchronous and waits
// for the resulting async task to complete. The task payload is discarded.
func (c *Client) patchAsync(endpoint string, pathParams map[string]string, body any) error {
	var raw helpers.AsyncResponse
	req := c.resty.R().SetBody(body).SetResult(&raw).SetError(&raw)
	for k, v := range pathParams {
		req = req.SetPathParam(k, v)
	}
	resp, err := req.Patch(endpoint)
	if err != nil {
		return err
	}
	_, err = helpers.ResolveAsync[any](c.resty, resp, raw, asyncTimeout)
	return err
}

// deleteAsync issues a DELETE on an endpoint known to be asynchronous and
// waits for the resulting async task to complete. The task payload is discarded.
func (c *Client) deleteAsync(endpoint string, pathParams map[string]string) error {
	var raw helpers.AsyncResponse
	req := c.resty.R().SetResult(&raw).SetError(&raw)
	for k, v := range pathParams {
		req = req.SetPathParam(k, v)
	}
	resp, err := req.Delete(endpoint)
	if err != nil {
		return err
	}
	_, err = helpers.ResolveAsync[any](c.resty, resp, raw, asyncTimeout)
	return err
}

// GetOpenrc fetches the openrc.sh file content for a user in the given region.
// The endpoint returns the file as the raw response body (202 + text content).
func (c *Client) GetOpenrc(publicCloudId, projectId, userId int64, region string) (string, error) {
	req := c.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(projectId)).
		SetPathParam("public_cloud_user_id", fmt.Sprint(userId))
	if region != "" {
		req = req.SetQueryParam("region", region)
	}
	resp, err := req.Get(EndpointUserOpenrc)
	if err != nil {
		return "", err
	}
	if resp.IsError() {
		return "", fmt.Errorf("openrc: HTTP %d: %s", resp.StatusCode(), resp.String())
	}
	return resp.String(), nil
}

// GetAuthentication fetches an authentication file (e.g. clouds.yaml) of a
// given type for a user in the given region.
func (c *Client) GetAuthentication(publicCloudId, projectId, userId int64, authType, region string) (string, error) {
	req := c.resty.R().
		SetPathParam("public_cloud_id", fmt.Sprint(publicCloudId)).
		SetPathParam("public_cloud_project_id", fmt.Sprint(projectId)).
		SetPathParam("public_cloud_user_id", fmt.Sprint(userId)).
		SetPathParam("type", authType)
	if region != "" {
		req = req.SetQueryParam("region", region)
	}
	resp, err := req.Get(EndpointUserAuthentication)
	if err != nil {
		return "", err
	}
	if resp.IsError() {
		return "", fmt.Errorf("authentication: HTTP %d: %s", resp.StatusCode(), resp.String())
	}
	return resp.String(), nil
}
