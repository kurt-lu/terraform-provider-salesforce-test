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

type accountResource struct {
	client *force.ForceApi
}

var _ resource.Resource = &accountResource{}

func (r *accountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "salesforce_account"
}

func (r *accountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Account Resource for the Salesforce Provider",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the resource.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the account.",
				Required:    true,
			},
			"account_number": schema.StringAttribute{
				Description: "Account number.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Account type (e.g., Customer, Prospect, Partner).",
				Optional:    true,
			},
			"industry": schema.StringAttribute{
				Description: "Industry (e.g., Technology, Healthcare, Finance).",
				Optional:    true,
			},
			"phone": schema.StringAttribute{
				Description: "Phone number.",
				Optional:    true,
			},
			"website": schema.StringAttribute{
				Description: "Website URL.",
				Optional:    true,
			},
		},
	}
}

type accountResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	AccountNumber types.String `tfsdk:"account_number"`
	Type          types.String `tfsdk:"type"`
	Industry      types.String `tfsdk:"industry"`
	Phone         types.String `tfsdk:"phone"`
	Website       types.String `tfsdk:"website"`
}

// Custom Account struct that implements force.SObject
type customAccount struct {
	Name            string `json:"Name"`
	AccountNumber   string `json:"AccountNumber,omitempty"`
	Type            string `json:"Type,omitempty"`
	Industry        string `json:"Industry,omitempty"`
	Phone           string `json:"Phone,omitempty"`
	Website         string `json:"Website,omitempty"`
}

func (a customAccount) ApiName() string {
	return "Account"
}

func (a customAccount) ExternalIdApiName() string {
	return ""
}

func (r *accountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data accountResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	account := customAccount{
		Name: data.Name.ValueString(),
	}
	if !data.AccountNumber.IsNull() {
		account.AccountNumber = data.AccountNumber.ValueString()
	}
	if !data.Type.IsNull() {
		account.Type = data.Type.ValueString()
	}
	if !data.Industry.IsNull() {
		account.Industry = data.Industry.ValueString()
	}
	if !data.Phone.IsNull() {
		account.Phone = data.Phone.ValueString()
	}
	if !data.Website.IsNull() {
		account.Website = data.Website.ValueString()
	}

	sfResp, err := r.client.InsertSObject(account)
	if err != nil {
		resp.Diagnostics.AddError("Error Inserting Account", err.Error())
		return
	}
	data.Id = types.StringValue(sfResp.Id)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *accountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data accountResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var account customAccount
	if err := r.client.GetSObject(data.Id.ValueString(), nil, &account); err != nil {
		resp.Diagnostics.AddError("Error Getting Account", err.Error())
		return
	}

	data.Name = types.StringValue(account.Name)
	data.AccountNumber = types.StringValue(account.AccountNumber)
	data.Type = types.StringValue(account.Type)
	data.Industry = types.StringValue(account.Industry)
	data.Phone = types.StringValue(account.Phone)
	data.Website = types.StringValue(account.Website)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *accountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data accountResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	account := customAccount{
		Name: data.Name.ValueString(),
	}
	if !data.AccountNumber.IsNull() {
		account.AccountNumber = data.AccountNumber.ValueString()
	}
	if !data.Type.IsNull() {
		account.Type = data.Type.ValueString()
	}
	if !data.Industry.IsNull() {
		account.Industry = data.Industry.ValueString()
	}
	if !data.Phone.IsNull() {
		account.Phone = data.Phone.ValueString()
	}
	if !data.Website.IsNull() {
		account.Website = data.Website.ValueString()
	}

	if err := r.client.UpdateSObject(data.Id.ValueString(), account); err != nil {
		resp.Diagnostics.AddError("Error Updating Account", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *accountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data accountResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSObject(data.Id.ValueString(), nil); err != nil {
		resp.Diagnostics.AddError("Error Deleting Account", err.Error())
		return
	}
} 