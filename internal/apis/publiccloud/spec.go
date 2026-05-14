package publiccloud

// Api is the contract exposed by all Public Cloud client implementations
// (real HTTP client + mock). Methods are added per phase as the surface grows.
type Api interface {
	// PublicCloud product. Account-scoped reads (list, config, accesses)
	// require the caller's account id because the API filters by it.
	ListPublicClouds(accountId int64) ([]*PublicCloud, error)
	GetPublicCloud(publicCloudId int64) (*PublicCloud, error)
	UpdatePublicCloud(input *PublicCloud) (bool, error)
	GetConfig(accountId int64) (*Config, error)
	GetAccesses(accountId int64) (*Accesses, error)

	// Projects
	GetProject(publicCloudId, projectId int64) (*Project, error)
	CreateProject(publicCloudId int64, input *ProjectCreate) (int64, error)
	UpdateProject(input *Project) (bool, error)
	DeleteProject(publicCloudId, projectId int64) (bool, error)

	// Users
	GetUser(publicCloudId, projectId, userId int64) (*User, error)
	CreateUser(publicCloudId, projectId int64, input *UserCreate) (int64, error)
	UpdateUser(publicCloudId, projectId, userId int64, input *UserUpdate) (bool, error)
	DeleteUser(publicCloudId, projectId, userId int64) (bool, error)

	GetOpenrc(publicCloudId, projectId, userId int64, region string) (string, error)
	GetAuthentication(publicCloudId, projectId, userId int64, authType, region string) (string, error)
}
