package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BackupAndDrApplicationResourceModel represents the resource model for a backup and DR application (policy)
type BackupAndDrApplicationResourceModel struct {
	// User-provided fields
	ID                  types.String   `tfsdk:"id"`                    // computed (policy_id)
	AccountID           types.String   `tfsdk:"account_id"`            // required
	PolicyName          types.String   `tfsdk:"policy_name"`           // required
	IntegrationID       types.String   `tfsdk:"integration_id"`        // required
	Region              types.String   `tfsdk:"region"`                // required
	ProviderType        types.String   `tfsdk:"provider_type"`         // required
	Description         types.String   `tfsdk:"description"`           // optional
	Schedule            *ScheduleModel `tfsdk:"schedule"`              // required block
	Scope               []ScopeModel   `tfsdk:"scope"`                 // optional list
	NotificationID      types.String   `tfsdk:"notification_id"`       // optional
	VCS                 *VCSModel      `tfsdk:"vcs"`                   // optional block
	RestoreInstructions types.String   `tfsdk:"restore_instructions"`  // optional
	BackupOnSave        types.Bool     `tfsdk:"backup_on_save"`        // optional, default true

	// Computed fields (never user-provided)
	Status               types.String `tfsdk:"status"`
	SnapshotsCount       types.Int64  `tfsdk:"snapshots_count"`
	LastBackupSnapshotID types.String `tfsdk:"last_backup_snapshot_id"`
	LastBackupTime       types.String `tfsdk:"last_backup_time"`
	LastBackupStatus     types.String `tfsdk:"last_backup_status"`
	NextBackupTime       types.String `tfsdk:"next_backup_time"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

// ScheduleModel represents the backup schedule configuration
type ScheduleModel struct {
	Frequency           types.String `tfsdk:"frequency"`               // required: One-time, Daily, Weekly, Monthly
	Hour                types.Int64  `tfsdk:"hour"`                    // optional: 0-23
	Minute              types.Int64  `tfsdk:"minute"`                  // optional: 0-59
	DaysOfWeek          types.List   `tfsdk:"days_of_week"`            // optional: required for Weekly
	MonthlyScheduleType types.String `tfsdk:"monthly_schedule_type"`   // optional: specific_day, specific_weekday, last_day
	DayOfMonth          types.Int64  `tfsdk:"day_of_month"`            // optional: 1-31, required for Monthly specific_day
	WeekdayOrdinal      types.String `tfsdk:"weekday_ordinal"`         // optional: First, Second, Third, Fourth, Last
	WeekdayName         types.String `tfsdk:"weekday_name"`            // optional: required for Monthly specific_weekday
	CronExpression      types.String `tfsdk:"cron_expression"`         // optional: alternative to explicit schedule
}

// ScopeModel represents a resource scope configuration
type ScopeModel struct {
	Type  types.String `tfsdk:"type"`  // required: tags, resource_group, asset_types, selected_resources
	Value types.List   `tfsdk:"value"` // required: list of strings, min 1
}

// VCSModel represents VCS integration configuration
type VCSModel struct {
	ProjectID        types.String `tfsdk:"project_id"`         // optional
	VCSIntegrationID types.String `tfsdk:"vcs_integration_id"` // optional
	RepoID           types.String `tfsdk:"repo_id"`            // optional
}

// BackupAndDrApplicationDataSourceModel represents the data source model for a single policy
// Note: Simplified schema without nested blocks (schedule, scope, vcs) to comply with data source limitations
type BackupAndDrApplicationDataSourceModel struct {
	PolicyID             types.String `tfsdk:"policy_id"`
	AccountID            types.String `tfsdk:"account_id"`
	PolicyName           types.String `tfsdk:"policy_name"`
	IntegrationID        types.String `tfsdk:"integration_id"`
	Region               types.String `tfsdk:"region"`
	ProviderType         types.String `tfsdk:"provider_type"`
	Description          types.String `tfsdk:"description"`
	ScheduleFrequency    types.String `tfsdk:"schedule_frequency"` // Simplified: only showing frequency
	NotificationID       types.String `tfsdk:"notification_id"`
	RestoreInstructions  types.String `tfsdk:"restore_instructions"`
	BackupOnSave         types.Bool   `tfsdk:"backup_on_save"`
	Status               types.String `tfsdk:"status"`
	SnapshotsCount       types.Int64  `tfsdk:"snapshots_count"`
	LastBackupSnapshotID types.String `tfsdk:"last_backup_snapshot_id"`
	LastBackupTime       types.String `tfsdk:"last_backup_time"`
	LastBackupStatus     types.String `tfsdk:"last_backup_status"`
	NextBackupTime       types.String `tfsdk:"next_backup_time"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

// BackupAndDrApplicationsDataSourceModel represents the data source model for listing policies
type BackupAndDrApplicationsDataSourceModel struct {
	ID            types.String                                 `tfsdk:"id"`
	AccountID     types.String                                 `tfsdk:"account_id"`
	Status        types.String                                 `tfsdk:"status"`
	IntegrationID types.String                                 `tfsdk:"integration_id"`
	Region        types.String                                 `tfsdk:"region"`
	ProviderType  types.String                                 `tfsdk:"provider_type"`
	Policies      []BackupAndDrApplicationDataSourceModel `tfsdk:"policies"`
}
