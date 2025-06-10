package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_user"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *unifi.Client
}

func (d *userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_user.UserDataSourceSchema(ctx)
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_user.UserModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get User
	user, err := d.client.GetUser(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read User",
			err.Error(),
		)
		return
	}

	parseUserDataSourceJson(*user, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseUserDataSourceJson(json unifi.User, model *datasource_user.UserModel) {
	model.Id = types.StringValue(json.ID)
	model.Blocked = types.BoolValue(json.Blocked)
	model.DevIdOverride = types.Int64Value(int64(json.DevIdOverride))
	model.FixedIp = types.StringValue(json.FixedIP)
	model.Hostname = types.StringValue(json.Hostname)
	model.Ip = types.StringValue(json.IP)
	model.LocalDnsRecord = types.StringValue(json.LocalDNSRecord)
	model.Mac = types.StringValue(json.MAC)
	model.Name = types.StringValue(json.Name)
	model.NetworkId = types.StringValue(json.NetworkID)
	model.Note = types.StringValue(json.Note)
	model.UserGroupId = types.StringValue(json.UserGroupID)
}
