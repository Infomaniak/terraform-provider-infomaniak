package provider

import (
	"os"
	mockKaas "terraform-provider-infomaniak/internal/apis/kaas/mock"
	"testing"
)

func TestMain(m *testing.M) {
	mockKaas.ResetCache()

	os.Exit(m.Run())
}
