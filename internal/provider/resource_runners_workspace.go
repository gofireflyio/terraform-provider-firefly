package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-firefly/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &runnersWorkspaceResource{}
	_ resource.ResourceWithConfigure   = &runnersWorkspaceResource{}
	_ resource.ResourceWithImportState = &runnersWorkspaceResource{}
)

// NewRunnersWorkspaceResource is a helper function to simplify the provider implementation
func NewRunnersWorkspaceResource() resource.Resource {
	return &runnersWorkspaceResource{}
}

// runnersWorkspaceResource is the resource implementation
type runnersWorkspaceResource struct {
	client *client.Client
}

// RunnersWorkspaceResourceModel describes the resource data model
type RunnersWorkspaceResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Repository           types.String `tfsdk:"repository"`
	WorkingDirectory     types.String `tfsdk:"working_directory"`
	VcsIntegrationID     types.String `tfsdk:"vcs_integration_id"`
	VcsType              types.String `tfsdk:"vcs_type"`
	DefaultBranch        types.String `tfsdk:"default_branch"`
	CronExecutionPattern types.String `tfsdk:"cron_execution_pattern"`
	IacType              types.String `tfsdk:"iac_type"`
	TerraformVersion     types.String `tfsdk:"terraform_version"`
	ApplyRule            types.String `tfsdk:"apply_rule"`
	Triggers             types.List   `tfsdk:"triggers"`
	Labels               types.List   `tfsdk:"labels"`
	Variables            types.List   `tfsdk:"variables"`
	ConsumedVariableSets types.List   `tfsdk:"consumed_variable_sets"`
	ProjectID            types.String `tfsdk:"project_id"`
	AccountID            types.String `tfsdk:"account_id"`
}

// IacProvisionerModel describes the IaC provisioner
type IacProvisionerModel struct {
	Type    types.String `tfsdk:"type"`
	Version types.String `tfsdk:"version"`
}

// Metadata returns the resource type name
func (r *runnersWorkspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runners_workspace"
}

