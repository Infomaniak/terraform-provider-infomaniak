package mock

import (
	"fmt"
	"terraform-provider-infomaniak/internal/apis/publiccloud"
	"time"
)

var _ publiccloud.Api = (*Client)(nil)

type Client struct{}

func New() *Client {
	return &Client{}
}

// fixedNow returns a deterministic timestamp so test assertions on
// created_at / updated_at fields stay stable across runs.
var fixedNow = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()

func (c *Client) ListPublicClouds(accountId int64) ([]*publiccloud.PublicCloud, error) {
	if accountId <= 0 {
		return nil, fmt.Errorf("account_id must be > 0")
	}
	return []*publiccloud.PublicCloud{
		c.syntheticPublicCloud(42),
		c.syntheticPublicCloud(43),
	}, nil
}

func (c *Client) GetPublicCloud(publicCloudId int64) (*publiccloud.PublicCloud, error) {
	if publicCloudId <= 0 {
		return nil, fmt.Errorf("public cloud id must be > 0")
	}
	key := (&publiccloud.PublicCloud{Id: publicCloudId}).Key()
	if cached, err := getFromCache[*publiccloud.PublicCloud](key); err == nil {
		return cached, nil
	}
	return c.syntheticPublicCloud(publicCloudId), nil
}

// UpdatePublicCloud applies PATCHable fields to a Public Cloud, persisting the
// result in the mock cache so subsequent reads return the new values.
func (c *Client) UpdatePublicCloud(input *publiccloud.PublicCloud) error {
	if input.Id <= 0 {
		return fmt.Errorf("public cloud id must be > 0")
	}

	current, err := getFromCache[*publiccloud.PublicCloud](input.Key())
	if err == ErrKeyNotFound {
		// Seed from the synthetic baseline when the cache is empty.
		current = c.syntheticPublicCloud(input.Id)
		if addErr := addToCache(current); addErr != nil {
			return addErr
		}
	} else if err != nil {
		return err
	}

	if input.CustomerName != "" {
		current.CustomerName = input.CustomerName
	}
	if input.Description != "" {
		current.Description = input.Description
	}
	if input.BillReference != "" {
		current.BillReference = input.BillReference
	}

	return updateCache(current)
}

// SeedPublicCloud lets tests insert a Public Cloud into the mock cache before
// running terraform import. Callers must call ResetCache() in TestMain to
// avoid leaking state between runs.
func (c *Client) SeedPublicCloud(pc *publiccloud.PublicCloud) error {
	return addToCache(pc)
}

func (c *Client) GetConfig(accountId int64) (*publiccloud.Config, error) {
	if accountId <= 0 {
		return nil, fmt.Errorf("account_id must be > 0")
	}
	return &publiccloud.Config{
		FreeTier:             100.0,
		FreeTierUsed:         12.5,
		AccountResourceLevel: 2,
		ProjectCount:         1,
		ValidFrom:            fixedNow,
		ValidTo:              fixedNow + int64((365 * 24 * time.Hour).Seconds()),
	}, nil
}

func (c *Client) GetAccesses(accountId int64) (*publiccloud.Accesses, error) {
	if accountId <= 0 {
		return nil, fmt.Errorf("account_id must be > 0")
	}
	return &publiccloud.Accesses{IsMaintenanceOngoing: false}, nil
}

func (c *Client) GetProject(publicCloudId, projectId int64) (*publiccloud.Project, error) {
	if publicCloudId <= 0 {
		return nil, fmt.Errorf("public cloud id must be > 0")
	}
	if projectId <= 0 {
		return nil, fmt.Errorf("project id must be > 0")
	}
	key := (&publiccloud.Project{PublicCloudId: publicCloudId, Id: projectId}).Key()
	if cached, err := getFromCache[*publiccloud.Project](key); err == nil {
		return cached, nil
	}
	return c.syntheticProject(publicCloudId, projectId), nil
}

