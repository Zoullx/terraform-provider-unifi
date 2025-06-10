package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_firewall_group"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &firewallGroupDataSource{}
)

func NewFirewallGroupDataSource() datasource.DataSource {
	return &firewallGroupDataSource{}
}

type firewallGroupDataSource struct {
	client *unifi.Client
}

func (d *firewallGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_group"
}

func (d *firewallGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_firewall_group.FirewallGroupDataSourceSchema(ctx)
}

func (d *firewallGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*unifi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *unifi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *firewallGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_firewall_group.FirewallGroupModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.Id.IsNull() || data.Id.IsUnknown()) && (data.Name.IsNull() || data.Name.IsUnknown()) {
		resp.Diagnostics.AddError(
			"Id and Name are null or unknown",
			"Id or Name is required to retrieve a Firewall Group.",
		)
		return
	}

	// Get Firewall Group
	var firewallGroup *unifi.FirewallGroup
	var err error
	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		firewallGroup, err = d.client.GetFirewallGroup(ctx, data.Site.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() && !data.Name.IsUnknown() {
		firewallGroup, err = d.client.GetFirewallGroupByName(ctx, data.Site.ValueString(), data.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Firewall Group",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseFirewallGroupDataSourceJson(ctx, *firewallGroup, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseFirewallGroupDataSourceJson(ctx context.Context, json unifi.FirewallGroup, model *datasource_firewall_group.FirewallGroupModel) diag.Diagnostics {
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
