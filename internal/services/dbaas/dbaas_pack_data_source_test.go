package dbaas

import (
	"regexp"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDbaasPackDatasource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"data_source.dbaas.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_dbaas_pack_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.infomaniak_dbaas_pack.pack", "type", "mysql"),
						resource.TestCheckResourceAttr("data.infomaniak_dbaas_pack.pack", "ram", "4"),
					),
				},
			},
		},
		"data_source.dbaas.search": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_dbaas_pack_search.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.infomaniak_dbaas_pack.pack", "type", "mysql"),
						resource.TestCheckResourceAttr("data.infomaniak_dbaas_pack.pack", "name", "essential-db-8"),
					),
				},
			},
		},
		"data_source.dbaas.too_many_results": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_dbaas_pack_too_many_results.tf"),
					ExpectError: regexp.MustCompile(`multiple packs found, please refine your search`),
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
