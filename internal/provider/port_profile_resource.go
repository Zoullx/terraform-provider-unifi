package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_port_profile"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &portProfileResource{}
	_ resource.ResourceWithConfigure   = &portProfileResource{}
	_ resource.ResourceWithImportState = &portProfileResource{}
)

func NewPortProfileResource() resource.Resource {
	return &portProfileResource{}
}

type portProfileResource struct {
	client unifi.Client
}

func (r *portProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_profile"
}

func (r *portProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_port_profile.PortProfileResourceSchema(ctx)
}

func (r *portProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *portProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *portProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_port_profile.PortProfileModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.PortProfile
	parsePortProfileResourceModel(data, &body)
	portProfile, err := r.client.CreatePortProfile(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Port Profile",
			"Could not create Port Profile, unexpected error: "+err.Error(),
		)
		return
	}

	parsePortProfileResourceJson(*portProfile, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_port_profile.PortProfileModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Port Profile value from Unifi
	portProfile, err := r.client.GetPortProfile(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Port Profile",
			"Could not read Port Profile ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	parsePortProfileResourceJson(*portProfile, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_port_profile.PortProfileModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.PortProfile
	parsePortProfileResourceModel(data, &body)
	portProfile, err := r.client.UpdatePortProfile(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Port Profile",
			"Could not create Port Profile, unexpected error: "+err.Error(),
		)
		return
	}

	parsePortProfileResourceJson(*portProfile, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_port_profile.PortProfileModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Port Profile
	err := r.client.DeletePortProfile(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Port Profile",
			"Could not delete Port Profile, unexpected error: "+err.Error(),
		)
		return
	}
}

func parsePortProfileResourceJson(json unifi.PortProfile, model *resource_port_profile.PortProfileModel) {
	model.Id = types.StringValue(json.ID)
	model.Name = types.StringValue(json.Name)
}

func parsePortProfileResourceModel(model resource_port_profile.PortProfileModel, json *unifi.PortProfile) {
	json.ID = model.Id.ValueString()
	json.Name = model.Name.ValueString()
}
