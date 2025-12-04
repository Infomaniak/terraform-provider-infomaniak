package dynamic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DynamicObjectValidator struct {
	Strict bool
}

func NewDynamicObjectValidator() *DynamicObjectValidator {
	return &DynamicObjectValidator{}
}

var _ validator.Dynamic = (*DynamicObjectValidator)(nil)

func (dv *DynamicObjectValidator) Description(context.Context) string {
	return ""
}

func (dv *DynamicObjectValidator) MarkdownDescription(context.Context) string {
	return ""
}

func (dv *DynamicObjectValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, res *validator.DynamicResponse) {
	if req.ConfigValue.IsNull() {
		return
	}

	obj, isObject := req.ConfigValue.UnderlyingValue().(basetypes.ObjectValue)
	if !isObject {
		res.Diagnostics.AddAttributeError(req.Path, "Type mismatch", fmt.Sprintf("Attribute should be an object and not a %v", obj.Type(ctx)))
		return
	}

	elems := obj.Attributes()
	if len(elems) == 0 {
		res.Diagnostics.AddAttributeError(req.Path, "configuration is empty", "configuration needs to have at least one element, delete the field if you do not want to configure it")
	}
	for _, value := range elems {
		_, isList := value.(basetypes.ListValue)
		_, isSet := value.(basetypes.SetValue)
		if isList || isSet {
			res.Diagnostics.AddAttributeError(req.Path, "Wrong type", "please use tuple when using a list or set inside a dynamic object")
		}
	}
}
