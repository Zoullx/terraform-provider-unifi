// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_network"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &networkDataSource{}
)

func NewNetworkDataSource() datasource.DataSource {
	return &networkDataSource{}
}

type networkDataSource struct {
	client unifi.Client
}

func (d *networkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (d *networkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_network.NetworkDataSourceSchema(ctx)
}

func (d *networkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *networkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_network.NetworkModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Network
	network, err := d.client.GetNetwork(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Network",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseNetworkDataSourceJson(ctx, *network, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseNetworkDataSourceJson(ctx context.Context, json unifi.Network, model *datasource_network.NetworkModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)
	model.SiteId = types.StringValue(json.SiteID)
	model.Name = types.StringValue(json.Name)
	model.AutoScaleEnabled = types.BoolValue(json.AutoScaleEnabled)
	model.DhcpStart = types.StringValue(json.DHCPDStart)
	model.DhcpStop = types.StringValue(json.DHCPDStop)
	model.DhcpEnabled = types.BoolValue(json.DHCPDEnabled)
	model.DhcpLeaseTime = types.Int64Value(int64(json.DHCPDLeaseTime))

	var dhcpDnsSlice []types.String
	for _, dns := range []string{
		json.DHCPDDNS1,
		json.DHCPDDNS2,
		json.DHCPDDNS3,
		json.DHCPDDNS4,
	} {
		if dns != "" {
			dhcpDnsSlice = append(dhcpDnsSlice, types.StringValue(dns))
		}
	}
	dhcpDnsList, diags := types.ListValueFrom(ctx, types.StringType, dhcpDnsSlice)
	if diags.HasError() {
		return diags
	}
	model.DhcpDns = dhcpDnsList

	model.DhcpDnsEnabled = types.BoolValue(json.DHCPDEnabled)
	model.DhcpGuardEnabled = types.BoolValue(json.DHCPguardEnabled)
	model.DhcpRelayEnabled = types.BoolValue(json.DHCPRelayEnabled)
	model.DhcpBootEnabled = types.BoolValue(json.DHCPDBootEnabled)
	model.DhcpBootServer = types.StringValue(json.DHCPDBootServer)
	model.DhcpBootFilename = types.StringValue(json.DHCPDBootFilename)
	model.DhcpConflictChecking = types.BoolValue(json.DHCPDConflictChecking)
	model.DhcpGatewayEnabled = types.BoolValue(json.DHCPDGatewayEnabled)
	model.DhcpNtpEnabled = types.BoolValue(json.DHCPDNtpEnabled)
	model.DhcpTftpServer = types.StringValue(json.DHCPDTFTPServer)
	model.DhcpTimeOffsetEnabled = types.BoolValue(json.DHCPDTimeOffsetEnabled)
	model.DhcpUnifiController = types.StringValue(json.DHCPDUnifiController)
	model.DhcpWinsEnabled = types.BoolValue(json.DHCPDWinsEnabled)
	model.DhcpWpadUrl = types.StringValue(json.DHCPDWPAdUrl)
	model.DhcpV6AllowSlaac = types.BoolValue(json.DHCPDV6AllowSlaac)

	var dhcpDnsv6Slice []types.String
	for _, dns := range []string{
		json.DHCPDV6DNS1,
		json.DHCPDV6DNS2,
		json.DHCPDV6DNS3,
		json.DHCPDV6DNS4,
	} {
		if dns != "" {
			dhcpDnsv6Slice = append(dhcpDnsv6Slice, types.StringValue(dns))
		}
	}
	model.DhcpV6Dns, diags = types.ListValueFrom(ctx, types.StringType, dhcpDnsv6Slice)
	if diags.HasError() {
		return diags
	}

	model.DhcpV6DnsAuto = types.BoolValue(json.DHCPDV6DNSAuto)
	model.DhcpV6Enabled = types.BoolValue(json.DHCPDV6Enabled)
	model.DhcpV6LeaseTime = types.Int64Value(int64(json.DHCPDV6LeaseTime))
	model.DhcpV6Start = types.StringValue(json.DHCPDV6Start)
	model.DhcpV6Stop = types.StringValue(json.DHCPDV6Stop)
	model.DomainName = types.StringValue(json.DomainName)
	model.Enabled = types.BoolValue(json.Enabled)
	model.GatewayType = types.StringValue(json.GatewayType)
	model.IgmpSnooping = types.BoolValue(json.IGMPSnooping)
	model.InternetAccessEnabled = types.BoolValue(json.InternetAccessEnabled)
	model.Ipv6PdAutoPrefixidEnabled = types.BoolValue(json.IPV6PDAutoPrefixidEnabled)
	model.Ipv6ClientAddressAssignment = types.StringValue(json.IPV6ClientAddressAssignment)
	model.Ipv6Enabled = types.BoolValue(json.IPV6Enabled)
	model.Ipv6InterfaceType = types.StringValue(json.IPV6InterfaceType)
	model.Ipv6StaticSubnet = types.StringValue(json.IPV6Subnet)
	model.Ipv6PdInterface = types.StringValue(json.IPV6PDInterface)
	model.Ipv6PdPrefixid = types.StringValue(json.IPV6PDPrefixid)
	model.Ipv6PdStart = types.StringValue(json.IPV6PDStart)
	model.Ipv6PdStop = types.StringValue(json.IPV6PDStop)
	model.Ipv6RaEnabled = types.BoolValue(json.IPV6RaEnabled)
	model.Ipv6RaPreferredLifetime = types.Int64Value(int64(json.IPV6RaPreferredLifetime))
	model.Ipv6RaPriority = types.StringValue(json.IPV6RaPriority)
	model.Ipv6RaValidLifetime = types.Int64Value(int64(json.IPV6RaValidLifetime))
	model.Ipv6SettingPreference = types.StringValue(json.IPV6SettingPreference)
	model.LteLanEnabled = types.BoolValue(json.LteLanEnabled)
	model.MulticastDnsEnabled = types.BoolValue(json.MdnsEnabled)

	var natIpAddresses = []datasource_network.NatOutboundIpAddressesValue{}
	for _, ipAddress := range json.NATOutboundIPAddresses {
		var addressPool []types.String
		for _, pool := range ipAddress.IPAddressPool {
			addressPool = append(addressPool, types.StringValue(pool))
		}
		addressPoolList, diags := types.ListValueFrom(ctx, types.StringType, addressPool)
		if diags.HasError() {
			return diags
		}

		natIpAddresses = append(natIpAddresses, datasource_network.NatOutboundIpAddressesValue{
			IpAddress:       types.StringValue(ipAddress.IPAddress),
			IpAddressPool:   addressPoolList,
			Mode:            types.StringValue(ipAddress.Mode),
			WanNetworkGroup: types.StringValue(ipAddress.WANNetworkGroup),
		})
	}
	model.NatOutboundIpAddresses, diags = types.ListValueFrom(ctx, datasource_network.NatOutboundIpAddressesValue{}.Type(ctx), natIpAddresses)
	if diags.HasError() {
		return diags
	}

	model.NetworkGroup = types.StringValue(json.NetworkGroup)
	model.NetworkIsolationEnabled = types.BoolValue(json.NetworkIsolationEnabled)
	model.Purpose = types.StringValue(json.Purpose)
	model.SettingPreference = types.StringValue(json.SettingPreference)
	model.Subnet = types.StringValue(json.IPSubnet)
	model.UpnpLanEnabled = types.BoolValue(json.UpnpLanEnabled)
	model.VlanId = types.Int64Value(int64(json.VLAN))
	model.VlanEnabled = types.BoolValue(json.VLANEnabled)
	model.WanIp = types.StringValue(json.WANIP)
	model.WanNetmask = types.StringValue(json.WANNetmask)
	model.WanGateway = types.StringValue(json.WANGateway)

	var wanDnsSlice []types.String
	for _, dns := range []string{
		json.WANDNS1,
		json.WANDNS2,
		json.WANDNS3,
		json.WANDNS4,
	} {
		if dns != "" {
			wanDnsSlice = append(wanDnsSlice, types.StringValue(dns))
		}
	}
	model.WanDns, diags = types.ListValueFrom(ctx, types.StringType, wanDnsSlice)
	if diags.HasError() {
		return diags
	}

	model.WanType = types.StringValue(json.WANType)
	model.WanNetworkGroup = types.StringValue(json.WANNetworkGroup)
	model.WanEgressQos = types.Int64Value(int64(json.WANEgressQOS))
	model.WanUsername = types.StringValue(json.WANUsername)
	model.WanPassword = types.StringValue(json.XWANPassword)
	model.WanTypeV6 = types.StringValue(json.WANTypeV6)
	model.WanDhcpV6PdSize = types.Int64Value(int64(json.WANDHCPv6PDSize))
	model.WanIpv6 = types.StringValue(json.WANIPV6)
	model.WanGatewayV6 = types.StringValue(json.WANGatewayV6)
	model.WanPrefixlen = types.Int64Value(int64(json.WANPrefixlen))

	return nil
}
