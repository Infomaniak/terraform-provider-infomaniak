package publiccloud

import (
	"fmt"
	"regexp"
	"terraform-provider-infomaniak/internal/apis/publiccloud"
	mockPublicCloud "terraform-provider-infomaniak/internal/apis/publiccloud/mock"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPublicCloudProjectResource_Schema(t *testing.T) {
	mockPublicCloud.ResetCache()

	testCases := map[string]resource.TestCase{
		"resource.public_cloud_project.direct.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "resource_project_direct_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("infomaniak_public_cloud_project.this", "name", "dev-project"),
						resource.TestCheckResourceAttr("infomaniak_public_cloud_project.this", "public_cloud_id", "42"),
						resource.TestCheckResourceAttr("infomaniak_public_cloud_project.this", "status", "ok"),
						resource.TestCheckResourceAttrSet("infomaniak_public_cloud_project.this", "id"),
						resource.TestCheckResourceAttrSet("infomaniak_public_cloud_project.this", "open_stack_name"),
					),
				},
			},
		},
		"resource.public_cloud_project.invite.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "resource_project_invite_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("infomaniak_public_cloud_project.this", "name", "invited-project"),
						resource.TestCheckResourceAttr("infomaniak_public_cloud_project.this", "invite", "true"),
					),
				},
			},
		},
		"resource.public_cloud_project.missing_name": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_project_missing_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "name" is required`),
				},
			},
		},
		"resource.public_cloud_project.invite_without_email": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_project_invite_without_email.tf"),
					ExpectError: regexp.MustCompile(`Missing user_email`),
				},
			},
		},
		"resource.public_cloud_project.invite_with_password": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_project_invite_with_password.tf"),
					ExpectError: regexp.MustCompile(`Unexpected user_password`),
				},
			},
		},
		"resource.public_cloud_project.direct_without_password": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_project_direct_without_password.tf"),
					ExpectError: regexp.MustCompile(`Missing user_password`),
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

func TestPublicCloudProjectResource_Plan(t *testing.T) {
	mockPublicCloud.ResetCache()

	testCases := map[string]resource.TestCase{
		"resource.public_cloud_project.no_changes": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{Config: test.MustGetTestFile("plan", "resource_project_plan_initial.tf")},
				{
					Config: test.MustGetTestFile("plan", "resource_project_plan_initial.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					},
				},
			},
		},
		"resource.public_cloud_project.rename_is_in_place_update": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{Config: test.MustGetTestFile("plan", "resource_project_plan_initial.tf")},
				{
					Config: test.MustGetTestFile("plan", "resource_project_plan_renamed.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("infomaniak_public_cloud_project.this", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		},
		"resource.public_cloud_project.change_public_cloud_id_forces_replace": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{Config: test.MustGetTestFile("plan", "resource_project_plan_initial.tf")},
				{
					Config: test.MustGetTestFile("plan", "resource_project_plan_change_cloud.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("infomaniak_public_cloud_project.this", plancheck.ResourceActionDestroyBeforeCreate),
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		tc.IsUnitTest = true
		t.Run(name, func(t *testing.T) {
			mockPublicCloud.ResetCache()
			resource.Test(t, tc)
		})
	}
}

func TestPublicCloudProjectResource_Import(t *testing.T) {
	mockPublicCloud.ResetCache()
	client := mockPublicCloud.New()

	const (
		seedCloudId int64 = 42
		seedId      int64 = 5454
	)
	if err := client.SeedProject(&publiccloud.Project{
		Id:            seedId,
		PublicCloudId: seedCloudId,
		Name:          "imported-project",
		OpenStackName: "PCP-imported",
		Status:        publiccloud.StatusOk,
		ResourceLevel: 2,
		UserCount:     1,
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ResourceName:  "infomaniak_public_cloud_project.this",
				Config:        test.MustGetTestFile("schema", "resource_project_import.tf"),
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%d,%d", seedCloudId, seedId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infomaniak_public_cloud_project.this", "id", fmt.Sprint(seedId)),
					resource.TestCheckResourceAttr("infomaniak_public_cloud_project.this", "public_cloud_id", fmt.Sprint(seedCloudId)),
					resource.TestCheckResourceAttr("infomaniak_public_cloud_project.this", "name", "imported-project"),
					resource.TestCheckResourceAttr("infomaniak_public_cloud_project.this", "open_stack_name", "PCP-imported"),
				),
			},
		},
	})
}
