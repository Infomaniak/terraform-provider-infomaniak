package kaas

import (
	"regexp"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestKaasDatasource_Schema(t *testing.T) {
	t.Parallel()

	testCases := map[string]resource.TestCase{
		"data_source.kaas.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "data_source_kaas_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.infomaniak_kaas.kluster", "region", "dc5"),
						resource.TestCheckResourceAttrSet("data.infomaniak_kaas.kluster", "kubeconfig"),
					),
				},
			},
		},
		"data_source.kaas.missing_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_kaas_missing_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		},
		"data_source.kaas.missing_public_cloud_project_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_kaas_missing_public_cloud_project_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "public_cloud_project_id" is required`),
				},
			},
		},
		"data_source.kaas.cant_specify_kubeconfig": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_kaas_cant_specify_kubeconfig.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*kubeconfig( )*=`),
				},
			},
		},
		"data_source.kaas.cant_specify_region": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "data_source_kaas_cant_specify_region.tf"),
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
