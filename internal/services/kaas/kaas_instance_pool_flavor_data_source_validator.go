package kaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type NameOrResourcesValidator struct {
}

func (v NameOrResourcesValidator) Description(ctx context.Context) string {
	return "either 'name' or at least one of 'cpu', 'ram', 'storage' must be specified"
}

func (v NameOrResourcesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v NameOrResourcesValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var data KaasInstancePoolFlavorModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameSet := !data.Name.IsNull() && !data.Name.IsUnknown() && data.Name.ValueString() != ""

	anyResourceSet :=
		(!data.Cpu.IsNull() && !data.Cpu.IsUnknown()) ||
			(!data.Ram.IsNull() && !data.Ram.IsUnknown()) ||
			(!data.Storage.IsNull() && !data.Storage.IsUnknown())

	if !nameSet && !anyResourceSet {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid configuration",
			v.Description(ctx),
		)
	}
}
