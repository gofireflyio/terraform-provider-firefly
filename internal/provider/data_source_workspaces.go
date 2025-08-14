package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-firefly/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &workspacesDataSource{}
	_ datasource.DataSourceWithConfigure = &workspacesDataSource{}
)

// NewWorkspacesDataSource is a helper function to simplify the provider implementation
func NewWorkspacesDataSource() datasource.DataSource {
	return &workspacesDataSource{}
}

// workspacesDataSource is the data source implementation
type workspacesDataSource struct {
	client *client.Client
}

// WorkspaceFiltersModel describes the workspace filters
type WorkspaceFiltersModel struct {
	WorkspaceName    types.List   `tfsdk:"workspace_name"`
	Repositories     types.List   `tfsdk:"repositories"`
	CITool           types.List   `tfsdk:"ci_tool"`
	Labels           types.List   `tfsdk:"labels"`
	Status           types.List   `tfsdk:"status"`
	IsManagedWorkflow types.Bool  `tfsdk:"is_managed_workflow"`
	VCSType          types.List   `tfsdk:"vcs_type"`
}

// WorkspaceDataModel describes a single workspace
type WorkspaceDataModel struct {
	ID                types.String `tfsdk:"id"`
	AccountID         types.String `tfsdk:"account_id"`
	WorkspaceID       types.String `tfsdk:"workspace_id"`
	WorkspaceName     types.String `tfsdk:"workspace_name"`
	Repo              types.String `tfsdk:"repo"`
	RepoURL           types.String `tfsdk:"repo_url"`
	VCSType           types.String `tfsdk:"vcs_type"`
	RunnerType        types.String `tfsdk:"runner_type"`
	LastRunStatus     types.String `tfsdk:"last_run_status"`
	LastApplyTime     types.String `tfsdk:"last_apply_time"`
	LastPlanTime      types.String `tfsdk:"last_plan_time"`
	LastRunTime       types.String `tfsdk:"last_run_time"`
	IACType           types.String `tfsdk:"iac_type"`
	IACTypeVersion    types.String `tfsdk:"iac_type_version"`
	Labels            types.List   `tfsdk:"labels"`
	RunsCount         types.Int64  `tfsdk:"runs_count"`
	IsWorkflowManaged types.Bool   `tfsdk:"is_workflow_managed"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

// WorkspacesDataSourceModel describes the data source data model
type WorkspacesDataSourceModel struct {
	Workspaces types.List          `tfsdk:"workspaces"`
	Filters    *WorkspaceFiltersModel `tfsdk:"filters"`
	SearchValue types.String       `tfsdk:"search_value"`
}

// Metadata returns the data source type name
func (d *workspacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspaces"
}

// Schema defines the schema for the data source
func (d *workspacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Firefly workspaces",
		Attributes: map[string]schema.Attribute{
			"workspaces": schema.ListNestedAttribute{
				Description: "List of workspaces",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the workspace",
							Computed:    true,
						},
						"account_id": schema.StringAttribute{
							Description: "Account ID associated with the workspace",
							Computed:    true,
						},
						"workspace_id": schema.StringAttribute{
							Description: "Workspace ID",
							Computed:    true,
						},
						"workspace_name": schema.StringAttribute{
							Description: "Name of the workspace",
							Computed:    true,
						},
						"repo": schema.StringAttribute{
							Description: "Repository associated with the workspace",
							Computed:    true,
						},
						"repo_url": schema.StringAttribute{
							Description: "Repository URL",
							Computed:    true,
						},
						"vcs_type": schema.StringAttribute{
							Description: "Version control system type",
							Computed:    true,
						},
						"runner_type": schema.StringAttribute{
							Description: "CI/CD runner type",
							Computed:    true,
						},
						"last_run_status": schema.StringAttribute{
							Description: "Status of the last run",
							Computed:    true,
						},
						"last_apply_time": schema.StringAttribute{
							Description: "Timestamp of the last apply",
							Computed:    true,
						},
						"last_plan_time": schema.StringAttribute{
							Description: "Timestamp of the last plan",
							Computed:    true,
						},
						"last_run_time": schema.StringAttribute{
							Description: "Timestamp of the last run",
							Computed:    true,
						},
						"iac_type": schema.StringAttribute{
							Description: "Type of Infrastructure as Code",
							Computed:    true,
						},
						"iac_type_version": schema.StringAttribute{
							Description: "Version of the IaC tool",
							Computed:    true,
						},
						"labels": schema.ListAttribute{
							Description: "List of labels associated with the workspace",
							Computed:    true,
							ElementType: types.StringType,
						},
						"runs_count": schema.Int64Attribute{
							Description: "Number of runs",
							Computed:    true,
						},
						"is_workflow_managed": schema.BoolAttribute{
							Description: "Whether the workflow is managed",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Timestamp when the workspace was created",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Timestamp when the workspace was last updated",
							Computed:    true,
						},
					},
				},
			},
			"search_value": schema.StringAttribute{
				Description: "Search value to filter workspaces",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"filters": schema.SingleNestedBlock{
				Description: "Filters for workspaces",
				Attributes: map[string]schema.Attribute{
					"workspace_name": schema.ListAttribute{
						Description: "Filter by workspace name",
						Optional:    true,
						ElementType: types.StringType,
					},
					"repositories": schema.ListAttribute{
						Description: "Filter by repositories",
						Optional:    true,
						ElementType: types.StringType,
					},
					"ci_tool": schema.ListAttribute{
						Description: "Filter by CI tool",
						Optional:    true,
						ElementType: types.StringType,
					},
					"labels": schema.ListAttribute{
						Description: "Filter by labels",
						Optional:    true,
						ElementType: types.StringType,
					},
					"status": schema.ListAttribute{
						Description: "Filter by status",
						Optional:    true,
						ElementType: types.StringType,
					},
					"is_managed_workflow": schema.BoolAttribute{
						Description: "Filter by whether the workflow is managed",
						Optional:    true,
					},
					"vcs_type": schema.ListAttribute{
						Description: "Filter by VCS type",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *workspacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

// Read refreshes the Terraform state with the latest data
func (d *workspacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkspacesDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request with filters
	request := &client.ListWorkspacesRequest{}
	
	// Add search value if provided
	if !data.SearchValue.IsNull() {
		request.SearchValue = data.SearchValue.ValueString()
	}
	
	// Add filters if provided
	if data.Filters != nil {
		filters := &client.WorkspaceFilters{}
		
		// Add workspace name filter
		if !data.Filters.WorkspaceName.IsNull() {
			var workspaceNames []string
			diags = data.Filters.WorkspaceName.ElementsAs(ctx, &workspaceNames, false)
			resp.Diagnostics.Append(diags...)
			if workspaceNames != nil {
				filters.WorkspaceName = workspaceNames
			}
		}
		
		// Add repositories filter
		if !data.Filters.Repositories.IsNull() {
			var repositories []string
			diags = data.Filters.Repositories.ElementsAs(ctx, &repositories, false)
			resp.Diagnostics.Append(diags...)
			if repositories != nil {
				filters.Repositories = repositories
			}
		}
		
		// Add CI tool filter
		if !data.Filters.CITool.IsNull() {
			var ciTools []string
			diags = data.Filters.CITool.ElementsAs(ctx, &ciTools, false)
			resp.Diagnostics.Append(diags...)
			if ciTools != nil {
				filters.CITool = ciTools
			}
		}
		
		// Add labels filter
		if !data.Filters.Labels.IsNull() {
			var labels []string
			diags = data.Filters.Labels.ElementsAs(ctx, &labels, false)
			resp.Diagnostics.Append(diags...)
			if labels != nil {
				filters.Labels = labels
			}
		}
		
		// Add status filter
		if !data.Filters.Status.IsNull() {
			var statuses []string
			diags = data.Filters.Status.ElementsAs(ctx, &statuses, false)
			resp.Diagnostics.Append(diags...)
			if statuses != nil {
				filters.Status = statuses
			}
		}
		
		// Add is managed workflow filter
		if !data.Filters.IsManagedWorkflow.IsNull() {
			isManagedWorkflow := data.Filters.IsManagedWorkflow.ValueBool()
			filters.IsManagedWorkflow = &isManagedWorkflow
		}
		
		// Add VCS type filter
		if !data.Filters.VCSType.IsNull() {
			var vcsTypes []string
			diags = data.Filters.VCSType.ElementsAs(ctx, &vcsTypes, false)
			resp.Diagnostics.Append(diags...)
			if vcsTypes != nil {
				filters.VCSType = vcsTypes
			}
		}
		
		request.Filters = filters
	}
	
	// Get workspaces from API
	workspaces, err := d.client.Workspaces.ListWorkspaces(request, 0, 100)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workspaces",
			fmt.Sprintf("Could not read workspaces: %s", err),
		)
		return
	}
	
	// Map response to model
	var workspaceModels []WorkspaceDataModel
	for _, workspace := range workspaces {
		workspaceModel := WorkspaceDataModel{
			ID:               types.StringValue(workspace.ID),
			AccountID:        types.StringValue(workspace.AccountID),
			WorkspaceID:      types.StringValue(workspace.WorkspaceID),
			WorkspaceName:    types.StringValue(workspace.WorkspaceName),
			Repo:             types.StringValue(workspace.Repo),
			RepoURL:          types.StringValue(workspace.RepoURL),
			VCSType:          types.StringValue(workspace.VCSType),
			RunnerType:       types.StringValue(workspace.RunnerType),
			LastRunStatus:    types.StringValue(workspace.LastRunStatus),
			LastApplyTime:    types.StringValue(workspace.LastApplyTime),
			LastPlanTime:     types.StringValue(workspace.LastPlanTime),
			LastRunTime:      types.StringValue(workspace.LastRunTime),
			IACType:          types.StringValue(workspace.IACType),
			IACTypeVersion:   types.StringValue(workspace.IACTypeVersion),
			RunsCount:        types.Int64Value(int64(workspace.RunsCount)),
			IsWorkflowManaged: types.BoolValue(workspace.IsWorkflowManaged),
			CreatedAt:        types.StringValue(workspace.CreatedAt),
			UpdatedAt:        types.StringValue(workspace.UpdatedAt),
		}
		
		// Set labels
		workspaceModel.Labels = types.ListValueMust(types.StringType, listToValues(workspace.Labels))
		
		workspaceModels = append(workspaceModels, workspaceModel)
	}
	
	// Set workspaces in the data model
	workspacesList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                 types.StringType,
			"account_id":         types.StringType,
			"workspace_id":       types.StringType,
			"workspace_name":     types.StringType,
			"repo":               types.StringType,
			"repo_url":           types.StringType,
			"vcs_type":           types.StringType,
			"runner_type":        types.StringType,
			"last_run_status":    types.StringType,
			"last_apply_time":    types.StringType,
			"last_plan_time":     types.StringType,
			"last_run_time":      types.StringType,
			"iac_type":           types.StringType,
			"iac_type_version":   types.StringType,
			"labels":             types.ListType{ElemType: types.StringType},
			"runs_count":         types.Int64Type,
			"is_workflow_managed": types.BoolType,
			"created_at":         types.StringType,
			"updated_at":         types.StringType,
		},
	}, workspaceModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	data.Workspaces = workspacesList
	
	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
