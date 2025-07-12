package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_user_groups"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &userGroupsDataSource{}
)

func NewUserGroupsDataSource() datasource.DataSource {
	return &userGroupsDataSource{}
}

type userGroupsDataSource struct {
	client unifi.Client
}

func (d *userGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_groups"
}

func (d *userGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_user_groups.UserGroupsDataSourceSchema(ctx)
}

func (d *userGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *userGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_user_groups.UserGroupsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get User Groups
	userGroups, err := d.client.ListUserGroup(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read User Groups",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseUserGroupsDataSourceJson(ctx, userGroups, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseUserGroupsDataSourceJson(ctx context.Context, json []unifi.UserGroup, model *datasource_user_groups.UserGroupsModel) diag.Diagnostics {
	userGroupList, diags := types.ListValueFrom(ctx, datasource_user_groups.UserGroupsValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.UserGroups = userGroupList

	return nil
}
