package publiccloud

import (
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPublicCloudsDataSource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"data_source.public_clouds.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_public_clouds_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.infomaniak_public_clouds.all", "account_id", "1"),
						resource.TestCheckResourceAttr("data.infomaniak_public_clouds.all", "public_clouds.#", "2"),
						resource.TestCheckResourceAttr("data.infomaniak_public_clouds.all", "public_clouds.0.id", "42"),
						resource.TestCheckResourceAttr("data.infomaniak_public_clouds.all", "public_clouds.1.id", "43"),
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
