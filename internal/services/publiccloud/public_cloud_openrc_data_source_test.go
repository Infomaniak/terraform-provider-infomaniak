package publiccloud

import (
	"regexp"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPublicCloudOpenrcDataSource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"data_source.public_cloud_openrc.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_public_cloud_openrc_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.infomaniak_public_cloud_openrc.this", "content"),
					),
				},
			},
		},
		"data_source.public_cloud_openrc.invalid_region": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_public_cloud_openrc_invalid_region.tf"),
					ExpectError: regexp.MustCompile(`Attribute region value must be one of`),
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
