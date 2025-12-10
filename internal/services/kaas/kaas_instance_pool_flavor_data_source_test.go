package kaas

import (
	"regexp"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestKaasInstancePoolFlavorDatasourceSchema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"data_source.kaas_instance_pool_flavor.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_kaas_instance_pool_flavor_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.infomaniak_kaas_instance_pool_flavor.my_flavor", "name"),
					),
				},
			},
		},
		"data_source.kaas_instance_pool.only_name": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_kaas_instance_pool_flavor_only_name.tf"),
					Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("data.infomaniak_kaas_instance_pool_flavor.my_flavor", "cpu")),
				},
			},
		},
		"data_source.kaas_instance_pool.empty": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_kaas_instance_pool_flavor_empty.tf"),
					ExpectError: regexp.MustCompile(`either 'name' or at least one of 'cpu', 'ram', 'storage' must be specified`),
				},
			},
		},
		"data_source.kaas_instance_pool.not_found": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_kaas_instance_pool_flavor_not_found.tf"),
					ExpectError: regexp.MustCompile(`flavor not found`),
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
