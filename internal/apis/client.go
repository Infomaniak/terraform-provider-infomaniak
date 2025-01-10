package apis

import (
	"terraform-provider-infomaniak/internal/apis/kaas"
	mock_kaas "terraform-provider-infomaniak/internal/apis/kaas/mock"
)

type Client struct {
	Kaas kaas.Api
}

// NewMockClient defines the mock client for Infomaniak's API,
// It is used for testing or dryrunning
func NewMockClient() *Client {
	return &Client{
		Kaas: mock_kaas.New(),
	}
}
