// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_dynamic_dnses"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &dynamicDnsesDataSource{}
)

func NewDynamicDnsesDataSource() datasource.DataSource {
	return &dynamicDnsesDataSource{}
}

type dynamicDnsesDataSource struct {
	client unifi.Client
}

func (d *dynamicDnsesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dynamic_dnses"
}

func (d *dynamicDnsesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_dynamic_dnses.DynamicDnsesDataSourceSchema(ctx)
}

func (d *dynamicDnsesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dynamicDnsesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_dynamic_dnses.DynamicDnsesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Dynamic DNSes
	dynamicDnses, err := d.client.ListDynamicDNS(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Dynamic DNSes",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseDynamicDnsesDataSourceJson(ctx, dynamicDnses, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseDynamicDnsesDataSourceJson(ctx context.Context, json []unifi.DynamicDNS, model *datasource_dynamic_dnses.DynamicDnsesModel) diag.Diagnostics {
	dynamicDnsesList, diags := types.ListValueFrom(ctx, datasource_dynamic_dnses.DynamicDnsesValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.DynamicDnses = dynamicDnsesList

	return nil
}
