package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AllOnlyUnanimousWarnings returns a validator which ensures that any configured attribute value
// attribute value validates against all the given validators.
//
// Use of AllOnlyUnanimousWarnings is only necessary when used in conjunction with Any or AnyWithAllWarnings
// as the Validators field automatically applies a logical AND.
func AllOnlyUnanimousWarnings(validators ...validator.String) validator.String {
	return allValidator{
		validators: validators,
	}
}

var _ validator.String = allValidator{}

// allValidator implements the validator.
type allValidator struct {
	validators []validator.String
}

// Description describes the validation in plain text formatting.
func (v allValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, subValidator := range v.validators {
		descriptions = append(descriptions, subValidator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy all of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v allValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v allValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var warningCount int
	for _, subValidator := range v.validators {
		validateResp := &validator.StringResponse{}

		subValidator.ValidateString(ctx, req, validateResp)

		resp.Diagnostics.Append(validateResp.Diagnostics...)

		if validateResp.Diagnostics.HasError() {
			return
		}

		if validateResp.Diagnostics.WarningsCount() > 0 {
			warningCount++
		}
	}

	if warningCount < len(v.validators) {
		resp.Diagnostics = resp.Diagnostics.Errors()
	}
}
