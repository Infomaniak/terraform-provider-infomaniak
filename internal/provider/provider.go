// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package provider

import (
	"context"
	"os"
	"terraform-provider-infomaniak/internal/apis"
	"terraform-provider-infomaniak/internal/provider/registry"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Environment variables used by the provider
const (
	INFOMANIAK_TOKEN = "INFOMANIAK_TOKEN"
	INFOMANIAK_HOST  = "INFOMANIAK_HOST"
)

// Ensure IkProvider satisfies various kaas interfaces.
var (
	_ provider.Provider              = &IkProvider{}
	_ provider.ProviderWithFunctions = &IkProvider{}

	DefaultHost = "https://api.infomaniak.com"
)

// IkProvider defines the kaas implementation.
type IkProvider struct {
	// version is set to the kaas version on release, "dev" when the
	// kaas is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	ik *IkProviderData
}

// IkProviderData defines the data associated with the provider
type IkProviderData struct {
	*apis.Client

	Data *IkProviderModel
}

type IkProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

func (p *IkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "infomaniak"
	resp.Version = p.version
}

func (p *IkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:            true,
				Description:         "The base endpoint for Infomaniak's API (including scheme).",
				MarkdownDescription: "The base endpoint for Infomaniak's API (including scheme).",
			},
			"token": schema.StringAttribute{
				Optional:            os.Getenv(INFOMANIAK_TOKEN) != "",
				Sensitive:           true,
				Description:         "The token used for authenticating against Infomaniak's API.",
				MarkdownDescription: "The token used for authenticating against Infomaniak's API.",
			},
		},
		Description:         "Infomaniak's provider.",
		MarkdownDescription: "Infomaniak's provider.",
	}
}

func (p *IkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Provider configuration started")

	if p.ik != nil {
		tflog.Debug(ctx, "Provider already present, skipping configuration")
		resp.DataSourceData = p.ik
		resp.ResourceData = p.ik
		return
	}

	var data IkProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Infomaniak API Host",
			"The provider cannot create the Infomaniak API client as there is an unknown configuration value for the Infomaniak API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INFOMANIAK_HOST environment variable.",
		)
	}

	if data.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Infomaniak API Token",
			"The provider cannot create the Infomaniak API client as there is an unknown configuration value for the Infomaniak API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INFOMANIAK_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv(INFOMANIAK_HOST)
	token := os.Getenv(INFOMANIAK_TOKEN)

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	if host == "" {
		host = DefaultHost
		data.Host = types.StringValue(host)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Infomaniak API Username",
			"The provider cannot create the Infomaniak API client as there is a missing or empty value for the Infomaniak API username. "+
				"Set the username value in the configuration or use the INFOMANIAK_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	p.ik = &IkProviderData{
		Client: apis.NewClient(host),
		Data:   &data,
	}

	resp.DataSourceData = p.ik
	resp.ResourceData = p.ik
}

func (p *IkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return registry.GetResources()
}

func (p *IkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return registry.GetDataSources()
}

func (p *IkProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IkProvider{
			version: version,
		}
	}
}

func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"infomaniak": providerserver.NewProtocol6WithError(New("test")()),
	}
}
