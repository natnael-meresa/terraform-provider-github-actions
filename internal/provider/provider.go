// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure GitHubActionsProvider satisfies various provider interfaces.
var _ provider.Provider = &GithubActionsProvider{}
var _ provider.ProviderWithFunctions = &GithubActionsProvider{}
var _ provider.ProviderWithEphemeralResources = &GithubActionsProvider{}
var _ provider.ProviderWithActions = &GithubActionsProvider{}

type GithubActionsProvider struct {
	version string
}

// GithubActionsProviderModel describes the provider data model.
type GithubActionsProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *GithubActionsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "github_actions"
	resp.Version = p.version
}

func (p *GithubActionsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *GithubActionsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data GithubActionsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *GithubActionsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVariableResource,
	}
}

func (p *GithubActionsProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *GithubActionsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVariableDataSource,
	}
}

func (p *GithubActionsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *GithubActionsProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GithubActionsProvider{
			version: version,
		}
	}
}
