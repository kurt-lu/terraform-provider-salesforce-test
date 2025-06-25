// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nimajalali/go-force/force"
	"github.com/nimajalali/go-force/sobjects"
)

type profileResource struct {
	client *force.ForceApi
}

var _ resource.Resource = &profileResource{}

func (r *profileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "salesforce_profile"
}

func (r *profileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Profile Resource for the Salesforce Provider. Please note that Users must have a Profile assigned to them, Profiles cannot be deleted if a User is assigned to it, and Salesforce does not allow the deletion of Users, only deactivation. Terraform will warn after destroy of a User that it has only been deactivated and now removed from state. A common issue with this pattern is a Profile and User created in tandem will fail to delete the Profile on destroy due to the lingering assignment. Should you wish to destroy a created Profile, it's advised that an apply that moves all affected Users to a static Profile be run first, after which the Profile can be safely destroyed.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the resource.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the profile.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the profile.",
				Optional:    true,
			},
			"user_license_id": schema.StringAttribute{
				Description: "ID of the UserLicense associated with this profile. Forces replacement if updated.",
				Required:    true,
			},
			"permissions": schema.MapAttribute{
				Description: "Map of permissions for the profile. At this time specific permissions can only be set, the comprehensive list will not be read from Salesforce. The keys should follow Salesforce 'SnakeCase' format however the 'Permissions' prefix should be omitted. Permissions will not import to state due to a technical limitation, you will need to run a subsequent apply if you have permissions set in config during import.",
				Optional:    true,
				ElementType: types.BoolType,
			},
		},
	}
}

type profileResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	UserLicenseId types.String `tfsdk:"user_license_id"`
	Permissions   types.Map    `tfsdk:"permissions"`
}

// Custom Profile struct that includes permissions
type customProfile struct {
	sobjects.Profile
	Permissions map[string]bool `json:"-"`
}

func (p customProfile) ApiName() string {
	return "Profile"
}

func (p customProfile) ExternalIdApiName() string {
	return ""
}

func (r *profileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data profileResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile := &sobjects.Profile{
		Name:          data.Name.ValueString(),
		UserLicenseId: data.UserLicenseId.ValueString(),
	}
	if !data.Description.IsNull() {
		profile.Description = data.Description.ValueString()
	}

	// Handle permissions separately since they're not in sobjects.Profile
	if !data.Permissions.IsNull() && data.Permissions.Elements() != nil {
		// Create a custom profile with permissions
		customProfile := customProfile{
			Profile: *profile,
		}
		customProfile.Permissions = make(map[string]bool)
		for k, v := range data.Permissions.Elements() {
			if b, ok := v.(types.Bool); ok {
				customProfile.Permissions["Permissions"+k] = b.ValueBool()
			}
		}
		sfResp, err := r.client.InsertSObject(customProfile)
		if err != nil {
			resp.Diagnostics.AddError("Error Inserting Profile", err.Error())
			return
		}
		data.Id = types.StringValue(sfResp.Id)
	} else {
		sfResp, err := r.client.InsertSObject(profile)
		if err != nil {
			resp.Diagnostics.AddError("Error Inserting Profile", err.Error())
			return
		}
		data.Id = types.StringValue(sfResp.Id)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *profileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data profileResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use custom profile to get permissions since they're not in sobjects.Profile
	var customProfile customProfile
	if err := r.client.GetSObject(data.Id.ValueString(), nil, &customProfile); err != nil {
		resp.Diagnostics.AddError("Error Getting Profile", err.Error())
		return
	}

	data.Name = types.StringValue(customProfile.Name)
	data.Description = types.StringValue(customProfile.Description)
	data.UserLicenseId = types.StringValue(customProfile.UserLicenseId)
	// expand permissions
	permissions := make(map[string]attr.Value)
	for k, v := range customProfile.Permissions {
		permissions[k] = types.BoolValue(v)
	}
	data.Permissions = types.MapValueMust(types.BoolType, permissions)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *profileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data profileResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile := &sobjects.Profile{
		Name:          data.Name.ValueString(),
		UserLicenseId: data.UserLicenseId.ValueString(),
	}
	if !data.Description.IsNull() {
		profile.Description = data.Description.ValueString()
	}

	// Handle permissions separately since they're not in sobjects.Profile
	if !data.Permissions.IsNull() && data.Permissions.Elements() != nil {
		// Create a custom profile with permissions
		customProfile := customProfile{
			Profile: *profile,
		}
		customProfile.Permissions = make(map[string]bool)
		for k, v := range data.Permissions.Elements() {
			if b, ok := v.(types.Bool); ok {
				customProfile.Permissions["Permissions"+k] = b.ValueBool()
			}
		}
		if err := r.client.UpdateSObject(data.Id.ValueString(), customProfile); err != nil {
			resp.Diagnostics.AddError("Error Updating Profile", err.Error())
			return
		}
	} else {
		if err := r.client.UpdateSObject(data.Id.ValueString(), profile); err != nil {
			resp.Diagnostics.AddError("Error Updating Profile", err.Error())
			return
		}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *profileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data profileResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSObject(data.Id.ValueString(), nil); err != nil {
		resp.Diagnostics.AddError("Error Deleting Profile", err.Error())
		return
	}
}
