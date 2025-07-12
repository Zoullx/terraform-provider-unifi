package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_setting_usg"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &settingUsgResource{}
	_ resource.ResourceWithConfigure   = &settingUsgResource{}
	_ resource.ResourceWithImportState = &settingUsgResource{}
)

func NewSettingUsgResource() resource.Resource {
	return &settingUsgResource{}
}

type settingUsgResource struct {
	client unifi.Client
}

func (r *settingUsgResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting_usg"
}

func (r *settingUsgResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_setting_usg.SettingUsgResourceSchema(ctx)
}

func (r *settingUsgResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingUsgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *settingUsgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_setting_usg.SettingUsgModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.SettingUsg
	resp.Diagnostics.Append(parseSettingUsgResourceModel(ctx, data, &body)...)
	settingUsg, err := r.client.UpdateSettingUsg(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Setting USG",
			"Could not create Setting USG, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseSettingUsgResourceJson(ctx, *settingUsg, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *settingUsgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_setting_usg.SettingUsgModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Setting USG value from Unifi
	settingUsg, err := r.client.GetSettingUsg(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Setting USG",
			"Could not read Setting USG ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseSettingUsgResourceJson(ctx, *settingUsg, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *settingUsgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_setting_usg.SettingUsgModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.SettingUsg
	resp.Diagnostics.Append(parseSettingUsgResourceModel(ctx, data, &body)...)
	settingUsg, err := r.client.UpdateSettingUsg(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Setting USG",
			"Could not create Setting USG, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseSettingUsgResourceJson(ctx, *settingUsg, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *settingUsgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_setting_usg.SettingUsgModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// // Delete existing AP Group
	// err := r.client.DeleteAPGroup(ctx, data.Site.ValueString(), data.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error deleting AP Group",
	// 		"Could not delete AP Group, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
}

func parseSettingUsgResourceJson(ctx context.Context, json unifi.SettingUsg, model *resource_setting_usg.SettingUsgModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)

	var dhcpRelayServerSlice []types.String
	for _, dns := range []string{
		json.DHCPRelayServer1,
		json.DHCPRelayServer2,
		json.DHCPRelayServer3,
		json.DHCPRelayServer4,
		json.DHCPRelayServer5,
	} {
		dhcpRelayServerSlice = append(dhcpRelayServerSlice, types.StringValue(dns))
	}
	dhcpRelayServerList, diags := types.ListValueFrom(ctx, types.StringType, dhcpRelayServerSlice)
	if diags.HasError() {
		return diags
	}
	model.DhcpRelayServers = dhcpRelayServerList

	model.MulticastDnsEnabled = types.BoolValue(json.MdnsEnabled)

	return nil
}

func parseSettingUsgResourceModel(ctx context.Context, model resource_setting_usg.SettingUsgModel, json *unifi.SettingUsg) diag.Diagnostics {
	json.ID = model.Id.ValueString()

	var dhcpRelayServerSlice []types.String
	if !model.DhcpRelayServers.IsUnknown() && !model.DhcpRelayServers.IsNull() {
		diags := model.DhcpRelayServers.ElementsAs(ctx, &dhcpRelayServerSlice, false)
		if diags.HasError() {
			return diags
		}
	}
	json.DHCPRelayServer1 = tfStringSliceValueAtIndex(dhcpRelayServerSlice, 0).ValueString()
	json.DHCPRelayServer2 = tfStringSliceValueAtIndex(dhcpRelayServerSlice, 1).ValueString()
	json.DHCPRelayServer3 = tfStringSliceValueAtIndex(dhcpRelayServerSlice, 2).ValueString()
	json.DHCPRelayServer4 = tfStringSliceValueAtIndex(dhcpRelayServerSlice, 3).ValueString()
	json.DHCPRelayServer5 = tfStringSliceValueAtIndex(dhcpRelayServerSlice, 4).ValueString()

	json.MdnsEnabled = model.MulticastDnsEnabled.ValueBool()

	return nil
}
