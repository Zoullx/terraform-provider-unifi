// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_dynamic_dns

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DynamicDnsDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host_name": schema.StringAttribute{
				Computed:            true,
				Description:         "The host name to update in the Dynamic DNS service.",
				MarkdownDescription: "The host name to update in the Dynamic DNS service.",
			},
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the Dynamic DNS to look up.",
				MarkdownDescription: "The ID of the Dynamic DNS to look up.",
			},
			"interface": schema.StringAttribute{
				Computed:            true,
				Description:         "The interface for the Dynamic DNS. Can be `wan` or `wan2`.",
				MarkdownDescription: "The interface for the Dynamic DNS. Can be `wan` or `wan2`.",
			},
			"login": schema.StringAttribute{
				Computed:            true,
				Description:         "The login username for the Dynamic DNS service.",
				MarkdownDescription: "The login username for the Dynamic DNS service.",
			},
			"password": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "The password for the Dynamic DNS service.",
				MarkdownDescription: "The password for the Dynamic DNS service.",
			},
			"server": schema.StringAttribute{
				Computed:            true,
				Description:         "The server for the Dynamic DNS service",
				MarkdownDescription: "The server for the Dynamic DNS service",
			},
			"service": schema.StringAttribute{
				Computed:            true,
				Description:         "The Dynamic DNS service provider, various values are supported (for example `dyndns`, etc.).",
				MarkdownDescription: "The Dynamic DNS service provider, various values are supported (for example `dyndns`, etc.).",
			},
			"site": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the site the Dynamic DNS is associated with.",
				MarkdownDescription: "The name of the site the Dynamic DNS is associated with.",
			},
		},
	}
}

type DynamicDnsModel struct {
	HostName  types.String `tfsdk:"host_name"`
	Id        types.String `tfsdk:"id"`
	Interface types.String `tfsdk:"interface"`
	Login     types.String `tfsdk:"login"`
	Password  types.String `tfsdk:"password"`
	Server    types.String `tfsdk:"server"`
	Service   types.String `tfsdk:"service"`
	Site      types.String `tfsdk:"site"`
}
