package provider

import (
	"fmt"
	"os"
	"terraform-provider-infomaniak/internal/apis"
)

func GetApiClient(providerData any) (*apis.Client, error) {
	data, ok := providerData.(*IkProviderData)
	if !ok {
		return nil, fmt.Errorf("expected *provider.IkProviderData, got: %T", providerData)
	}

	mocked := os.Getenv("MOCKED")
	if data.Version.ValueString() == "dev" && mocked == "true" {
		return apis.NewMockClient(), nil
	}

	client := apis.NewClient(data.Data.Host.ValueString(), data.Data.Token.ValueString(), data.Version.ValueString())

	return client, nil
}
