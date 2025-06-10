package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zoullx/terraform-provider-unifi/internal/resource_port_forward"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zoullx/unifi-go/unifi"
)

var (
	_ resource.Resource                = &portForwardResource{}
	_ resource.ResourceWithConfigure   = &portForwardResource{}
	_ resource.ResourceWithImportState = &portForwardResource{}
)

func NewPortForwardResource() resource.Resource {
	return &portForwardResource{}
}

type portForwardResource struct {
	client *unifi.Client
}

func (r *portForwardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forward"
}

func (r *portForwardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_port_forward.PortForwardResourceSchema(ctx)
}

func (r *portForwardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *portForwardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *portForwardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_port_forward.PortForwardModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.PortForward
	parsePortForwardResourceModel(data, &body)
	portForward, err := r.client.CreatePortForward(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Port Forward",
			"Could not create Port Forward, unexpected error: "+err.Error(),
		)
		return
	}

	parsePortForwardResourceJson(*portForward, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portForwardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_port_forward.PortForwardModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Port Forward value from Unifi
	portForward, err := r.client.GetPortForward(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Port Forward",
			"Could not read Port Forward ID "+data.Id.ValueString()+"; "+err.Error(),
		)
		return
	}

	parsePortForwardResourceJson(*portForward, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portForwardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_port_forward.PortForwardModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var body unifi.PortForward
	parsePortForwardResourceModel(data, &body)
	portForward, err := r.client.UpdatePortForward(ctx, data.Site.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Port Forward",
			"Could not create Port Forward, unexpected error: "+err.Error(),
		)
		return
	}

	parsePortForwardResourceJson(*portForward, &data)

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portForwardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_port_forward.PortForwardModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Port Forward
	err := r.client.DeletePortForward(ctx, data.Site.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Port Forward",
			"Could not delete Port Forward, unexpected error: "+err.Error(),
		)
		return
	}
}

func parsePortForwardResourceJson(json unifi.PortForward, model *resource_port_forward.PortForwardModel) {
	model.Id = types.StringValue(json.ID)
	model.SiteId = types.StringValue(json.SiteID)
	model.DstPort = types.StringValue(json.DstPort)
	model.Enabled = types.BoolValue(json.Enabled)
	model.FwdIp = types.StringValue(json.Fwd)
	model.FwdPort = types.StringValue(json.FwdPort)
	model.Log = types.BoolValue(json.Log)
	model.Name = types.StringValue(json.Name)
	model.PortForwardInterface = types.StringValue(json.PfwdInterface)
	model.Protocol = types.StringValue(json.Proto)
	model.SrcIp = types.StringValue(json.Src)
}

func parsePortForwardResourceModel(model resource_port_forward.PortForwardModel, json *unifi.PortForward) {
	json.ID = model.Id.ValueString()
	json.SiteID = model.SiteId.ValueString()
	json.DstPort = model.DstPort.ValueString()
	json.Enabled = model.Enabled.ValueBool()
	json.Fwd = model.FwdIp.ValueString()
	json.FwdPort = model.FwdPort.ValueString()
	json.Log = model.Log.ValueBool()
	json.Name = model.Name.ValueString()
	json.PfwdInterface = model.PortForwardInterface.ValueString()
	json.Proto = model.Protocol.ValueString()
	json.Src = model.SrcIp.ValueString()
}
