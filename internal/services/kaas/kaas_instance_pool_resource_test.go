package kaas

import (
	"fmt"
	"regexp"
	"terraform-provider-infomaniak/internal/apis/kaas"
	mockKaas "terraform-provider-infomaniak/internal/apis/kaas/mock"
	"terraform-provider-infomaniak/internal/provider"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestKaasInstancePoolResource_Config(t *testing.T) {
	t.Parallel()

	testCases := map[string]resource.TestCase{
		"resource.kaas_instance_pool.good": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("schema", "resource_kaas_instance_pool_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("infomaniak_kaas.kluster", "id"),
						resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "name", "coucou"),
						resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "flavor_name", "test"),
						resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "min_instances", "3"),
						// resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "max_instances", "6"),
					),
				},
			},
		},
		"resource.kaas_instance_pool.missing_flavor_name": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_kaas_instance_pool_missing_flavor_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "flavor_name" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas_instance_pool.missing_kaas_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_kaas_instance_pool_missing_kaas_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "kaas_id" is required, but no definition was found.`),
				},
			},
		},
		// "resource.kaas_instance_pool.missing_max_instances": {
		// 	ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
		// 	Steps: []resource.TestStep{
		// 		{
		// 			Config:      test.MustGetTestFile("schema", "resource_kaas_instance_pool_missing_max_instances.tf"),
		// 			ExpectError: regexp.MustCompile(`The argument "max_instances" is required, but no definition was found.`),
		// 		},
		// 	},
		// },
		"resource.kaas_instance_pool.missing_min_instances": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_kaas_instance_pool_missing_min_instances.tf"),
					ExpectError: regexp.MustCompile(`The argument "min_instances" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas_instance_pool.missing_name": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_kaas_instance_pool_missing_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas_instance_pool.missing_public_cloud_project_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_kaas_instance_pool_missing_public_cloud_project_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "public_cloud_project_id" is required`),
				},
			},
		},
		"resource.kaas_instance_pool.cant_specify_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_kaas_instance_pool_cant_specify_id.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*id( )*=`),
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

func TestKaasInstancePoolResource_Plan(t *testing.T) {
	t.Parallel()

	testCases := map[string]resource.TestCase{
		"resource.kaas_instance_pool.no_changes": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("plan", "resource_kaas_instance_pool_test_no_changes.tf"),
				},
				{
					Config: test.MustGetTestFile("plan", "resource_kaas_instance_pool_test_no_changes.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		},
		"resource.kaas_instance_pool.change_kaas_id_causes_replace": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("plan", "resource_kaas_instance_pool_test_change_kaas_id_1.tf"),
				},
				{
					Config: test.MustGetTestFile("plan", "resource_kaas_instance_pool_test_change_kaas_id_2.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("infomaniak_kaas_instance_pool.instance_pool", plancheck.ResourceActionDestroyBeforeCreate),
						},
					},
				},
			},
		},
		"resource.kaas_instance_pool.change_public_cloud_project_id_causes_error": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("plan", "resource_kaas_instance_pool_test_change_public_cloud_project_id_1.tf"),
				},
				{
					Config:      test.MustGetTestFile("plan", "resource_kaas_instance_pool_test_change_public_cloud_project_id_2.tf"),
					ExpectError: regexp.MustCompile(`key not found`),
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

func TestKaasInstancePoolResource_Import(t *testing.T) {
	client := mockKaas.New()

	kaasId, err := client.CreateKaas(&kaas.Kaas{
		Project: kaas.KaasProject{
			PublicCloudId: 536,
			ProjectId:     451,
		},
		PackId: 1,
		Region: "das45",
	})
	if err != nil {
		t.Fatalf("Could not create Kaas for import test, got : %v", err)
	}

	kaasObject, err := client.GetKaas(536, 451, kaasId)
	if err != nil {
		t.Fatalf("Could not get Kaas for import test, got : %v", err)
	}

	defer func() {
		// Kaas project should be deleted after instance pool.
		_, err = client.DeleteKaas(kaasObject.Project.PublicCloudId, kaasObject.Project.ProjectId, kaasObject.Id)
		if err != nil {
			t.Fatalf("Could not delete Kaas in import test, got : %v", err)
		}
	}()

	instancePoolId, err := client.CreateInstancePool(kaasObject.Project.PublicCloudId, kaasObject.Project.ProjectId, &kaas.InstancePool{
		KaasId:       kaasObject.Id,
		Name:         "supername",
		FlavorName:   "superflavorname",
		MinInstances: 1,
		// MaxInstances: 99,
		Labels: map[string]string{"node-role.kubernetes.io/worker": "high"},
	})
	if err != nil {
		t.Fatalf("Could not create instance pool for import test, got : %v", err)
	}

	instancePool, err := client.GetInstancePool(kaasObject.Project.PublicCloudId, kaasObject.Project.ProjectId, kaasObject.Id, instancePoolId)
	if err != nil {
		t.Fatalf("Could not get instance pool for import test, got : %v", err)
	}

	defer func() {
		_, err = client.DeleteInstancePool(kaasObject.Project.PublicCloudId, kaasObject.Project.ProjectId, kaasObject.Id, instancePool.Id)
		if err != nil {
			t.Fatalf("Could not delete Kaas in import test, got : %v", err)
		}
	}()

	resourcePublicCloudId := kaasObject.Project.PublicCloudId
	resourcePublicCloudProjectId := kaasObject.Project.ProjectId
	resourceKaasId := instancePool.KaasId
	resourceId := instancePool.Id

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ResourceName:  "infomaniak_kaas_instance_pool.instance_pool",
				Config:        test.MustGetTestFile("plan", "resource_kaas_instance_pool_test_no_changes.tf"),
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%d,%d,%d,%d", resourcePublicCloudId, resourcePublicCloudProjectId, resourceKaasId, resourceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "id", fmt.Sprint(resourceId)),
					resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "public_cloud_id", fmt.Sprint(resourcePublicCloudId)),
					resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "public_cloud_project_id", fmt.Sprint(resourcePublicCloudProjectId)),
					resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "kaas_id", fmt.Sprint(resourceKaasId)),
				),
			},
		},
	})
}
