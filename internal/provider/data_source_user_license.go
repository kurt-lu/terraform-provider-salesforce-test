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

type userLicenseDataSource struct {
	client *force.ForceApi
}

var _ datasource.DataSource = &userLicenseDataSource{}

func (d *userLicenseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "salesforce_user_license"
}

func (d *userLicenseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User License Data Source for the Salesforce Provider",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the resource.",
				Computed:    true,
			},
			"license_definition_key": schema.StringAttribute{
				Description: "A string that uniquely identifies a particular user license. Valid options vary depending on organization type and configuration. For a complete list see https://developer.salesforce.com/docs/atlas.en-us.api.meta/api/sforce_api_objects_userlicense.htm",
				Required:    true,
			},
		},
	}
}

type userLicenseDataModel struct {
	Id                   types.String `tfsdk:"id"`
	LicenseDefinitionKey types.String `tfsdk:"license_definition_key"`
}

type userLicenseQueryResponse struct {
	sobjects.BaseQuery
	Records []struct {
		Id                   string `json:"Id"`
		LicenseDefinitionKey string `json:"LicenseDefinitionKey"`
	}
}

func (d *userLicenseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userLicenseDataModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var query userLicenseQueryResponse
	licenseDefinitionKeyFilter := fmt.Sprintf("LicenseDefinitionKey = '%s'", data.LicenseDefinitionKey.ValueString())
	if err := d.client.Query(force.BuildQuery("Id, LicenseDefinitionKey", "UserLicense", []string{licenseDefinitionKeyFilter}), &query); err != nil {
		resp.Diagnostics.AddError("Error Getting User License", err.Error())
		return
	}
	if len(query.Records) == 0 {
		resp.Diagnostics.AddError("Error Getting User License", fmt.Sprintf("No User License where %s", licenseDefinitionKeyFilter))
		return
	}

	record := query.Records[0]
	data.Id = types.StringValue(record.Id)
	data.LicenseDefinitionKey = types.StringValue(record.LicenseDefinitionKey)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
