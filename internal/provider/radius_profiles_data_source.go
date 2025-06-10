package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_radius_profiles"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &radiusProfilesDataSource{}
)

func NewRadiusProfilesDataSource() datasource.DataSource {
	return &radiusProfilesDataSource{}
}

type radiusProfilesDataSource struct {
	client *unifi.Client
}

func (d *radiusProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_radius_profiles"
}

func (d *radiusProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_radius_profiles.RadiusProfilesDataSourceSchema(ctx)
}

func (d *radiusProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *radiusProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_radius_profiles.RadiusProfilesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get RADIUS Profiles
	radiusProfiles, err := d.client.ListRADIUSProfile(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read RADIUS Profiles",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseRadiusProfilesDataSourceJson(ctx, radiusProfiles, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseRadiusProfilesDataSourceJson(ctx context.Context, json []unifi.RADIUSProfile, model *datasource_radius_profiles.RadiusProfilesModel) diag.Diagnostics {
	radiusProfileList, diags := types.ListValueFrom(ctx, datasource_radius_profiles.RadiusProfilesValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.RadiusProfiles = radiusProfileList

	return nil
}
