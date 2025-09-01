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

var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *client.Client
}

type ProjectSingleDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Path                 types.String `tfsdk:"path"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Labels               types.List   `tfsdk:"labels"`
	CronExecutionPattern types.String `tfsdk:"cron_execution_pattern"`
	ParentID             types.String `tfsdk:"parent_id"`
	AccountID            types.String `tfsdk:"account_id"`
	MembersCount         types.Int64  `tfsdk:"members_count"`
	WorkspaceCount       types.Int64  `tfsdk:"workspace_count"`
}

func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows_project"
}

func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Firefly project by ID or path",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the project",
				Optional:    true,
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "The path of the project (alternative to id)",
				Optional:    true,
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
	}
}

func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectSingleDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that either ID or path is provided, but not both
	hasID := !data.ID.IsNull() && data.ID.ValueString() != ""
	hasPath := !data.Path.IsNull() && data.Path.ValueString() != ""

	if !hasID && !hasPath {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'path' must be specified to identify the project",
		)
		return
	}

	if hasID && hasPath {
		resp.Diagnostics.AddError(
			"Conflicting Attributes",
			"Only one of 'id' or 'path' should be specified, not both",
		)
		return
	}

	var project *client.Project
	var err error

	if hasID {
		projectID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading project by ID", map[string]interface{}{"id": projectID})
		project, err = d.client.Projects.GetProject(projectID)
	} else {
		// Search for project by path (name)
		projectPath := data.Path.ValueString()
		tflog.Debug(ctx, "Reading project by path", map[string]interface{}{"path": projectPath})
		
		// Use ListProjects with search to find project by name/path
		projects, err := d.client.Projects.ListProjects(100, 0, projectPath)
		if err != nil {
			resp.Diagnostics.AddError("Error Searching Projects", fmt.Sprintf("Could not search for project path %s: %s", projectPath, err))
			return
		}

		// Find exact match by name
		var foundProject *client.Project
		for _, p := range projects.Data {
			if p.Name == projectPath {
				foundProject = &p
				break
			}
		}

		if foundProject == nil {
			resp.Diagnostics.AddError("Project Not Found", fmt.Sprintf("Could not find project with path %s", projectPath))
			return
		}

		project = foundProject
	}

	if err != nil {
		resp.Diagnostics.AddError("Error Reading Project", fmt.Sprintf("Could not read project: %s", err))
		return
	}

	// Convert labels
	var labelsList types.List
	if len(project.Labels) > 0 {
		labelValues := make([]types.String, len(project.Labels))
		for i, label := range project.Labels {
			labelValues[i] = types.StringValue(label)
		}
		labelsList = types.ListValueMust(types.StringType, labelListToValues(labelValues))
	} else {
		labelsList = types.ListValueMust(types.StringType, []attr.Value{})
	}

	data.ID = types.StringValue(project.ID)
	data.Path = types.StringValue(project.Name) // Path is the project name
	data.Name = types.StringValue(project.Name)
	data.Description = types.StringValue(project.Description)
	data.Labels = labelsList
	data.CronExecutionPattern = types.StringValue(project.CronExecutionPattern)
	data.ParentID = types.StringValue(project.ParentID)
	data.AccountID = types.StringValue(project.AccountID)
	data.MembersCount = types.Int64Value(int64(project.MembersCount))
	data.WorkspaceCount = types.Int64Value(int64(project.WorkspaceCount))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}