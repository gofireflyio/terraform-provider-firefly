package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &workspaceLabelsResource{}
	_ resource.ResourceWithConfigure   = &workspaceLabelsResource{}
	_ resource.ResourceWithImportState = &workspaceLabelsResource{}
)

// NewWorkspaceLabelsResource is a helper function to simplify the provider implementation
func NewWorkspaceLabelsResource() resource.Resource {
	return &workspaceLabelsResource{}
}

// workspaceLabelsResource is the resource implementation
type workspaceLabelsResource struct {
	client *client.Client
}

// WorkspaceLabelsResourceModel describes the workspace labels resource data model
type WorkspaceLabelsResourceModel struct {
	ID           types.String `tfsdk:"id"`
	WorkspaceID  types.String `tfsdk:"workspace_id"`
	WorkspaceName types.String `tfsdk:"workspace_name"`
	Labels       types.List   `tfsdk:"labels"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name
func (r *workspaceLabelsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_labels"
}

// Schema defines the schema for the resource
func (r *workspaceLabelsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages labels for a Firefly workspace",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the workspace",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "The ID of the workspace to manage labels for",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"workspace_name": schema.StringAttribute{
				Description: "The name of the workspace",
				Computed:    true,
			},
			"labels": schema.ListAttribute{
				Description: "List of labels to assign to the workspace",
				Required:    true,
				ElementType: types.StringType,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the labels were last updated",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *workspaceLabelsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates or updates workspace labels
func (r *workspaceLabelsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan WorkspaceLabelsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the list of labels from the plan
	var labels []string
	diags = plan.Labels.ElementsAs(ctx, &labels, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the workspace labels
	updateResp, err := r.client.Workspaces.UpdateWorkspaceLabels(plan.WorkspaceID.ValueString(), labels)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Workspace Labels",
			fmt.Sprintf("Could not update labels for workspace ID %s: %s", plan.WorkspaceID.ValueString(), err),
		)
		return
	}

	// Update the plan with values from the response
	plan.ID = types.StringValue(updateResp.ID)
	plan.WorkspaceName = types.StringValue(updateResp.WorkspaceName)
	plan.UpdatedAt = types.StringValue(updateResp.UpdatedAt)
	
	// Set labels from the response
	plan.Labels = types.ListValueMust(types.StringType, listToValues(updateResp.Labels))

	// Set state to fully populated plan
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data
func (r *workspaceLabelsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state WorkspaceLabelsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Search for workspace by ID
	workspaceID := state.WorkspaceID.ValueString()
	
	// List workspaces (filtering will be done client-side since there's no direct get endpoint)
	workspaces, err := r.client.Workspaces.ListWorkspaces(&client.ListWorkspacesRequest{}, 0, 100)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workspace",
			fmt.Sprintf("Could not list workspaces: %s", err),
		)
		return
	}

	// Find the workspace with matching ID
	var workspace *client.Workspace
	for _, w := range workspaces {
		if w.WorkspaceID == workspaceID {
			workspace = &w
			break
		}
	}

	if workspace == nil {
		resp.Diagnostics.AddError(
			"Workspace Not Found",
			fmt.Sprintf("Workspace with ID %s not found", workspaceID),
		)
		return
	}

	// Update state with found workspace data
	state.ID = types.StringValue(workspace.ID)
	state.WorkspaceName = types.StringValue(workspace.WorkspaceName)
	
	// Set labels from the workspace
	state.Labels = types.ListValueMust(types.StringType, listToValues(workspace.Labels))
	
	// Set updated timestamp
	state.UpdatedAt = types.StringValue(workspace.UpdatedAt)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the workspace labels
func (r *workspaceLabelsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan WorkspaceLabelsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the list of labels from the plan
	var labels []string
	diags = plan.Labels.ElementsAs(ctx, &labels, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the workspace labels
	updateResp, err := r.client.Workspaces.UpdateWorkspaceLabels(plan.WorkspaceID.ValueString(), labels)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Workspace Labels",
			fmt.Sprintf("Could not update labels for workspace ID %s: %s", plan.WorkspaceID.ValueString(), err),
		)
		return
	}

	// Update the plan with values from the response
	plan.ID = types.StringValue(updateResp.ID)
	plan.WorkspaceName = types.StringValue(updateResp.WorkspaceName)
	plan.UpdatedAt = types.StringValue(updateResp.UpdatedAt)
	
	// Set labels from the response
	plan.Labels = types.ListValueMust(types.StringType, listToValues(updateResp.Labels))

	// Set state to fully populated plan
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete removes the labels from the workspace by setting an empty list
func (r *workspaceLabelsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state WorkspaceLabelsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Clear labels from workspace by setting an empty list
	_, err := r.client.Workspaces.UpdateWorkspaceLabels(state.WorkspaceID.ValueString(), []string{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Workspace Labels",
			fmt.Sprintf("Could not clear labels for workspace ID %s: %s", state.WorkspaceID.ValueString(), err),
		)
		return
	}

	// No additional state setting required for deletion
}

// ImportState imports a workspace labels resource by workspace ID
func (r *workspaceLabelsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Simply set workspace_id, the Read method will do the rest
	resource.ImportStatePassthroughID(ctx, path.Root("workspace_id"), req, resp)
}
