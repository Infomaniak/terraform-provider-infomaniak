package helpers

import (
	"fmt"
	"testing"
)

func Test_ClientErrorPrint(t *testing.T) {
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

	fmt.Println(err.Error())
}
