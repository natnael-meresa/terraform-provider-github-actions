package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/natnael-meresa/terraform-provider-github-actions/internal/client"
)

var _ resource.Resource = &VariableResource{}
var _ resource.ResourceWithImportState = &VariableResource{}

type VariableResource struct {
	client *client.Client
}

type VariableResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Repository types.String `tfsdk:"repository"`
	Name       types.String `tfsdk:"name"`
	Value      types.String `tfsdk:"value"`
}

func NewVariableResource() resource.Resource {
	return &VariableResource{}
}

func (r *VariableResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (r *VariableResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"repository": schema.StringAttribute{
				Required:    true,
				Description: "Github repository in owner/repository format.",
			},

			"name": schema.SetAttribute{
				Required:    true,
				Description: "Github Actions variable name.",
			},

			"value": schema.SetAttribute{
				Required:    true,
				Description: "Value of the Github Actions variable.",
			},
		},
	}
}

func (r *VariableResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *VariableResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data VariableResourceModel

	resp.Diagnostics.Append(
		req.Plan.Get(ctx, &data)...,
	)

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
			"Invalid repository",
			"Expected repository to be in owner/repository format.",
		)
		return
	}

	owner := parts[0]
	repo := parts[1]

	err := r.client.CreateVariable(
		ctx,
		owner,
		repo,
		client.Variable{
			Name:  data.Name.ValueString(),
			Value: data.Value.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create github actions variable",
			err.Error(),
		)

		return
	}

	tflog.Trace(ctx, "created a resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *VariableResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data VariableResourceModel

	resp.Diagnostics.Append(
		req.State.Get(ctx, &data)...,
	)

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

	variable, err := r.client.GetVariable(
		ctx,
		parts[0],
		parts[1],
		data.Name.ValueString(),
	)

	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok &&
			apiErr.StatusCode == http.StatusNotFound {

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"Unable to Read GitHub Actions Variable",
			err.Error(),
		)

		return
	}

	data.Value = types.StringValue(variable.Value)

	resp.Diagnostics.Append(
		resp.State.Set(ctx, &data)...,
	)
}

func (r *VariableResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data VariableResourceModel

	resp.Diagnostics.Append(
		req.Plan.Get(ctx, &data)...,
	)

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

	err := r.client.UpdateVariable(
		ctx,
		parts[0],
		parts[1],
		data.Name.ValueString(),
		data.Value.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update GitHub Actions Variable",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(ctx, &data)...,
	)
}

func (r *VariableResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data VariableResourceModel

	resp.Diagnostics.Append(
		req.State.Get(ctx, &data)...,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	parts := strings.SplitN(
		data.Repository.ValueString(),
		"/",
		2,
	)

	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid repository",
			"Expected repository to be in owner/repository format.",
		)
		return
	}

	err := r.client.DeleteVariable(
		ctx,
		parts[0],
		parts[1],
		data.Name.ValueString(),
	)

	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok &&
			apiErr.StatusCode == http.StatusNotFound {
			return
		}

		resp.Diagnostics.AddError(
			"Unable to delete github actions variable",
			err.Error(),
		)

		return
	}
}

func (r *VariableResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
