package mock

import (
	"fmt"
	"math/rand/v2"
	"terraform-provider-infomaniak/internal/apis/dbaas"
)

// Ensure that our client implements Api
var (
	_ dbaas.Api = (*Client)(nil)
)

type Client struct{}

func (c *Client) GetPacks() (map[string][]*dbaas.DBaaSPack, error) {
	return map[string][]*dbaas.DBaaSPack{
		"mysql": {
			{
				Id:   1,
				Name: "essential-1",
			},
			{
				Id:   2,
				Name: "essential-2",
			},
		},
	}, nil
}

// CreateDBaaS implements dbaas.Api.
func (c *Client) CreateDBaaS(input *dbaas.DBaaS) (*dbaas.DBaaSCreateInfo, error) {
	// Checks
	if input.Project.PublicCloudId == 0 {
		return nil, fmt.Errorf("dbaas is missing public cloud project id")
	}

	if input.Region == "" {
		return nil, fmt.Errorf("dbaas is missing region")
	}
	if input.PackId == 0 {
		return nil, fmt.Errorf("dbaas is missing pack id")
	}

	packs, _ := c.GetPacks()
	dbPacks, ok := packs[input.Type]

	if !ok {
		return nil, fmt.Errorf("dbType not found")
	}

	var pack *dbaas.DBaaSPack

	for _, dbPack := range dbPacks {
		if dbPack.Id == input.PackId {
			pack = dbPack
			break
		}
	}

	if pack == nil {
		return nil, fmt.Errorf("dbaas pack not found")
	}

	var obj = dbaas.DBaaS{
		Project: input.Project,
		Region:  input.Region,
		Type:    input.Type,
		Version: input.Version,
		PackId:  input.PackId,
		Pack:    pack,
		Name:    input.Name,
		Connection: &dbaas.DBaaSConnectionInfo{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "p@sSw0rd",
			Ca:       "This is totally a valid CA",
		},
		KubernetesIdentifier: "pck-zbeet07",
	}
	obj.Id = rand.Int64()

	var createInfo = dbaas.DBaaSCreateInfo{
		Id:             obj.Id,
		RootPassword:   obj.Connection.Password,
		KubeIdentifier: obj.KubernetesIdentifier,
	}
	return &createInfo, addToCache(&obj)
}

// CreateDBaasScheduleBackup implements dbaas.Api.
func (c *Client) CreateDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, backupSchedules *dbaas.DBaasBackupSchedule) (int64, error) {
	return 0, nil
}

// DeleteDBaaS implements dbaas.Api.
func (c *Client) DeleteDBaaS(publicCloudId int64, publicCloudProjectId int64, DBaaSId int64) (bool, error) {
	return true, nil
}

// DeleteDBaasScheduleBackup implements dbaas.Api.
func (c *Client) DeleteDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64) (bool, error) {
	return true, nil
}

// FindPack implements dbaas.Api.
func (c *Client) FindPack(dbType string, name string) (*dbaas.DBaaSPack, error) {
	packs, _ := c.GetPacks()
	dbPacks, ok := packs[dbType]

	if !ok {
		return nil, fmt.Errorf("dbType not found")
	}

	for _, pack := range dbPacks {
		if pack.Name == name {
			return pack, nil
		}
	}

	return nil, fmt.Errorf("pack not found")
}

// GetConfiguration implements dbaas.Api.
func (c *Client) GetConfiguration(publicCloudId int64, publicCloudProjectId int64, dbaasId int64) (map[string]any, error) {
	return map[string]any{
		"max_connections": 200,
	}, nil
}

// GetDBaaS implements dbaas.Api.
func (c *Client) GetDBaaS(publicCloudId int64, publicCloudProjectId int64, DBaaSId int64) (*dbaas.DBaaS, error) {

	key := fmt.Sprintf("%d-%d-%d", publicCloudId, publicCloudProjectId, DBaaSId)
	obj, err := getFromCache[*dbaas.DBaaS](key)
	if err != nil {
		return nil, err
	}

	obj.Status = "ready"

	return obj, nil
}

// GetDBaasScheduleBackup implements dbaas.Api.
func (c *Client) GetDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64) (*dbaas.DBaasBackupSchedule, error) {
	return nil, nil
}

// GetDbaasPack implements dbaas.Api.
func (c *Client) GetDbaasPack(params dbaas.PackFilter) (*dbaas.Pack, error) {
	return nil, nil
}

// GetDbaasRegions implements dbaas.Api.
func (c *Client) GetDbaasRegions() ([]string, error) {
	return []string{"dc4-a", "dc5-a"}, nil
}

// GetDbaasTypes implements dbaas.Api.
func (c *Client) GetDbaasTypes() ([]*dbaas.DbaasType, error) {
	return []*dbaas.DbaasType{
		{
			Name:     "mysql",
			Versions: []string{"8.0"},
		},
	}, nil
}

// GetIpFilters implements dbaas.Api.
func (c *Client) GetIpFilters(publicCloudId int64, publicCloudProjectId int64, dbaasId int64) ([]string, error) {
	return []string{"0.0.0.0/0"}, nil
}

// PatchIpFilters implements dbaas.Api.
func (c *Client) PatchIpFilters(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, filters dbaas.AllowedCIDRs) (bool, error) {
	return true, nil
}

// PutConfiguration implements dbaas.Api.
func (c *Client) PutConfiguration(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, configuration map[string]any) (bool, error) {
	return true, nil
}

// UpdateDBaaS implements dbaas.Api.
func (c *Client) UpdateDBaaS(input *dbaas.DBaaS) (bool, error) {
	return true, nil
}

// UpdateDBaasScheduleBackup implements dbaas.Api.
func (c *Client) UpdateDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64, backupSchedules *dbaas.DBaasBackupSchedule) (bool, error) {
	return true, nil
}

func New() *Client {
	return &Client{}
}
