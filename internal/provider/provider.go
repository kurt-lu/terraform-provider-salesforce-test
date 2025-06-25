// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-salesforce/internal/auth"
	"github.com/nimajalali/go-force/force"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

type salesforceProvider struct {
	client *force.ForceApi
}

var _ provider.Provider = &salesforceProvider{}

func New() provider.Provider {
	return &salesforceProvider{}
}

func (p *salesforceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "salesforce"
}

func (p *salesforceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Provider for managing a Salesforce Organization",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "Client ID of the connected app. Corresponds to Consumer Key in the user interface. Can be specified with the environment variable SALESFORCE_CLIENT_ID.",
				Optional:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "Private Key associated to the public certificate that was uploaded to the connected app. This may point to a file location or be set directly. This should not be confused with the Consumer Secret in the user interface. Can be specified with the environment variable SALESFORCE_PRIVATE_KEY.",
				Optional:    true,
				Sensitive:   true,
			},
			"api_version": schema.StringAttribute{
				Description: "API version of the salesforce org in the format: MAJOR.MINOR (please omit any leading 'v'). The provider requires at least version 53.0. Can be specified with the environment variable SALESFORCE_API_VERSION.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Salesforce Username of a System Administrator like user for the provider to authenticate as. Can be specified with the environment variable SALESFORCE_USERNAME.",
				Optional:    true,
			},
			"login_url": schema.StringAttribute{
				Description: "Directs the authentication request, defaults to the production endpoint https://login.salesforce.com, should be set to https://test.salesforce.com for sandbox organizations. Can be specified with the environment variable SALESFORCE_LOGIN_URL.",
				Optional:    true,
			},
		},
	}
}

type providerDataModel struct {
	ClientId   types.String `tfsdk:"client_id"`
	PrivateKey types.String `tfsdk:"private_key"`
	ApiVersion types.String `tfsdk:"api_version"`
	Username   types.String `tfsdk:"username"`
	LoginUrl   types.String `tfsdk:"login_url"`
}

func (p *salesforceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if unset, fallback to env
	if config.ClientId.IsNull() {
		config.ClientId = types.StringValue(os.Getenv("SALESFORCE_CLIENT_ID"))
	}
	if config.PrivateKey.IsNull() {
		config.PrivateKey = types.StringValue(os.Getenv("SALESFORCE_PRIVATE_KEY"))
	}
	if config.ApiVersion.IsNull() {
		config.ApiVersion = types.StringValue(os.Getenv("SALESFORCE_API_VERSION"))
	}
	if config.Username.IsNull() {
		config.Username = types.StringValue(os.Getenv("SALESFORCE_USERNAME"))
	}
	if config.LoginUrl.IsNull() {
		config.LoginUrl = types.StringValue(os.Getenv("SALESFORCE_LOGIN_URL"))
	}

	// required if still unset
	if config.ClientId.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			pathRoot("client_id"),
			"Invalid provider config",
			"client_id must be set.",
		)
		return
	}
	if config.PrivateKey.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			pathRoot("private_key"),
			"Invalid provider config",
			"private_key must be set.",
		)
		return
	}
	if config.ApiVersion.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			pathRoot("api_version"),
			"Invalid provider config",
			"api_version must be set.",
		)
		return
	}
	if config.Username.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			pathRoot("username"),
			"Invalid provider config",
			"username must be set.",
		)
		return
	}

	client, err := auth.Client(auth.Config{
		ApiVersion: config.ApiVersion.ValueString(),
		Username:   config.Username.ValueString(),
		ClientId:   config.ClientId.ValueString(),
		PrivateKey: config.PrivateKey.ValueString(),
		LoginUrl:   config.LoginUrl.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating salesforce client", err.Error())
		return
	}
	p.client = client
}

func (p *salesforceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &profileDataSource{client: p.client} },
		func() datasource.DataSource { return &accountDataSource{client: p.client} },
		func() datasource.DataSource { return &userLicenseDataSource{client: p.client} },
	}
}

func (p *salesforceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &accountResource{client: p.client} },
		func() resource.Resource { return &profileResource{client: p.client} },
		func() resource.Resource { return &userResource{client: p.client} },
		func() resource.Resource { return &userRoleResource{client: p.client} },
	}
}

// Helper for attribute error paths
func pathRoot(attr string) path.Path {
	return path.Root(attr)
}
