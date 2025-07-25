package provider

import (
	"fmt"
	"terraform-provider-infomaniak/internal/apis"
)

func GetApiClient(providerData any) (*apis.Client, error) {
	data, ok := providerData.(*IkProviderData)
	if !ok {
		return nil, fmt.Errorf("expected *provider.IkProviderData, got: %T", providerData)
	}

	client := apis.NewClient(data.Data.Host.ValueString(), data.Data.Token.ValueString(), data.Version.ValueString())
	if data.Version.ValueString() == "test" {
		client = apis.NewMockClient()
	}

	return client, nil
}
