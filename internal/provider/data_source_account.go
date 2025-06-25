// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nimajalali/go-force/force"
	"github.com/nimajalali/go-force/sobjects"
)

type accountDataSource struct {
	client *force.ForceApi
}

var _ datasource.DataSource = &accountDataSource{}

func (d *accountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "salesforce_account"
}

func (d *accountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Account Data Source for the Salesforce Provider",
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
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Account type (e.g., Customer, Prospect, Partner).",
				Computed:    true,
			},
			"industry": schema.StringAttribute{
				Description: "Industry (e.g., Technology, Healthcare, Finance).",
				Computed:    true,
			},
			"phone": schema.StringAttribute{
				Description: "Phone number.",
				Computed:    true,
			},
			"website": schema.StringAttribute{
				Description: "Website URL.",
				Computed:    true,
			},
		},
	}
}

type accountDataModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	AccountNumber types.String `tfsdk:"account_number"`
	Type          types.String `tfsdk:"type"`
	Industry      types.String `tfsdk:"industry"`
	Phone         types.String `tfsdk:"phone"`
	Website       types.String `tfsdk:"website"`
}

type accountQueryResponse struct {
	sobjects.BaseQuery
	Records []struct {
		Id            string `json:"Id"`
		Name          string `json:"Name"`
		AccountNumber string `json:"AccountNumber"`
		Type          string `json:"Type"`
		Industry      string `json:"Industry"`
		Phone         string `json:"Phone"`
		Website       string `json:"Website"`
	}
}

func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data accountDataModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var query accountQueryResponse
	nameFilter := fmt.Sprintf("Name = '%s'", data.Name.ValueString())
	if err := d.client.Query(force.BuildQuery("Id, Name, AccountNumber, Type, Industry, Phone, Website", "Account", []string{nameFilter}), &query); err != nil {
		resp.Diagnostics.AddError("Error Getting Account", err.Error())
		return
	}
	if len(query.Records) == 0 {
		resp.Diagnostics.AddError("Error Getting Account", fmt.Sprintf("No Account where %s", nameFilter))
		return
	}

	record := query.Records[0]
	data.Id = types.StringValue(record.Id)
	data.Name = types.StringValue(record.Name)
	data.AccountNumber = types.StringValue(record.AccountNumber)
	data.Type = types.StringValue(record.Type)
	data.Industry = types.StringValue(record.Industry)
	data.Phone = types.StringValue(record.Phone)
	data.Website = types.StringValue(record.Website)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
} 