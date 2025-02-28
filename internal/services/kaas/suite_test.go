package kaas

import (
	"os"
	mockKaas "terraform-provider-infomaniak/internal/apis/kaas/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	mockKaas.ResetCache()
	Register()

	os.Exit(m.Run())
}

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "KaaS Service Suite")
}
