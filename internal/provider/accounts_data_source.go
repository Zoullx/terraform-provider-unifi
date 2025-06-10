package provider

import (
	"context"
	"fmt"

	"github.com/zoullx/terraform-provider-unifi/internal/datasource_accounts"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ datasource.DataSource = &accountsDataSource{}
)

func NewAccountsDataSource() datasource.DataSource {
	return &accountsDataSource{}
}

type accountsDataSource struct {
	client *unifi.Client
}

func (d *accountsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accounts"
}

func (d *accountsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_accounts.AccountsDataSourceSchema(ctx)
}

func (d *accountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_accounts.AccountsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Accounts
	accounts, err := d.client.ListAccount(ctx, data.Site.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Accounts",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(parseAccountsDataSourceJson(ctx, accounts, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseAccountsDataSourceJson(ctx context.Context, json []unifi.Account, model *datasource_accounts.AccountsModel) diag.Diagnostics {
	// var accounts []datasource_accounts.AccountsValue
	// for _, account := range json {
	// 	accounts = append(accounts, datasource_accounts.AccountsValue{
	// 		Id:               types.StringValue(account.ID),
	// 		Name:             types.StringValue(account.Name),
	// 		Password:         types.StringValue(account.XPassword),
	// 		TunnelType:       types.Int64Value(int64(account.TunnelType)),
	// 		TunnelMediumType: types.Int64Value(int64(account.TunnelMediumType)),
	// 		NetworkId:        types.StringValue(account.NetworkID),
	// 	})
	// }

	accountsList, diags := types.ListValueFrom(ctx, datasource_accounts.AccountsValue{}.Type(ctx), json)
	if diags.HasError() {
		return diags
	}
	model.Accounts = accountsList

	return nil
}
