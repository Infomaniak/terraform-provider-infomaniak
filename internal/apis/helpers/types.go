package helpers

import (
	"fmt"
	"strings"
)

type NormalizedApiResponse[K any] struct {
	Result string    `json:"result"`
	Data   K         `json:"data"`
	Error  *ApiError `json:"error"`
}

type ApiError struct {
	Description string          `json:"description"`
	Errors      []*ApiError     `json:"errors"`
	Context     ApiErrorContext `json:"context"`
}

type ApiErrorContext struct {
	Attribute string `json:"attribute"`
	Values    []any  `json:"values"`
}

func (apiError *ApiError) Error() string {
	var builder strings.Builder

	builder.WriteString(apiError.Description)

	if len(apiError.Context.Values) > 0 {
		builder.WriteString(fmt.Sprintf(" (possible values: %v)", apiError.Context.Values))
	}

	if len(apiError.Errors) > 0 {
		builder.WriteString(":\n")
	}

	for _, err := range apiError.Errors {
		tabulated := "  " + strings.ReplaceAll(err.Error(), "\n", "\n  ")
		builder.WriteString(tabulated + "\n")
	}

	return strings.TrimSuffix(builder.String(), "\n")
}
