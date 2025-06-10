package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_firewall_groups"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &firewallGroupsDataSource{}
)

func NewFirewallGroupsDataSource() datasource.DataSource {
	return &firewallGroupsDataSource{}
}

type firewallGroupsDataSource struct {
	client *unifi.Client
}

func (d *firewallGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_groups"
}

func (d *firewallGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_firewall_groups.FirewallGroupsDataSourceSchema(ctx)
}

func (d *firewallGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *firewallGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_firewall_groups.FirewallGroupsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Firewall Groups
	firewallGroups, err := d.client.ListFirewallGroup(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Firewall Groups",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseFirewallGroupsDataSourceJson(ctx, firewallGroups, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseFirewallGroupsDataSourceJson(ctx context.Context, json []unifi.FirewallGroup, model *datasource_firewall_groups.FirewallGroupsModel) diag.Diagnostics {
	firewallGroupList, diags := types.ListValueFrom(ctx, datasource_firewall_groups.FirewallGroupsValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.FirewallGroups = firewallGroupList

	return nil
}
