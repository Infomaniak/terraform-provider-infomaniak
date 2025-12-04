package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"terraform-provider-infomaniak/internal/dynamic"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ObjectStateManager will keep the state of stateEffective up to date with newEffective.
// It will also keep userDefined up to date (prevents changes when API set default values)
func ObjectStateManager(ctx context.Context, newEffective types.Dynamic, stateEffective types.Dynamic, userDefined types.Dynamic) (types.Dynamic, types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics

	incomingFromApi, d := ConvertDynamicObjectToTerraformMap(newEffective)
	diags.Append(d...)

	incomingFromState, d := ConvertDynamicObjectToTerraformMap(stateEffective)
	diags.Append(d...)

	local, d := ConvertDynamicObjectToTerraformMap(userDefined)
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
			stateEffectiveTfValue, err := stateEffectiveValue.ToTerraformValue(ctx)
			if err != nil {
				diags.AddError("could not get terraform value", "could not get a state effective value")
			}
			incomingEffectiveTfValue, err := incomingEffectiveValue.ToTerraformValue(ctx)
			if err != nil {
				diags.AddError("could not get terraform value", "could not get an incoming from api effective value")
			}
			if !stateEffectiveTfValue.Equal(incomingEffectiveTfValue) {
				local[incomingEffectiveKey] = incomingEffectiveValue
			}
		}
		incomingFromState[incomingEffectiveKey] = incomingEffectiveValue
	}

	localMap, d := ConvertMapToDynamicObject(ctx, local)
	diags.Append(d...)

	effectiveMap, d := ConvertMapToDynamicObject(ctx, incomingFromState)
	diags.Append(d...)

	return effectiveMap, localMap, diags
}

// ConvertDynamicObjectToTerraformMap will transform a Dynamic terraform object into a map of terraform values
func ConvertDynamicObjectToTerraformMap(dyn types.Dynamic) (map[string]attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if dyn.IsNull() || dyn.IsUnknown() {
		return make(map[string]attr.Value), diags
	}

	objVal, ok := dyn.UnderlyingValue().(basetypes.ObjectValue)
	if !ok {
		diags.AddError(
			"Invalid type",
			fmt.Sprintf("dynamic should be an object, got %T", dyn.UnderlyingValue()),
		)
		return nil, diags
	}
	elems := objVal.Attributes()
	return elems, diags
}

// ConvertDynamicObjectToMapAny will transform a Dynamic terraform object into a map of any go values
func ConvertDynamicObjectToMapAny(dyn types.Dynamic) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	converted := make(map[string]any)

	if dyn.IsNull() || dyn.IsUnknown() {
		return converted, diags
	}

	body, err := dynamic.ToJSON(dyn)
	if err != nil {
		diags.AddError("json error", fmt.Sprintf("could not convert dynamic to json: %v", err))
	}

	err = json.Unmarshal(body, &converted)
	if err != nil {
		diags.AddError("json error", fmt.Sprintf("could not unmarshall json: %v", err))
	}

	return converted, diags
}

// ConvertMapToDynamicObject will transform a map of terraform values into a terraform dynamic object
func ConvertMapToDynamicObject(ctx context.Context, toconvert map[string]attr.Value) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(toconvert) == 0 {
		return types.DynamicNull(), diags
	}

	typesMap := make(map[string]attr.Type)
	valuesMap := make(map[string]attr.Value)

	for k, v := range toconvert {
		typesMap[k] = v.Type(ctx)
		valuesMap[k] = v
	}

	obj, d := types.ObjectValue(typesMap, valuesMap)

	diags.Append(d...)
	dyn := types.DynamicValue(obj)

	return dyn, diags
}

// ConvertIntsToStrings will transform comparable types
// E.g: 100 = "100"
func ConvertIntsToStrings(input map[string]any) map[string]any {
	output := make(map[string]any)
	for key, value := range input {
		switch typedValue := value.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			output[key] = fmt.Sprint(typedValue)
		case map[string]any:
			newOutput := make(map[string]any)
			ConvertIntsToStrings(typedValue)
			output[key] = newOutput
		default:
			output[key] = typedValue
		}
	}
	return output
}
