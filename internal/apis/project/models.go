package project

type CreateProject struct {
	PublicCloudId   int    `json:"public_cloud_id,omitempty"`
	Name            string `json:"project_name,omitempty"`
	UserDescription string `json:"user_description"`
	UserEmail       string `json:"user_email"`
	UserPassword    string `json:"user_password,omitempty"`
}

type UpdateProject struct {
	Name string `json:"name,omitempty"`
}

type Project struct {
	ProjectId      int          `json:"public_cloud_project_id,omitempty"`
	PublicCloudId  int          `json:"public_cloud_id,omitempty"`
	OpenStackName  string       `json:"open_stack_name,omitempty"`
	Name           string       `json:"name,omitempty"`
	Price          float32      `json:"price,omitempty"`
	ResourceLevel  int          `json:"resource_level,omitempty"`
	Status         string       `json:"status,omitempty"`
	PriceUpdatedAt int          `json:"price_updated_at,omitempty"`
	BillingStartAt int          `json:"billing_start_at,omitempty"`
	BillingEndAt   int          `json:"billing_end_at,omitempty"`
	CreatedAt      int          `json:"created_at,omitempty"`
	UpdatedAt      int          `json:"updated_at,omitempty"`
	DeletedAt      int          `json:"deleted_at"`
	UserCount      int          `json:"user_count,omitempty"`
	Tags           []ProjectTag `json:tags,omitempty"`
}

type ProjectTag struct {
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Color        int    `json:"color,omitempty"`
	ProductCount int    `json:"product_count"`
}
