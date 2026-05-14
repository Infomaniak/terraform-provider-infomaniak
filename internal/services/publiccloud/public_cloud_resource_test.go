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
)

func TestPublicCloudResource_Schema(t *testing.T) {
	testCases := map[string]resource.TestCase{
		"resource.public_cloud.create_is_rejected": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_public_cloud_create_rejected.tf"),
					ExpectError: regexp.MustCompile(`Public Cloud cannot be created via Terraform`),
				},
			},
		},
		"resource.public_cloud.cant_specify_id": {
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.MustGetTestFile("schema", "resource_public_cloud_cant_specify_id.tf"),
					ExpectError: regexp.MustCompile(`[0-9]+:( )*id( )*=`),
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

func TestPublicCloudResource_Import(t *testing.T) {
	mockPublicCloud.ResetCache()
	client := mockPublicCloud.New()

	const seedId int64 = 4242
	if err := client.SeedPublicCloud(&publiccloud.PublicCloud{
		Id:           seedId,
		AccountId:    1,
		ServiceId:    140,
		ServiceName:  "public_cloud",
		CustomerName: "imported-cloud",
		InternalName: "PC-imported",
		Description:  "managed via terraform import",
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ResourceName:  "infomaniak_public_cloud.this",
				Config:        test.MustGetTestFile("schema", "resource_public_cloud_import.tf"),
				ImportState:   true,
				ImportStateId: fmt.Sprint(seedId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infomaniak_public_cloud.this", "id", fmt.Sprint(seedId)),
					resource.TestCheckResourceAttr("infomaniak_public_cloud.this", "customer_name", "imported-cloud"),
					resource.TestCheckResourceAttr("infomaniak_public_cloud.this", "internal_name", "PC-imported"),
					resource.TestCheckResourceAttr("infomaniak_public_cloud.this", "description", "managed via terraform import"),
				),
			},
		},
	})
}
