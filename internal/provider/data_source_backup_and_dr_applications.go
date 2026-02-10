package provider

import (
	"context"
	"fmt"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &backupAndDRApplicationsDataSource{}

func NewBackupAndDRApplicationsDataSource() datasource.DataSource {
	return &backupAndDRApplicationsDataSource{}
}

type backupAndDRApplicationsDataSource struct {
	client *client.Client
}

type BackupAndDRApplicationsDataSourceModel struct {
	ID           types.String                              `tfsdk:"id"`
	ProviderType types.String                              `tfsdk:"provider_type"`
	Status       types.String                              `tfsdk:"status"`
	Applications []BackupAndDRApplicationDataSourceModel   `tfsdk:"applications"`
}

type BackupAndDRApplicationDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	PolicyName     types.String `tfsdk:"policy_name"`
	Description    types.String `tfsdk:"description"`
	IntegrationID  types.String `tfsdk:"integration_id"`
	Region         types.String `tfsdk:"region"`
	ProviderType   types.String `tfsdk:"provider_type"`
	Status         types.String `tfsdk:"status"`
	SnapshotsCount types.Int64  `tfsdk:"snapshots_count"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func (d *backupAndDRApplicationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_and_dr_applications"
}

func (d *backupAndDRApplicationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for retrieving Firefly Backup & DR applications (backup policies)",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The data source identifier",
				Computed:            true,
			},
			"provider_type": schema.StringAttribute{
				MarkdownDescription: "Filter by cloud provider type (e.g., 'aws', 'azure', 'gcp')",
				Optional:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Filter by status ('Active' or 'Inactive')",
				Optional:            true,
			},
			"applications": schema.ListNestedAttribute{
				MarkdownDescription: "List of backup & DR applications",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the backup policy",
							Computed:            true,
						},
						"policy_name": schema.StringAttribute{
							MarkdownDescription: "The name of the backup policy",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the backup policy",
							Computed:            true,
						},
						"integration_id": schema.StringAttribute{
							MarkdownDescription: "The cloud integration ID",
							Computed:            true,
						},
						"region": schema.StringAttribute{
							MarkdownDescription: "The cloud region",
							Computed:            true,
						},
						"provider_type": schema.StringAttribute{
							MarkdownDescription: "The cloud provider type",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "The status of the backup policy",
							Computed:            true,
						},
						"snapshots_count": schema.Int64Attribute{
							MarkdownDescription: "The number of snapshots",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "The creation timestamp",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "The last update timestamp",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *backupAndDRApplicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *backupAndDRApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BackupAndDRApplicationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading backup & DR applications")

	result, err := d.client.BackupAndDR.List(1, 100)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup & DR applications",
			fmt.Sprintf("Could not list backup policies: %s", err),
		)
		return
	}

	// Apply client-side filters
	filterProvider := ""
	if !data.ProviderType.IsNull() && !data.ProviderType.IsUnknown() {
		filterProvider = data.ProviderType.ValueString()
	}
	filterStatus := ""
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		filterStatus = data.Status.ValueString()
	}

	var applications []BackupAndDRApplicationDataSourceModel
	for _, policy := range result.Policies {
		if filterProvider != "" && policy.ProviderType != filterProvider {
			continue
		}
		if filterStatus != "" && policy.Status != filterStatus {
			continue
		}

		app := BackupAndDRApplicationDataSourceModel{
			ID:             types.StringValue(policy.ID),
			PolicyName:     types.StringValue(policy.PolicyName),
			IntegrationID:  types.StringValue(policy.IntegrationID),
			Region:         types.StringValue(policy.Region),
			ProviderType:   types.StringValue(policy.ProviderType),
			Status:         types.StringValue(policy.Status),
			SnapshotsCount: types.Int64Value(int64(policy.SnapshotsCount)),
		}

		if policy.Description != "" {
			app.Description = types.StringValue(policy.Description)
		} else {
			app.Description = types.StringNull()
		}

		if policy.CreatedAt != "" {
			app.CreatedAt = types.StringValue(policy.CreatedAt)
		} else {
			app.CreatedAt = types.StringNull()
		}

		if policy.UpdatedAt != "" {
			app.UpdatedAt = types.StringValue(policy.UpdatedAt)
		} else {
			app.UpdatedAt = types.StringNull()
		}

		applications = append(applications, app)
	}

	data.Applications = applications
	data.ID = types.StringValue("backup-and-dr-applications")

	tflog.Debug(ctx, "Read backup & DR applications", map[string]interface{}{
		"count": len(data.Applications),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
