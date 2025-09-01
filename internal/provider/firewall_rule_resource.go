// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_firewall_rule"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &firewallRuleResource{}
	_ resource.ResourceWithConfigure   = &firewallRuleResource{}
	_ resource.ResourceWithImportState = &firewallRuleResource{}
)

func NewFirewallRuleResource() resource.Resource {
	return &firewallRuleResource{}
}

type firewallRuleResource struct {
	client unifi.Client
}

func (r *firewallRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_rule"
}

func (r *firewallRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_firewall_rule.FirewallRuleResourceSchema(ctx)
}

func (r *firewallRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *firewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *firewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_firewall_rule.FirewallRuleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.FirewallRule
	resp.Diagnostics.Append(parseFirewallRuleResourceModel(ctx, data, &body)...)
	firewallRule, err := r.client.CreateFirewallRule(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Firewall Rule",
			"Could not create Firewall Rule, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseFirewallRuleResourceJson(ctx, *firewallRule, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_firewall_rule.FirewallRuleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Firewall Rule value from Unifi
	firewallRule, err := r.client.GetFirewallRule(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Firewall Rule",
			"Could not read Firewall Rule ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseFirewallRuleResourceJson(ctx, *firewallRule, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_firewall_rule.FirewallRuleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.FirewallRule
	resp.Diagnostics.Append(parseFirewallRuleResourceModel(ctx, data, &body)...)
	firewallRule, err := r.client.UpdateFirewallRule(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Firewall Rule",
			"Could not create Firewall Rule, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseFirewallRuleResourceJson(ctx, *firewallRule, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_firewall_rule.FirewallRuleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Firewall Rule
	err := r.client.DeleteFirewallRule(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Firewall Rule",
			"Could not delete Firewall Rule, unexpected error: "+err.Error(),
		)
		return
	}
}

func parseFirewallRuleResourceJson(ctx context.Context, json unifi.FirewallRule, model *resource_firewall_rule.FirewallRuleModel) diag.Diagnostics {
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

func parseFirewallRuleResourceModel(ctx context.Context, model resource_firewall_rule.FirewallRuleModel, json *unifi.FirewallRule) diag.Diagnostics {
	json.ID = model.Id.ValueString()
	json.SiteID = model.SiteId.ValueString()
	json.Action = model.Action.ValueString()
	json.DstAddress = model.DstAddress.ValueString()
	json.DstAddressIPV6 = model.DstAddressIpv6.ValueString()

	if !model.DstFirewallGroupIds.IsUnknown() && !model.DstFirewallGroupIds.IsNull() {
		diags := model.DstFirewallGroupIds.ElementsAs(ctx, &json.DstFirewallGroupIDs, false)
		if diags.HasError() {
			return diags
		}
	}

	json.DstNetworkID = model.DstNetworkId.ValueString()
	json.DstNetworkType = model.DstNetworkType.ValueString()
	json.DstPort = model.DstPort.ValueString()
	json.Enabled = model.Enabled.ValueBool()
	json.ICMPTypename = model.IcmpTypename.ValueString()
	json.ICMPv6Typename = model.IcmpV6Typename.ValueString()
	json.IPSec = model.IpSec.ValueString()
	json.Logging = model.Logging.ValueBool()
	json.Name = model.Name.ValueString()
	json.Protocol = model.Protocol.ValueString()
	json.ProtocolV6 = model.ProtocolV6.ValueString()
	json.ProtocolMatchExcepted = model.ProtocolMatchExcepted.ValueBool()
	json.RuleIndex = int(model.RuleIndex.ValueInt64())
	json.Ruleset = model.Ruleset.ValueString()
	json.SettingPreference = model.SettingPreference.ValueString()
	json.SrcAddress = model.SrcAddress.ValueString()
	json.SrcAddressIPV6 = model.SrcAddressIpv6.ValueString()

	if !model.SrcFirewallGroupIds.IsUnknown() && !model.SrcFirewallGroupIds.IsNull() {
		diags := model.SrcFirewallGroupIds.ElementsAs(ctx, &json.SrcFirewallGroupIDs, false)
		if diags.HasError() {
			return diags
		}
	}

	json.SrcMACAddress = model.SrcMac.ValueString()
	json.SrcNetworkID = model.SrcNetworkId.ValueString()
	json.SrcNetworkType = model.SrcNetworkType.ValueString()
	json.SrcPort = model.SrcPort.ValueString()
	json.StateEstablished = model.StateEstablished.ValueBool()
	json.StateInvalid = model.StateInvalid.ValueBool()
	json.StateNew = model.StateNew.ValueBool()
	json.StateRelated = model.StateRelated.ValueBool()

	return nil
}
