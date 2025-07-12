package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_firewall_group"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &firewallGroupResource{}
	_ resource.ResourceWithConfigure   = &firewallGroupResource{}
	_ resource.ResourceWithImportState = &firewallGroupResource{}
)

func NewFirewallGroupResource() resource.Resource {
	return &firewallGroupResource{}
}

type firewallGroupResource struct {
	client unifi.Client
}

func (r *firewallGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_group"
}

func (r *firewallGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_firewall_group.FirewallGroupResourceSchema(ctx)
}

func (r *firewallGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *firewallGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *firewallGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_firewall_group.FirewallGroupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.FirewallGroup
	resp.Diagnostics.Append(parseFirewallGroupResourceModel(ctx, data, &body)...)
	firewallGroup, err := r.client.CreateFirewallGroup(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Firewall Group",
			"Could not create Firewall Group, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseFirewallGroupResourceJson(ctx, *firewallGroup, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_firewall_group.FirewallGroupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Firewall Group value from Unifi
	firewallGroup, err := r.client.GetFirewallGroup(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Firewall Group",
			"Could not read Firewall Group ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseFirewallGroupResourceJson(ctx, *firewallGroup, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_firewall_group.FirewallGroupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.FirewallGroup
	resp.Diagnostics.Append(parseFirewallGroupResourceModel(ctx, data, &body)...)
	firewallGroup, err := r.client.UpdateFirewallGroup(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Firewall Group",
			"Could not create Firewall Group, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseFirewallGroupResourceJson(ctx, *firewallGroup, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_firewall_group.FirewallGroupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Firewall Group
	err := r.client.DeleteFirewallGroup(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Firewall Group",
			"Could not delete Firewall Group, unexpected error: "+err.Error(),
		)
		return
	}
}

func parseFirewallGroupResourceJson(ctx context.Context, json unifi.FirewallGroup, model *resource_firewall_group.FirewallGroupModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)
	model.SiteId = types.StringValue(json.SiteID)
	model.Name = types.StringValue(json.Name)
	model.Type = types.StringValue(json.GroupType)

	memberList, diags := types.ListValueFrom(ctx, types.StringType, json.GroupMembers)
	if diags.HasError() {
		return diags
	}
	model.Members = memberList

	return nil
}

func parseFirewallGroupResourceModel(ctx context.Context, model resource_firewall_group.FirewallGroupModel, json *unifi.FirewallGroup) diag.Diagnostics {
	json.ID = model.Id.ValueString()
	json.SiteID = model.SiteId.ValueString()
	json.Name = model.Name.ValueString()
	json.GroupType = model.Type.ValueString()

	if !model.Members.IsUnknown() && !model.Members.IsNull() {
		diags := model.Members.ElementsAs(ctx, &json.GroupMembers, false)
		if diags.HasError() {
			return diags
		}
	}

	return nil
}
