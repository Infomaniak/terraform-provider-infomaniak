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

	localMap, d := ConvertMapToDynamicObject(local)
	diags.Append(d...)

	effectiveMap, d := ConvertMapToDynamicObject(incomingFromState)
	diags.Append(d...)

	return effectiveMap, localMap, diags
}

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

func ConvertDynamicObjectToMapAny(dyn types.Dynamic) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	converted := make(map[string]any)

	if dyn.IsNull() || dyn.IsUnknown() {
		return converted, diags
	}

	body, err := dynamic.ToJSON(dyn)
	if err != nil {
		diags.AddError("error", "error")
	}

	err = json.Unmarshal(body, &converted)
	if err != nil {
		diags.AddError("error", "error")
	}

	return converted, diags
}

func ConvertMapToDynamicObject(toconvert map[string]attr.Value) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(toconvert) == 0 {
		return types.DynamicNull(), diags
	}

	typesMap := make(map[string]attr.Type)
	valuesMap := make(map[string]attr.Value)

	for k, v := range toconvert {
		val := types.DynamicValue(v)
		typesMap[k] = types.DynamicType
		valuesMap[k] = val
	}

	obj, d := types.ObjectValue(typesMap, valuesMap)

	diags.Append(d...)
	dyn := types.DynamicValue(obj)

	return dyn, diags
}
