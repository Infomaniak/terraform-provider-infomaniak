package dbaasmigration

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func v0ObjectType() tftypes.Type {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"public_cloud_id":         tftypes.Number,
			"public_cloud_project_id": tftypes.Number,
			"id":                      tftypes.Number,
			"kube_identifier":         tftypes.String,
			"name":                    tftypes.String,
			"pack_name":               tftypes.String,
			"region":                  tftypes.String,
			"type":                    tftypes.String,
			"version":                 tftypes.String,
			"host":                    tftypes.String,
			"port":                    tftypes.String,
			"user":                    tftypes.String,
			"password":                tftypes.String,
			"ca":                      tftypes.String,
			"allowed_cidrs":           tftypes.List{ElementType: tftypes.String},
			"configuration":           tftypes.Map{ElementType: tftypes.String},
			"effective_configuration": tftypes.Map{ElementType: tftypes.String},
		},
	}
}

func v1ObjectType() tftypes.Type {
	objectType := v0ObjectType().(tftypes.Object)
	objectType.AttributeTypes["configuration"] = tftypes.DynamicPseudoType
	objectType.AttributeTypes["effective_configuration"] = tftypes.DynamicPseudoType
	return objectType
}

func v0Value(configuration, effectiveConfiguration tftypes.Value) tftypes.Value {
	return tftypes.NewValue(v0ObjectType(), map[string]tftypes.Value{
		"public_cloud_id":         tftypes.NewValue(tftypes.Number, big.NewFloat(1)),
		"public_cloud_project_id": tftypes.NewValue(tftypes.Number, big.NewFloat(2)),
		"id":                      tftypes.NewValue(tftypes.Number, big.NewFloat(3)),
		"kube_identifier":         tftypes.NewValue(tftypes.String, "dbaas-abc123"),
		"name":                    tftypes.NewValue(tftypes.String, "test-dbaas"),
		"pack_name":               tftypes.NewValue(tftypes.String, "db-xxs-1"),
		"region":                  tftypes.NewValue(tftypes.String, "dc3-a"),
		"type":                    tftypes.NewValue(tftypes.String, "mysql"),
		"version":                 tftypes.NewValue(tftypes.String, "8.0"),
		"host":                    tftypes.NewValue(tftypes.String, "test-dbaas.infomaniak.cloud"),
		"port":                    tftypes.NewValue(tftypes.String, "3306"),
		"user":                    tftypes.NewValue(tftypes.String, "root"),
		"password":                tftypes.NewValue(tftypes.String, "secret"),
		"ca":                      tftypes.NewValue(tftypes.String, "ca-content"),
		"allowed_cidrs": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
			tftypes.NewValue(tftypes.String, "10.0.0.0/8"),
		}),
		"configuration":           configuration,
		"effective_configuration": effectiveConfiguration,
	})
}

func v1Schema() schema.Schema {
	return schema.Schema{
		Version: 1,
		Attributes: map[string]schema.Attribute{
			"public_cloud_id":         schema.Int64Attribute{},
			"public_cloud_project_id": schema.Int64Attribute{},
			"id":                      schema.Int64Attribute{},
			"kube_identifier":         schema.StringAttribute{},
			"name":                    schema.StringAttribute{},
			"pack_name":               schema.StringAttribute{},
			"region":                  schema.StringAttribute{},
			"type":                    schema.StringAttribute{},
			"version":                 schema.StringAttribute{},
			"host":                    schema.StringAttribute{},
			"port":                    schema.StringAttribute{},
			"user":                    schema.StringAttribute{},
			"password":                schema.StringAttribute{},
			"ca":                      schema.StringAttribute{},
			"allowed_cidrs":           schema.ListAttribute{ElementType: types.StringType},
			"configuration":           schema.DynamicAttribute{},
			"effective_configuration": schema.DynamicAttribute{},
		},
	}
}

func emptyV1State() *tfsdk.State {
	objectType := v1ObjectType().(tftypes.Object)
	attributes := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		attributes[name] = tftypes.NewValue(attributeType, nil)
	}
	return &tfsdk.State{Raw: tftypes.NewValue(v1ObjectType(), attributes), Schema: v1Schema()}
}

func runUpgrader(t *testing.T, configuration, effectiveConfiguration tftypes.Value) (resource.UpgradeStateResponse, DBaasModelV1) {
	t.Helper()

	request := resource.UpgradeStateRequest{
		State: &tfsdk.State{
			Raw:    v0Value(configuration, effectiveConfiguration),
			Schema: GetV0Schema(),
		},
	}
	response := resource.UpgradeStateResponse{State: *emptyV1State()}

	StateUpgrader(context.Background(), request, &response)

	var upgraded DBaasModelV1
	if !response.Diagnostics.HasError() {
		diags := response.State.Get(context.Background(), &upgraded)
		if diags.HasError() {
			t.Fatalf("failed to decode upgraded state: %v", diags)
		}
	}

	return response, upgraded
}

