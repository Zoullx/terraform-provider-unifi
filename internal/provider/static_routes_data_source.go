package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_static_routes"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &staticRoutesDataSource{}
)

func NewStaticRoutesDataSource() datasource.DataSource {
	return &staticRoutesDataSource{}
}

type staticRoutesDataSource struct {
	client *unifi.Client
}

func (d *staticRoutesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_routes"
}

func (d *staticRoutesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_static_routes.StaticRoutesDataSourceSchema(ctx)
}

func (d *staticRoutesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *staticRoutesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_static_routes.StaticRoutesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Static Routes
	staticRoutes, err := d.client.ListRouting(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Static Routes",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseStaticRoutesDataSourceJson(ctx, staticRoutes, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseStaticRoutesDataSourceJson(ctx context.Context, json []unifi.Routing, model *datasource_static_routes.StaticRoutesModel) diag.Diagnostics {
	staticRouteList, diags := types.ListValueFrom(ctx, datasource_static_routes.StaticRoutesValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.StaticRoutes = staticRouteList

	return nil
}
