package dbaas

import (
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDbaasDatasource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"data_source.dbaas.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_dbaas_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.infomaniak_dbaas.db", "type", "mysql"),
						resource.TestCheckResourceAttr("data.infomaniak_dbaas.db", "effective_configuration.max_connections", "200"),
						resource.TestCheckResourceAttrSet("data.infomaniak_dbaas.db", "ca"),
					),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, tc)
		})
	}
}