func dynamicStringAttributes(t *testing.T, dyn types.Dynamic) map[string]string {
	t.Helper()

	obj, ok := dyn.UnderlyingValue().(types.Object)
	if !ok {
		t.Fatalf("expected dynamic object, got %T", dyn.UnderlyingValue())
	}
	attributes := make(map[string]string, len(obj.Attributes()))
	for name, value := range obj.Attributes() {
		str, ok := value.(types.String)
		if !ok {
			t.Fatalf("expected string attribute %q, got %T", name, value)
		}
		attributes[name] = str.ValueString()
	}
	return attributes
}

func TestStateUpgrader_CopiesScalarState(t *testing.T) {
	response, upgraded := runUpgrader(t,
		tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
			"maintenance_dow":  tftypes.NewValue(tftypes.String, "monday"),
			"maintenance_time": tftypes.NewValue(tftypes.String, "03:00:00"),
		}),
		tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
			"maintenance_dow":  tftypes.NewValue(tftypes.String, "monday"),
			"maintenance_time": tftypes.NewValue(tftypes.String, "03:00:00"),
		}),
	)

	if response.Diagnostics.HasError() {
		t.Fatalf("unexpected upgrade diagnostics: %v", response.Diagnostics)
	}
	if upgraded.PublicCloudId.ValueInt64() != 1 {
		t.Errorf("public_cloud_id not copied, got %v", upgraded.PublicCloudId.ValueInt64())
	}
	if upgraded.Id.ValueInt64() != 3 {
		t.Errorf("id not copied, got %v", upgraded.Id.ValueInt64())
	}
	if upgraded.Name.ValueString() != "test-dbaas" {
		t.Errorf("name not copied, got %v", upgraded.Name.ValueString())
	}

	expectedConfiguration := map[string]string{"maintenance_dow": "monday", "maintenance_time": "03:00:00"}
	if diff := compareMaps(expectedConfiguration, dynamicStringAttributes(t, upgraded.Configuration)); diff != "" {
		t.Errorf("configuration not preserved during upgrade:\n%s", diff)
	}
	if diff := compareMaps(expectedConfiguration, dynamicStringAttributes(t, upgraded.EffectiveConfiguration)); diff != "" {
		t.Errorf("effective_configuration not preserved during upgrade:\n%s", diff)
	}
}

func TestStateUpgrader_UnknownConfigurationWritesState(t *testing.T) {
	response, upgraded := runUpgrader(t,
		tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, tftypes.UnknownValue),
		tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{}),
	)

	if response.Diagnostics.HasError() {
		t.Fatalf("unexpected upgrade diagnostics: %v", response.Diagnostics)
	}
	if !upgraded.Configuration.IsUnknown() {
		t.Errorf("configuration should stay unknown, got %v", upgraded.Configuration)
	}
	if upgraded.PublicCloudId.ValueInt64() != 1 {
		t.Errorf("state was not written when configuration is unknown, public_cloud_id got %v", upgraded.PublicCloudId.ValueInt64())
	}
}

func TestStateUpgrader_NullConfigurationStaysNull(t *testing.T) {
	response, upgraded := runUpgrader(t,
		tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
	)

	if response.Diagnostics.HasError() {
		t.Fatalf("unexpected upgrade diagnostics: %v", response.Diagnostics)
	}
	if !upgraded.Configuration.IsNull() {
		t.Errorf("configuration should stay null, got %v", upgraded.Configuration)
	}
	if !upgraded.EffectiveConfiguration.IsNull() {
		t.Errorf("effective_configuration should stay null, got %v", upgraded.EffectiveConfiguration)
	}
	if upgraded.PublicCloudId.ValueInt64() != 1 {
		t.Errorf("state was not written when configuration is null, public_cloud_id got %v", upgraded.PublicCloudId.ValueInt64())
	}
}

func compareMaps(expected, actual map[string]string) string {
	if len(expected) != len(actual) {
		return "length mismatch"
	}
	for key, expectedValue := range expected {
		actualValue, ok := actual[key]
		if !ok {
			return "missing key: " + key
		}
		if expectedValue != actualValue {
			return "key " + key + ": expected " + expectedValue + ", got " + actualValue
		}
	}
	return ""
}
