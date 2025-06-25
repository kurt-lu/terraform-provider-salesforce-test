// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nimajalali/go-force/force"
)

type userRoleResource struct {
	client *force.ForceApi
}

var _ resource.Resource = &userRoleResource{}

func (r *userRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "salesforce_user_role"
}

func (r *userRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User Role Resource for the Salesforce Provider",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the resource.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the role. Corresponds to Label on the user interface.",
				Required:    true,
			},
			"developer_name": schema.StringAttribute{
				Description: "The unique name of the object in the API. This name can contain only underscores and alphanumeric characters, and must be unique in your org. It must begin with a letter, not include spaces, not end with an underscore, and not contain two consecutive underscores. In managed packages, this field prevents naming conflicts on package installations. With this field, a developer can change the object's name in a managed package and the changes are reflected in a subscriber's organization. Corresponds to Role Name in the user interface.",
				Required:    true,
			},
			"parent_role_id": schema.StringAttribute{
				Description: "The ID of the parent role.",
				Optional:    true,
			},
		},
	}
}

type userRoleResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DeveloperName types.String `tfsdk:"developer_name"`
	ParentRoleId  types.String `tfsdk:"parent_role_id"`
}

// Custom UserRole struct that implements force.SObject
type customUserRole struct {
	Name          string `json:"Name"`
	DeveloperName string `json:"DeveloperName"`
	ParentRoleId  string `json:"ParentRoleId,omitempty"`
}

func (r customUserRole) ApiName() string {
	return "UserRole"
}

func (r customUserRole) ExternalIdApiName() string {
	return ""
}

func (r *userRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userRoleResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := customUserRole{
		Name:          data.Name.ValueString(),
		DeveloperName: data.DeveloperName.ValueString(),
	}
	if !data.ParentRoleId.IsNull() {
		role.ParentRoleId = data.ParentRoleId.ValueString()
	}

	sfResp, err := r.client.InsertSObject(role)
	if err != nil {
		resp.Diagnostics.AddError("Error Inserting User Role", err.Error())
		return
	}
	data.Id = types.StringValue(sfResp.Id)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *userRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userRoleResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var role customUserRole
	if err := r.client.GetSObject(data.Id.ValueString(), nil, &role); err != nil {
		resp.Diagnostics.AddError("Error Getting User Role", err.Error())
		return
	}

	data.Name = types.StringValue(role.Name)
	data.DeveloperName = types.StringValue(role.DeveloperName)
	data.ParentRoleId = types.StringValue(role.ParentRoleId)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *userRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userRoleResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := customUserRole{
		Name:          data.Name.ValueString(),
		DeveloperName: data.DeveloperName.ValueString(),
	}
	if !data.ParentRoleId.IsNull() {
		role.ParentRoleId = data.ParentRoleId.ValueString()
	}

	if err := r.client.UpdateSObject(data.Id.ValueString(), role); err != nil {
		resp.Diagnostics.AddError("Error Updating User Role", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *userRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userRoleResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSObject(data.Id.ValueString(), nil); err != nil {
		resp.Diagnostics.AddError("Error Deleting User Role", err.Error())
		return
	}
}
