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
func Int64AllOnlyUnanimousWarnings(validators ...validator.Int64) validator.Int64 {
	return int64AllValidator{
		validators: validators,
	}
}

var _ validator.Int64 = int64AllValidator{}

// int64AllValidator implements the validator.
type int64AllValidator struct {
	validators []validator.Int64
}

// Description describes the validation in plain text formatting.
func (v int64AllValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, subValidator := range v.validators {
		descriptions = append(descriptions, subValidator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy all of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v int64AllValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v int64AllValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	var warningCount int
	for _, subValidator := range v.validators {
		validateResp := &validator.Int64Response{}

		subValidator.ValidateInt64(ctx, req, validateResp)

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
