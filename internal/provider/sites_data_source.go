package provider

import (
	"context"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_sites"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &sitesDataSource{}
)

func NewSitesDataSource() datasource.DataSource {
	return &sitesDataSource{}
}

type sitesDataSource struct {
	client *unifi.Client
}

func (d *sitesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sites"
}

func (d *sitesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_sites.SitesDataSourceSchema(ctx)
}

func (d *sitesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_sites.SitesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Sites
	sites, err := d.client.ListSites(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Sites",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseSitesDataSourceJson(ctx, sites, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseSitesDataSourceJson(ctx context.Context, json []unifi.Site, model *datasource_sites.SitesModel) diag.Diagnostics {
	siteList, diags := types.ListValueFrom(ctx, datasource_sites.SitesValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.Sites = siteList

	return nil
}
