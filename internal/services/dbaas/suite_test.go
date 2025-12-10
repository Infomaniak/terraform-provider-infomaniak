package dbaas

import (
	"os"
	mockDBaas "terraform-provider-infomaniak/internal/apis/dbaas/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	mockDBaas.ResetCache()
	Register()

	os.Exit(m.Run())
}

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "DBaaS Service Suite")
}
