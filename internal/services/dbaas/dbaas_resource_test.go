package dbaas

import (
	"regexp"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDbaasResource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"resource.dbaas.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "resource_dbaas_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("infomaniak_dbaas.db", "type", "mysql"),
						resource.TestCheckResourceAttrSet("infomaniak_dbaas.db", "id"),
						resource.TestCheckResourceAttrSet("infomaniak_dbaas.db", "ca"),
					),
				},
			},
		},
		"resource.dbaas.missing_region": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_missing_region.tf"),
					ExpectError: regexp.MustCompile(`The argument "region" is required, but no definition was found.`),
				},
			},
		},
		"resource.dbaas.missing_type": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_missing_type.tf"),
					ExpectError: regexp.MustCompile(`The argument "type" is required, but no definition was found.`),
				},
			},
		},
		"resource.dbaas.missing_version": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_missing_version.tf"),
					ExpectError: regexp.MustCompile(`The argument "version" is required, but no definition was found.`),
				},
			},
		},
		"resource.dbaas.missing_pack_name": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_missing_pack_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "pack_name" is required, but no definition was found.`),
				},
			},
		},
		"resource.dbaas.missing_name": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_missing_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
				},
			},
		},
		"resource.dbaas.invalid_region": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_invalid_region.tf"),
					ExpectError: regexp.MustCompile(`The selected region is invalid.`),
				},
			},
		},
		"resource.dbaas.invalid_type": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_invalid_type.tf"),
					ExpectError: regexp.MustCompile(`The selected filter.type is invalid.`),
				},
			},
		},
		"resource.dbaas.invalid_version": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_invalid_version.tf"),
					ExpectError: regexp.MustCompile(`The selected version is invalid.`),
				},
			},
		},
		"resource.dbaas.invalid_pack": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_invalid_pack.tf"),
					ExpectError: regexp.MustCompile(`pack not found`),
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
