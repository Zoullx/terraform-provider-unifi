package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_devices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &devicesDataSource{}
)

func NewDevicesDataSource() datasource.DataSource {
	return &devicesDataSource{}
}

type devicesDataSource struct {
	client unifi.Client
}

func (d *devicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

func (d *devicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_devices.DevicesDataSourceSchema(ctx)
}

func (d *devicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *devicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_devices.DevicesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Devices
	devices, err := d.client.ListDevice(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Devices",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseDevicesDataSourceJson(ctx, devices, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseDevicesDataSourceJson(ctx context.Context, json []unifi.Device, model *datasource_devices.DevicesModel) diag.Diagnostics {
	devicesList, diags := types.ListValueFrom(ctx, datasource_devices.DevicesValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.Devices = devicesList

	return nil
}
