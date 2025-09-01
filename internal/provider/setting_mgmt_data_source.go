// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_setting_mgmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &settingMgmtDataSource{}
)

func NewSettingMgmtDataSource() datasource.DataSource {
	return &settingMgmtDataSource{}
}

type settingMgmtDataSource struct {
	client unifi.Client
}

func (d *settingMgmtDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting_mgmt"
}

func (d *settingMgmtDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_setting_mgmt.SettingMgmtDataSourceSchema(ctx)
}

func (d *settingMgmtDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *settingMgmtDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_setting_mgmt.SettingMgmtModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Setting Mgmt
	settingMgmt, err := d.client.GetSettingMgmt(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Setting Mgmt",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseSettingMgmtDataSourceJson(ctx, *settingMgmt, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseSettingMgmtDataSourceJson(ctx context.Context, json unifi.SettingMgmt, model *datasource_setting_mgmt.SettingMgmtModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)
	model.AutoUpgrade = types.BoolValue(json.AutoUpgrade)
	model.SshEnabled = types.BoolValue(json.XSshEnabled)

	sshKeyList, diags := types.ListValueFrom(ctx, datasource_setting_mgmt.SshKeysValue{}.Type(ctx), json.XSshKeys)
	if diags.HasError() {
		return diags
	}
	model.SshKeys = sshKeyList

	return nil
}
