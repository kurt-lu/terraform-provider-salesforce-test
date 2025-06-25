// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/nimajalali/go-force/force"
	"github.com/nimajalali/go-force/sobjects"
)

var userDefaults = resourceDefaults{
	defaults: map[string]attr.Value{
		tftypes.NewAttributePath().WithAttributeName("email_encoding_key").String(): types.StringValue("UTF-8"),
		tftypes.NewAttributePath().WithAttributeName("language_locale_key").String(): types.StringValue("en_US"),
		tftypes.NewAttributePath().WithAttributeName("locale_sid_key").String(): types.StringValue("en_US"),
		tftypes.NewAttributePath().WithAttributeName("time_zone_sid_key").String(): types.StringValue("America/New_York"),
	},
}

type userResource struct {
	client *force.ForceApi
}

var _ resource.Resource = &userResource{}

func (r *userResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "salesforce_user"
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User Resource for the Salesforce Provider",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the resource.",
				Computed:    true,
			},
			"alias": schema.StringAttribute{
				Description: "The user's alias. For example, jsmith.",
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "The user's email address.",
				Required:    true,
			},
			"email_encoding_key": schema.StringAttribute{
				Description: "The email encoding for the user, such as ISO-8859-1 or UTF-8. Defaults to UTF-8.",
				Optional:    true,
				Computed:    true,
			},
			"language_locale_key": schema.StringAttribute{
				Description: "The user's language. Defaults to en_US.",
				Optional:    true,
				Computed:    true,
			},
			"last_name": schema.StringAttribute{
				Description: "The user's last name.",
				Required:    true,
			},
			"locale_sid_key": schema.StringAttribute{
				Description: "The value of the field affects formatting and parsing of values, especially numeric values, in the user interface. It doesn't affect the API. The field values are named according to the language, and the country if necessary, using two-letter ISO codes. The set of names is based on the ISO standard. You can also manually set a user's locale in the user interface, and then use that value for inserting or updating other users via the API. Defaults to en_US.",
				Optional:    true,
				Computed:    true,
			},
			"profile_id": schema.StringAttribute{
				Description: "ID of the user's Profile. Use this value to cache metadata based on profile.",
				Required:    true,
			},
			"time_zone_sid_key": schema.StringAttribute{
				Description: "A User time zone affects the offset used when displaying or entering times in the user interface. But the API doesn't use a User time zone when querying or setting values. Values for this field are named using region and key city, according to ISO standards. You can also manually set one User time zone in the user interface, and then use that value for creating or updating other User records via the API. Defaults to America/New_York.",
				Optional:    true,
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Contains the name that a user enters to log in to the API or the user interface. The value for this field must be in the form of an email address, using all lowercase characters. It must also be unique across all organizations. If you try to create or update a User with a duplicate value for this field, the operation is rejected. Each inserted User also counts as a license. Every organization has a maximum number of licenses. If you attempt to exceed the maximum number of licenses by inserting User records, the create request is rejected.",
				Required:    true,
			},
			"user_role_id": schema.StringAttribute{
				Description: "ID of the user's UserRole.",
				Optional:    true,
			},
			"reset_password": schema.BoolAttribute{
				Description: "Reset password and send an email to the user. No reset is performed if this field is omitted, is false, or was true and remained true on subsequent apply. Please set to false and then true in subsequent applies, or have it set to true on create to trigger the reset.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

type userResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Alias             types.String `tfsdk:"alias"`
	Email             types.String `tfsdk:"email"`
	EmailEncodingKey  types.String `tfsdk:"email_encoding_key"`
	LanguageLocaleKey types.String `tfsdk:"language_locale_key"`
	LastName          types.String `tfsdk:"last_name"`
	LocaleSidKey      types.String `tfsdk:"locale_sid_key"`
	ProfileID         types.String `tfsdk:"profile_id"`
	TimeZoneSidKey    types.String `tfsdk:"time_zone_sid_key"`
	Username          types.String `tfsdk:"username"`
	UserRoleId        types.String `tfsdk:"user_role_id"`
	ResetPassword     types.Bool   `tfsdk:"reset_password"`
}

