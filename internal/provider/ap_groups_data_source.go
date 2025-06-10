package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_ap_groups"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &apGroupsDataSource{}
)

func NewApGroupsDataSource() datasource.DataSource {
	return &apGroupsDataSource{}
}

type apGroupsDataSource struct {
	client *unifi.Client
}

func (d *apGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ap_groups"
}

func (d *apGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_ap_groups.ApGroupsDataSourceSchema(ctx)
}

func (d *apGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *apGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_ap_groups.ApGroupsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get AP Groups
	apGroups, err := d.client.ListAPGroup(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read AP Groups",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseApGroupsDataSourceJson(ctx, apGroups, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseApGroupsDataSourceJson(ctx context.Context, json []unifi.APGroup, model *datasource_ap_groups.ApGroupsModel) diag.Diagnostics {
	apGroupsList, diags := types.ListValueFrom(ctx, datasource_ap_groups.ApGroupsValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.ApGroups = apGroupsList

	return nil
}
