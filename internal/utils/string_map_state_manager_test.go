package utils_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-infomaniak/internal/utils"
)

var _ = Describe("StringMapStateManager", func() {
	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("when there are no user defined settings", func() {
		It("should return empty maps", func() {
			newEffective := types.MapNull(types.StringType)
			stateEffective := types.MapNull(types.StringType)
			userDefined := types.MapNull(types.StringType)

			effectiveMap, localMap, diags := utils.StringMapStateManager(ctx, newEffective, stateEffective, userDefined)

			Expect(diags.HasError()).To(BeFalse())
			Expect(effectiveMap.IsNull()).To(BeTrue())
			Expect(localMap.IsNull()).To(BeTrue())
		})
	})

	Context("when there are user defined settings that match API values", func() {
		It("should preserve user defined settings", func() {
			newEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("value1"),
				"setting2": types.StringValue("value2"),
			}
			stateEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("value1"),
				"setting2": types.StringValue("value2"),
			}
			userDefinedMap := map[string]attr.Value{
				"setting1": types.StringValue("value1"),
				"setting2": types.StringValue("value2"),
			}

			newEffective, _ := types.MapValue(types.StringType, newEffectiveMap)
			stateEffective, _ := types.MapValue(types.StringType, stateEffectiveMap)
			userDefined, _ := types.MapValue(types.StringType, userDefinedMap)

			effectiveMap, localMap, diags := utils.StringMapStateManager(ctx, newEffective, stateEffective, userDefined)

			Expect(diags.HasError()).To(BeFalse())

			// Check localMap (userDefined)
			localElements := localMap.Elements()
			Expect(localElements).To(HaveLen(2))
			Expect(localElements["setting1"]).To(Equal(types.StringValue("value1")))
			Expect(localElements["setting2"]).To(Equal(types.StringValue("value2")))

			// Check effectiveMap (stateEffective updated)
			effectiveElements := effectiveMap.Elements()
			Expect(effectiveElements).To(HaveLen(2))
			Expect(effectiveElements["setting1"]).To(Equal(types.StringValue("value1")))
			Expect(effectiveElements["setting2"]).To(Equal(types.StringValue("value2")))
		})
	})

	Context("when API returns new values for user managed settings", func() {
		It("should update user defined settings with API values", func() {
			newEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("new_value1"), // Changed
				"setting2": types.StringValue("value2"),
			}
			stateEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("old_value1"), // Old value
				"setting2": types.StringValue("value2"),
			}
			userDefinedMap := map[string]attr.Value{
				"setting1": types.StringValue("old_value1"), // User's old value
				"setting2": types.StringValue("value2"),
			}

			newEffective, _ := types.MapValue(types.StringType, newEffectiveMap)
			stateEffective, _ := types.MapValue(types.StringType, stateEffectiveMap)
			userDefined, _ := types.MapValue(types.StringType, userDefinedMap)

			effectiveMap, localMap, diags := utils.StringMapStateManager(ctx, newEffective, stateEffective, userDefined)

			Expect(diags.HasError()).To(BeFalse())

			// Check localMap (userDefined) - should be updated with new API values
			localElements := localMap.Elements()
			Expect(localElements).To(HaveLen(2))
			Expect(localElements["setting1"]).To(Equal(types.StringValue("new_value1"))) // Updated
			Expect(localElements["setting2"]).To(Equal(types.StringValue("value2")))

			// Check effectiveMap (stateEffective updated)
			effectiveElements := effectiveMap.Elements()
			Expect(effectiveElements).To(HaveLen(2))
			Expect(effectiveElements["setting1"]).To(Equal(types.StringValue("new_value1"))) // Updated
			Expect(effectiveElements["setting2"]).To(Equal(types.StringValue("value2")))
		})
	})

	Context("when there are settings not managed by user", func() {
		It("should preserve API values for non-user managed settings", func() {
			newEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("value1"), // User managed
				"setting2": types.StringValue("value2"), // Not user managed
				"setting3": types.StringValue("value3"), // Not user managed
			}
			stateEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("value1"),
				"setting2": types.StringValue("value2"),
				"setting3": types.StringValue("value3"),
			}
			userDefinedMap := map[string]attr.Value{
				"setting1": types.StringValue("value1"), // Only this is user managed
			}

			newEffective, _ := types.MapValue(types.StringType, newEffectiveMap)
			stateEffective, _ := types.MapValue(types.StringType, stateEffectiveMap)
			userDefined, _ := types.MapValue(types.StringType, userDefinedMap)

			effectiveMap, localMap, diags := utils.StringMapStateManager(ctx, newEffective, stateEffective, userDefined)

			Expect(diags.HasError()).To(BeFalse())

			// Check localMap (userDefined) - should only contain user managed settings
			localElements := localMap.Elements()
			Expect(localElements).To(HaveLen(1)) // Only setting1
			Expect(localElements["setting1"]).To(Equal(types.StringValue("value1")))

			// Check effectiveMap (stateEffective updated)
			effectiveElements := effectiveMap.Elements()
			Expect(effectiveElements).To(HaveLen(3)) // All settings
			Expect(effectiveElements["setting1"]).To(Equal(types.StringValue("value1")))
			Expect(effectiveElements["setting2"]).To(Equal(types.StringValue("value2")))
			Expect(effectiveElements["setting3"]).To(Equal(types.StringValue("value3")))
		})
	})

	Context("when user adds new settings that weren't previously managed", func() {
		It("should include new user settings in the result", func() {
			newEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("value1"),
				"setting2": types.StringValue("value2"),
			}
			stateEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("value1"),
				"setting2": types.StringValue("value2"),
			}
			userDefinedMap := map[string]attr.Value{
				"setting1": types.StringValue("value1"),
				"setting2": types.StringValue("value2"),
				"setting3": types.StringValue("user_value3"), // New setting
			}

			newEffective, _ := types.MapValue(types.StringType, newEffectiveMap)
			stateEffective, _ := types.MapValue(types.StringType, stateEffectiveMap)
			userDefined, _ := types.MapValue(types.StringType, userDefinedMap)

			effectiveMap, localMap, diags := utils.StringMapStateManager(ctx, newEffective, stateEffective, userDefined)

			Expect(diags.HasError()).To(BeFalse())

			// Check localMap (userDefined) - should include the new setting
			localElements := localMap.Elements()
			Expect(localElements).To(HaveLen(3))
			Expect(localElements["setting1"]).To(Equal(types.StringValue("value1")))
			Expect(localElements["setting2"]).To(Equal(types.StringValue("value2")))
			Expect(localElements["setting3"]).To(Equal(types.StringValue("user_value3")))

			// Check effectiveMap (stateEffective updated)
			effectiveElements := effectiveMap.Elements()
			Expect(effectiveElements).To(HaveLen(2)) // Only API known settings
			Expect(effectiveElements["setting1"]).To(Equal(types.StringValue("value1")))
			Expect(effectiveElements["setting2"]).To(Equal(types.StringValue("value2")))
		})
	})

	Context("when effective value changed from elsewhere (not from Terraform)", func() {
		It("should update both effectiveMap and userDefined map", func() {
			newEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("api_changed_value"), // Changed from API
				"setting2": types.StringValue("value2"),
			}
			stateEffectiveMap := map[string]attr.Value{
				"setting1": types.StringValue("terraform_value"), // Old value from Terraform
				"setting2": types.StringValue("value2"),
			}
			userDefinedMap := map[string]attr.Value{
				"setting2": types.StringValue("value2"),
			}

			newEffective, _ := types.MapValue(types.StringType, newEffectiveMap)
			stateEffective, _ := types.MapValue(types.StringType, stateEffectiveMap)
			userDefined, _ := types.MapValue(types.StringType, userDefinedMap)

			effectiveMap, localMap, diags := utils.StringMapStateManager(ctx, newEffective, stateEffective, userDefined)

			Expect(diags.HasError()).To(BeFalse())

			// Check localMap (userDefined) - should be updated with API changed value
			localElements := localMap.Elements()
			Expect(localElements).To(HaveLen(2))
			Expect(localElements["setting1"]).To(Equal(types.StringValue("api_changed_value"))) // Updated
			Expect(localElements["setting2"]).To(Equal(types.StringValue("value2")))

			// Check effectiveMap (stateEffective updated) - should also be updated with API changed value
			effectiveElements := effectiveMap.Elements()
			Expect(effectiveElements).To(HaveLen(2))
			Expect(effectiveElements["setting1"]).To(Equal(types.StringValue("api_changed_value"))) // Updated
			Expect(effectiveElements["setting2"]).To(Equal(types.StringValue("value2")))
		})
	})
})