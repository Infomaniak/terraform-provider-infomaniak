package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringMapStateManager(ctx context.Context, newEffective types.Map, stateEffective types.Map, userDefined types.Map) (types.Map, types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	incomingFromApi := make(map[string]string)
	diags.Append(newEffective.ElementsAs(ctx, &incomingFromApi, false)...)

	incomingFromState := make(map[string]string)
	diags.Append(stateEffective.ElementsAs(ctx, &incomingFromState, false)...)

	local := make(map[string]string)
	diags.Append(userDefined.ElementsAs(ctx, &local, false)...)

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
			if stateEffectiveValue != incomingEffectiveValue {
				local[incomingEffectiveKey] = incomingEffectiveValue
			}
			incomingFromState[incomingEffectiveKey] = incomingEffectiveValue
		}
	}

	localMap, diags := types.MapValueFrom(ctx, types.StringType, local)
	diags.Append(diags...)

	effectiveMap, diags := types.MapValueFrom(ctx, types.StringType, incomingFromState)
	diags.Append(diags...)

	return effectiveMap, localMap, diags
}
