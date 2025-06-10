package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_ap_group"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &apGroupDataSource{}
)

func NewApGroupDataSource() datasource.DataSource {
	return &apGroupDataSource{}
}

type apGroupDataSource struct {
	client *unifi.Client
}

func (d *apGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ap_group"
}

func (d *apGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_ap_group.ApGroupDataSourceSchema(ctx)
}

func (d *apGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *apGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_ap_group.ApGroupModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.Id.IsNull() || data.Id.IsUnknown()) && (data.Name.IsNull() || data.Name.IsUnknown()) {
		resp.Diagnostics.AddError(
			"Id and Name are null or unknown",
			"Id or Name is required to retrieve an AP Group.",
		)
		return
	}

	// Get AP Group
	var apGroup *unifi.APGroup
	var err error
	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		apGroup, err = d.client.GetAPGroup(ctx, data.Site.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() && !data.Name.IsUnknown() {
		apGroup, err = d.client.GetAPGroupByName(ctx, data.Site.ValueString(), data.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read AP Group",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseApGroupDataSourceJson(ctx, *apGroup, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseApGroupDataSourceJson(ctx context.Context, json unifi.APGroup, model *datasource_ap_group.ApGroupModel) diag.Diagnostics {
	model.Id = types.StringValue(json.ID)
	model.Name = types.StringValue(json.Name)

	deviceMacList, diags := types.ListValueFrom(ctx, types.StringType, json.DeviceMACs)
	if diags.HasError() {
		return diags
	}
	model.DeviceMacs = deviceMacList

	return nil
}
