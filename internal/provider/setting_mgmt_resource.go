package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_setting_mgmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &settingMgmtResource{}
	_ resource.ResourceWithConfigure   = &settingMgmtResource{}
	_ resource.ResourceWithImportState = &settingMgmtResource{}
)

func NewSettingMgmtResource() resource.Resource {
	return &settingMgmtResource{}
}

type settingMgmtResource struct {
	client unifi.Client
}

func (r *settingMgmtResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting_mgmt"
}

func (r *settingMgmtResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_setting_mgmt.SettingMgmtResourceSchema(ctx)
}

func (r *settingMgmtResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingMgmtResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *settingMgmtResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_setting_mgmt.SettingMgmtModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.SettingMgmt
	resp.Diagnostics.Append(parseSettingMgmtResourceModel(ctx, data, &body)...)
	settingMgmt, err := r.client.UpdateSettingMgmt(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Setting Mgmt",
			"Could not updating Setting Mgmt, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseSettingMgmtResourceJson(ctx, *settingMgmt, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *settingMgmtResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_setting_mgmt.SettingMgmtModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Setting Mgmt value from Unifi
	settingMgmt, err := r.client.GetSettingMgmt(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Setting Mgmt",
			"Could not read Setting Mgmt ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseSettingMgmtResourceJson(ctx, *settingMgmt, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *settingMgmtResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_setting_mgmt.SettingMgmtModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.SettingMgmt
	resp.Diagnostics.Append(parseSettingMgmtResourceModel(ctx, data, &body)...)
	settingMgmt, err := r.client.UpdateSettingMgmt(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Setting Mgmt",
			"Could not create Setting Mgmt, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseSettingMgmtResourceJson(ctx, *settingMgmt, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *settingMgmtResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_setting_mgmt.SettingMgmtModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// // Delete existing Setting Mgmt
	// err := r.client.DeleteSett(ctx, data.Site.ValueString(), data.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error deleting Setting Mgmt",
	// 		"Could not delete Setting Mgmt, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
}

func parseSettingMgmtResourceJson(ctx context.Context, json unifi.SettingMgmt, model *resource_setting_mgmt.SettingMgmtModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)
	model.AutoUpgrade = types.BoolValue(json.AutoUpgrade)
	model.SshEnabled = types.BoolValue(json.XSshEnabled)

	sshKeyList, diags := types.ListValueFrom(ctx, resource_setting_mgmt.SshKeysValue{}.Type(ctx), json.XSshKeys)
	if diags.HasError() {
		return diags
	}
	model.SshKeys = sshKeyList

	return nil
}

func parseSettingMgmtResourceModel(ctx context.Context, model resource_setting_mgmt.SettingMgmtModel, json *unifi.SettingMgmt) diag.Diagnostics {
	json.ID = model.Id.ValueString()
	json.AutoUpgrade = model.AutoUpgrade.ValueBool()
	json.XSshEnabled = model.SshEnabled.ValueBool()

	if !model.SshKeys.IsUnknown() && !model.SshKeys.IsNull() {
		diags := model.SshKeys.ElementsAs(ctx, &json.XSshKeys, false)
		if diags.HasError() {
			return diags
		}
	}

	return nil
}
