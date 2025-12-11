package mock

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"time"

	"github.com/samber/lo"
)

// Ensure that our client implements Api
var (
	_ dbaas.Api = (*Client)(nil)
)

type Client struct{}

func (c *Client) getPacks() map[string][]*dbaas.Pack {
	return map[string][]*dbaas.Pack{
		"mysql": {
			{
				ID:        1,
				Type:      "mysql",
				Name:      "essential-db-4",
				Group:     "essential-db",
				Instances: 1,
				CPU:       1,
				RAM:       4,
				Storage:   80,
			},
			{
				ID:        2,
				Type:      "mysql",
				Name:      "essential-db-8",
				Group:     "essential-db",
				Instances: 1,
				CPU:       2,
				RAM:       8,
				Storage:   160,
			},
		},
	}
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

	// check if the type is valid
	dbTypes, _ := c.GetDbaasTypes()
	dbType, ok := lo.Find(dbTypes, func(dbType *dbaas.DbaasType) bool {
		return dbType.Name == input.Type
	})

	if !ok {
		return nil, fmt.Errorf("The selected filter.type is invalid.")
	}

	// check if the version is valid
	if !slices.Contains(dbType.Versions, input.Version) {
		return nil, fmt.Errorf("The selected version is invalid.")
	}

	// check if the region is valid
	regions, _ := c.GetDbaasRegions()
	if !slices.Contains(regions, input.Region) {
		return nil, fmt.Errorf("The selected region is invalid.")
	}

	// get the requested pack
	pack, ok := lo.Find(c.getPacks()[input.Type], func(item *dbaas.Pack) bool {
		return item.ID == input.PackId
	})

	if !ok {
		return nil, fmt.Errorf("pack not found")
	}

	var obj = dbaas.DBaaS{
		Project: input.Project,
		Region:  input.Region,
		Type:    input.Type,
		Version: input.Version,
		PackId:  input.PackId,
		Pack: &dbaas.DBaaSPack{
			Id:   pack.ID,
			Name: pack.Name,
		},
		Name: input.Name,
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
	if backupSchedules.ScheduledAt == nil {
		return 0, fmt.Errorf("dbaas backup schedule is missing schedule_at")
	}

	if backupSchedules.Retention == nil {
		return 0, fmt.Errorf("dbaas backup schedule is retention schedule_at")
	}

	_, err := time.Parse("15:04", *backupSchedules.ScheduledAt)
	if err != nil {
		return 0, fmt.Errorf("The scheduled at does not match the format H:i.")
	}

	id := rand.Int64()
	backupSchedules.Id = &id

	name := fmt.Sprintf("%d", rand.Float64())
	backupSchedules.Name = &name

	return id, addToCache(backupSchedules)
}

// DeleteDBaaS implements dbaas.Api.
func (c *Client) DeleteDBaaS(publicCloudId int64, publicCloudProjectId int64, DBaaSId int64) (bool, error) {
	var obj = dbaas.DBaaS{
		Project: dbaas.DBaaSProject{
			PublicCloudId: publicCloudId,
			ProjectId:     publicCloudProjectId,
		},
		Id: DBaaSId,
	}

	return true, removeFromCache(&obj)
}

// DeleteDBaasScheduleBackup implements dbaas.Api.
func (c *Client) DeleteDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64) (bool, error) {
	var obj = dbaas.DBaasBackupSchedule{
		Id: &id,
	}

	return true, removeFromCache(&obj)
}

// FindPack implements dbaas.Api.
func (c *Client) FindPack(dbType string, name string) (*dbaas.DBaaSPack, error) {
	packs, ok := c.getPacks()[dbType]
	if !ok {
		return nil, fmt.Errorf("The selected filter.type is invalid.")
	}

	dbPack, ok := lo.Find(packs, func(item *dbaas.Pack) bool {
		return item.Name == name
	})

	if !ok {
		return nil, fmt.Errorf("pack not found")
	}

	return &dbaas.DBaaSPack{
		Id:   dbPack.ID,
		Name: dbPack.Name,
	}, nil
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
	key := fmt.Sprintf("%d", id)
	obj, err := getFromCache[*dbaas.DBaasBackupSchedule](key)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// GetDbaasPack implements dbaas.Api.
func (c *Client) GetDbaasPack(params dbaas.PackFilter) (*dbaas.Pack, error) {
	packs, ok := c.getPacks()[params.DbType]
	if !ok {
		return nil, fmt.Errorf("The selected filter.type is invalid.")
	}

	filteredPacks := lo.Filter(packs, func(pack *dbaas.Pack, _ int) bool {
		found := true

		if params.Group != nil && pack.Group != *params.Group {
			found = false
		}

		if params.Name != nil && pack.Name != *params.Name {
			found = false
		}

		if params.Instances != nil && pack.Instances != *params.Instances {
			found = false
		}

		if params.Cpu != nil && pack.CPU != *params.Cpu {
			found = false
		}

		if params.Ram != nil && pack.RAM != *params.Ram {
			found = false
		}

		if params.Storage != nil && pack.Storage != *params.Storage {
			found = false
		}

		return found
	})

	if len(filteredPacks) == 0 {
		return nil, fmt.Errorf("pack not found")
	}

	if len(filteredPacks) > 1 {
		return nil, fmt.Errorf("multiple packs found, please refine your search")
	}

	return filteredPacks[0], nil
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
	// Checks
	if input.Project.PublicCloudId == 0 {
		return false, fmt.Errorf("dbaas is missing public cloud project id")
	}
	if input.Id == 0 {
		return false, fmt.Errorf("dbaas is missing kaas id")
	}
	if input.PackId == 0 {
		return false, fmt.Errorf("dbaas is missing pack id")
	}
	if input.Region != "" {
		return false, fmt.Errorf("client cannot update region")
	}

	var obj = dbaas.DBaaS{
		Id:      input.Id,
		Project: input.Project,

		Name:   input.Name,
		PackId: input.PackId,
	}

	return true, updateCache(&obj)
}

// UpdateDBaasScheduleBackup implements dbaas.Api.
func (c *Client) UpdateDBaasScheduleBackup(publicCloudId int64, publicCloudProjectId int64, dbaasId int64, id int64, backupSchedules *dbaas.DBaasBackupSchedule) (bool, error) {
	// Checks
	if publicCloudProjectId == 0 {
		return false, fmt.Errorf("dbaas backup schedule is missing public cloud project id")
	}
	if id == 0 {
		return false, fmt.Errorf("dbaas backup schedule is missing kaas id")
	}
	if backupSchedules.ScheduledAt == nil {
		return false, fmt.Errorf("dbaas backup schedule is missing scheduled_at")
	}

	if backupSchedules.Retention == nil {
		return false, fmt.Errorf("dbaas backup schedule is retention schedule_at")
	}

	return true, updateCache(backupSchedules)
}

func New() *Client {
	return &Client{}
}
