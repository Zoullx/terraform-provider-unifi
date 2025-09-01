// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_ap_group"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &apGroupResource{}
	_ resource.ResourceWithConfigure   = &apGroupResource{}
	_ resource.ResourceWithImportState = &apGroupResource{}
)

func NewApGroupResource() resource.Resource {
	return &apGroupResource{}
}

type apGroupResource struct {
	client unifi.Client
}

func (r *apGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ap_group"
}

func (r *apGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_ap_group.ApGroupResourceSchema(ctx)
}

func (r *apGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *apGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *apGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_ap_group.ApGroupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.APGroup
	resp.Diagnostics.Append(parseApGroupResourceModel(ctx, data, &body)...)
	apGroup, err := r.client.CreateAPGroup(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AP Group",
			"Could not create AP Group, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseApGroupResourceJson(ctx, *apGroup, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_ap_group.ApGroupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed AP Group value from Unifi
	apGroup, err := r.client.GetAPGroup(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AP Group",
			"Could not read AP Group ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseApGroupResourceJson(ctx, *apGroup, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_ap_group.ApGroupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.APGroup
	resp.Diagnostics.Append(parseApGroupResourceModel(ctx, data, &body)...)
	apGroup, err := r.client.UpdateAPGroup(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AP Group",
			"Could not create AP Group, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseApGroupResourceJson(ctx, *apGroup, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_ap_group.ApGroupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing AP Group
	err := r.client.DeleteAPGroup(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting AP Group",
			"Could not delete AP Group, unexpected error: "+err.Error(),
		)
		return
	}
}

func parseApGroupResourceJson(ctx context.Context, json unifi.APGroup, model *resource_ap_group.ApGroupModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)
	model.Name = types.StringValue(json.Name)

	deviceMacList, diags := types.ListValueFrom(ctx, types.StringType, json.DeviceMACs)
	if diags.HasError() {
		return diags
	}
	model.DeviceMacs = deviceMacList

	return nil
}

func parseApGroupResourceModel(ctx context.Context, model resource_ap_group.ApGroupModel, json *unifi.APGroup) diag.Diagnostics {
	json.ID = model.Id.ValueString()
	json.Name = model.Name.ValueString()

	// var deviceMacSlice []string
	if !model.DeviceMacs.IsUnknown() && !model.DeviceMacs.IsNull() {
		diags := model.DeviceMacs.ElementsAs(ctx, &json.DeviceMACs, false)
		if diags.HasError() {
			return diags
		}
	}
	// json.DeviceMACs = deviceMacSlice

	return nil
}
