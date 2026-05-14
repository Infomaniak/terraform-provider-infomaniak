package publiccloud

import (
	"regexp"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPublicCloudProjectDataSource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"data_source.public_cloud_project.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_public_cloud_project_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.infomaniak_public_cloud_project.this", "id", "54"),
						resource.TestCheckResourceAttr("data.infomaniak_public_cloud_project.this", "public_cloud_id", "42"),
						resource.TestCheckResourceAttr("data.infomaniak_public_cloud_project.this", "status", "ok"),
					),
				},
			},
		},
		"data_source.public_cloud_project.missing_public_cloud_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_public_cloud_project_missing_public_cloud_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "public_cloud_id" is required`),
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
