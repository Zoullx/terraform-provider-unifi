// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_network"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &networkResource{}
	_ resource.ResourceWithConfigure   = &networkResource{}
	_ resource.ResourceWithImportState = &networkResource{}
)

func NewNetworkResource() resource.Resource {
	return &networkResource{}
}

type networkResource struct {
	client unifi.Client
}

func (r *networkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (r *networkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_network.NetworkResourceSchema(ctx)
}

func (r *networkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nill check when handling ProviderData because Terraform
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

	r.client = client
}

func (r *networkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: site/id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("site"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *networkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_network.NetworkModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.Network
	resp.Diagnostics.Append(parseNetworkResourceModel(ctx, data, &body)...)
	network, err := r.client.CreateNetwork(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Network",
			"Could not create Network, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseNetworkResourceJson(ctx, *network, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *networkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_network.NetworkModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Network value from Unifi
	network, err := r.client.GetNetwork(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Network",
			"Could not read Network ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseNetworkResourceJson(ctx, *network, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *networkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_network.NetworkModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.Network
	resp.Diagnostics.Append(parseNetworkResourceModel(ctx, data, &body)...)
	network, err := r.client.UpdateNetwork(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Network",
			"Could not create Network, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseNetworkResourceJson(ctx, *network, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *networkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_network.NetworkModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Network
	err := r.client.DeleteNetwork(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Network",
			"Could not delete Network, unexpected error: "+err.Error(),
		)
		return
	}
}

func parseNetworkResourceJson(ctx context.Context, json unifi.Network, model *resource_network.NetworkModel) diag.Diagnostics {
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

	var natIpAddresses = []resource_network.NatOutboundIpAddressesValue{}
	for _, ipAddress := range json.NATOutboundIPAddresses {
		var addressPool []types.String
		for _, pool := range ipAddress.IPAddressPool {
			addressPool = append(addressPool, types.StringValue(pool))
		}
		addressPoolList, diags := types.ListValueFrom(ctx, types.StringType, addressPool)
		if diags.HasError() {
			return diags
		}

		natIpAddresses = append(natIpAddresses, resource_network.NatOutboundIpAddressesValue{
			IpAddress:       types.StringValue(ipAddress.IPAddress),
			IpAddressPool:   addressPoolList,
			Mode:            types.StringValue(ipAddress.Mode),
			WanNetworkGroup: types.StringValue(ipAddress.WANNetworkGroup),
		})
	}
	model.NatOutboundIpAddresses, diags = types.ListValueFrom(ctx, resource_network.NatOutboundIpAddressesValue{}.Type(ctx), natIpAddresses)
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

func parseNetworkResourceModel(ctx context.Context, model resource_network.NetworkModel, json *unifi.Network) diag.Diagnostics {
	json.ID = model.Id.ValueString()
	json.SiteID = model.SiteId.ValueString()
	json.Name = model.Name.ValueString()
	json.AutoScaleEnabled = model.AutoScaleEnabled.ValueBool()
	json.DHCPDStart = model.DhcpStart.ValueString()
	json.DHCPDStop = model.DhcpStop.ValueString()
	json.DHCPDEnabled = model.DhcpEnabled.ValueBool()
	json.DHCPDLeaseTime = int(model.DhcpLeaseTime.ValueInt64())

	var dhcpDnsSlice []types.String
	if !model.DhcpDns.IsUnknown() && !model.DhcpDns.IsNull() {
		diags := model.DhcpDns.ElementsAs(ctx, &dhcpDnsSlice, false)
		if diags.HasError() {
			return diags
		}
	}
	json.DHCPDDNS1 = tfStringSliceValueAtIndex(dhcpDnsSlice, 0).ValueString()
	json.DHCPDDNS2 = tfStringSliceValueAtIndex(dhcpDnsSlice, 1).ValueString()
	json.DHCPDDNS3 = tfStringSliceValueAtIndex(dhcpDnsSlice, 2).ValueString()
	json.DHCPDDNS4 = tfStringSliceValueAtIndex(dhcpDnsSlice, 3).ValueString()

	json.DHCPDDNSEnabled = model.DhcpDnsEnabled.ValueBool()
	json.DHCPguardEnabled = model.DhcpGuardEnabled.ValueBool()
	json.DHCPRelayEnabled = model.DhcpRelayEnabled.ValueBool()
	json.DHCPDBootEnabled = model.DhcpBootEnabled.ValueBool()
	json.DHCPDBootServer = model.DhcpBootServer.ValueString()
	json.DHCPDBootFilename = model.DhcpBootFilename.ValueString()
	json.DHCPDConflictChecking = model.DhcpConflictChecking.ValueBool()
	json.DHCPDGatewayEnabled = model.DhcpGatewayEnabled.ValueBool()
	json.DHCPDNtpEnabled = model.DhcpNtpEnabled.ValueBool()
	json.DHCPDTFTPServer = model.DhcpTftpServer.ValueString()
	json.DHCPDTimeOffsetEnabled = model.DhcpTimeOffsetEnabled.ValueBool()
	json.DHCPDUnifiController = model.DhcpUnifiController.ValueString()
	json.DHCPDWinsEnabled = model.DhcpWinsEnabled.ValueBool()
	json.DHCPDWPAdUrl = model.DhcpWpadUrl.ValueString()
	json.DHCPDV6AllowSlaac = model.DhcpV6AllowSlaac.ValueBool()

	var dhcpv6DnsSlice []types.String
	if !model.DhcpV6Dns.IsUnknown() && !model.DhcpV6Dns.IsNull() {
		diags := model.DhcpV6Dns.ElementsAs(ctx, &dhcpv6DnsSlice, false)
		if diags.HasError() {
			return diags
		}
	}
	json.DHCPDV6DNS1 = tfStringSliceValueAtIndex(dhcpv6DnsSlice, 0).ValueString()
	json.DHCPDV6DNS2 = tfStringSliceValueAtIndex(dhcpv6DnsSlice, 1).ValueString()
	json.DHCPDV6DNS3 = tfStringSliceValueAtIndex(dhcpv6DnsSlice, 2).ValueString()
	json.DHCPDV6DNS4 = tfStringSliceValueAtIndex(dhcpv6DnsSlice, 3).ValueString()

	json.DHCPDV6DNSAuto = model.DhcpV6DnsAuto.ValueBool()
	json.DHCPDV6Enabled = model.DhcpV6Enabled.ValueBool()
	json.DHCPDV6LeaseTime = int(model.DhcpV6LeaseTime.ValueInt64())
	json.DHCPDV6Start = model.DhcpV6Start.ValueString()
	json.DHCPDV6Stop = model.DhcpV6Stop.ValueString()
	json.DomainName = model.DomainName.ValueString()
	json.Enabled = model.Enabled.ValueBool()
	json.GatewayType = model.GatewayType.ValueString()
	json.IGMPSnooping = model.IgmpSnooping.ValueBool()
	json.InternetAccessEnabled = model.InternetAccessEnabled.ValueBool()
	json.IPV6PDAutoPrefixidEnabled = model.Ipv6PdAutoPrefixidEnabled.ValueBool()
	json.IPV6ClientAddressAssignment = model.Ipv6ClientAddressAssignment.ValueString()
	json.IPV6Enabled = model.Ipv6Enabled.ValueBool()
	json.IPV6InterfaceType = model.Ipv6InterfaceType.ValueString()
	json.IPV6Subnet = model.Ipv6StaticSubnet.ValueString()
	json.IPV6PDInterface = model.Ipv6PdInterface.ValueString()
	json.IPV6PDPrefixid = model.Ipv6PdPrefixid.ValueString()
	json.IPV6PDStart = model.Ipv6PdStart.ValueString()
	json.IPV6PDStop = model.Ipv6PdStop.ValueString()
	json.IPV6RaEnabled = model.Ipv6RaEnabled.ValueBool()
	json.IPV6RaPreferredLifetime = int(model.Ipv6RaPreferredLifetime.ValueInt64())
	json.IPV6RaPriority = model.Ipv6RaPriority.ValueString()
	json.IPV6RaValidLifetime = int(model.Ipv6RaValidLifetime.ValueInt64())
	json.SettingPreference = model.SettingPreference.ValueString()
	json.LteLanEnabled = model.LteLanEnabled.ValueBool()
	json.MdnsEnabled = model.MulticastDnsEnabled.ValueBool()

	var natIpAddresses []unifi.NetworkNATOutboundIPAddresses
	if !model.NatOutboundIpAddresses.IsUnknown() && !model.NatOutboundIpAddresses.IsNull() {
		diags := model.NatOutboundIpAddresses.ElementsAs(ctx, natIpAddresses, false)
		if diags.HasError() {
			return diags
		}
	}
	json.NATOutboundIPAddresses = natIpAddresses

	json.NetworkGroup = model.NetworkGroup.ValueString()
	json.NetworkIsolationEnabled = model.NetworkIsolationEnabled.ValueBool()
	json.Purpose = model.Purpose.ValueString()
	json.SettingPreference = model.SettingPreference.ValueString()
	json.IPSubnet = model.Subnet.ValueString()
	json.UpnpLanEnabled = model.UpnpLanEnabled.ValueBool()
	json.VLAN = int(model.VlanId.ValueInt64())
	json.VLANEnabled = model.VlanEnabled.ValueBool()
	json.WANIP = model.WanIp.ValueString()
	json.WANNetmask = model.WanNetmask.ValueString()
	json.WANGateway = model.WanGateway.ValueString()

	var wanDnsSlice []types.String
	if !model.WanDns.IsUnknown() && !model.WanDns.IsNull() {
		diags := model.WanDns.ElementsAs(ctx, &wanDnsSlice, false)
		if diags.HasError() {
			return diags
		}
	}
	json.WANDNS1 = tfStringSliceValueAtIndex(wanDnsSlice, 0).ValueString()
	json.WANDNS2 = tfStringSliceValueAtIndex(wanDnsSlice, 1).ValueString()
	json.WANDNS3 = tfStringSliceValueAtIndex(wanDnsSlice, 2).ValueString()
	json.WANDNS4 = tfStringSliceValueAtIndex(wanDnsSlice, 3).ValueString()

	json.WANType = model.WanType.ValueString()
	json.WANNetworkGroup = model.WanNetworkGroup.ValueString()
	json.WANEgressQOS = int(model.WanEgressQos.ValueInt64())
	json.WANUsername = model.WanUsername.ValueString()
	json.XWANPassword = model.WanPassword.ValueString()
	json.WANTypeV6 = model.WanTypeV6.ValueString()
	json.WANDHCPv6PDSize = int(model.WanDhcpV6PdSize.ValueInt64())
	json.WANIPV6 = model.WanIpv6.ValueString()
	json.WANGatewayV6 = model.WanGatewayV6.ValueString()
	json.WANPrefixlen = int(model.WanPrefixlen.ValueInt64())

	return nil
}

func tfStringSliceValueAtIndex(s []types.String, i int) types.String {
	if s != nil && i <= (len(s)-1) {
		return s[i]
	}
	return types.StringValue("")
}
