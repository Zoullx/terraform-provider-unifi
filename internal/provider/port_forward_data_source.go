// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_port_forward"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &portForwardDataSource{}
)

func NewPortForwardDataSource() datasource.DataSource {
	return &portForwardDataSource{}
}

type portForwardDataSource struct {
	client unifi.Client
}

func (d *portForwardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forward"
}

func (d *portForwardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_port_forward.PortForwardDataSourceSchema(ctx)
}

func (d *portForwardDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *portForwardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_port_forward.PortForwardModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Port Forward
	portForward, err := d.client.GetPortForward(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Port Forward",
			err.Error(),
		)
		return
	}

	parsePortForwardDataSourceJson(*portForward, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parsePortForwardDataSourceJson(json unifi.PortForward, model *datasource_port_forward.PortForwardModel) {
	model.Id = types.StringValue(json.ID)
	model.SiteId = types.StringValue(json.SiteID)
	model.DstPort = types.StringValue(json.DstPort)
	model.Enabled = types.BoolValue(json.Enabled)
	model.FwdIp = types.StringValue(json.Fwd)
	model.FwdPort = types.StringValue(json.FwdPort)
	model.Log = types.BoolValue(json.Log)
	model.Name = types.StringValue(json.Name)
	model.PortForwardInterface = types.StringValue(json.PfwdInterface)
	model.Protocol = types.StringValue(json.Proto)
	model.SrcIp = types.StringValue(json.Src)
}
