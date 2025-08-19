package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &projectsDataSource{}
	_ datasource.DataSourceWithConfigure = &projectsDataSource{}
)

// NewProjectsDataSource is a helper function to simplify the provider implementation
func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

// projectsDataSource is the data source implementation
type projectsDataSource struct {
	client *client.Client
}

// ProjectsDataSourceModel describes the data source data model
type ProjectsDataSourceModel struct {
	SearchQuery types.String         `tfsdk:"search_query"`
	Projects    []ProjectDataModel   `tfsdk:"projects"`
}

// ProjectDataModel describes a single project in the data source
type ProjectDataModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Labels               types.List   `tfsdk:"labels"`
	CronExecutionPattern types.String `tfsdk:"cron_execution_pattern"`
	ParentID             types.String `tfsdk:"parent_id"`
	AccountID            types.String `tfsdk:"account_id"`
	MembersCount         types.Int64  `tfsdk:"members_count"`
	WorkspaceCount       types.Int64  `tfsdk:"workspace_count"`
}

// Metadata returns the data source type name
func (d *projectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows_projects"
}

// Schema defines the schema for the data source
func (d *projectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of Firefly projects",
		Attributes: map[string]schema.Attribute{
			"search_query": schema.StringAttribute{
				Description: "Optional search query to filter projects",
				Optional:    true,
			},
			"projects": schema.ListNestedAttribute{
				Description: "List of projects",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the project",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the project",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the project",
							Computed:    true,
						},
						"labels": schema.ListAttribute{
							Description: "Labels assigned to the project",
							Computed:    true,
							ElementType: types.StringType,
						},
						"cron_execution_pattern": schema.StringAttribute{
							Description: "Cron pattern for scheduled executions",
							Computed:    true,
						},
						"parent_id": schema.StringAttribute{
							Description: "ID of the parent project",
							Computed:    true,
						},
						"account_id": schema.StringAttribute{
							Description: "ID of the account the project belongs to",
							Computed:    true,
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
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *projectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsDataSourceModel

	// Read configuration
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get search query
	searchQuery := ""
	if !data.SearchQuery.IsNull() {
		searchQuery = data.SearchQuery.ValueString()
	}

	tflog.Debug(ctx, "Reading projects", map[string]interface{}{
		"search_query": searchQuery,
	})

	// Get projects from API
	projectsResp, err := d.client.Projects.ListProjects(100, 0, searchQuery) // Default to first 100
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Projects",
			fmt.Sprintf("Could not read projects: %s", err),
		)
		return
	}

	// Map response to model
	projects := make([]ProjectDataModel, len(projectsResp.Data))
	for i, project := range projectsResp.Data {
		// Convert labels
		var labelsList types.List
		if len(project.Labels) > 0 {
			labelValues := make([]types.String, len(project.Labels))
			for j, label := range project.Labels {
				labelValues[j] = types.StringValue(label)
			}
			labelsList = types.ListValueMust(types.StringType, labelListToValues(labelValues))
		} else {
			labelsList = types.ListValueMust(types.StringType, []attr.Value{})
		}

		projects[i] = ProjectDataModel{
			ID:                   types.StringValue(project.ID),
			Name:                 types.StringValue(project.Name),
			Description:          types.StringValue(project.Description),
			Labels:               labelsList,
			CronExecutionPattern: types.StringValue(project.CronExecutionPattern),
			ParentID:             types.StringValue(project.ParentID),
			AccountID:            types.StringValue(project.AccountID),
			MembersCount:         types.Int64Value(int64(project.MembersCount)),
			WorkspaceCount:       types.Int64Value(int64(project.WorkspaceCount)),
		}
	}

	data.Projects = projects

	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}