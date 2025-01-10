package provider

import (
	"fmt"
	"regexp"
	"terraform-provider-infomaniak/internal/apis/kaas"
	mockKaas "terraform-provider-infomaniak/internal/apis/kaas/mock"
	"terraform-provider-infomaniak/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestKaasInstancePoolResource_Config(t *testing.T) {
	t.Parallel()

	testCases := map[string]resource.TestCase{
		"resource.kaas_instance_pool.good": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "schema", "resource_kaas_instance_pool_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("infomaniak_kaas.kluster", "id"),
						resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "name", "coucou"),
						resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "flavor_name", "test"),
						resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "min_instances", "3"),
						resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "max_instances", "6"),
					),
				},
			},
		},
		"resource.kaas_instance_pool.missing_flavor_name": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_instance_pool_missing_flavor_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "flavor_name" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas_instance_pool.missing_kaas_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_instance_pool_missing_kaas_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "kaas_id" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas_instance_pool.missing_max_instances": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_instance_pool_missing_max_instances.tf"),
					ExpectError: regexp.MustCompile(`The argument "max_instances" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas_instance_pool.missing_min_instances": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_instance_pool_missing_min_instances.tf"),
					ExpectError: regexp.MustCompile(`The argument "min_instances" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas_instance_pool.missing_name": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_instance_pool_missing_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas_instance_pool.missing_pcp_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_instance_pool_missing_pcp_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "pcp_id" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas_instance_pool.cant_specify_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_instance_pool_cant_specify_id.tf"),
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
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_instance_pool_test_no_changes.tf"),
				},
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_instance_pool_test_no_changes.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		},
		"resource.kaas_instance_pool.change_kaas_id_causes_replace": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_instance_pool_test_change_kaas_id_1.tf"),
				},
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_instance_pool_test_change_kaas_id_2.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("infomaniak_kaas_instance_pool.instance_pool", plancheck.ResourceActionDestroyBeforeCreate),
						},
					},
				},
			},
		},
		"resource.kaas_instance_pool.change_pcp_id_causes_error": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_instance_pool_test_change_pcp_id_1.tf"),
				},
				{
					Config:      test.MustGetTestFile("kaas", "plan", "resource_kaas_instance_pool_test_change_pcp_id_2.tf"),
					ExpectError: regexp.MustCompile(`key '(.*)' not found`),
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
	k, err := client.CreateKaas(&kaas.Kaas{
		PcpId:  "451",
		Region: "das45",
	})
	if err != nil {
		t.Fatalf("Could not create Kaas for import test, got : %v", err)
	}

	defer func() {
		// Kaas project should be deleted after instance pool.
		err = client.DeleteKaas(k.PcpId, k.Id)
		if err != nil {
			t.Fatalf("Could not delete Kaas in import test, got : %v", err)
		}
	}()

	instancePool, err := client.CreateInstancePool(&kaas.InstancePool{
		PcpId:        k.PcpId,
		KaasId:       k.Id,
		Name:         "supername",
		FlavorName:   "superflavorname",
		MinInstances: 1,
		MaxInstances: 99,
	})
	if err != nil {
		t.Fatalf("Could not create Kaas for import test, got : %v", err)
	}

	defer func() {
		err = client.DeleteInstancePool(k.PcpId, k.Id, instancePool.Id)
		if err != nil {
			t.Fatalf("Could not delete Kaas in import test, got : %v", err)
		}
	}()

	resourcePcpId := instancePool.PcpId
	resourceKaasId := instancePool.KaasId
	resourceId := instancePool.Id

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ResourceName:  "infomaniak_kaas_instance_pool.instance_pool",
				Config:        test.MustGetTestFile("kaas", "plan", "resource_kaas_instance_pool_test_no_changes.tf"),
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s,%s,%s", resourcePcpId, resourceKaasId, resourceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "id", resourceId),
					resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "pcp_id", resourcePcpId),
					resource.TestCheckResourceAttr("infomaniak_kaas_instance_pool.instance_pool", "kaas_id", resourceKaasId),
				),
			},
		},
	})
}
