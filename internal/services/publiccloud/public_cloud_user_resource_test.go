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

func TestPublicCloudUserResource_Schema(t *testing.T) {
	mockPublicCloud.ResetCache()

	testCases := map[string]resource.TestCase{
		"resource.public_cloud_user.direct.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "resource_user_direct_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "public_cloud_id", "42"),
						resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "public_cloud_project_id", "54"),
						resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "description", "ops user"),
						resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "status", "ok"),
						resource.TestCheckResourceAttrSet("infomaniak_public_cloud_user.this", "id"),
						resource.TestCheckResourceAttrSet("infomaniak_public_cloud_user.this", "open_stack_name"),
					),
				},
			},
		},
		"resource.public_cloud_user.invite.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "resource_user_invite_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "invite", "true"),
					),
				},
			},
		},
		"resource.public_cloud_user.missing_project_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_user_missing_project_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "public_cloud_project_id" is required`),
				},
			},
		},
		"resource.public_cloud_user.invite_without_email": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_user_invite_without_email.tf"),
					ExpectError: regexp.MustCompile(`Missing email`),
				},
			},
		},
		"resource.public_cloud_user.direct_without_password": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_user_direct_without_password.tf"),
					ExpectError: regexp.MustCompile(`Missing password`),
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

func TestPublicCloudUserResource_Plan(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"resource.public_cloud_user.no_changes": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{Config: test.MustGetTestFile("plan", "resource_user_plan_initial.tf")},
				{
					Config: test.MustGetTestFile("plan", "resource_user_plan_initial.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					},
				},
			},
		},
		"resource.public_cloud_user.description_change_is_in_place_update": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{Config: test.MustGetTestFile("plan", "resource_user_plan_initial.tf")},
				{
					Config: test.MustGetTestFile("plan", "resource_user_plan_description_changed.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("infomaniak_public_cloud_user.this", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		},
		"resource.public_cloud_user.invite_change_forces_replace": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{Config: test.MustGetTestFile("plan", "resource_user_plan_initial.tf")},
				{
					Config: test.MustGetTestFile("plan", "resource_user_plan_change_invite.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("infomaniak_public_cloud_user.this", plancheck.ResourceActionDestroyBeforeCreate),
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

func TestPublicCloudUserResource_Import(t *testing.T) {
	mockPublicCloud.ResetCache()
	client := mockPublicCloud.New()

	const (
		seedCloudId   int64 = 42
		seedProjectId int64 = 54
		seedUserId    int64 = 7777
	)
	if err := client.SeedUser(&publiccloud.User{
		Id:                   seedUserId,
		PublicCloudId:        seedCloudId,
		PublicCloudProjectId: seedProjectId,
		OpenStackName:        "PCU-imported",
		Description:          "imported user",
		Status:               publiccloud.StatusOk,
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ResourceName:  "infomaniak_public_cloud_user.this",
				Config:        test.MustGetTestFile("schema", "resource_user_import.tf"),
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%d,%d,%d", seedCloudId, seedProjectId, seedUserId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "id", fmt.Sprint(seedUserId)),
					resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "public_cloud_id", fmt.Sprint(seedCloudId)),
					resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "public_cloud_project_id", fmt.Sprint(seedProjectId)),
					resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "open_stack_name", "PCU-imported"),
					resource.TestCheckResourceAttr("infomaniak_public_cloud_user.this", "description", "imported user"),
				),
			},
		},
	})
}
