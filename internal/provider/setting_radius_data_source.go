package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_setting_radius"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &settingRadiusDataSource{}
)

func NewSettingRadiusDataSource() datasource.DataSource {
	return &settingRadiusDataSource{}
}

type settingRadiusDataSource struct {
	client unifi.Client
}

func (d *settingRadiusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting_radius"
}

func (d *settingRadiusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_setting_radius.SettingRadiusDataSourceSchema(ctx)
}

func (d *settingRadiusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *settingRadiusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_setting_radius.SettingRadiusModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Setting RADIUS
	settingRadius, err := d.client.GetSettingRadius(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Setting RADIUS",
			err.Error(),
		)
		return
	}

	parseSettingRadiusDataSourceJson(*settingRadius, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseSettingRadiusDataSourceJson(json unifi.SettingRadius, model *datasource_setting_radius.SettingRadiusModel) {
	model.Id = types.StringValue(json.ID)
	model.AccountingEnabled = types.BoolValue(json.AccountingEnabled)
	model.AccountingPort = types.Int64Value(int64(json.AcctPort))
	model.AuthPort = types.Int64Value(int64(json.AuthPort))
	model.Enabled = types.BoolValue(json.Enabled)
	model.InterimUpdateInterval = types.Int64Value(int64(json.InterimUpdateInterval))
	model.Secret = types.StringValue(json.XSecret)
	model.TunneledReply = types.BoolValue(json.TunneledReply)
}