// Schema defines the schema for the resource
func (r *runnersWorkspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Firefly runners workspace",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the workspace",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the workspace",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the workspace",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"repository": schema.StringAttribute{
				Description: "Repository URL or name",
				Required:    true,
			},
			"working_directory": schema.StringAttribute{
				Description: "Working directory within the repository",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"vcs_integration_id": schema.StringAttribute{
				Description: "VCS integration ID",
				Required:    true,
			},
			"vcs_type": schema.StringAttribute{
				Description: "VCS type (e.g., github, gitlab)",
				Required:    true,
			},
			"default_branch": schema.StringAttribute{
				Description: "Default branch for the workspace",
				Required:    true,
			},
			"cron_execution_pattern": schema.StringAttribute{
				Description: "Cron pattern for scheduled executions",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"iac_type": schema.StringAttribute{
				Description: "Infrastructure as Code type (terraform, opentofu)",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("terraform"),
			},
			"terraform_version": schema.StringAttribute{
				Description: "Terraform version to use",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("1.5.7"),
			},
			"apply_rule": schema.StringAttribute{
				Description: "Apply rule (manual or auto)",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("manual"),
			},
			"triggers": schema.ListAttribute{
				Description: "List of triggers for the workspace",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"labels": schema.ListAttribute{
				Description: "Labels to assign to the workspace",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"consumed_variable_sets": schema.ListAttribute{
				Description: "List of variable set IDs that this workspace consumes",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "Project ID for workspace assignment",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"account_id": schema.StringAttribute{
				Description: "Account ID that the workspace belongs to",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"variables": schema.ListNestedBlock{
				Description: "Variables associated with the workspace",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "The variable key",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The variable value",
							Required:    true,
							Sensitive:   true,
						},
						"sensitivity": schema.StringAttribute{
							Description: "The sensitivity of the variable (string or secret)",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("string"),
						},
						"destination": schema.StringAttribute{
							Description: "The destination of the variable (env or iac)",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("env"),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *runnersWorkspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new runners workspace
func (r *runnersWorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan RunnersWorkspaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert triggers to string slice
	var triggers []string
	if !plan.Triggers.IsNull() && !plan.Triggers.IsUnknown() {
		diags = plan.Triggers.ElementsAs(ctx, &triggers, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert labels to string slice
	var labels []string
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		diags = plan.Labels.ElementsAs(ctx, &labels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert consumed variable sets to string slice
	var consumedVariableSets []string
	if !plan.ConsumedVariableSets.IsNull() && !plan.ConsumedVariableSets.IsUnknown() {
		diags = plan.ConsumedVariableSets.ElementsAs(ctx, &consumedVariableSets, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert variables to API format
	var variables []client.Variable
	if !plan.Variables.IsNull() && !plan.Variables.IsUnknown() {
		var varModels []ProjectVariableModel
		diags = plan.Variables.ElementsAs(ctx, &varModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, v := range varModels {
			variables = append(variables, client.Variable{
				Key:         v.Key.ValueString(),
				Value:       v.Value.ValueString(),
				Sensitivity: v.Sensitivity.ValueString(),
				Destination: v.Destination.ValueString(),
			})
		}
	}

	// Handle project ID (can be nil)
	var projectID *string
	if !plan.ProjectID.IsNull() && plan.ProjectID.ValueString() != "" {
		pid := plan.ProjectID.ValueString()
		projectID = &pid
	}

	// Create the workspace
	createReq := client.CreateRunnersWorkspaceRequest{
		RunnerType:           "firefly_runners",
		IacType:              plan.IacType.ValueString(),
		WorkspaceName:        plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Labels:               labels,
		VcsID:                plan.VcsIntegrationID.ValueString(),
		Repo:                 plan.Repository.ValueString(),
		DefaultBranch:        plan.DefaultBranch.ValueString(),
		VcsType:              plan.VcsType.ValueString(),
		WorkDir:              plan.WorkingDirectory.ValueString(),
		Variables:            variables,
		ConsumedVariableSets: consumedVariableSets,
		Execution: client.ExecutionConfig{
			Triggers:         triggers,
			ApplyRule:        plan.ApplyRule.ValueString(),
			TerraformVersion: plan.TerraformVersion.ValueString(),
		},
		Project: projectID,
	}

	tflog.Debug(ctx, "Creating runners workspace", map[string]interface{}{
		"name": createReq.WorkspaceName,
	})

	workspace, err := r.client.RunnersWorkspaces.CreateRunnersWorkspace(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Runners Workspace",
			fmt.Sprintf("Could not create runners workspace: %s", err),
		)
		return
	}

	// Map response to model
	plan.ID = types.StringValue(workspace.ID)
	plan.AccountID = types.StringValue(workspace.AccountID)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data
func (r *runnersWorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state RunnersWorkspaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get workspace from API
	workspace, err := r.client.RunnersWorkspaces.GetRunnersWorkspace(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Runners Workspace",
			fmt.Sprintf("Could not read runners workspace ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Map response to model
	state.Name = types.StringValue(workspace.Name)
	state.Description = types.StringValue(workspace.Description)
	state.Repository = types.StringValue(workspace.Repository)
	state.WorkingDirectory = types.StringValue(workspace.WorkingDirectory)
	state.VcsIntegrationID = types.StringValue(workspace.VcsIntegrationID)
	state.VcsType = types.StringValue(workspace.Vcs)
	state.DefaultBranch = types.StringValue(workspace.DefaultBranch)
	state.CronExecutionPattern = types.StringValue(workspace.CronExecutionPattern)
	state.AccountID = types.StringValue(workspace.AccountID)

	// Handle IaC provisioner
	if workspace.IacProvisioner != nil {
		state.IacType = types.StringValue(workspace.IacProvisioner.Type)
		state.TerraformVersion = types.StringValue(workspace.IacProvisioner.Version)
	}

	// Convert labels
	if len(workspace.Labels) > 0 {
		labelList := make([]types.String, len(workspace.Labels))
		for i, label := range workspace.Labels {
			labelList[i] = types.StringValue(label)
		}
		state.Labels = types.ListValueMust(types.StringType, labelListToValues(labelList))
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *runnersWorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan RunnersWorkspaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert labels to string slice
	var labels []string
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		diags = plan.Labels.ElementsAs(ctx, &labels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert consumed variable sets to string slice
	var consumedVariableSets []string
	if !plan.ConsumedVariableSets.IsNull() && !plan.ConsumedVariableSets.IsUnknown() {
		diags = plan.ConsumedVariableSets.ElementsAs(ctx, &consumedVariableSets, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert variables to API format
	var variables []client.Variable
	if !plan.Variables.IsNull() && !plan.Variables.IsUnknown() {
		var varModels []ProjectVariableModel
		diags = plan.Variables.ElementsAs(ctx, &varModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, v := range varModels {
			variables = append(variables, client.Variable{
				Key:         v.Key.ValueString(),
				Value:       v.Value.ValueString(),
				Sensitivity: v.Sensitivity.ValueString(),
				Destination: v.Destination.ValueString(),
			})
		}
	}

	// Update the workspace
	updateReq := client.UpdateRunnersWorkspaceRequest{
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Labels:               labels,
		VcsIntegrationID:     plan.VcsIntegrationID.ValueString(),
		Repository:           plan.Repository.ValueString(),
		DefaultBranch:        plan.DefaultBranch.ValueString(),
		WorkingDirectory:     plan.WorkingDirectory.ValueString(),
		CronExecutionPattern: plan.CronExecutionPattern.ValueString(),
		IacProvisioner: &client.IacProvisioner{
			Type:    plan.IacType.ValueString(),
			Version: plan.TerraformVersion.ValueString(),
		},
		Variables:            variables,
		ConsumedVariableSets: consumedVariableSets,
	}

	tflog.Debug(ctx, "Updating runners workspace", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": updateReq.Name,
	})

	workspace, err := r.client.RunnersWorkspaces.UpdateRunnersWorkspace(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Runners Workspace",
			fmt.Sprintf("Could not update runners workspace ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Map response to model
	plan.AccountID = types.StringValue(workspace.AccountID)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *runnersWorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state RunnersWorkspaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete workspace
	tflog.Debug(ctx, "Deleting runners workspace", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := r.client.RunnersWorkspaces.DeleteRunnersWorkspace(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Runners Workspace",
			fmt.Sprintf("Could not delete runners workspace ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

// ImportState imports a resource by ID
func (r *runnersWorkspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Save the import ID as the resource ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}