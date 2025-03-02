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

func TestKaasResource_Schema(t *testing.T) {
	t.Parallel()

	testCases := map[string]resource.TestCase{
		"resource.kaas.good": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "schema", "resource_kaas_good.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("infomaniak_kaas.kluster", "region", "dc5"),
						resource.TestCheckResourceAttrSet("infomaniak_kaas.kluster", "id"),
						resource.TestCheckResourceAttrSet("infomaniak_kaas.kluster", "kubeconfig"),
					),
				},
			},
		},
		"resource.kaas.missing_region": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_missing_region.tf"),
					ExpectError: regexp.MustCompile(`The argument "region" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas.missing_pcp_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_missing_pcp_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "pcp_id" is required, but no definition was found.`),
				},
			},
		},
		"resource.kaas.cant_specify_id": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_cant_specify_id.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*id( )*=`),
				},
			},
		},
		"resource.kaas.cant_specify_kubeconfig": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("kaas", "schema", "resource_kaas_cant_specify_kubeconfig.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*kubeconfig( )*=`),
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

func TestKaasResource_Plan(t *testing.T) {
	t.Parallel()

	testCases := map[string]resource.TestCase{
		"resource.kaas.no_changes": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_test_no_changes.tf"),
				},
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_test_no_changes.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		},
		"resource.kaas.change_pcp_id_causes_replace": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_test_change_pcp_id_1.tf"),
				},
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_test_change_pcp_id_2.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("infomaniak_kaas.kluster", plancheck.ResourceActionDestroyBeforeCreate),
						},
					},
				},
			},
		},
		"resource.kaas.change_region_causes_replace": {
			ProtoV6ProviderFactories: protoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_test_change_region_1.tf"),
				},
				{
					Config: test.MustGetTestFile("kaas", "plan", "resource_kaas_test_change_region_2.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("infomaniak_kaas.kluster", plancheck.ResourceActionDestroyBeforeCreate),
						},
					},
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

func TestKaasResource_Import(t *testing.T) {
	var resourcePcpId, resourceId string

	client := mockKaas.New()
	k, err := client.CreateKaas(&kaas.Kaas{
		PcpId:  "451",
		Region: "das45",
	})
	if err != nil {
		t.Fatalf("Could not create Kaas for import test, got : %v", err)
	}
	defer func() {
		err = client.DeleteKaas(k.PcpId, k.Id)
		if err != nil {
			t.Fatalf("Could not delete Kaas in import test, got : %v", err)
		}
	}()

	resourceId = k.Id
	resourcePcpId = k.PcpId

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ResourceName:  "infomaniak_kaas.kluster",
				Config:        test.MustGetTestFile("kaas", "plan", "resource_kaas_test_no_changes.tf"),
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s,%s", resourcePcpId, resourceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infomaniak_kaas.kluster", "id", resourceId),
					resource.TestCheckResourceAttr("infomaniak_kaas.kluster", "pcp_id", resourcePcpId),
					resource.TestCheckResourceAttr("infomaniak_kaas.kluster", "region", "das45"),
					resource.TestCheckResourceAttr("infomaniak_kaas.kluster", "kubeconfig", k.Kubeconfig),
				),
			},
		},
	})
}
