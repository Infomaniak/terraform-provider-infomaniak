package publiccloud

import (
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPublicCloudConfigDataSource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"data_source.public_cloud_config.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_public_cloud_config_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.infomaniak_public_cloud_config.this", "free_tier"),
						resource.TestCheckResourceAttrSet("data.infomaniak_public_cloud_config.this", "account_resource_level"),
					),
				},
			},
		},
	}

	for name, tc := range testCases {
		tc.IsUnitTest = true
		t.Run(name, func(t *testing.T) {
			resource.Test(t, tc)
		})
	}
}
