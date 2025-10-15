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

	if data.Version.ValueString() == "dev" {
		return apis.NewMockClient(), nil
	}

	client := apis.NewClient(data.Data.Host.ValueString(), data.Data.Token.ValueString(), data.Version.ValueString())

	return client, nil
}
