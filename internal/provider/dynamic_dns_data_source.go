// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_dynamic_dns"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &dynamicDnsDataSource{}
)

func NewDynamicDnsDataSource() datasource.DataSource {
	return &dynamicDnsDataSource{}
}

type dynamicDnsDataSource struct {
	client unifi.Client
}

func (d *dynamicDnsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dynamic_dns"
}

func (d *dynamicDnsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_dynamic_dns.DynamicDnsDataSourceSchema(ctx)
}

func (d *dynamicDnsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dynamicDnsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_dynamic_dns.DynamicDnsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Dynamic DNS
	dynamicDns, err := d.client.GetDynamicDNS(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Dynamic DNS",
			err.Error(),
		)
		return
	}

	parseDynamicDnsDataSourceJson(*dynamicDns, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseDynamicDnsDataSourceJson(json unifi.DynamicDNS, model *datasource_dynamic_dns.DynamicDnsModel) {
	model.Id = types.StringValue(json.ID)
	model.HostName = types.StringValue(json.HostName)
	model.Interface = types.StringValue(json.Interface)
	model.Login = types.StringValue(json.Login)
	model.Password = types.StringValue(json.XPassword)
	model.Server = types.StringValue(json.Server)
	model.Service = types.StringValue(json.Service)
}
