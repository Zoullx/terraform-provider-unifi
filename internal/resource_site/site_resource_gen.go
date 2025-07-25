// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_site

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func SiteResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The description of the Site.",
				MarkdownDescription: "The description of the Site.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the Site to look up.",
				MarkdownDescription: "The ID of the Site to look up.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_update": schema.StringAttribute{
				Computed:            true,
				Description:         "Timestamp of the last Terraform update of the Site.",
				MarkdownDescription: "Timestamp of the last Terraform update of the Site.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the Site.",
				MarkdownDescription: "The name of the Site.",
			},
		},
	}
}

type SiteModel struct {
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	LastUpdate  types.String `tfsdk:"last_update"`
	Name        types.String `tfsdk:"name"`
}
