// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_dynamic_dns"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &dynamicDnsResource{}
	_ resource.ResourceWithConfigure   = &dynamicDnsResource{}
	_ resource.ResourceWithImportState = &dynamicDnsResource{}
)

func NewDynamicDnsResource() resource.Resource {
	return &dynamicDnsResource{}
}

type dynamicDnsResource struct {
	client unifi.Client
}

func (r *dynamicDnsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dynamic_dns"
}

func (r *dynamicDnsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_dynamic_dns.DynamicDnsResourceSchema(ctx)
}

func (r *dynamicDnsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dynamicDnsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: site/id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("site"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *dynamicDnsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_dynamic_dns.DynamicDnsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.DynamicDNS
	parseDynamicDnsResourceModel(data, &body)
	dynamicDns, err := r.client.CreateDynamicDNS(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Dynamic DNS",
			"Could not create Dynamic DNS, unexpected error: "+err.Error(),
		)
		return
	}

	parseDynamicDnsResourceJson(*dynamicDns, &data)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dynamicDnsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_dynamic_dns.DynamicDnsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Dynamic DNS value from Unifi
	dynamicDns, err := r.client.GetDynamicDNS(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Dynamic DNS",
			"Could not read Dynamic DNS ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	parseDynamicDnsResourceJson(*dynamicDns, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dynamicDnsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_dynamic_dns.DynamicDnsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.DynamicDNS
	parseDynamicDnsResourceModel(data, &body)
	dynamicDns, err := r.client.UpdateDynamicDNS(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Dynamic DNS",
			"Could not create Dynamic DNS, unexpected error: "+err.Error(),
		)
		return
	}

	parseDynamicDnsResourceJson(*dynamicDns, &data)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dynamicDnsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_dynamic_dns.DynamicDnsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Dyanmic DNS
	err := r.client.DeleteDynamicDNS(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Dynamic DNS",
			"Could not delete Dynamic DNS, unexpected error: "+err.Error(),
		)
		return
	}
}

func parseDynamicDnsResourceJson(json unifi.DynamicDNS, model *resource_dynamic_dns.DynamicDnsModel) {
	model.Id = types.StringValue(json.ID)
	model.HostName = types.StringValue(json.HostName)
	model.Interface = types.StringValue(json.Interface)
	model.Login = types.StringValue(json.Login)
	model.Password = types.StringValue(json.XPassword)
	model.Server = types.StringValue(json.Server)
	model.Service = types.StringValue(json.Service)
}

func parseDynamicDnsResourceModel(model resource_dynamic_dns.DynamicDnsModel, json *unifi.DynamicDNS) {
	json.ID = model.Id.ValueString()
	json.HostName = model.HostName.ValueString()
	json.Interface = model.Interface.ValueString()
	json.Login = model.HostName.ValueString()
	json.XPassword = model.Password.ValueString()
	json.Server = model.Server.ValueString()
	json.Service = model.Service.ValueString()
}