// Custom User struct that includes UserRoleId
type customUser struct {
	sobjects.User
	UserRoleId string `json:"UserRoleId,omitempty"`
}

func (u customUser) ApiName() string {
	return "User"
}

func (u customUser) ExternalIdApiName() string {
	return ""
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := &sobjects.User{
		Alias:             data.Alias.ValueString(),
		Email:             data.Email.ValueString(),
		EmailEncodingKey:  data.EmailEncodingKey.ValueString(),
		LanguageLocaleKey: data.LanguageLocaleKey.ValueString(),
		LastName:          data.LastName.ValueString(),
		LocaleSidKey:      data.LocaleSidKey.ValueString(),
		ProfileId:         data.ProfileID.ValueString(),
		TimeZoneSidKey:    data.TimeZoneSidKey.ValueString(),
		Username:          data.Username.ValueString(),
	}
	
	// UserRoleId is not in sobjects.User, so we need to use a custom struct
	if !data.UserRoleId.IsNull() {
		customUser := customUser{
			User:       *user,
			UserRoleId: data.UserRoleId.ValueString(),
		}
		sfResp, err := r.client.InsertSObject(customUser)
		if err != nil {
			resp.Diagnostics.AddError("Error Inserting User", err.Error())
			return
		}
		data.Id = types.StringValue(sfResp.Id)
	} else {
		sfResp, err := r.client.InsertSObject(user)
		if err != nil {
			resp.Diagnostics.AddError("Error Inserting User", err.Error())
			return
		}
		data.Id = types.StringValue(sfResp.Id)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use custom user to get UserRoleId which is not in sobjects.User
	var customUser customUser
	if err := r.client.GetSObject(data.Id.ValueString(), nil, &customUser); err != nil {
		resp.Diagnostics.AddError("Error Getting User", err.Error())
		return
	}

	data.Alias = types.StringValue(customUser.Alias)
	data.Email = types.StringValue(customUser.Email)
	data.EmailEncodingKey = types.StringValue(customUser.EmailEncodingKey)
	data.LanguageLocaleKey = types.StringValue(customUser.LanguageLocaleKey)
	data.LastName = types.StringValue(customUser.LastName)
	data.LocaleSidKey = types.StringValue(customUser.LocaleSidKey)
	data.ProfileID = types.StringValue(customUser.ProfileId)
	data.TimeZoneSidKey = types.StringValue(customUser.TimeZoneSidKey)
	data.Username = types.StringValue(customUser.Username)
	data.UserRoleId = types.StringValue(customUser.UserRoleId)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := &sobjects.User{
		Alias:             data.Alias.ValueString(),
		Email:             data.Email.ValueString(),
		EmailEncodingKey:  data.EmailEncodingKey.ValueString(),
		LanguageLocaleKey: data.LanguageLocaleKey.ValueString(),
		LastName:          data.LastName.ValueString(),
		LocaleSidKey:      data.LocaleSidKey.ValueString(),
		ProfileId:         data.ProfileID.ValueString(),
		TimeZoneSidKey:    data.TimeZoneSidKey.ValueString(),
		Username:          data.Username.ValueString(),
	}

	// Always use custom user since UserRoleId is not in sobjects.User
	customUser := customUser{
		User: *user,
	}
	if !data.UserRoleId.IsNull() {
		customUser.UserRoleId = data.UserRoleId.ValueString()
	}

	if err := r.client.UpdateSObject(data.Id.ValueString(), customUser); err != nil {
		resp.Diagnostics.AddError("Error Updating User", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSObject(data.Id.ValueString(), nil); err != nil {
		resp.Diagnostics.AddError("Error Deleting User", err.Error())
		return
	}
}
