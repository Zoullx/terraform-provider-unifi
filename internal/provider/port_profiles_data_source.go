// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_port_profiles"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &portProfilesDataSource{}
)

func NewPortProfilesDataSource() datasource.DataSource {
	return &portProfilesDataSource{}
}

type portProfilesDataSource struct {
	client unifi.Client
}

func (d *portProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_profiles"
}

func (d *portProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_port_profiles.PortProfilesDataSourceSchema(ctx)
}

func (d *portProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *portProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_port_profiles.PortProfilesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Port Profiles
	portProfiles, err := d.client.ListPortProfile(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Port Profiles",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parsePortProfilesDataSourceJson(ctx, portProfiles, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parsePortProfilesDataSourceJson(ctx context.Context, json []unifi.PortProfile, model *datasource_port_profiles.PortProfilesModel) diag.Diagnostics {
	portProfileList, diags := types.ListValueFrom(ctx, datasource_port_profiles.PortProfilesValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.PortProfiles = portProfileList

	return nil
}
