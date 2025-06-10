// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/zoullx/terraform-provider-unifi/internal/provider_unifi"

	"github.com/zoullx/unifi-go/unifi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &UnifiProvider{}

type UnifiProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *UnifiProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "unifi"
	resp.Version = p.version
}

func (p *UnifiProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = provider_unifi.UnifiProviderSchema(ctx)
}

func (p *UnifiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Unifi client")

	// Retrieve provider data from configuration
	var data provider_unifi.UnifiModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Unifi Host",
			"The provider cannot create the Unifi API client as there is an unknown configuration value for the Unifi host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the UNIFI_HOST environment variable.",
		)
	}

	if data.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Unifi API Key",
			"The provider cannot create the Unifi API client as there is an unknown configuration value for the Unifi API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the UNIFI_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("UNIFI_HOST")
	apiKey := os.Getenv("UNIFI_API_KEY")
	insecure := false

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.ApiKey.IsNull() {
		apiKey = data.ApiKey.ValueString()
	}

	if !data.AllowInsecure.IsNull() {
		insecure = data.AllowInsecure.ValueBool()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Unifi Host",
			"The provider cannot create the Unifi API client as there is a missing or empty value for the Unifi host. "+
				"Set the host value in the configuration or use the UNIFI_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Unifi API Key",
			"The provider cannot create the Unifi API client as there is a missing or empty value for the Unifi API key. "+
				"Set the API key value in the configuration or use the UNIFI_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "unifi_host", host)
	ctx = tflog.SetField(ctx, "unifi_api_key", apiKey)
	ctx = tflog.SetField(ctx, "unifi_insecure", insecure)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "unifi_api_key")

	tflog.Debug(ctx, "Creating Unifi client")

	// Create a new unifi client using the configuration values
	client, err := unifi.NewClient(ctx, host, apiKey, insecure)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Unifi API Client",
			"An unexpected error occurred when creating the Unifi API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Unifi Client Error: "+err.Error(),
		)
		return
	}

	// Make the Unifi client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Unifi client", map[string]any{"success": true})
}

func (p *UnifiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccountResource,
		NewApGroupResource,
		NewDeviceResource,
		NewDynamicDnsResource,
		NewFirewallGroupResource,
		NewFirewallRuleResource,
		NewNetworkResource,
		NewPortForwardResource,
		NewPortProfileResource,
		NewRadiusProfileResource,
		NewSettingMgmtResource,
		NewSettingRadiusResource,
		NewSettingUsgResource,
		NewSiteResource,
		NewStaticRouteResource,
		NewUserResource,
		NewUserGroupResource,
		NewWlanResource,
	}
}

func (p *UnifiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountDataSource,
		NewAccountsDataSource,
		NewApGroupDataSource,
		NewApGroupsDataSource,
		NewDeviceDataSource,
		NewDevicesDataSource,
		NewDynamicDnsDataSource,
		NewDynamicDnsesDataSource,
		NewFirewallGroupDataSource,
		NewFirewallGroupsDataSource,
		NewFirewallRuleDataSource,
		NewFirewallRulesDataSource,
		NewNetworkDataSource,
		NewNetworksDataSource,
		NewPortForwardDataSource,
		NewPortForwardsDataSource,
		NewPortProfileDataSource,
		NewPortProfilesDataSource,
		NewRadiusProfileDataSource,
		NewRadiusProfilesDataSource,
		NewSettingMgmtDataSource,
		NewSettingRadiusDataSource,
		NewSettingUsgDataSource,
		NewSiteDataSource,
		NewSitesDataSource,
		NewStaticRouteDataSource,
		NewStaticRoutesDataSource,
		NewUserDataSource,
		NewUsersDataSource,
		NewUserGroupDataSource,
		NewUserGroupsDataSource,
		NewWlanDataSource,
		NewWlansDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UnifiProvider{
			version: version,
		}
	}
}
