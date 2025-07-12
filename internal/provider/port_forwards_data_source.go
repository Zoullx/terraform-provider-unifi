package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_port_forwards"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &portForwardsDataSource{}
)

func NewPortForwardsDataSource() datasource.DataSource {
	return &portForwardsDataSource{}
}

type portForwardsDataSource struct {
	client unifi.Client
}

func (d *portForwardsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forwards"
}

func (d *portForwardsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_port_forwards.PortForwardsDataSourceSchema(ctx)
}

func (d *portForwardsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *portForwardsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_port_forwards.PortForwardsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Port Forwards
	portForwards, err := d.client.ListPortForward(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Port Forwards",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parsePortForwardsDataSourceJson(ctx, portForwards, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parsePortForwardsDataSourceJson(ctx context.Context, json []unifi.PortForward, model *datasource_port_forwards.PortForwardsModel) diag.Diagnostics {
	portForwardList, diags := types.ListValueFrom(ctx, datasource_port_forwards.PortForwardsValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.PortForwards = portForwardList

	return nil
}
