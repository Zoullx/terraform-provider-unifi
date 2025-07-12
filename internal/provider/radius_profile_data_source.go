package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_radius_profile"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &radiusProfileDataSource{}
)

func NewRadiusProfileDataSource() datasource.DataSource {
	return &radiusProfileDataSource{}
}

type radiusProfileDataSource struct {
	client unifi.Client
}

func (d *radiusProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_radius_profile"
}

func (d *radiusProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_radius_profile.RadiusProfileDataSourceSchema(ctx)
}

func (d *radiusProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *radiusProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_radius_profile.RadiusProfileModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get RADIUS Profile
	radiusProfile, err := d.client.GetRADIUSProfile(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read RADIUS Profile",
			err.Error(),
		)
		return
	}

	parseRadiusProfileDataSourceJson(*radiusProfile, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseRadiusProfileDataSourceJson(json unifi.RADIUSProfile, model *datasource_radius_profile.RadiusProfileModel) {
	model.Id = types.StringValue(json.ID)
	model.Name = types.StringValue(json.Name)
}
