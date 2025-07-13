package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_setting_usg"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &settingUsgDataSource{}
)

func NewSettingUsgDataSource() datasource.DataSource {
	return &settingUsgDataSource{}
}

type settingUsgDataSource struct {
	client unifi.Client
}

func (d *settingUsgDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting_usg"
}

func (d *settingUsgDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_setting_usg.SettingUsgDataSourceSchema(ctx)
}

func (d *settingUsgDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(unifi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *unifi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *settingUsgDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_setting_usg.SettingUsgModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Setting USG
	settingUsg, err := d.client.GetSettingUsg(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Setting USG",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseSettingUsgDataSourceJson(ctx, *settingUsg, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseSettingUsgDataSourceJson(ctx context.Context, json unifi.SettingUsg, model *datasource_setting_usg.SettingUsgModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)
	model.SiteId = types.StringValue(json.SiteID)

	var dhcpRelayServerSlice []types.String
	for _, server := range []string{
		json.DHCPRelayServer1,
		json.DHCPRelayServer2,
		json.DHCPRelayServer3,
		json.DHCPRelayServer4,
		json.DHCPRelayServer5,
	} {
		if server != "" {
			dhcpRelayServerSlice = append(dhcpRelayServerSlice, types.StringValue(server))
		}
	}
	dhcpRelayServerList, diags := types.ListValueFrom(ctx, types.StringType, dhcpRelayServerSlice)
	if diags.HasError() {
		return diags
	}
	model.DhcpRelayServers = dhcpRelayServerList

	model.MulticastDnsEnabled = types.BoolValue(json.MdnsEnabled)

	return nil
}
