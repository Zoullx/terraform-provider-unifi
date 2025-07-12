package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_wlan"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &wlanResource{}
	_ resource.ResourceWithConfigure   = &wlanResource{}
	_ resource.ResourceWithImportState = &wlanResource{}
)

func NewWlanResource() resource.Resource {
	return &wlanResource{}
}

type wlanResource struct {
	client unifi.Client
}

func (r *wlanResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wlan"
}

func (r *wlanResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_wlan.WlanResourceSchema(ctx)
}

func (r *wlanResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *wlanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *wlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_wlan.WlanModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.WLAN
	resp.Diagnostics.Append(parseWlanResourceModel(ctx, data, &body)...)
	wlan, err := r.client.CreateWLAN(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating WLAN",
			"Could not create WLAN, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseWlanResourceJson(ctx, *wlan, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *wlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_wlan.WlanModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed WLAN value from Unifi
	wlan, err := r.client.GetWLAN(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading WLAN",
			"Could not read WLAN ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseWlanResourceJson(ctx, *wlan, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *wlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_wlan.WlanModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.WLAN
	resp.Diagnostics.Append(parseWlanResourceModel(ctx, data, &body)...)
	wlan, err := r.client.UpdateWLAN(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating WLAN",
			"Could not create WLAN, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseWlanResourceJson(ctx, *wlan, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *wlanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_wlan.WlanModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing WLAN
	err := r.client.DeleteWLAN(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting WLAN",
			"Could not delete WLAN, unexpected error: "+err.Error(),
		)
		return
	}
}

func parseWlanResourceJson(ctx context.Context, json unifi.WLAN, model *resource_wlan.WlanModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)
	model.SiteId = types.StringValue(json.SiteID)

	apGroupIdList, diags := types.ListValueFrom(ctx, types.StringType, json.ApGroupIDs)
	if diags.HasError() {
		return diags
	}
	model.ApGroupIds = apGroupIdList

	model.ApGroupMode = types.StringValue(json.ApGroupMode)
	model.BSupported = types.BoolValue(json.BSupported)

	broadcastFilterList, diags := types.ListValueFrom(ctx, types.StringType, json.BroadcastFilterList)
	if diags.HasError() {
		return diags
	}
	model.BroadcastFilterList = broadcastFilterList

	model.BssTransition = types.BoolValue(json.BssTransition)
	model.Dtim6e = types.Int64Value(int64(json.DTIM6E))
	model.Dtim2g = types.Int64Value(int64(json.DTIMNg))
	model.Dtim5g = types.Int64Value(int64(json.DTIMNa))
	model.DtimMode = types.StringValue(json.DTIMMode)
	model.Enabled = types.BoolValue(json.Enabled)
	model.EnhancedIot = types.BoolValue(json.EnhancedIot)
	model.FastRoamingEnabled = types.BoolValue(json.FastRoamingEnabled)
	model.GroupRekey = types.Int64Value(int64(json.GroupRekey))
	model.HideSsid = types.BoolValue(json.HideSSID)
	model.Hotspot2confEnabled = types.BoolValue(json.Hotspot2ConfEnabled)
	model.IsGuest = types.BoolValue(json.IsGuest)
	model.IappEnabled = types.BoolValue(json.IappEnabled)
	model.IappKey = types.StringValue(json.XIappKey)
	model.L2Isolation = types.BoolValue(json.L2Isolation)
	model.MacFilterEnabled = types.BoolValue(json.MACFilterEnabled)

	model.MacFilterList, diags = types.ListValueFrom(ctx, types.StringType, json.MACFilterList)
	if diags.HasError() {
		return diags
	}

	model.MacFilterPolicy = types.StringValue(json.MACFilterPolicy)
	model.Minimum2gAdvertisingRates = types.BoolValue(json.MinrateNgAdvertisingRates)
	model.Minimum2gDataRateEnabled = types.BoolValue(json.MinrateNgEnabled)
	model.Minimum2gDataRateKbps = types.Int64Value(int64(json.MinrateNgDataRateKbps))
	model.Minimum5gAdvertisingRates = types.BoolValue(json.MinrateNaAdvertisingRates)
	model.Minimum5gDataRateEnabled = types.BoolValue(json.MinrateNaEnabled)
	model.Minimum5gDataRateKbps = types.Int64Value(int64(json.MinrateNaDataRateKbps))
	model.MinimumDataRateSettingPreference = types.StringValue(json.MinrateSettingPreference)
	model.MloEnabled = types.BoolValue(json.MloEnabled)
	model.MulticastEnhanceEnabled = types.BoolValue(json.MulticastEnhanceEnabled)
	model.Name = types.StringValue(json.Name)
	model.NetworkId = types.StringValue(json.NetworkID)
	model.No2ghzOui = types.BoolValue(json.No2GhzOui)
	model.OptimizeIotWifiConnectivity = types.BoolValue(json.OptimizeIotWifiConnectivity)
	model.Passphrase = types.StringValue(json.XPassphrase)
	model.PassphraseAutogenerated = types.BoolValue(json.PassphraseAutogenerated)
	model.PmfMode = types.StringValue(json.PMFMode)

	privatePresharedKeyList, diags := types.ListValueFrom(ctx, resource_wlan.PrivatePresharedKeysValue{}.Type(ctx), json.PrivatePresharedKeys)
	if diags.HasError() {
		return diags
	}
	model.PrivatePresharedKeys = privatePresharedKeyList

	model.PrivatePresharedKeysEnabled = types.BoolValue(json.PrivatePresharedKeysEnabled)
	model.ProxyArp = types.BoolValue(json.ProxyArp)
	model.RadiusDasEnabled = types.BoolValue(json.RADIUSDasEnabled)
	model.RadiusMacAuthEnabled = types.BoolValue(json.RADIUSMACAuthEnabled)
	model.RadiusMacAclFormat = types.StringValue(json.RADIUSMACaclFormat)
	model.RadiusProfileId = types.StringValue(json.RADIUSProfileID)

	saeGroupList, diags := types.ListValueFrom(ctx, types.Int64Type, json.SaeGroups)
	if diags.HasError() {
		return diags
	}
	model.SaeGroups = saeGroupList

	saePskList, diags := types.ListValueFrom(ctx, resource_wlan.SaePsksValue{}.Type(ctx), json.SaePsk)
	if diags.HasError() {
		return diags
	}
	model.SaePsks = saePskList

	model.Schedule, diags = types.ListValueFrom(ctx, resource_wlan.ScheduleValue{}.Type(ctx), json.Schedule)
	if diags.HasError() {
		return diags
	}

	model.Security = types.StringValue(json.Security)
	model.SettingPreference = types.StringValue(json.SettingPreference)
	model.UapsdEnabled = types.BoolValue(json.UapsdEnabled)
	model.UserGroupId = types.StringValue(json.UserGroupID)
	model.WlanBand = types.StringValue(json.WLANBand)

	wlanBandList, diags := types.ListValueFrom(ctx, types.StringType, json.WLANBands)
	if diags.HasError() {
		return diags
	}
	model.WlanBands = wlanBandList

	model.WpaEnc = types.StringValue(json.WPAEnc)
	model.WpaMode = types.StringValue(json.WPAMode)
	model.Wpa3Enhanced192 = types.BoolValue(json.WPA3Enhanced192)
	model.Wpa3FastRoaming = types.BoolValue(json.WPA3FastRoaming)
	model.Wpa3Support = types.BoolValue(json.WPA3Support)
	model.Wpa3Transition = types.BoolValue(json.WPA3Transition)

	return nil
}

func parseWlanResourceModel(ctx context.Context, model resource_wlan.WlanModel, json *unifi.WLAN) diag.Diagnostics {
	json.ID = model.Id.ValueString()
	json.SiteID = model.SiteId.ValueString()

	if !model.ApGroupIds.IsUnknown() && !model.ApGroupIds.IsNull() {
		diags := model.ApGroupIds.ElementsAs(ctx, &json.ApGroupIDs, false)
		if diags.HasError() {
			return diags
		}
	}

	json.ApGroupMode = model.ApGroupMode.ValueString()
	json.BSupported = model.BSupported.ValueBool()

	if !model.BroadcastFilterList.IsUnknown() && !model.BroadcastFilterList.IsNull() {
		diags := model.BroadcastFilterList.ElementsAs(ctx, &json.BroadcastFilterList, false)
		if diags.HasError() {
			return diags
		}
	}

	json.BssTransition = model.BssTransition.ValueBool()
	json.DTIM6E = int(model.Dtim6e.ValueInt64())
	json.DTIMNg = int(model.Dtim2g.ValueInt64())
	json.DTIMNa = int(model.Dtim5g.ValueInt64())
	json.DTIMMode = model.DtimMode.ValueString()
	json.Enabled = model.Enabled.ValueBool()
	json.EnhancedIot = model.EnhancedIot.ValueBool()
	json.FastRoamingEnabled = model.FastRoamingEnabled.ValueBool()
	json.GroupRekey = int(model.GroupRekey.ValueInt64())
	json.HideSSID = model.HideSsid.ValueBool()
	json.Hotspot2ConfEnabled = model.Hotspot2confEnabled.ValueBool()
	json.IsGuest = model.IsGuest.ValueBool()
	json.IappEnabled = model.IappEnabled.ValueBool()
	json.XIappKey = model.IappKey.ValueString()
	json.L2Isolation = model.L2Isolation.ValueBool()
	json.MACFilterEnabled = model.MacFilterEnabled.ValueBool()

	if !model.MacFilterList.IsUnknown() && !model.MacFilterList.IsNull() {
		diags := model.MacFilterList.ElementsAs(ctx, &json.MACFilterList, false)
		if diags.HasError() {
			return diags
		}
	}

	json.MACFilterPolicy = model.MacFilterPolicy.ValueString()
	json.MinrateNgAdvertisingRates = model.Minimum2gAdvertisingRates.ValueBool()
	json.MinrateNgEnabled = model.Minimum2gDataRateEnabled.ValueBool()
	json.MinrateNgDataRateKbps = int(model.Minimum2gDataRateKbps.ValueInt64())
	json.MinrateNaAdvertisingRates = model.Minimum5gAdvertisingRates.ValueBool()
	json.MinrateNaEnabled = model.Minimum5gDataRateEnabled.ValueBool()
	json.MinrateNaDataRateKbps = int(model.Minimum5gDataRateKbps.ValueInt64())
	json.MinrateSettingPreference = model.MinimumDataRateSettingPreference.ValueString()
	json.MloEnabled = model.MloEnabled.ValueBool()
	json.MulticastEnhanceEnabled = model.MulticastEnhanceEnabled.ValueBool()
	json.Name = model.Name.ValueString()
	json.NetworkID = model.NetworkId.ValueString()
	json.No2GhzOui = model.No2ghzOui.ValueBool()
	json.OptimizeIotWifiConnectivity = model.OptimizeIotWifiConnectivity.ValueBool()
	json.XPassphrase = model.Passphrase.ValueString()
	json.PassphraseAutogenerated = model.PassphraseAutogenerated.ValueBool()
	json.PMFMode = model.PmfMode.ValueString()

	if !model.PrivatePresharedKeys.IsUnknown() && !model.PrivatePresharedKeys.IsNull() {
		diags := model.PrivatePresharedKeys.ElementsAs(ctx, &json.PrivatePresharedKeys, false)
		if diags.HasError() {
			return diags
		}
	}

	json.PrivatePresharedKeysEnabled = model.PrivatePresharedKeysEnabled.ValueBool()
	json.ProxyArp = model.ProxyArp.ValueBool()
	json.RADIUSDasEnabled = model.RadiusDasEnabled.ValueBool()
	json.RADIUSMACAuthEnabled = model.RadiusMacAuthEnabled.ValueBool()
	json.RADIUSMACaclFormat = model.RadiusMacAclFormat.ValueString()
	json.RADIUSProfileID = model.RadiusProfileId.ValueString()

	if !model.SaeGroups.IsUnknown() && !model.SaeGroups.IsNull() {
		diags := model.SaeGroups.ElementsAs(ctx, &json.SaeGroups, false)
		if diags.HasError() {
			return diags
		}
	}

	if !model.SaePsks.IsUnknown() && !model.SaePsks.IsNull() {
		diags := model.SaePsks.ElementsAs(ctx, &json.SaePsk, false)
		if diags.HasError() {
			return diags
		}
	}

	if !model.Schedule.IsUnknown() && !model.Schedule.IsNull() {
		diags := model.Schedule.ElementsAs(ctx, &json.Schedule, false)
		if diags.HasError() {
			return diags
		}
	}

	json.Security = model.Security.ValueString()
	json.SettingPreference = model.SettingPreference.ValueString()
	json.UapsdEnabled = model.UapsdEnabled.ValueBool()
	json.UserGroupID = model.UserGroupId.ValueString()
	json.WLANBand = model.WlanBand.ValueString()

	if !model.WlanBands.IsUnknown() && !model.WlanBands.IsNull() {
		diags := model.WlanBands.ElementsAs(ctx, &json.WLANBands, false)
		if diags.HasError() {
			return diags
		}
	}

	json.WPAEnc = model.WpaEnc.ValueString()
	json.WPAMode = model.WpaMode.ValueString()
	json.WPA3Enhanced192 = model.Wpa3Enhanced192.ValueBool()
	json.WPA3FastRoaming = model.Wpa3FastRoaming.ValueBool()
	json.WPA3Support = model.Wpa3Support.ValueBool()
	json.WPA3Transition = model.Wpa3Transition.ValueBool()

	return nil
}
