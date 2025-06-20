// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nimajalali/go-force/force"
)

type accountType struct {
}

func (accountType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Account Resource for the Salesforce Provider",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "ID of the resource.",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					staticComputed{},
				},
			},
			"name": {
				Description: "The name of the account.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					notEmptyString{},
				},
			},
			"account_number": {
				Description: "Account number.",
				Type:        types.StringType,
				Optional:    true,
			},
			"type": {
				Description: "Account type (e.g., Customer, Prospect, Partner).",
				Type:        types.StringType,
				Optional:    true,
			},
			"industry": {
				Description: "Industry (e.g., Technology, Healthcare, Finance).",
				Type:        types.StringType,
				Optional:    true,
			},
			"phone": {
				Description: "Phone number.",
				Type:        types.StringType,
				Optional:    true,
			},
			"website": {
				Description: "Website URL.",
				Type:        types.StringType,
				Optional:    true,
			},
		},
	}, nil
}

func (a accountType) NewResource(_ context.Context, prov tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, ok := prov.(*provider)
	if !ok {
		return nil, diag.Diagnostics{errorConvertingProvider(a)}
	}
	return &accountResource{
		client: provider.client,
	}, nil
}

type accountResource struct {
	client *force.ForceApi
}

type accountResourceData struct {
	Name          string       `tfsdk:"name"`
	AccountNumber *string      `tfsdk:"account_number"`
	Type          *string      `tfsdk:"type"`
	Industry      *string      `tfsdk:"industry"`
	Phone         *string      `tfsdk:"phone"`
	Website       *string      `tfsdk:"website"`
	Id            types.String `tfsdk:"id"`
}

func (a accountResourceData) ToMap(exclude ...string) accountMap {
	aMap := make(accountMap)
	if a.Name != "" {
		aMap["Name"] = a.Name
	}
	if a.AccountNumber != nil && *a.AccountNumber != "" {
		aMap["AccountNumber"] = *a.AccountNumber
	}
	if a.Type != nil && *a.Type != "" {
		aMap["Type"] = *a.Type
	}
	if a.Industry != nil && *a.Industry != "" {
		aMap["Industry"] = *a.Industry
	}
	if a.Phone != nil && *a.Phone != "" {
		aMap["Phone"] = *a.Phone
	}
	if a.Website != nil && *a.Website != "" {
		aMap["Website"] = *a.Website
	}
	// exclude keys, useful for update
	for _, k := range exclude {
		delete(aMap, k)
	}
	return aMap
}

type accountMap map[string]interface{}

func (a accountMap) ToStateData() accountResourceData {
	data := accountResourceData{
		Name: a["Name"].(string),
	}
	if accountNumber, ok := a["AccountNumber"]; ok && accountNumber != nil {
		accountNumberStr := accountNumber.(string)
		data.AccountNumber = &accountNumberStr
	}
	if accountType, ok := a["Type"]; ok && accountType != nil {
		accountTypeStr := accountType.(string)
		data.Type = &accountTypeStr
	}
	if industry, ok := a["Industry"]; ok && industry != nil {
		industryStr := industry.(string)
		data.Industry = &industryStr
	}
	if phone, ok := a["Phone"]; ok && phone != nil {
		phoneStr := phone.(string)
		data.Phone = &phoneStr
	}
	if website, ok := a["Website"]; ok && website != nil {
		websiteStr := website.(string)
		data.Website = &websiteStr
	}
	return data
}

func (accountMap) ApiName() string {
	return "Account"
}

func (accountMap) ExternalIdApiName() string {
	return ""
}

func (a *accountResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data accountResourceData
	if diags := req.Plan.Get(ctx, &data); diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	sfResp, err := a.client.InsertSObject(data.ToMap())
	if err != nil {
		resp.Diagnostics.AddError("Error Inserting Account", err.Error())
		return
	}
	data.Id = types.String{Value: sfResp.Id}

	resp.Diagnostics = resp.State.Set(ctx, &data)
}

func (a *accountResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data accountResourceData
	if diags := req.State.Get(ctx, &data); diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	var aMap accountMap
	if err := a.client.GetSObject(data.Id.Value, nil, &aMap); err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error Getting Account", err.Error())
		}
		return
	}

	d := aMap.ToStateData()
	// copy the ID back over
	d.Id = data.Id

	resp.Diagnostics = resp.State.Set(ctx, &d)
}

func (a *accountResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data accountResourceData
	if diags := req.Plan.Get(ctx, &data); diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	if err := a.client.UpdateSObject(data.Id.Value, data.ToMap()); err != nil {
		resp.Diagnostics.AddError("Error Updating Account", err.Error())
		return
	}

	resp.Diagnostics = resp.State.Set(ctx, &data)
}

func (a *accountResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data accountResourceData
	if diags := req.State.Get(ctx, &data); diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	if err := a.client.DeleteSObject(data.Id.Value); err != nil {
		resp.Diagnostics.AddError("Error Deleting Account", err.Error())
		return
	}
}

func (a *accountResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tfsdk.NewAttributePath().WithAttributeName("id"), req, resp)
} 