package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_user_group"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &userGroupDataSource{}
)

func NewUserGroupDataSource() datasource.DataSource {
	return &userGroupDataSource{}
}

type userGroupDataSource struct {
	client unifi.Client
}

func (d *userGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (d *userGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_user_group.UserGroupDataSourceSchema(ctx)
}

func (d *userGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *userGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_user_group.UserGroupModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.Id.IsNull() || data.Id.IsUnknown()) && (data.Name.IsNull() || data.Name.IsUnknown()) {
		resp.Diagnostics.AddError(
			"Id and Name are null or unknown",
			"Id or Name is required to retrieve a User Group.",
		)
		return
	}

	// Get User Group
	var userGroup *unifi.UserGroup
	var err error
	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		userGroup, err = d.client.GetUserGroup(ctx, data.Site.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() && !data.Name.IsUnknown() {
		userGroup, err = d.client.GetUserGroupByName(ctx, data.Site.ValueString(), data.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read User Group",
			err.Error(),
		)
		return
	}

	parseUserGroupDataSourceJson(*userGroup, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseUserGroupDataSourceJson(json unifi.UserGroup, model *datasource_user_group.UserGroupModel) {
	model.Id = types.StringValue(json.ID)
	model.Name = types.StringValue(json.Name)
	model.QosRateMaxDown = types.Int64Value(int64(json.QOSRateMaxDown))
	model.QosRateMaxUp = types.Int64Value(int64(json.QOSRateMaxUp))
}
