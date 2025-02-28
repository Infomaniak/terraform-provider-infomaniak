package helpers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	BaseURI = "https://example.com"
)

var _ = Describe("Helpers Tests", func() {
	Context("Test Api Error Formatting", func() {
		err := ApiError{
			Description: "validation",
			Errors: []*ApiError{
				{
					Description: "required",
					Errors: []*ApiError{
						{
							Description: "missing field",
							Context: ApiErrorContext{
								Values: []any{"tete", "tata"},
							},
						},
					},
				},
				{
					Description: "required",
					Errors: []*ApiError{
						{
							Description: "missing field",
						},
					},
				},
			},
		}

		It("should tell the user about missing fields", func() {
			Expect(err.Error()).To(ContainSubstring("tete"))
			Expect(err.Error()).To(ContainSubstring("tata"))
		})
	})
})