func (c *Client) CreateProject(publicCloudId int64, input *publiccloud.ProjectCreate) (int64, error) {
	if publicCloudId <= 0 {
		return 0, fmt.Errorf("public cloud id must be > 0")
	}
	if input.Name == "" {
		return 0, fmt.Errorf("project name is required")
	}
	if input.Invite {
		if input.UserEmail == "" {
			return 0, fmt.Errorf("user_email is required when invite=true")
		}
	} else {
		if input.UserPassword == "" {
			return 0, fmt.Errorf("user_password is required when invite=false")
		}
	}

	project := &publiccloud.Project{
		Id:             genId(),
		PublicCloudId:  publicCloudId,
		Name:           input.Name,
		OpenStackName:  fmt.Sprintf("PCP-mock-%d", time.Now().UnixNano()%100000),
		Status:         publiccloud.StatusOk,
		ResourceLevel:  2,
		UserCount:      1,
		CreatedAt:      fixedNow,
		UpdatedAt:      fixedNow,
		BillingStartAt: fixedNow,
		BillingEndAt:   fixedNow + int64((30 * 24 * time.Hour).Seconds()),
		PriceUpdatedAt: fixedNow,
	}
	if err := addToCache(project); err != nil {
		return 0, err
	}
	return project.Id, nil
}

func (c *Client) UpdateProject(input *publiccloud.Project) error {
	if input.PublicCloudId <= 0 || input.Id <= 0 {
		return fmt.Errorf("public_cloud_id and project id must be > 0")
	}

	current, err := getFromCache[*publiccloud.Project](input.Key())
	if err == ErrKeyNotFound {
		current = c.syntheticProject(input.PublicCloudId, input.Id)
		if addErr := addToCache(current); addErr != nil {
			return addErr
		}
	} else if err != nil {
		return err
	}

	if input.Name != "" {
		current.Name = input.Name
	}
	current.UpdatedAt = fixedNow

	return updateCache(current)
}

func (c *Client) DeleteProject(publicCloudId, projectId int64) error {
	if publicCloudId <= 0 || projectId <= 0 {
		return fmt.Errorf("public_cloud_id and project id must be > 0")
	}
	obj := &publiccloud.Project{PublicCloudId: publicCloudId, Id: projectId}
	if err := removeFromCache(obj); err != nil && err != ErrKeyNotFound {
		return err
	}
	return nil
}

// SeedProject inserts a Project into the mock cache for tests.
func (c *Client) SeedProject(p *publiccloud.Project) error {
	return addToCache(p)
}

func (c *Client) syntheticProject(publicCloudId, projectId int64) *publiccloud.Project {
	return &publiccloud.Project{
		Id:             projectId,
		PublicCloudId:  publicCloudId,
		Name:           fmt.Sprintf("mock-project-%d", projectId),
		OpenStackName:  fmt.Sprintf("PCP-mock-%d", projectId),
		Status:         publiccloud.StatusOk,
		Price:          0.0,
		ResourceLevel:  2,
		UserCount:      1,
		CreatedAt:      fixedNow,
		UpdatedAt:      fixedNow,
		BillingStartAt: fixedNow,
		BillingEndAt:   fixedNow + int64((30 * 24 * time.Hour).Seconds()),
		PriceUpdatedAt: fixedNow,
	}
}

func (c *Client) GetUser(publicCloudId, projectId, userId int64) (*publiccloud.User, error) {
	if publicCloudId <= 0 || projectId <= 0 || userId <= 0 {
		return nil, fmt.Errorf("public_cloud_id, project_id and user_id must be > 0")
	}
	key := (&publiccloud.User{PublicCloudId: publicCloudId, PublicCloudProjectId: projectId, Id: userId}).Key()
	if cached, err := getFromCache[*publiccloud.User](key); err == nil {
		return cached, nil
	}
	return c.syntheticUser(publicCloudId, projectId, userId), nil
}

