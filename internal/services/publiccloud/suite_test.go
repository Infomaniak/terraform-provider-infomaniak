package publiccloud

import (
	"os"
	mockPublicCloud "terraform-provider-infomaniak/internal/apis/publiccloud/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	mockPublicCloud.ResetCache()
	Register()
	os.Exit(m.Run())
}

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Public Cloud Service Suite")
}
