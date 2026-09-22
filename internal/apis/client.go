package apis

import (
	"terraform-provider-infomaniak/internal/apis/dbaas"
	implem_dbaas "terraform-provider-infomaniak/internal/apis/dbaas/implementation"
	mock_dbaas "terraform-provider-infomaniak/internal/apis/dbaas/mock"
	"terraform-provider-infomaniak/internal/apis/domain"
	"terraform-provider-infomaniak/internal/apis/kaas"

	implem_kaas "terraform-provider-infomaniak/internal/apis/kaas/implementation"
	mock_kaas "terraform-provider-infomaniak/internal/apis/kaas/mock"

	implem_domain "terraform-provider-infomaniak/internal/apis/domain/implementation"

	"terraform-provider-infomaniak/internal/apis/publiccloud"
	implem_publiccloud "terraform-provider-infomaniak/internal/apis/publiccloud/implementation"
	mock_publiccloud "terraform-provider-infomaniak/internal/apis/publiccloud/mock"
)

type Client struct {
	Kaas        kaas.Api
	Domain      domain.Api
	DBaas       dbaas.Api
	PublicCloud publiccloud.Api
}

// NewMockClient defines the mock client for Infomaniak's API,
// It is used for testing or dryrunning
func NewMockClient() *Client {
	return &Client{
		Kaas:        mock_kaas.New(),
		DBaas:       mock_dbaas.New(),
		PublicCloud: mock_publiccloud.New(),
	}
}

// NewClient defines the client for Infomaniak's API
func NewClient(baseUri, token, version string) *Client {
	return &Client{
		Kaas:        implem_kaas.New(baseUri, token, version),
		DBaas:       implem_dbaas.New(baseUri, token, version),
		Domain:      implem_domain.New(baseUri, token, version),
		PublicCloud: implem_publiccloud.New(baseUri, token, version),
	}
}
