package provider

import (
	"regexp"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestKaasInstancePoolDatasource_Schema(t *testing.T) {
	t.Parallel()

	testCases := map[string]resource.TestCase{
		"data_source.kaas_instance_pool.good": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "schema", "data_source_kaas_instance_pool_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.infomaniak_kaas_instance_pool.instance_pool", "id"),
						resource.TestCheckResourceAttr("data.infomaniak_kaas_instance_pool.instance_pool", "name", "coucou"),
						resource.TestCheckResourceAttr("data.infomaniak_kaas_instance_pool.instance_pool", "flavor_name", "test"),
						resource.TestCheckResourceAttr("data.infomaniak_kaas_instance_pool.instance_pool", "min_instances", "3"),
						resource.TestCheckResourceAttr("data.infomaniak_kaas_instance_pool.instance_pool", "max_instances", "6"),
					),
				},
			},
		},
		"data_source.kaas_instance_pool.missing_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_instance_pool_missing_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		},
		"data_source.kaas_instance_pool.missing_pcp_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_instance_pool_missing_pcp_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "pcp_id" is required, but no definition was found.`),
				},
			},
		},
		"data_source.kaas_instance_pool.missing_kaas_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_instance_pool_missing_kaas_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "kaas_id" is required, but no definition was found.`),
				},
			},
		},
		"data_source.kaas_instance_pool.cant_specify_flavor_name": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_instance_pool_cant_specify_flavor_name.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*flavor_name( )*=`),
				},
			},
		},
		"data_source.kaas_instance_pool.cant_specify_max_instances": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_instance_pool_cant_specify_max_instances.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*max_instances( )*=`),
				},
			},
		},
		"data_source.kaas_instance_pool.cant_specify_min_instances": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_instance_pool_cant_specify_min_instances.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*min_instances( )*=`),
				},
			},
		},
		"data_source.kaas_instance_pool.cant_specify_name": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "data_source_kaas_instance_pool_cant_specify_name.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*name( )*=`),
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
