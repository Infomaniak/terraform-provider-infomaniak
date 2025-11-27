package utils

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ObjectStateManager(ctx context.Context, newEffective types.Dynamic, stateEffective types.Dynamic, userDefined types.Dynamic) (types.Dynamic, types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics

	incomingFromApi, d := ConvertToMap(newEffective)
	diags.Append(d...)

	incomingFromState, d := ConvertToMap(stateEffective)
	diags.Append(d...)

	local, d := ConvertToMap(userDefined)
	diags.Append(d...)

	for incomingEffectiveKey, incomingEffectiveValue := range incomingFromApi {
		_, localManagerUseKey := local[incomingEffectiveKey]
		if localManagerUseKey {
			local[incomingEffectiveKey] = incomingEffectiveValue
		}
	}

	for incomingEffectiveKey, incomingEffectiveValue := range incomingFromApi {
		stateEffectiveValue, stateEffectiveUseKey := incomingFromState[incomingEffectiveKey]
		if stateEffectiveUseKey {
			// The user changed the value from an other source than terraform
			if !reflect.DeepEqual(stateEffectiveValue, incomingEffectiveValue) {
				local[incomingEffectiveKey] = incomingEffectiveValue
			}
		}
		incomingFromState[incomingEffectiveKey] = incomingEffectiveValue
	}

	localMap, d := ConvertMap(ctx, local)
	diags.Append(d...)

	effectiveMap, d := ConvertMap(ctx, incomingFromState)
	diags.Append(d...)

	return effectiveMap, localMap, diags
}

func ConvertToMap(dyn types.Dynamic) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	if dyn.IsNull() || dyn.IsUnknown() {
		return nil, diags
	}

	objVal, ok := dyn.UnderlyingValue().(basetypes.ObjectValue)
	if !ok {
		diags.AddError(
			"Invalid type",
			fmt.Sprintf("dynamic should be an object, got %T", dyn.UnderlyingValue()),
		)
		return nil, diags
	}

	converted := make(map[string]any)

	elems := objVal.Attributes()

	for k, v := range elems {
		decoded, err := DecodeValue(v)
		if err != nil {
			diags.AddError(
				"Failed to decode dynamic object",
				err.Error(),
			)
		}
		converted[k] = decoded
	}

	return converted, diags
}

func ConvertMap(ctx context.Context, toconvert map[string]any) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics
	typesMap := make(map[string]attr.Type)
	valuesMap := make(map[string]attr.Value)

	for k, v := range toconvert {
		val, d := convertToDynamic(ctx, v)
		diags.Append(d...)
		typesMap[k] = types.DynamicType
		valuesMap[k] = val
	}

	obj, d := types.ObjectValue(typesMap, valuesMap)

	diags.Append(d...)
	dynamic := types.DynamicValue(obj)

	return dynamic, diags
}

func convertToDynamic(ctx context.Context, v any) (types.Dynamic, diag.Diagnostics) {
	switch t := v.(type) {
	case nil:
		return types.DynamicNull(), nil
	case string:
		return types.DynamicValue(types.StringValue(t)), nil
	case bool:
		return types.DynamicValue(types.BoolValue(t)), nil
	case int:
		return types.DynamicValue(types.Int64Value(int64(t))), nil
	case int64:
		return types.DynamicValue(types.Int64Value(t)), nil
	case float64:
		return types.DynamicValue(types.Float64Value(t)), nil
	case []any:
		lv, diags := types.ListValueFrom(ctx, types.DynamicType, t)
		return types.DynamicValue(lv), diags
	case map[string]any:
		attrTypes := map[string]attr.Type{}
		for k := range t {
			attrTypes[k] = types.DynamicType
		}
		ov, diags := types.ObjectValueFrom(ctx, attrTypes, t)
		if diags.HasError() {
			return types.Dynamic{}, diags
		}
		return types.DynamicValue(ov), nil
	default:
		return types.DynamicValue(types.StringValue(fmt.Sprint(v))), nil
	}
}

func DecodeValue(v attr.Value) (any, error) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	switch val := v.(type) {

	case types.Dynamic:
		underlying := val.UnderlyingValue()
		return DecodeValue(underlying)

	case types.String:
		return val.ValueString(), nil

	case types.Int64:
		return val.ValueInt64(), nil

	case types.Float64:
		return val.ValueFloat64(), nil

	case types.Bool:
		return val.ValueBool(), nil

	case types.Number:
		return val.ValueBigFloat(), nil

	case types.Tuple:
		var result []any
		elems := val.Elements()
		for _, e := range elems {
			decoded, err := DecodeValue(e)
			if err != nil {
				return nil, err
			}
			result = append(result, decoded)
		}
		return result, nil

	case types.List:
		var result []any
		elems := val.Elements()
		for _, e := range elems {
			decoded, err := DecodeValue(e)
			if err != nil {
				return nil, err
			}
			result = append(result, decoded)
		}
		return result, nil

	case types.Map:
		result := make(map[string]any)
		elems := val.Elements()
		for key, e := range elems {
			decoded, err := DecodeValue(e)
			if err != nil {
				return nil, err
			}
			result[key] = decoded
		}
		return result, nil

	default:
		return nil, fmt.Errorf("unsupported attr.Value type: %T", v)
	}
}
