package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/femnad/terraform-provider-omglol/internal/omglol"
)

const (
	apiKeyEnv   = "OMGLOL_API_KEY"
	usernameEnv = "OMGLOL_USERNAME"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &omglolProvider{}
)

type omglolProvider struct {
	version string
}

type omglolProviderModel struct {
	ApiKey   types.String `tfsdk:"api_key"`
	Username types.String `tfsdk:"username"`
}

func New() func() provider.Provider {
	return func() provider.Provider {
		return omglolProvider{}
	}
}

func (p omglolProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config omglolProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown omg.lol API key",
			"The provider cannot create omg.lol API client as there is an unknown value for API key")
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown omg.lol username",
			"The provider cannot create omg.lol API client as there is an unknown value for username")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv(apiKeyEnv)
	username := os.Getenv(usernameEnv)

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing omg.lol API key",
			"The provider cannot create omg.lol API client as there is a missing value for API key"+
				fmt.Sprintf("Set the API key in the provider configuration or use the %s environment var",
					apiKeyEnv))
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing omg.lol username",
			"The provider cannot create omg.lol API client as there is a missing value for API key"+
				fmt.Sprintf("Set the username in the provider configuration or use the %s environment var",
					usernameEnv))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := omglol.NewClient(username, apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable create omg.lol API client",
			fmt.Sprintf("Unexpected error: %s", err.Error()))
	}

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p omglolProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "omglol"
	resp.Version = p.version
}

func (p omglolProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
		}}
}

func (p omglolProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p omglolProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDNSResource,
	}
}
