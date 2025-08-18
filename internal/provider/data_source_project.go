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
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Firefly project by ID",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the project",
				Required:    true,
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

	projectID := data.ID.ValueString()
	tflog.Debug(ctx, "Reading project", map[string]interface{}{"id": projectID})

	project, err := d.client.Projects.GetProject(projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Project", fmt.Sprintf("Could not read project ID %s: %s", projectID, err))
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