package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_user"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client unifi.Client
}

func (r *userResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_user.UserResourceSchema(ctx)
}

func (r *userResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nill check when handling ProviderData because Terraform
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

	r.client = client
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: site/id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("site"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_user.UserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.User
	parseUserResourceModel(data, &body)
	user, err := r.client.CreateUser(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating User",
			"Could not create User, unexpected error: "+err.Error(),
		)
		return
	}

	parseUserResourceJson(*user, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_user.UserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed User value from Unifi
	user, err := r.client.GetUser(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading User",
			"Could not read User ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	parseUserResourceJson(*user, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_user.UserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.User
	parseUserResourceModel(data, &body)
	user, err := r.client.UpdateUser(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating User",
			"Could not create User, unexpected error: "+err.Error(),
		)
		return
	}

	parseUserResourceJson(*user, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_user.UserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing User
	err := r.client.DeleteUser(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting User",
			"Could not delete User, unexpected error: "+err.Error(),
		)
		return
	}
}

func parseUserResourceJson(json unifi.User, model *resource_user.UserModel) {
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

func parseUserResourceModel(model resource_user.UserModel, json *unifi.User) {
	json.ID = model.Id.ValueString()
	json.Blocked = model.Blocked.ValueBool()
	json.DevIdOverride = int(model.DevIdOverride.ValueInt64())
	json.FixedIP = model.FixedIp.ValueString()
	json.Hostname = model.Hostname.ValueString()
	json.IP = model.Ip.ValueString()
	json.LocalDNSRecord = model.LocalDnsRecord.ValueString()
	json.MAC = model.Mac.ValueString()
	json.Name = model.Name.ValueString()
	json.NetworkID = model.NetworkId.ValueString()
	json.Note = model.Note.ValueString()
	json.UserGroupID = model.UserGroupId.ValueString()
}
