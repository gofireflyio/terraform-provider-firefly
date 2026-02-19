package provider

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &BackupAndDrApplicationsDataSource{}

// NewBackupAndDrApplicationsDataSource creates a new backup and DR applications data source
func NewBackupAndDrApplicationsDataSource() datasource.DataSource {
	return &BackupAndDrApplicationsDataSource{}
}

// BackupAndDrApplicationsDataSource defines the data source implementation
type BackupAndDrApplicationsDataSource struct {
	client *client.Client
}

func (d *BackupAndDrApplicationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_and_dr_applications"
}

func (d *BackupAndDrApplicationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for retrieving Firefly Backup & DR applications (backup policies)",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The data source identifier",
				Computed:            true,
			},
			"account_id": schema.StringAttribute{
				MarkdownDescription: "The account ID to list policies for",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Filter by policy status (Active/Inactive)",
				Optional:            true,
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "Filter by integration ID",
				Optional:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Filter by cloud region",
				Optional:            true,
			},
			"provider_type": schema.StringAttribute{
				MarkdownDescription: "Filter by cloud provider type",
				Optional:            true,
			},
			"policies": schema.ListNestedAttribute{
				MarkdownDescription: "List of backup policies",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy_id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the backup policy",
							Computed:            true,
						},
						"account_id": schema.StringAttribute{
							MarkdownDescription: "The account ID",
							Computed:            true,
						},
						"policy_name": schema.StringAttribute{
							MarkdownDescription: "The name of the backup policy",
							Computed:            true,
						},
						"integration_id": schema.StringAttribute{
							MarkdownDescription: "The integration ID",
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
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the backup policy",
							Computed:            true,
						},
						"schedule_frequency": schema.StringAttribute{
							MarkdownDescription: "Backup schedule frequency (One-time, Daily, Weekly, Monthly)",
							Computed:            true,
						},
						"notification_id": schema.StringAttribute{
							MarkdownDescription: "Notification channel ID",
							Computed:            true,
						},
						"restore_instructions": schema.StringAttribute{
							MarkdownDescription: "Restore instructions",
							Computed:            true,
						},
						"backup_on_save": schema.BoolAttribute{
							MarkdownDescription: "Whether to backup on save",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Current status of the policy",
							Computed:            true,
						},
						"snapshots_count": schema.Int64Attribute{
							MarkdownDescription: "Number of snapshots",
							Computed:            true,
						},
						"last_backup_snapshot_id": schema.StringAttribute{
							MarkdownDescription: "ID of the last backup snapshot",
							Computed:            true,
						},
						"last_backup_time": schema.StringAttribute{
							MarkdownDescription: "Timestamp of the last backup",
							Computed:            true,
						},
						"last_backup_status": schema.StringAttribute{
							MarkdownDescription: "Status of the last backup",
							Computed:            true,
						},
						"next_backup_time": schema.StringAttribute{
							MarkdownDescription: "Timestamp of the next backup",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "Creation timestamp",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "Last update timestamp",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *BackupAndDrApplicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BackupAndDrApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BackupAndDrApplicationsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	accountID := data.AccountID.ValueString()

	// Build filters
	filters := &client.PolicyListFilters{}

	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		filters.Status = data.Status.ValueString()
	}

	if !data.IntegrationID.IsNull() && !data.IntegrationID.IsUnknown() {
		filters.IntegrationID = data.IntegrationID.ValueString()
	}

	if !data.Region.IsNull() && !data.Region.IsUnknown() {
		filters.Region = data.Region.ValueString()
	}

	if !data.ProviderType.IsNull() && !data.ProviderType.IsUnknown() {
		filters.ProviderType = data.ProviderType.ValueString()
	}

	tflog.Debug(ctx, "Reading backup policies", map[string]interface{}{
		"account_id":     accountID,
		"status":         filters.Status,
		"integration_id": filters.IntegrationID,
		"region":         filters.Region,
		"provider_type":  filters.ProviderType,
	})

	// Get policies
	policies, err := d.client.BackupAndDr.List(filters)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup policies",
			fmt.Sprintf("Could not read backup policies: %s", err),
		)
		return
	}

	// Map policies to data source model
	data.Policies = make([]BackupAndDrApplicationDataSourceModel, len(policies.Data))
	for i, policy := range policies.Data {
		policyModel := BackupAndDrApplicationDataSourceModel{
			PolicyID:          types.StringValue(policy.PolicyID),
			AccountID:         types.StringValue(policy.AccountID),
			PolicyName:        types.StringValue(policy.PolicyName),
			IntegrationID:     types.StringValue(policy.IntegrationID),
			Region:            types.StringValue(policy.Region),
			ProviderType:      types.StringValue(policy.ProviderType),
			ScheduleFrequency: types.StringValue(policy.Schedule.Frequency),
			BackupOnSave:      types.BoolValue(policy.BackupOnSave),
			Status:            types.StringValue(policy.Status),
			SnapshotsCount:    types.Int64Value(int64(policy.SnapshotsCount)),
			CreatedAt:         types.StringValue(policy.CreatedAt),
			UpdatedAt:         types.StringValue(policy.UpdatedAt),
		}

		// Optional string fields
		policyModel.Description = StringValueOrNull(policy.Description)

		policyModel.NotificationID = StringValueOrNull(policy.NotificationID)

		policyModel.RestoreInstructions = StringValueOrNull(policy.RestoreInstructions)

		policyModel.LastBackupSnapshotID = StringValueOrNull(policy.LastBackupSnapshotID)

		policyModel.LastBackupTime = StringValueOrNull(policy.LastBackupTime)

		policyModel.LastBackupStatus = StringValueOrNull(policy.LastBackupStatus)

		policyModel.NextBackupTime = StringValueOrNull(policy.NextBackupTime)

		data.Policies[i] = policyModel
	}

	// Generate unique ID based on account_id and filters
	idStr := fmt.Sprintf("backup-dr-apps-%s-%s-%s-%s-%s",
		accountID,
		filters.Status,
		filters.IntegrationID,
		filters.Region,
		filters.ProviderType,
	)
	hash := sha256.Sum256([]byte(idStr))
	data.ID = types.StringValue(fmt.Sprintf("%x", hash[:8]))

	tflog.Debug(ctx, "Read backup policies", map[string]interface{}{
		"count": len(data.Policies),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
