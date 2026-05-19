package publiccloud

import (
	"regexp"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPublicCloudDataSource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"data_source.public_cloud.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_public_cloud_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.infomaniak_public_cloud.this", "id", "42"),
						resource.TestCheckResourceAttr("data.infomaniak_public_cloud.this", "service_name", "public_cloud"),
						resource.TestCheckResourceAttrSet("data.infomaniak_public_cloud.this", "customer_name"),
					),
				},
			},
		},
		"data_source.public_cloud.missing_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_public_cloud_missing_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "id" is required`),
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
