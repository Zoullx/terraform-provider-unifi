// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_ap_groups

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ApGroupsDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ap_groups": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_macs": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							Description:         "The MAC addresses of the APs associated with the AP Group.",
							MarkdownDescription: "The MAC addresses of the APs associated with the AP Group.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "The ID of the AP Group.",
							MarkdownDescription: "The ID of the AP Group.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "The name of the AP Group.",
							MarkdownDescription: "The name of the AP Group.",
						},
					},
					CustomType: ApGroupsType{
						ObjectType: types.ObjectType{
							AttrTypes: ApGroupsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed:            true,
				Description:         "The list of AP Groups associated with the site.",
				MarkdownDescription: "The list of AP Groups associated with the site.",
			},
			"site": schema.StringAttribute{
				Optional:            true,
				Description:         "The name of the site the AP Groups are associated with.",
				MarkdownDescription: "The name of the site the AP Groups are associated with.",
			},
		},
	}
}

type ApGroupsModel struct {
	ApGroups types.List   `tfsdk:"ap_groups"`
	Site     types.String `tfsdk:"site"`
}

var _ basetypes.ObjectTypable = ApGroupsType{}

type ApGroupsType struct {
	basetypes.ObjectType
}

func (t ApGroupsType) Equal(o attr.Type) bool {
	other, ok := o.(ApGroupsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ApGroupsType) String() string {
	return "ApGroupsType"
}

func (t ApGroupsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	deviceMacsAttribute, ok := attributes["device_macs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`device_macs is missing from object`)

		return nil, diags
	}

	deviceMacsVal, ok := deviceMacsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`device_macs expected to be basetypes.ListValue, was: %T`, deviceMacsAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return nil, diags
	}

	idVal, ok := idAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.StringValue, was: %T`, idAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return nil, diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ApGroupsValue{
		DeviceMacs: deviceMacsVal,
		Id:         idVal,
		Name:       nameVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewApGroupsValueNull() ApGroupsValue {
	return ApGroupsValue{
		state: attr.ValueStateNull,
	}
}

func NewApGroupsValueUnknown() ApGroupsValue {
	return ApGroupsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewApGroupsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ApGroupsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ApGroupsValue Attribute Value",
				"While creating a ApGroupsValue value, a missing attribute value was detected. "+
					"A ApGroupsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ApGroupsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ApGroupsValue Attribute Type",
				"While creating a ApGroupsValue value, an invalid attribute value was detected. "+
					"A ApGroupsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ApGroupsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ApGroupsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ApGroupsValue Attribute Value",
				"While creating a ApGroupsValue value, an extra attribute value was detected. "+
					"A ApGroupsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ApGroupsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewApGroupsValueUnknown(), diags
	}

	deviceMacsAttribute, ok := attributes["device_macs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`device_macs is missing from object`)

		return NewApGroupsValueUnknown(), diags
	}

	deviceMacsVal, ok := deviceMacsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`device_macs expected to be basetypes.ListValue, was: %T`, deviceMacsAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewApGroupsValueUnknown(), diags
	}

	idVal, ok := idAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.StringValue, was: %T`, idAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewApGroupsValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	if diags.HasError() {
		return NewApGroupsValueUnknown(), diags
	}

	return ApGroupsValue{
		DeviceMacs: deviceMacsVal,
		Id:         idVal,
		Name:       nameVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewApGroupsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ApGroupsValue {
	object, diags := NewApGroupsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewApGroupsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ApGroupsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewApGroupsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewApGroupsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewApGroupsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewApGroupsValueMust(ApGroupsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ApGroupsType) ValueType(ctx context.Context) attr.Value {
	return ApGroupsValue{}
}

var _ basetypes.ObjectValuable = ApGroupsValue{}

type ApGroupsValue struct {
	DeviceMacs basetypes.ListValue   `tfsdk:"device_macs"`
	Id         basetypes.StringValue `tfsdk:"id"`
	Name       basetypes.StringValue `tfsdk:"name"`
	state      attr.ValueState
}

func (v ApGroupsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["device_macs"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.DeviceMacs.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["device_macs"] = val

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v ApGroupsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ApGroupsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ApGroupsValue) String() string {
	return "ApGroupsValue"
}

func (v ApGroupsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var deviceMacsVal basetypes.ListValue
	switch {
	case v.DeviceMacs.IsUnknown():
		deviceMacsVal = types.ListUnknown(types.StringType)
	case v.DeviceMacs.IsNull():
		deviceMacsVal = types.ListNull(types.StringType)
	default:
		var d diag.Diagnostics
		deviceMacsVal, d = types.ListValue(types.StringType, v.DeviceMacs.Elements())
		diags.Append(d...)
	}

	if diags.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			"device_macs": basetypes.ListType{
				ElemType: types.StringType,
			},
			"id":   basetypes.StringType{},
			"name": basetypes.StringType{},
		}), diags
	}

	attributeTypes := map[string]attr.Type{
		"device_macs": basetypes.ListType{
			ElemType: types.StringType,
		},
		"id":   basetypes.StringType{},
		"name": basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"device_macs": deviceMacsVal,
			"id":          v.Id,
			"name":        v.Name,
		})

	return objVal, diags
}

func (v ApGroupsValue) Equal(o attr.Value) bool {
	other, ok := o.(ApGroupsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.DeviceMacs.Equal(other.DeviceMacs) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	return true
}

func (v ApGroupsValue) Type(ctx context.Context) attr.Type {
	return ApGroupsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ApGroupsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"device_macs": basetypes.ListType{
			ElemType: types.StringType,
		},
		"id":   basetypes.StringType{},
		"name": basetypes.StringType{},
	}
}
