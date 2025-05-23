package validators

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = ipValidator{}

type ipValidator struct {
	Type string
}

func (validator ipValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string must be an %s", validator.Type)
}

func (validator ipValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v ipValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	ip, err := netip.ParseAddr(value)
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			err.Error(),
		))

		return
	}

	verifFunc := ip.Is4
	if v.Type == "IPv6" {
		verifFunc = ip.Is6
	}

	if !verifFunc() {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("string must be an %s", v.Type),
		))

		return
	}
}

func IsIPv4() ipValidator {
	return ipValidator{
		Type: "IPv4",
	}
}

func IsIPv6() ipValidator {
	return ipValidator{
		Type: "IPv6",
	}
}
