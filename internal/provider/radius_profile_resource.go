package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_radius_profile"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &radiusProfileResource{}
	_ resource.ResourceWithConfigure   = &radiusProfileResource{}
	_ resource.ResourceWithImportState = &radiusProfileResource{}
)

func NewRadiusProfileResource() resource.Resource {
	return &radiusProfileResource{}
}

type radiusProfileResource struct {
	client *unifi.Client
}

func (r *radiusProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_radius_profile"
}

func (r *radiusProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_radius_profile.RadiusProfileResourceSchema(ctx)
}

func (r *radiusProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nill check when handling ProviderData because Terraform
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

	r.client = client
}

func (r *radiusProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *radiusProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_radius_profile.RadiusProfileModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.RADIUSProfile
	parseRadiusProfileResourceModel(data, &body)
	radiusProfile, err := r.client.CreateRADIUSProfile(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating RADIUS Profile",
			"Could not create RADIUS Profile, unexpected error: "+err.Error(),
		)
		return
	}

	parseRadiusProfileResourceJson(*radiusProfile, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *radiusProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_radius_profile.RadiusProfileModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed RADIUS Profile value from Unifi
	radiusProfile, err := r.client.GetRADIUSProfile(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading RADIUS Profile",
			"Could not read RADIUS Profile ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	parseRadiusProfileResourceJson(*radiusProfile, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *radiusProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_radius_profile.RadiusProfileModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.RADIUSProfile
	parseRadiusProfileResourceModel(data, &body)
	radiusProfile, err := r.client.UpdateRADIUSProfile(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating RADIUS Profile",
			"Could not create RADIUS Profile, unexpected error: "+err.Error(),
		)
		return
	}

	parseRadiusProfileResourceJson(*radiusProfile, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *radiusProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_radius_profile.RadiusProfileModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing RADIUS Profile
	err := r.client.DeleteRADIUSProfile(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting RADIUS Profile",
			"Could not delete RADIUS Profile, unexpected error: "+err.Error(),
		)
		return
	}
}

func parseRadiusProfileResourceJson(json unifi.RADIUSProfile, model *resource_radius_profile.RadiusProfileModel) {
	model.Id = types.StringValue(json.ID)
	model.Name = types.StringValue(json.Name)
}

func parseRadiusProfileResourceModel(model resource_radius_profile.RadiusProfileModel, json *unifi.RADIUSProfile) {
	json.ID = model.Id.ValueString()
	json.Name = model.Name.ValueString()
}
