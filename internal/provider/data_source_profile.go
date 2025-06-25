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

type profileDataSource struct {
	client *force.ForceApi
}

var _ datasource.DataSource = &profileDataSource{}

func (d *profileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "salesforce_profile"
}

func (d *profileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Profile Data Source for the Salesforce Provider",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the resource.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the profile.",
				Required:    true,
			},
		},
	}
}

type profileDataModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type profileQueryResponse struct {
	sobjects.BaseQuery
	Records []struct {
		Id   string `json:"Id"`
		Name string `json:"Name"`
	}
}

func (d *profileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data profileDataModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var query profileQueryResponse
	nameFilter := fmt.Sprintf("Name = '%s'", data.Name.ValueString())
	if err := d.client.Query(force.BuildQuery("Id, Name", "Profile", []string{nameFilter}), &query); err != nil {
		resp.Diagnostics.AddError("Error Getting Profile", err.Error())
		return
	}
	if len(query.Records) == 0 {
		resp.Diagnostics.AddError("Error Getting Profile", fmt.Sprintf("No Profile where %s", nameFilter))
		return
	}

	record := query.Records[0]
	data.Id = types.StringValue(record.Id)
	data.Name = types.StringValue(record.Name)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