func (c *Client) CreateUser(publicCloudId, projectId int64, input *publiccloud.UserCreate) (int64, error) {
	if publicCloudId <= 0 || projectId <= 0 {
		return 0, fmt.Errorf("public_cloud_id and project_id must be > 0")
	}
	if input.Invite {
		if input.Email == "" {
			return 0, fmt.Errorf("email is required when invite=true")
		}
	} else {
		if input.Password == "" {
			return 0, fmt.Errorf("password is required when invite=false")
		}
	}

	user := &publiccloud.User{
		Id:                   genId(),
		PublicCloudId:        publicCloudId,
		PublicCloudProjectId: projectId,
		OpenStackName:        fmt.Sprintf("PCU-mock-%d", time.Now().UnixNano()%100000),
		Description:          input.Description,
		Status:               publiccloud.StatusOk,
		CreatedAt:            fixedNow,
		UpdatedAt:            fixedNow,
	}
	if err := addToCache(user); err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (c *Client) UpdateUser(publicCloudId, projectId, userId int64, input *publiccloud.UserUpdate) error {
	if publicCloudId <= 0 || projectId <= 0 || userId <= 0 {
		return fmt.Errorf("public_cloud_id, project_id and user_id must be > 0")
	}

	key := (&publiccloud.User{PublicCloudId: publicCloudId, PublicCloudProjectId: projectId, Id: userId}).Key()
	current, err := getFromCache[*publiccloud.User](key)
	if err == ErrKeyNotFound {
		current = c.syntheticUser(publicCloudId, projectId, userId)
		if addErr := addToCache(current); addErr != nil {
			return addErr
		}
	} else if err != nil {
		return err
	}

	if input.Description != "" {
		current.Description = input.Description
	}
	current.UpdatedAt = fixedNow

	return updateCache(current)
}

func (c *Client) DeleteUser(publicCloudId, projectId, userId int64) error {
	if publicCloudId <= 0 || projectId <= 0 || userId <= 0 {
		return fmt.Errorf("public_cloud_id, project_id and user_id must be > 0")
	}
	obj := &publiccloud.User{PublicCloudId: publicCloudId, PublicCloudProjectId: projectId, Id: userId}
	if err := removeFromCache(obj); err != nil && err != ErrKeyNotFound {
		return err
	}
	return nil
}

// SeedUser inserts a User into the mock cache for tests.
func (c *Client) SeedUser(u *publiccloud.User) error {
	return addToCache(u)
}

func (c *Client) syntheticUser(publicCloudId, projectId, userId int64) *publiccloud.User {
	return &publiccloud.User{
		Id:                   userId,
		PublicCloudId:        publicCloudId,
		PublicCloudProjectId: projectId,
		OpenStackName:        fmt.Sprintf("PCU-mock-%d", userId),
		Description:          "mock user",
		Status:               publiccloud.StatusOk,
		CreatedAt:            fixedNow,
		UpdatedAt:            fixedNow,
	}
}

func (c *Client) GetOpenrc(publicCloudId, projectId, userId int64, region string) (string, error) {
	if region == "" {
		region = publiccloud.RegionPub1
	}
	return fmt.Sprintf("#!/usr/bin/env bash\n# mock openrc for cloud=%d project=%d user=%d region=%s\nexport OS_USERNAME=mock\n",
		publicCloudId, projectId, userId, region), nil
}

func (c *Client) GetAuthentication(publicCloudId, projectId, userId int64, authType, region string) (string, error) {
	if authType == "" {
		return "", fmt.Errorf("authentication type is required")
	}
	if region == "" {
		region = publiccloud.RegionPub1
	}
	return fmt.Sprintf("# mock %s auth file for cloud=%d project=%d user=%d region=%s\n",
		authType, publicCloudId, projectId, userId, region), nil
}

func (c *Client) syntheticPublicCloud(id int64) *publiccloud.PublicCloud {
	return &publiccloud.PublicCloud{
		Id:                     id,
		AccountId:              1,
		ServiceId:              140,
		ServiceName:            "public_cloud",
		CustomerName:           fmt.Sprintf("mock-customer-%d", id),
		InternalName:           fmt.Sprintf("PC-mock-%d", id),
		Description:            "mock public cloud",
		BillReference:          "",
		CreatedAt:              fixedNow,
		ExpiredAt:              fixedNow + int64((365 * 24 * time.Hour).Seconds()),
		IsFree:                 false,
		IsZeroPrice:            false,
		IsTrial:                false,
		IsLocked:               false,
		HasMaintenance:         false,
		HasOperationInProgress: false,
	}
}
