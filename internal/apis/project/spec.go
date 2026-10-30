package project

type Api interface {
	GetProject(publicCloudId int, publicCloudProjectId int) (*Project, error)
	CreateProject(input *CreateProject) (int, error)
	UpdateProject(publicCloudId int, publicCloudProjectId int, input *UpdateProject) (bool, error)
	DeleteProject(publicCloudId int, publicCloudProjectId int) (bool, error)
}
