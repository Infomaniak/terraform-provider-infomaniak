package dbaas

import (
	"regexp"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDbaasBacketScheduleResource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"resource.dbaas_backup_schedule.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "resource_dbaas_backup_schedule_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("infomaniak_dbaas_backup_schedule.my_schedule", "name"),
					),
				},
			},
		},
		"resource.dbaas_backup_schedule.missing_scheduled_at": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_backup_schedule_missing_scheduled_at.tf"),
					ExpectError: regexp.MustCompile(`The argument "scheduled_at" is required, but no definition was found.`),
				},
			},
		},
		"resource.dbaas_backup_schedule.invalid_scheduled_at": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_dbaas_backup_schedule_invalid_scheduled_at.tf"),
					ExpectError: regexp.MustCompile(`The scheduled at does not match the format H:i.`),
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
