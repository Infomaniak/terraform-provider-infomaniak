package dbaas

import (
	"fmt"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider"
)

func GetApiClient(providerData any) (*apis.Client, error) {
	data, ok := providerData.(*provider.IkProviderData)
	if !ok {
		return nil, fmt.Errorf("expected *provider.IkProviderData, got: %T", providerData)
	}

	client := apis.NewClient(data.Data.Host.ValueString(), data.Data.Token.ValueString(), data.Version.ValueString())
	if data.Version.ValueString() == "test" {
		client = apis.NewMockClient()
	}

	return client, nil
}
