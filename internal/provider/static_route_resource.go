package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_static_route"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &staticRouteResource{}
	_ resource.ResourceWithConfigure   = &staticRouteResource{}
	_ resource.ResourceWithImportState = &staticRouteResource{}
)

func NewStaticRouteResource() resource.Resource {
	return &staticRouteResource{}
}

type staticRouteResource struct {
	client unifi.Client
}

func (r *staticRouteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_route"
}

func (r *staticRouteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_static_route.StaticRouteResourceSchema(ctx)
}

func (r *staticRouteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *staticRouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *staticRouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_static_route.StaticRouteModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.Routing
	parseStaticRouteResourceModel(data, &body)
	staticRoute, err := r.client.CreateRouting(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Static Route",
			"Could not create Static Route, unexpected error: "+err.Error(),
		)
		return
	}

	parseStaticRouteResourceJson(*staticRoute, &data)

	// data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staticRouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_static_route.StaticRouteModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Static Route value from Unifi
	staticRoute, err := r.client.GetRouting(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Static Route",
			"Could not read Static Route ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	parseStaticRouteResourceJson(*staticRoute, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staticRouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_static_route.StaticRouteModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.Routing
	parseStaticRouteResourceModel(data, &body)
	staticRoute, err := r.client.UpdateRouting(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Static Route",
			"Could not create Static Route, unexpected error: "+err.Error(),
		)
		return
	}

	parseStaticRouteResourceJson(*staticRoute, &data)

	// data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staticRouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_static_route.StaticRouteModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Static Route
	err := r.client.DeleteRouting(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Static Route",
			"Could not delete Static Route, unexpected error: "+err.Error(),
		)
		return
	}
}

func parseStaticRouteResourceJson(json unifi.Routing, model *resource_static_route.StaticRouteModel) {
	model.Id = types.StringValue(json.ID)
	model.Distance = types.Int64Value(int64(json.StaticRouteDistance))
	model.Interface = types.StringValue(json.StaticRouteInterface)
	model.Name = types.StringValue(json.Name)
	model.Network = types.StringValue(json.StaticRouteNetwork)
	model.NextHop = types.StringValue(json.StaticRouteNexthop)
	model.Type = types.StringValue(json.StaticRouteType)
}

func parseStaticRouteResourceModel(model resource_static_route.StaticRouteModel, json *unifi.Routing) {
	json.ID = model.Id.ValueString()
	json.StaticRouteDistance = int(model.Distance.ValueInt64())
	json.StaticRouteInterface = model.Interface.ValueString()
	json.Name = model.Name.ValueString()
	json.StaticRouteNetwork = model.Network.ValueString()
	json.StaticRouteNexthop = model.NextHop.ValueString()
	json.StaticRouteType = model.Type.ValueString()
}
