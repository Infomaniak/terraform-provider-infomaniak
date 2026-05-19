package publiccloud

import (
	"regexp"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPublicCloudUserDataSource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"data_source.public_cloud_user.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_public_cloud_user_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.infomaniak_public_cloud_user.this", "id", "7"),
						resource.TestCheckResourceAttr("data.infomaniak_public_cloud_user.this", "status", "ok"),
						resource.TestCheckResourceAttrSet("data.infomaniak_public_cloud_user.this", "open_stack_name"),
					),
				},
			},
		},
		"data_source.public_cloud_user.missing_project_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_public_cloud_user_missing_project_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "public_cloud_project_id" is required`),
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
