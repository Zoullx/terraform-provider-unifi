package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_setting_radius"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &settingRadiusResource{}
	_ resource.ResourceWithConfigure   = &settingRadiusResource{}
	_ resource.ResourceWithImportState = &settingRadiusResource{}
)

func NewSettingRadiusResource() resource.Resource {
	return &settingRadiusResource{}
}

type settingRadiusResource struct {
	client *unifi.Client
}

func (r *settingRadiusResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting_radius"
}

func (r *settingRadiusResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_setting_radius.SettingRadiusResourceSchema(ctx)
}

func (r *settingRadiusResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nill check when handling ProviderData because Terraform
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

	r.client = client
}

func (r *settingRadiusResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *settingRadiusResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_setting_radius.SettingRadiusModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.SettingRadius
	parseSettingRadiusResourceModel(data, &body)
	settingRadius, err := r.client.UpdateSettingRadius(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Setting RADIUS",
			"Could not create Setting RADIUS, unexpected error: "+err.Error(),
		)
		return
	}

	parseSettingRadiusResourceJson(*settingRadius, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *settingRadiusResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_setting_radius.SettingRadiusModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Setting RADIUS value from Unifi
	settingRadius, err := r.client.GetSettingRadius(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Setting RADIUS",
			"Could not read Setting RADIUS ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	parseSettingRadiusResourceJson(*settingRadius, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *settingRadiusResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_setting_radius.SettingRadiusModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.SettingRadius
	parseSettingRadiusResourceModel(data, &body)
	settingRadius, err := r.client.UpdateSettingRadius(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Setting RADIUS",
			"Could not create Setting RADIUS, unexpected error: "+err.Error(),
		)
		return
	}

	parseSettingRadiusResourceJson(*settingRadius, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *settingRadiusResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_setting_radius.SettingRadiusModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// // Delete existing Port Forward
	// err := r.client.DeletePortForward(ctx, data.Site.ValueString(), data.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error deleting Port Forward",
	// 		"Could not delete Port Forward, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
}

func parseSettingRadiusResourceJson(json unifi.SettingRadius, model *resource_setting_radius.SettingRadiusModel) {
	model.Id = types.StringValue(json.ID)
	model.AccountingEnabled = types.BoolValue(json.AccountingEnabled)
	model.AccountingPort = types.Int64Value(int64(json.AcctPort))
	model.AuthPort = types.Int64Value(int64(json.AuthPort))
	model.Enabled = types.BoolValue(json.Enabled)
	model.InterimUpdateInterval = types.Int64Value(int64(json.InterimUpdateInterval))
	model.Secret = types.StringValue(json.XSecret)
	model.TunneledReply = types.BoolValue(json.TunneledReply)
}

func parseSettingRadiusResourceModel(model resource_setting_radius.SettingRadiusModel, json *unifi.SettingRadius) {
	json.ID = model.Id.ValueString()
	json.AccountingEnabled = model.AccountingEnabled.ValueBool()
	json.AcctPort = int(model.AccountingPort.ValueInt64())
	json.AuthPort = int(model.AuthPort.ValueInt64())
	json.Enabled = model.Enabled.ValueBool()
	json.InterimUpdateInterval = int(model.InterimUpdateInterval.ValueInt64())
	json.XSecret = model.Secret.ValueString()
	json.TunneledReply = model.TunneledReply.ValueBool()
}
