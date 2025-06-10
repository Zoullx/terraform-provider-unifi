package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_firewall_rule"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &firewallRuleDataSource{}
)

func NewFirewallRuleDataSource() datasource.DataSource {
	return &firewallRuleDataSource{}
}

type firewallRuleDataSource struct {
	client *unifi.Client
}

func (d *firewallRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_rule"
}

func (d *firewallRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_firewall_rule.FirewallRuleDataSourceSchema(ctx)
}

func (d *firewallRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*unifi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *unifi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *firewallRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_firewall_rule.FirewallRuleModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.Id.IsNull() || data.Id.IsUnknown()) && (data.Name.IsNull() || data.Name.IsUnknown()) {
		resp.Diagnostics.AddError(
			"Id and Name are null or unknown",
			"Id or Name is required to retrieve a Firewall Rule.",
		)
		return
	}

	// Get Firewall Rule
	var firewallRule *unifi.FirewallRule
	var err error
	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		firewallRule, err = d.client.GetFirewallRule(ctx, data.Site.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() && !data.Name.IsUnknown() {
		firewallRule, err = d.client.GetFirewallRuleByName(ctx, data.Site.ValueString(), data.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Firewall Rule",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseFirewallRuleDataSourceJson(ctx, *firewallRule, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseFirewallRuleDataSourceJson(ctx context.Context, json unifi.FirewallRule, model *datasource_firewall_rule.FirewallRuleModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)
	model.SiteId = types.StringValue(json.SiteID)
	model.Action = types.StringValue(json.Action)
	model.DstAddress = types.StringValue(json.DstAddress)
	model.DstAddressIpv6 = types.StringValue(json.DstAddressIPV6)

	dstFirewallGroupIdList, diags := types.ListValueFrom(ctx, types.StringType, json.DstFirewallGroupIDs)
	if diags.HasError() {
		return diags
	}
	model.DstFirewallGroupIds = dstFirewallGroupIdList

	model.DstNetworkId = types.StringValue(json.DstNetworkID)
	model.DstNetworkType = types.StringValue(json.DstNetworkType)
	model.DstPort = types.StringValue(json.DstPort)
	model.Enabled = types.BoolValue(json.Enabled)
	model.IcmpTypename = types.StringValue(json.ICMPTypename)
	model.IcmpV6Typename = types.StringValue(json.ICMPv6Typename)
	model.IpSec = types.StringValue(json.IPSec)
	model.Logging = types.BoolValue(json.Logging)
	model.Name = types.StringValue(json.Name)
	model.Protocol = types.StringValue(json.Protocol)
	model.ProtocolV6 = types.StringValue(json.ProtocolV6)
	model.ProtocolMatchExcepted = types.BoolValue(json.ProtocolMatchExcepted)
	model.RuleIndex = types.Int64Value(int64(json.RuleIndex))
	model.Ruleset = types.StringValue(json.Ruleset)
	model.SettingPreference = types.StringValue(json.SettingPreference)
	model.SrcAddress = types.StringValue(json.SrcAddress)
	model.SrcAddressIpv6 = types.StringValue(json.SrcAddressIPV6)

	srcFirewallGroupIdList, diags := types.ListValueFrom(ctx, types.StringType, json.SrcFirewallGroupIDs)
	if diags.HasError() {
		return diags
	}
	model.SrcFirewallGroupIds = srcFirewallGroupIdList

	model.SrcMac = types.StringValue(json.SrcMACAddress)
	model.SrcNetworkId = types.StringValue(json.SrcNetworkID)
	model.SrcNetworkType = types.StringValue(json.SrcNetworkType)
	model.SrcPort = types.StringValue(json.SrcPort)
	model.StateEstablished = types.BoolValue(json.StateEstablished)
	model.StateInvalid = types.BoolValue(json.StateInvalid)
	model.StateNew = types.BoolValue(json.StateNew)
	model.StateRelated = types.BoolValue(json.StateRelated)

	return nil
}
