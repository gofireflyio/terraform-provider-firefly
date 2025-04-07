package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-firefly/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &workspaceRunsDataSource{}
	_ datasource.DataSourceWithConfigure = &workspaceRunsDataSource{}
)

// NewWorkspaceRunsDataSource is a helper function to simplify the provider implementation
func NewWorkspaceRunsDataSource() datasource.DataSource {
	return &workspaceRunsDataSource{}
}

// workspaceRunsDataSource is the data source implementation
type workspaceRunsDataSource struct {
	client *client.Client
}

// WorkspaceRunFiltersModel describes the workspace run filters
type WorkspaceRunFiltersModel struct {
	RunID       types.List   `tfsdk:"run_id"`
	RunName     types.List   `tfsdk:"run_name"`
	Status      types.List   `tfsdk:"status"`
	Branch      types.List   `tfsdk:"branch"`
	CommitID    types.List   `tfsdk:"commit_id"`
	CITool      types.List   `tfsdk:"ci_tool"`
	VCSType     types.List   `tfsdk:"vcs_type"`
	Repository  types.List   `tfsdk:"repository"`
}

// WorkspaceRunDataModel describes a single workspace run
type WorkspaceRunDataModel struct {
	ID            types.String `tfsdk:"id"`
	WorkspaceID   types.String `tfsdk:"workspace_id"`
	WorkspaceName types.String `tfsdk:"workspace_name"`
	RunID         types.String `tfsdk:"run_id"`
	RunName       types.String `tfsdk:"run_name"`
	Status        types.String `tfsdk:"status"`
	Branch        types.String `tfsdk:"branch"`
	CommitID      types.String `tfsdk:"commit_id"`
	CommitURL     types.String `tfsdk:"commit_url"`
	RunnerType    types.String `tfsdk:"runner_type"`
	BuildID       types.String `tfsdk:"build_id"`
	BuildURL      types.String `tfsdk:"build_url"`
	BuildName     types.String `tfsdk:"build_name"`
	VCSType       types.String `tfsdk:"vcs_type"`
	Repo          types.String `tfsdk:"repo"`
	RepoURL       types.String `tfsdk:"repo_url"`
	Title         types.String `tfsdk:"title"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

// WorkspaceRunsDataSourceModel describes the data source data model
type WorkspaceRunsDataSourceModel struct {
	WorkspaceID  types.String            `tfsdk:"workspace_id"`
	Runs         types.List              `tfsdk:"runs"`
	Filters      *WorkspaceRunFiltersModel `tfsdk:"filters"`
	SearchValue  types.String            `tfsdk:"search_value"`
}

// Metadata returns the data source type name
func (d *workspaceRunsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_runs"
}

// Schema defines the schema for the data source
func (d *workspaceRunsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of runs for a Firefly workspace",
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				Description: "ID of the workspace to fetch runs for",
				Required:    true,
			},
			"runs": schema.ListNestedAttribute{
				Description: "List of workspace runs",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the run",
							Computed:    true,
						},
						"workspace_id": schema.StringAttribute{
							Description: "ID of the workspace",
							Computed:    true,
						},
						"workspace_name": schema.StringAttribute{
							Description: "Name of the workspace",
							Computed:    true,
						},
						"run_id": schema.StringAttribute{
							Description: "Run ID",
							Computed:    true,
						},
						"run_name": schema.StringAttribute{
							Description: "Name of the run",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of the run",
							Computed:    true,
						},
						"branch": schema.StringAttribute{
							Description: "Branch associated with the run",
							Computed:    true,
						},
						"commit_id": schema.StringAttribute{
							Description: "Commit ID",
							Computed:    true,
						},
						"commit_url": schema.StringAttribute{
							Description: "URL of the commit",
							Computed:    true,
						},
						"runner_type": schema.StringAttribute{
							Description: "Type of CI/CD runner",
							Computed:    true,
						},
						"build_id": schema.StringAttribute{
							Description: "Build ID",
							Computed:    true,
						},
						"build_url": schema.StringAttribute{
							Description: "URL of the build",
							Computed:    true,
						},
						"build_name": schema.StringAttribute{
							Description: "Name of the build",
							Computed:    true,
						},
						"vcs_type": schema.StringAttribute{
							Description: "Type of version control system",
							Computed:    true,
						},
						"repo": schema.StringAttribute{
							Description: "Repository associated with the run",
							Computed:    true,
						},
						"repo_url": schema.StringAttribute{
							Description: "URL of the repository",
							Computed:    true,
						},
						"title": schema.StringAttribute{
							Description: "Title of the run",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Timestamp when the run was created",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Timestamp when the run was last updated",
							Computed:    true,
						},
					},
				},
			},
			"search_value": schema.StringAttribute{
				Description: "Search value to filter runs",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"filters": schema.SingleNestedBlock{
				Description: "Filters for workspace runs",
				Attributes: map[string]schema.Attribute{
					"run_id": schema.ListAttribute{
						Description: "Filter by run ID",
						Optional:    true,
						ElementType: types.StringType,
					},
					"run_name": schema.ListAttribute{
						Description: "Filter by run name",
						Optional:    true,
						ElementType: types.StringType,
					},
					"status": schema.ListAttribute{
						Description: "Filter by status",
						Optional:    true,
						ElementType: types.StringType,
					},
					"branch": schema.ListAttribute{
						Description: "Filter by branch",
						Optional:    true,
						ElementType: types.StringType,
					},
					"commit_id": schema.ListAttribute{
						Description: "Filter by commit ID",
						Optional:    true,
						ElementType: types.StringType,
					},
					"ci_tool": schema.ListAttribute{
						Description: "Filter by CI tool",
						Optional:    true,
						ElementType: types.StringType,
					},
					"vcs_type": schema.ListAttribute{
						Description: "Filter by VCS type",
						Optional:    true,
						ElementType: types.StringType,
					},
					"repository": schema.ListAttribute{
						Description: "Filter by repository",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *workspaceRunsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *workspaceRunsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkspaceRunsDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the workspace ID
	workspaceID := data.WorkspaceID.ValueString()
	if workspaceID == "" {
		resp.Diagnostics.AddError(
			"Missing Workspace ID",
			"Workspace ID is required",
		)
		return
	}

	// Prepare the request with filters
	request := &client.ListWorkspaceRunsRequest{}
	
	// Add search value if provided
	if !data.SearchValue.IsNull() {
		request.SearchValue = data.SearchValue.ValueString()
	}
	
	// Add filters if provided
	if data.Filters != nil {
		filters := &client.WorkspaceRunFilters{}
		
		// Add run ID filter
		if !data.Filters.RunID.IsNull() {
			var runIDs []string
			diags = data.Filters.RunID.ElementsAs(ctx, &runIDs, false)
			resp.Diagnostics.Append(diags...)
			if runIDs != nil {
				filters.RunID = runIDs
			}
		}
		
		// Add run name filter
		if !data.Filters.RunName.IsNull() {
			var runNames []string
			diags = data.Filters.RunName.ElementsAs(ctx, &runNames, false)
			resp.Diagnostics.Append(diags...)
			if runNames != nil {
				filters.RunName = runNames
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
		
		// Add branch filter
		if !data.Filters.Branch.IsNull() {
			var branches []string
			diags = data.Filters.Branch.ElementsAs(ctx, &branches, false)
			resp.Diagnostics.Append(diags...)
			if branches != nil {
				filters.Branch = branches
			}
		}
		
		// Add commit ID filter
		if !data.Filters.CommitID.IsNull() {
			var commitIDs []string
			diags = data.Filters.CommitID.ElementsAs(ctx, &commitIDs, false)
			resp.Diagnostics.Append(diags...)
			if commitIDs != nil {
				filters.CommitID = commitIDs
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
		
		// Add VCS type filter
		if !data.Filters.VCSType.IsNull() {
			var vcsTypes []string
			diags = data.Filters.VCSType.ElementsAs(ctx, &vcsTypes, false)
			resp.Diagnostics.Append(diags...)
			if vcsTypes != nil {
				filters.VCSType = vcsTypes
			}
		}
		
		// Add repository filter
		if !data.Filters.Repository.IsNull() {
			var repositories []string
			diags = data.Filters.Repository.ElementsAs(ctx, &repositories, false)
			resp.Diagnostics.Append(diags...)
			if repositories != nil {
				filters.Repository = repositories
			}
		}
		
		request.Filters = filters
	}
	
	// Get runs from API
	runs, err := d.client.Workspaces.ListWorkspaceRuns(workspaceID, request, 0, 100)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workspace Runs",
			fmt.Sprintf("Could not read runs for workspace ID %s: %s", workspaceID, err),
		)
		return
	}
	
	// Map response to model
	var runModels []WorkspaceRunDataModel
	for _, run := range runs {
		runModel := WorkspaceRunDataModel{
			ID:            types.StringValue(run.ID),
			WorkspaceID:   types.StringValue(run.WorkspaceID),
			WorkspaceName: types.StringValue(run.WorkspaceName),
			RunID:         types.StringValue(run.RunID),
			RunName:       types.StringValue(run.RunName),
			Status:        types.StringValue(run.Status),
			Branch:        types.StringValue(run.Branch),
			CommitID:      types.StringValue(run.CommitID),
			CommitURL:     types.StringValue(run.CommitURL),
			RunnerType:    types.StringValue(run.RunnerType),
			BuildID:       types.StringValue(run.BuildID),
			BuildURL:      types.StringValue(run.BuildURL),
			BuildName:     types.StringValue(run.BuildName),
			VCSType:       types.StringValue(run.VCSType),
			Repo:          types.StringValue(run.Repo),
			RepoURL:       types.StringValue(run.RepoURL),
			Title:         types.StringValue(run.Title),
			CreatedAt:     types.StringValue(run.CreatedAt),
			UpdatedAt:     types.StringValue(run.UpdatedAt),
		}
		
		runModels = append(runModels, runModel)
	}
	
	// Set runs in the data model
	runsList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]types.Type{
			"id":             types.StringType,
			"workspace_id":   types.StringType,
			"workspace_name": types.StringType,
			"run_id":         types.StringType,
			"run_name":       types.StringType,
			"status":         types.StringType,
			"branch":         types.StringType,
			"commit_id":      types.StringType,
			"commit_url":     types.StringType,
			"runner_type":    types.StringType,
			"build_id":       types.StringType,
			"build_url":      types.StringType,
			"build_name":     types.StringType,
			"vcs_type":       types.StringType,
			"repo":           types.StringType,
			"repo_url":       types.StringType,
			"title":          types.StringType,
			"created_at":     types.StringType,
			"updated_at":     types.StringType,
		},
	}, runModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	data.Runs = runsList
	
	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
