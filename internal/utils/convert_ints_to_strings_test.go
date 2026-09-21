package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"terraform-provider-infomaniak/internal/utils"
)

var _ = Describe("ConvertIntsToStrings", func() {
	It("should convert top-level numbers to strings", func() {
		input := map[string]any{
			"max_connections": 100,
			"ratio":           0.5,
			"name":            "test",
		}

		output := utils.ConvertIntsToStrings(input)

		Expect(output).To(Equal(map[string]any{
			"max_connections": "100",
			"ratio":           "0.5",
			"name":            "test",
		}))
	})

	It("should convert numbers inside nested maps", func() {
		input := map[string]any{
			"mysql": map[string]any{
				"max_connections": 100,
				"slow_query_log":  true,
			},
		}

		output := utils.ConvertIntsToStrings(input)

		Expect(output).To(Equal(map[string]any{
			"mysql": map[string]any{
				"max_connections": "100",
				"slow_query_log":  true,
			},
		}))
	})

	It("should convert numbers in deeply nested maps", func() {
		input := map[string]any{
			"mysql": map[string]any{
				"options": map[string]any{
					"max_allowed_packet": 67108864,
				},
			},
		}

		output := utils.ConvertIntsToStrings(input)

		Expect(output).To(Equal(map[string]any{
			"mysql": map[string]any{
				"options": map[string]any{
					"max_allowed_packet": "67108864",
				},
			},
		}))
	})
})
