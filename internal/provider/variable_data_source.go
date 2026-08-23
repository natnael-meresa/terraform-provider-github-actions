package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/natnael-meresa/terraform-provider-github-actions/internal/client"
)

var _ datasource.DataSource = &VariableDataSource{}

type VariableDataSource struct {
	client *client.Client
}

type VariableDataSourceModel struct {
	Repository types.String `tfsdk:"repository"`
	Name       types.String `tfsdk:"name"`
	Value      types.String `tfsdk:"value"`
}

func NewVariableDataSource() datasource.DataSource {
	return &VariableDataSource{}
}

func (d *VariableDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (d *VariableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"repository": schema.StringAttribute{
				Required:    true,
				Description: "Github repository in owner/repository format.",
			},

			"name": schema.StringAttribute{
				Required:    true,
				Description: "Github Actions variable name.",
			},

			"value": schema.StringAttribute{
				Computed:    true,
				Description: "Value of the Github Actions variable.",
			},
		},
	}
}

func (d *VariableDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *VariableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VariableDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	parts := strings.SplitN(
		data.Repository.ValueString(),
		"/",
		2,
	)

	if len(parts) != 2 {
		resp.Diagnostics.AddAttributeError(
			path.Root("repository"),
			"Invalid Repository",
			"Repository must be in owner/repository format.",
		)

		return
	}

	variable, err := d.client.GetVariable(
		ctx,
		parts[0],
		parts[1],
		data.Name.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read variable",
			err.Error(),
		)
		return
	}

	data.Value = types.StringValue(variable.Value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
