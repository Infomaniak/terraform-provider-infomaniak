package provider

import (
	"regexp"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestKaasDatasource_Schema(t *testing.T) {
	t.Parallel()

	testCases := map[string]resource.TestCase{
		"data_source.kaas.good": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "schema", "data_source_kaas_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.infomaniak_kaas.kluster", "region", "dc5"),
						resource.TestCheckResourceAttrSet("data.infomaniak_kaas.kluster", "kubeconfig"),
					),
				},
			},
		},
		"data_source.kaas.missing_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_missing_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		},
		"data_source.kaas.missing_pcp_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_missing_pcp_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "pcp_id" is required, but no definition was found.`),
				},
			},
		},
		"data_source.kaas.cant_specify_kubeconfig": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_cant_specify_kubeconfig.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*kubeconfig( )*=`),
				},
			},
		},
		"data_source.kaas.cant_specify_region": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_cant_specify_region.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*region( )*=`),
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
