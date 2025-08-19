package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

// NewProjectResource is a helper function to simplify the provider implementation
func NewProjectResource() resource.Resource {
	return &projectResource{}
}

// projectResource is the resource implementation
type projectResource struct {
	client *client.Client
}

// ProjectResourceModel describes the resource data model
type ProjectResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Labels               types.List   `tfsdk:"labels"`
	CronExecutionPattern types.String `tfsdk:"cron_execution_pattern"`
	Variables            types.List   `tfsdk:"variables"`
	ParentID             types.String `tfsdk:"parent_id"`
	AccountID            types.String `tfsdk:"account_id"`
	MembersCount         types.Int64  `tfsdk:"members_count"`
	WorkspaceCount       types.Int64  `tfsdk:"workspace_count"`
}

// ProjectVariableModel describes a project variable
type ProjectVariableModel struct {
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	Sensitivity types.String `tfsdk:"sensitivity"`
	Destination types.String `tfsdk:"destination"`
}

// Metadata returns the resource type name
func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows_project"
}

// Schema defines the schema for the resource
func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Firefly project",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the project",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the project",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the project",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"labels": schema.ListAttribute{
				Description: "Labels to assign to the project",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"cron_execution_pattern": schema.StringAttribute{
				Description: "Cron pattern for scheduled executions",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"parent_id": schema.StringAttribute{
				Description: "ID of the parent project (for nested projects)",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"account_id": schema.StringAttribute{
				Description: "ID of the account the project belongs to",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"members_count": schema.Int64Attribute{
				Description: "Number of members assigned to the project",
				Computed:    true,
			},
			"workspace_count": schema.Int64Attribute{
				Description: "Number of workspaces in the project",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"variables": schema.ListNestedBlock{
				Description: "Variables associated with the project",
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
func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new project
func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ProjectResourceModel
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
				Sensitivity: client.VariableSensitivity(v.Sensitivity.ValueString()),
				Destination: client.VariableDestination(v.Destination.ValueString()),
			})
		}
	}

	// Get the root project ID if no parent is specified
	parentID := plan.ParentID.ValueString()
	if parentID == "" {
		rootProjectID, err := r.getRootProjectID(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Finding Root Project",
				fmt.Sprintf("Could not find root project: %s", err),
			)
			return
		}
		parentID = rootProjectID
	}

	// Create the project
	createReq := client.CreateProjectRequest{
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Labels:               labels,
		CronExecutionPattern: plan.CronExecutionPattern.ValueString(),
		Variables:            variables,
		ParentID:             parentID,
	}

	tflog.Debug(ctx, "Creating project", map[string]interface{}{
		"name": createReq.Name,
	})

	project, err := r.client.Projects.CreateProject(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Project",
			fmt.Sprintf("Could not create project: %s", err),
		)
		return
	}

	// Map response to model
	plan.ID = types.StringValue(project.ID)
	plan.AccountID = types.StringValue(project.AccountID)
	plan.MembersCount = types.Int64Value(int64(project.MembersCount))
	plan.WorkspaceCount = types.Int64Value(int64(project.WorkspaceCount))

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data
func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get project from API
	project, err := r.client.Projects.GetProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			fmt.Sprintf("Could not read project ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Map response to model
	state.Name = types.StringValue(project.Name)
	state.Description = types.StringValue(project.Description)
	state.CronExecutionPattern = types.StringValue(project.CronExecutionPattern)
	state.ParentID = types.StringValue(project.ParentID)
	state.AccountID = types.StringValue(project.AccountID)
	state.MembersCount = types.Int64Value(int64(project.MembersCount))
	state.WorkspaceCount = types.Int64Value(int64(project.WorkspaceCount))

	// Convert labels
	labelList := make([]types.String, len(project.Labels))
	for i, label := range project.Labels {
		labelList[i] = types.StringValue(label)
	}
	state.Labels = types.ListValueMust(types.StringType, labelListToValues(labelList))

	// Convert variables
	if len(project.Variables) > 0 {
		varList := make([]ProjectVariableModel, len(project.Variables))
		for i, v := range project.Variables {
			varList[i] = ProjectVariableModel{
				Key:         types.StringValue(v.Key),
				Value:       types.StringValue(v.Value),
				Sensitivity: types.StringValue(string(v.Sensitivity)),
				Destination: types.StringValue(string(v.Destination)),
			}
		}
		state.Variables = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":         types.StringType,
				"value":       types.StringType,
				"sensitivity": types.StringType,
				"destination": types.StringType,
			},
		}, projectVariablesToValues(ctx, varList))
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ProjectResourceModel
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
				Sensitivity: client.VariableSensitivity(v.Sensitivity.ValueString()),
				Destination: client.VariableDestination(v.Destination.ValueString()),
			})
		}
	}

	// Update the project
	updateReq := client.UpdateProjectRequest{
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Labels:               labels,
		CronExecutionPattern: plan.CronExecutionPattern.ValueString(),
		Variables:            variables,
	}

	tflog.Debug(ctx, "Updating project", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": updateReq.Name,
	})

	project, err := r.client.Projects.UpdateProject(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Project",
			fmt.Sprintf("Could not update project ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Map response to model
	plan.AccountID = types.StringValue(project.AccountID)
	plan.MembersCount = types.Int64Value(int64(project.MembersCount))
	plan.WorkspaceCount = types.Int64Value(int64(project.WorkspaceCount))

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete project
	tflog.Debug(ctx, "Deleting project", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := r.client.Projects.DeleteProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Project",
			fmt.Sprintf("Could not delete project ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

// ImportState imports a resource by ID
func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Save the import ID as the resource ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper functions
func labelListToValues(labels []types.String) []attr.Value {
	values := make([]attr.Value, len(labels))
	for i, label := range labels {
		values[i] = label
	}
	return values
}

func projectVariablesToValues(ctx context.Context, variables []ProjectVariableModel) []attr.Value {
	values := make([]attr.Value, len(variables))
	for i, v := range variables {
		obj, _ := types.ObjectValue(map[string]attr.Type{
			"key":         types.StringType,
			"value":       types.StringType,
			"sensitivity": types.StringType,
			"destination": types.StringType,
		}, map[string]attr.Value{
			"key":         v.Key,
			"value":       v.Value,
			"sensitivity": v.Sensitivity,
			"destination": v.Destination,
		})
		values[i] = obj
	}
	return values
}

// getRootProjectID finds the root project ID by listing projects and finding one without a parent
func (r *projectResource) getRootProjectID(ctx context.Context) (string, error) {
	// List projects to find the root project
	projects, err := r.client.Projects.ListProjects(100, 0, "")
	if err != nil {
		return "", fmt.Errorf("failed to list projects: %w", err)
	}

	// Find the project without a parent (root project)
	for _, project := range projects.Data {
		if project.ParentID == "" {
			tflog.Debug(ctx, "Found root project", map[string]interface{}{
				"root_project_id": project.ID,
				"root_project_name": project.Name,
			})
			return project.ID, nil
		}
	}

	return "", fmt.Errorf("no root project found")
}