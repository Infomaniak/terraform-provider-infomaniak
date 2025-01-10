package provider

import (
	"os"
	mockKaas "terraform-provider-infomaniak/internal/apis/kaas/mock"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestMain(m *testing.M) {
	mockKaas.ResetCache()

	os.Exit(m.Run())
}

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"infomaniak": providerserver.NewProtocol6WithError(New("test")()),
	}
}
