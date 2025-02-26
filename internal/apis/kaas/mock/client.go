package mock

import (
	"fmt"
	"regexp"
	"terraform-provider-infomaniak/internal/apis/kaas"
)

// Ensure that our client implements Api
var (
	_         kaas.Api = (*Client)(nil)
	dnsRegexp          = regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) GetPacks() ([]*kaas.KaasPack, error) {
	return []*kaas.KaasPack{
		{
			Id:          1,
			Name:        "standard",
			Description: "Standard Cluster",
		},
		{
			Id:          2,
			Name:        "pro",
			Description: "Pro Cluster",
		},
	}, nil
}

func (c *Client) GetVersions() ([]string, error) {
	return []string{"1.29", "1.30", "1.31"}, nil
}

func (c *Client) GetKaas(publicCloudId int, publicCloudProjectId int, kaasId int) (*kaas.Kaas, error) {
	key := fmt.Sprintf("%d-%d-%d", publicCloudId, publicCloudProjectId, kaasId)
	obj, err := getFromCache[*kaas.Kaas](key)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (c *Client) CreateKaas(input *kaas.Kaas) (*kaas.Kaas, error) {
	// Checks
	if input.Project.PublicCloudId == 0 {
		return nil, fmt.Errorf("kaas is missing public cloud project id")
	}
	if input.Region == "" {
		return nil, fmt.Errorf("kaas is missing region")
	}

	var obj = kaas.Kaas{
		Project:    input.Project,
		Region:     input.Region,
		Kubeconfig: genKubeconfig(),
	}
	obj.Id = genId()

	return &obj, addToCache(&obj)
}

func (c *Client) UpdateKaas(input *kaas.Kaas) (*kaas.Kaas, error) {
	// Checks
	if input.Project.PublicCloudId == 0 {
		return nil, fmt.Errorf("kaas is missing public cloud project id")
	}
	if input.Id == 0 {
		return nil, fmt.Errorf("kaas is missing kaas id")
	}
	if input.Region != "" {
		return nil, fmt.Errorf("client cannot update region")
	}
	if input.Kubeconfig != "" {
		return nil, fmt.Errorf("client cannot update kubeconfig")
	}

	found, err := c.GetKaas(input.Project.PublicCloudId, input.Project.ProjectId, input.Id)
	if err != nil {
		return nil, err
	}

	var obj = kaas.Kaas{
		Id:      input.Id,
		Project: input.Project,

		Region:     input.Region,
		Kubeconfig: found.Kubeconfig,
	}

	return &obj, updateCache(&obj)
}

func (c *Client) DeleteKaas(publicCloudId int, publicCloudProjectId int, kaasId int) error {
	var obj = kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: publicCloudId,
			ProjectId:     publicCloudProjectId,
		},
		Id: kaasId,
	}

	return removeFromCache(&obj)
}

func (c *Client) GetInstancePool(publicCloudId int, publicCloudProjectId int, kaasId int, instancePoolId int) (*kaas.InstancePool, error) {
	_, err := c.GetKaas(publicCloudId, publicCloudProjectId, kaasId)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%d-%d", kaasId, instancePoolId)
	obj, err := getFromCache[*kaas.InstancePool](key)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (c *Client) CreateInstancePool(publicCloudId int, publicCloudProjectId int, input *kaas.InstancePool) (*kaas.InstancePool, error) {
	// Checks
	if publicCloudId == 0 {
		return nil, fmt.Errorf("instance pool is missing public cloud id")
	}
	if publicCloudProjectId == 0 {
		return nil, fmt.Errorf("instance pool is missing public cloud project id")
	}
	if input.KaasId == 0 {
		return nil, fmt.Errorf("instance pool is missing kaas id")
	}
	if !dnsRegexp.MatchString(input.Name) {
		return nil, fmt.Errorf("instance pool name should be a dns name according to RFC 1123")
	}
	if input.FlavorName == "" {
		return nil, fmt.Errorf("instance pool is missing flavor name")
	}
	if input.MinInstances < 0 {
		return nil, fmt.Errorf("instance pool min instances should be greater than 0")
	}
	// if input.MaxInstances < 0 {
	// 	return nil, fmt.Errorf("instance pool max instances should be greater than 0")
	// }
	// if input.MinInstances > input.MaxInstances {
	// 	return nil, fmt.Errorf("instance pool min instance should be lesser than (or equal) max")
	// }

	_, err := c.GetKaas(publicCloudId, publicCloudProjectId, input.KaasId)
	if err != nil {
		return nil, err
	}

	var obj = kaas.InstancePool{
		Id:     genId(),
		KaasId: input.KaasId,

		Name:         input.Name,
		FlavorName:   input.FlavorName,
		MinInstances: input.MinInstances,
		// MaxInstances: input.MaxInstances,
	}

	return &obj, addToCache(&obj)
}

func (c *Client) UpdateInstancePool(publicCloudId int, publicCloudProjectId int, input *kaas.InstancePool) (*kaas.InstancePool, error) {
	// Checks
	if publicCloudId == 0 {
		return nil, fmt.Errorf("instance pool is missing public cloud id")
	}
	if publicCloudProjectId == 0 {
		return nil, fmt.Errorf("instance pool is missing public cloud project id")
	}
	if input.KaasId == 0 {
		return nil, fmt.Errorf("instance pool is missing kaas id")
	}
	if input.Id == 0 {
		return nil, fmt.Errorf("instance pool is instance pool id")
	}
	if !dnsRegexp.MatchString(input.Name) {
		return nil, fmt.Errorf("instance pool name should be a dns name according to RFC 1123")
	}
	if input.FlavorName == "" {
		return nil, fmt.Errorf("instance pool is missing flavor name")
	}
	if input.MinInstances < 0 {
		return nil, fmt.Errorf("instance pool min instances should be greater than 0")
	}
	// if input.MaxInstances < 0 {
	// 	return nil, fmt.Errorf("instance pool max instances should be greater than 0")
	// }
	// if input.MinInstances > input.MaxInstances {
	// 	return nil, fmt.Errorf("instance pool min instance should be lesser than (or equal) max")
	// }

	_, err := c.GetKaas(publicCloudId, publicCloudProjectId, input.KaasId)
	if err != nil {
		return nil, err
	}

	_, err = c.GetInstancePool(publicCloudId, publicCloudProjectId, input.KaasId, input.Id)
	if err != nil {
		return nil, err
	}

	var obj = kaas.InstancePool{
		KaasId: input.KaasId,
		Id:     input.Id,

		Name:         input.Name,
		FlavorName:   input.FlavorName,
		MinInstances: input.MinInstances,
		// MaxInstances: input.MaxInstances,
	}

	return &obj, updateCache(&obj)
}

func (c *Client) DeleteInstancePool(publicCloudId int, publicCloudProjectId int, kaasId int, instancePoolId int) error {
	_, err := c.GetKaas(publicCloudId, publicCloudProjectId, kaasId)
	if err != nil {
		return err
	}

	var obj = kaas.InstancePool{
		KaasId: kaasId,
		Id:     instancePoolId,
	}

	return removeFromCache(&obj)
}
