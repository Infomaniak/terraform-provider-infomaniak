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

	obj.Status = "Active"

	return obj, nil
}

func (client *Client) GetKubeconfig(publicCloudId int, publicCloudProjectId int, kaasId int) (string, error) {
	return genKubeconfig(), nil
}

func (c *Client) CreateKaas(input *kaas.Kaas) (int, error) {
	// Checks
	if input.Project.PublicCloudId == 0 {
		return 0, fmt.Errorf("kaas is missing public cloud project id")
	}
	if input.Region == "" {
		return 0, fmt.Errorf("kaas is missing region")
	}

	var obj = kaas.Kaas{
		Project:           input.Project,
		Region:            input.Region,
		PackId:            input.PackId,
		KubernetesVersion: input.KubernetesVersion,
		Name:              input.Name,
	}
	obj.Id = genId()

	return obj.Id, addToCache(&obj)
}

func (c *Client) UpdateKaas(input *kaas.Kaas) (bool, error) {
	// Checks
	if input.Project.PublicCloudId == 0 {
		return false, fmt.Errorf("kaas is missing public cloud project id")
	}
	if input.Id == 0 {
		return false, fmt.Errorf("kaas is missing kaas id")
	}
	if input.Region != "" {
		return false, fmt.Errorf("client cannot update region")
	}

	var obj = kaas.Kaas{
		Id:      input.Id,
		Project: input.Project,

		Name:              input.Name,
		Region:            input.Region,
		PackId:            input.PackId,
		KubernetesVersion: input.KubernetesVersion,
	}

	return true, updateCache(&obj)
}

func (c *Client) DeleteKaas(publicCloudId int, publicCloudProjectId int, kaasId int) (bool, error) {
	var obj = kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: publicCloudId,
			ProjectId:     publicCloudProjectId,
		},
		Id: kaasId,
	}

	return true, removeFromCache(&obj)
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

	obj.Status = "Active"

	return obj, nil
}

func (c *Client) CreateInstancePool(publicCloudId int, publicCloudProjectId int, input *kaas.InstancePool) (int, error) {
	// Checks
	if publicCloudId == 0 {
		return 0, fmt.Errorf("instance pool is missing public cloud id")
	}
	if publicCloudProjectId == 0 {
		return 0, fmt.Errorf("instance pool is missing public cloud project id")
	}
	if input.KaasId == 0 {
		return 0, fmt.Errorf("instance pool is missing kaas id")
	}
	if !dnsRegexp.MatchString(input.Name) {
		return 0, fmt.Errorf("instance pool name should be a dns name according to RFC 1123")
	}
	if input.FlavorName == "" {
		return 0, fmt.Errorf("instance pool is missing flavor name")
	}
	if input.MinInstances < 0 {
		return 0, fmt.Errorf("instance pool min instances should be greater than 0")
	}
	// if input.MaxInstances < 0 {
	// 	return nil, fmt.Errorf("instance pool max instances should be greater than 0")
	// }
	// if input.MinInstances > input.MaxInstances {
	// 	return nil, fmt.Errorf("instance pool min instance should be lesser than (or equal) max")
	// }

	_, err := c.GetKaas(publicCloudId, publicCloudProjectId, input.KaasId)
	if err != nil {
		return 0, err
	}

	var obj = kaas.InstancePool{
		Id:     genId(),
		KaasId: input.KaasId,

		Name:             input.Name,
		FlavorName:       input.FlavorName,
		AvailabilityZone: input.AvailabilityZone,
		MinInstances:     input.MinInstances,
		// MaxInstances: input.MaxInstances,
	}

	return obj.Id, addToCache(&obj)
}

func (c *Client) UpdateInstancePool(publicCloudId int, publicCloudProjectId int, input *kaas.InstancePool) (bool, error) {
	// Checks
	if publicCloudId == 0 {
		return false, fmt.Errorf("instance pool is missing public cloud id")
	}
	if publicCloudProjectId == 0 {
		return false, fmt.Errorf("instance pool is missing public cloud project id")
	}
	if input.KaasId == 0 {
		return false, fmt.Errorf("instance pool is missing kaas id")
	}
	if input.Id == 0 {
		return false, fmt.Errorf("instance pool is instance pool id")
	}
	if !dnsRegexp.MatchString(input.Name) {
		return false, fmt.Errorf("instance pool name should be a dns name according to RFC 1123")
	}
	if input.FlavorName == "" {
		return false, fmt.Errorf("instance pool is missing flavor name")
	}
	if input.MinInstances < 0 {
		return false, fmt.Errorf("instance pool min instances should be greater than 0")
	}
	// if input.MaxInstances < 0 {
	// 	return nil, fmt.Errorf("instance pool max instances should be greater than 0")
	// }
	// if input.MinInstances > input.MaxInstances {
	// 	return nil, fmt.Errorf("instance pool min instance should be lesser than (or equal) max")
	// }

	_, err := c.GetKaas(publicCloudId, publicCloudProjectId, input.KaasId)
	if err != nil {
		return false, err
	}

	_, err = c.GetInstancePool(publicCloudId, publicCloudProjectId, input.KaasId, input.Id)
	if err != nil {
		return false, err
	}

	var obj = kaas.InstancePool{
		KaasId: input.KaasId,
		Id:     input.Id,

		Name:             input.Name,
		FlavorName:       input.FlavorName,
		AvailabilityZone: input.AvailabilityZone,
		MinInstances:     input.MinInstances,
		// MaxInstances: input.MaxInstances,
	}

	return true, updateCache(&obj)
}

func (c *Client) DeleteInstancePool(publicCloudId int, publicCloudProjectId int, kaasId int, instancePoolId int) (bool, error) {
	_, err := c.GetKaas(publicCloudId, publicCloudProjectId, kaasId)
	if err != nil {
		return false, err
	}

	var obj = kaas.InstancePool{
		KaasId: kaasId,
		Id:     instancePoolId,
	}

	return true, removeFromCache(&obj)
}
