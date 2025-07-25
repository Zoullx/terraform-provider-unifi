// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_setting_radius

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func SettingRadiusDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"accounting_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Enable RADIUS accounting.",
				MarkdownDescription: "Enable RADIUS accounting.",
			},
			"accounting_port": schema.Int64Attribute{
				Computed:            true,
				Description:         "The port for accounting communications.",
				MarkdownDescription: "The port for accounting communications.",
			},
			"auth_port": schema.Int64Attribute{
				Computed:            true,
				Description:         "The port for authentication communications.",
				MarkdownDescription: "The port for authentication communications.",
			},
			"enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "RADIUS server enabled.",
				MarkdownDescription: "RADIUS server enabled.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the Setting RADIUS.",
				MarkdownDescription: "The ID of the Setting RADIUS.",
			},
			"interim_update_interval": schema.Int64Attribute{
				Computed:            true,
				Description:         "Statistics will be collected from connected clients at this interval.",
				MarkdownDescription: "Statistics will be collected from connected clients at this interval.",
			},
			"secret": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "RADIUS secret passphrase.",
				MarkdownDescription: "RADIUS secret passphrase.",
			},
			"site": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the site the Setting RADIUS is associated with.",
				MarkdownDescription: "The name of the site the Setting RADIUS is associated with.",
			},
			"tunneled_reply": schema.BoolAttribute{
				Computed:            true,
				Description:         "Encrypt communication between the server and the client.",
				MarkdownDescription: "Encrypt communication between the server and the client.",
			},
		},
	}
}

type SettingRadiusModel struct {
	AccountingEnabled     types.Bool   `tfsdk:"accounting_enabled"`
	AccountingPort        types.Int64  `tfsdk:"accounting_port"`
	AuthPort              types.Int64  `tfsdk:"auth_port"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	Id                    types.String `tfsdk:"id"`
	InterimUpdateInterval types.Int64  `tfsdk:"interim_update_interval"`
	Secret                types.String `tfsdk:"secret"`
	Site                  types.String `tfsdk:"site"`
	TunneledReply         types.Bool   `tfsdk:"tunneled_reply"`
}
