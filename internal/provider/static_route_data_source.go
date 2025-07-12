package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_static_route"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &staticRouteDataSource{}
)

func NewStaticRouteDataSource() datasource.DataSource {
	return &staticRouteDataSource{}
}

type staticRouteDataSource struct {
	client unifi.Client
}

func (d *staticRouteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_route"
}

func (d *staticRouteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_static_route.StaticRouteDataSourceSchema(ctx)
}

func (d *staticRouteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *staticRouteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_static_route.StaticRouteModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Static Route
	staticRoute, err := d.client.GetRouting(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Static Route",
			err.Error(),
		)
		return
	}

	parseStaticRouteDataSourceJson(*staticRoute, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseStaticRouteDataSourceJson(json unifi.Routing, model *datasource_static_route.StaticRouteModel) {
	model.Id = types.StringValue(json.ID)
	model.Distance = types.Int64Value(int64(json.StaticRouteDistance))
	model.Interface = types.StringValue(json.StaticRouteInterface)
	model.Name = types.StringValue(json.Name)
	model.Network = types.StringValue(json.StaticRouteNetwork)
	model.NextHop = types.StringValue(json.StaticRouteNexthop)
	model.Type = types.StringValue(json.StaticRouteType)
}
