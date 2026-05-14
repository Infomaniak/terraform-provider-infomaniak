package publiccloud

import "strconv"

// Status values returned by Project / User endpoints.
const (
	StatusCreating  = "creating"
	StatusDeleting  = "deleting"
	StatusDisabled  = "disabled"
	StatusDisabling = "disabling"
	StatusError     = "error"
	StatusOk        = "ok"
	StatusUpdating  = "updating"
)

// Region values accepted by the openrc / authentication endpoints.
const (
	RegionPub1 = "pub1"
	RegionPub2 = "pub2"
)

// PublicCloud represents an Infomaniak Public Cloud product. Only fields
// exposed through Terraform are kept; deeply nested product-status objects
// (maintenance, locked) are omitted.
type PublicCloud struct {
	Id                     int64  `json:"public_cloud_id,omitempty"`
	AccountId              int64  `json:"account_id,omitempty"`
	ServiceId              int64  `json:"service_id,omitempty"`
	ServiceName            string `json:"service_name,omitempty"`
	CustomerName           string `json:"customer_name,omitempty"`
	InternalName           string `json:"internal_name,omitempty"`
	Description            string `json:"description,omitempty"`
	BillReference          string `json:"bill_reference,omitempty"`
	CreatedAt              int64  `json:"created_at,omitempty"`
	ExpiredAt              int64  `json:"expired_at,omitempty"`
	IsFree                 bool   `json:"is_free,omitempty"`
	IsZeroPrice            bool   `json:"is_zero_price,omitempty"`
	IsTrial                bool   `json:"is_trial,omitempty"`
	IsLocked               bool   `json:"is_locked,omitempty"`
	HasMaintenance         bool   `json:"has_maintenance,omitempty"`
	HasOperationInProgress bool   `json:"has_operation_in_progress,omitempty"`
}

// Key uniquely identifies a Public Cloud product in the mock cache.
func (p *PublicCloud) Key() string {
	return strconv.FormatInt(p.Id, 10)
}

// Project represents a Public Cloud project (OpenStack tenant).
type Project struct {
	Id             int64   `json:"public_cloud_project_id,omitempty"`
	PublicCloudId  int64   `json:"public_cloud_id,omitempty"`
	Name           string  `json:"name,omitempty"`
	OpenStackName  string  `json:"open_stack_name,omitempty"`
	Status         string  `json:"status,omitempty"`
	Price          float64 `json:"price,omitempty"`
	ResourceLevel  int64   `json:"resource_level,omitempty"`
	UserCount      int64   `json:"user_count,omitempty"`
	CreatedAt      int64   `json:"created_at,omitempty"`
	UpdatedAt      int64   `json:"updated_at,omitempty"`
	BillingStartAt int64   `json:"billing_start_at,omitempty"`
	BillingEndAt   int64   `json:"billing_end_at,omitempty"`
	PriceUpdatedAt int64   `json:"price_updated_at,omitempty"`
}

// Key uniquely identifies a Project in the mock cache.
func (p *Project) Key() string {
	return strconv.FormatInt(p.PublicCloudId, 10) + "-" + strconv.FormatInt(p.Id, 10)
}

// ProjectCreate is the body shape for POST /projects (and /projects/invite).
// Invite=true switches the endpoint and changes which user fields are required.
type ProjectCreate struct {
	Name            string `json:"project_name"`
	UserEmail       string `json:"user_email,omitempty"`
	UserPassword    string `json:"user_password,omitempty"`
	UserDescription string `json:"user_description,omitempty"`
	Invite          bool   `json:"-"`
}

// User represents an OpenStack user inside a Public Cloud project.
// Password and email are write-only on the API and never returned by GET.
type User struct {
	Id                   int64  `json:"public_cloud_user_id,omitempty"`
	PublicCloudId        int64  `json:"-"`
	PublicCloudProjectId int64  `json:"public_cloud_project_id,omitempty"`
	OpenStackName        string `json:"open_stack_name,omitempty"`
	Description          string `json:"description,omitempty"`
	Status               string `json:"status,omitempty"`
	CreatedAt            int64  `json:"created_at,omitempty"`
	UpdatedAt            int64  `json:"updated_at,omitempty"`
}

// Key uniquely identifies a User in the mock cache.
func (u *User) Key() string {
	return strconv.FormatInt(u.PublicCloudId, 10) + "-" +
		strconv.FormatInt(u.PublicCloudProjectId, 10) + "-" +
		strconv.FormatInt(u.Id, 10)
}

// UserCreate is the body shape for POST /users (and /users/invite). Invite=true
// switches the endpoint; password is sent in the non-invite case, email in the
// invite case.
type UserCreate struct {
	Description string `json:"description,omitempty"`
	Password    string `json:"password,omitempty"`
	Email       string `json:"email,omitempty"`
	Invite      bool   `json:"-"`
}

// UserUpdate carries the PATCH-able fields of a user.
type UserUpdate struct {
	Description string `json:"description,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
}

// Config is the per-account configuration returned by /1/public_clouds/config.
// valid_from / valid_to are nullable on the wire; they decode to 0 when null.
type Config struct {
	FreeTier             float64 `json:"free_tier"`
	FreeTierUsed         float64 `json:"free_tier_used"`
	AccountResourceLevel int64   `json:"account_resource_level"`
	ProjectCount         int64   `json:"project_count"`
	ValidFrom            int64   `json:"valid_from"`
	ValidTo              int64   `json:"valid_to"`
}

// Accesses is the shape of /1/public_clouds/accesses. Despite the endpoint
// name, the response carries maintenance status for the API surface, not a
// list of accessible regions.
type Accesses struct {
	IsMaintenanceOngoing bool `json:"is_maintenance_ongoing"`
}
