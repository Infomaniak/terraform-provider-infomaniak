package domain

import (
	"context"
	"fmt"
	"reflect"
	"terraform-provider-infomaniak/internal/apis/domain"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ validator.String = NewAttributeValidatorFor(domain.RecordTypeA)
var _ validator.Int64 = NewAttributeValidatorFor(domain.RecordTypeA)

func NewAttributeValidatorFor[T domain.RecordConstraint](recordType T) attributeValidatorFor[T] {
	return attributeValidatorFor[T]{
		recordType: recordType,
	}
}

type attributeValidatorFor[T domain.RecordConstraint] struct {
	recordType T
}

func (validator attributeValidatorFor[T]) Description(_ context.Context) string {
	return ""
}

func (validator attributeValidatorFor[T]) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v attributeValidatorFor[T]) Validate(ctx context.Context, path path.Path, config tfsdk.Config, configVal tftypes.Value, response diag.Diagnostics) diag.Diagnostics {
	var configModel RecordModel

	diag := config.Get(ctx, &configModel)
	if diag.ErrorsCount() > 0 {
		return append(response, diag...)
	}

	if !configModel.RawTarget.IsNull() {
		if !configVal.IsNull() {
			response.AddAttributeWarning(
				path,
				"field will not be used",
				fmt.Sprintf("field %v will not be used when raw_target is specified", path),
			)
			return response
		}
		return response
	}

	effectiveRecordType := configModel.Type.ValueString()
	expectedRecordType := reflect.ValueOf(v.recordType).Field(0).String()

	if effectiveRecordType == expectedRecordType {
		// In this case the field must be set
		if configVal.IsNull() {
			response.Append(validatordiag.InvalidAttributeValueDiagnostic(
				path,
				"field is required",
				fmt.Sprintf("field must be set when Record Type is %s", effectiveRecordType),
			))
			return response
		}
	} else {
		// In this case the field must be unset
		if !configVal.IsNull() {
			response.AddAttributeWarning(
				path,
				"field will not be used",
				fmt.Sprintf("field %v will not be used because your record type is %s", path, effectiveRecordType),
			)
			return response
		}
	}

	return response
}

func (v attributeValidatorFor[T]) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	val, err := request.ConfigValue.ToTerraformValue(ctx)
	if err != nil {
		response.Diagnostics.AddError(
			"got error when converting value",
			err.Error(),
		)
		return
	}

	response.Diagnostics = append(response.Diagnostics, v.Validate(ctx, request.Path, request.Config, val, response.Diagnostics)...)
}

func (v attributeValidatorFor[T]) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	val, err := request.ConfigValue.ToTerraformValue(ctx)
	if err != nil {
		response.Diagnostics.AddError(
			"got error when converting value",
			err.Error(),
		)
		return
	}

	response.Diagnostics = append(response.Diagnostics, v.Validate(ctx, request.Path, request.Config, val, response.Diagnostics)...)
}
