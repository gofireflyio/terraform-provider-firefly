package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BackupAndDrApplicationResourceModel represents the resource model for a backup and DR application
type BackupAndDrApplicationResourceModel struct {
	// User-provided fields
	ID                  types.String `tfsdk:"id"`                   // computed (application_id)
	AccountID           types.String `tfsdk:"account_id"`           // required
	ApplicationName     types.String `tfsdk:"application_name"`     // required
	IntegrationID       types.String `tfsdk:"integration_id"`       // required
	Region              types.String `tfsdk:"region"`               // required
	ProviderType        types.String `tfsdk:"provider_type"`        // required
	Description         types.String `tfsdk:"description"`          // optional
	Frequency           types.Int64  `tfsdk:"frequency"`            // optional: 4, 8, 16, or 24 hours
	Scope               []ScopeModel `tfsdk:"scope"`                // optional list
	NotificationID      types.String `tfsdk:"notification_id"`      // optional
	VCS                 *VCSModel    `tfsdk:"vcs"`                  // optional block
	RestoreInstructions types.String `tfsdk:"restore_instructions"` // optional
	BackupOnSave        types.Bool   `tfsdk:"backup_on_save"`       // optional, default true
	TargetAccount       types.String `tfsdk:"target_account"`       // optional
	TargetRegion        types.String `tfsdk:"target_region"`        // optional
	AutoCreatePR        types.Bool   `tfsdk:"auto_create_pr"`       // optional
	ResilienceEnabled   types.Bool   `tfsdk:"resilience_enabled"`   // optional

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

// ScopeModel represents a resource scope configuration
type ScopeModel struct {
	Type  types.String `tfsdk:"type"`  // required: tags, resource_group, asset_types, excluded_asset_types, selected_resources, excluded_resources
	Value types.List   `tfsdk:"value"` // required: list of strings, min 1
}

// VCSModel represents VCS integration configuration
type VCSModel struct {
	VCSIntegrationID types.String `tfsdk:"vcs_integration_id"` // optional
	RepoID           types.String `tfsdk:"repo_id"`            // optional
}

// BackupAndDrApplicationDataSourceModel represents the data source model for a single application
// Note: Simplified schema without nested blocks (scope, vcs) to comply with data source limitations
type BackupAndDrApplicationDataSourceModel struct {
	ApplicationID        types.String `tfsdk:"application_id"`
	AccountID            types.String `tfsdk:"account_id"`
	ApplicationName      types.String `tfsdk:"application_name"`
	IntegrationID        types.String `tfsdk:"integration_id"`
	Region               types.String `tfsdk:"region"`
	ProviderType         types.String `tfsdk:"provider_type"`
	Description          types.String `tfsdk:"description"`
	Frequency            types.Int64  `tfsdk:"frequency"`
	NotificationID       types.String `tfsdk:"notification_id"`
	RestoreInstructions  types.String `tfsdk:"restore_instructions"`
	TargetAccount        types.String `tfsdk:"target_account"`
	TargetRegion         types.String `tfsdk:"target_region"`
	AutoCreatePR         types.Bool   `tfsdk:"auto_create_pr"`
	ResilienceEnabled    types.Bool   `tfsdk:"resilience_enabled"`
	Status               types.String `tfsdk:"status"`
	SnapshotsCount       types.Int64  `tfsdk:"snapshots_count"`
	LastBackupSnapshotID types.String `tfsdk:"last_backup_snapshot_id"`
	LastBackupTime       types.String `tfsdk:"last_backup_time"`
	LastBackupStatus     types.String `tfsdk:"last_backup_status"`
	NextBackupTime       types.String `tfsdk:"next_backup_time"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

// BackupAndDrApplicationsDataSourceModel represents the data source model for listing applications
type BackupAndDrApplicationsDataSourceModel struct {
	ID            types.String                            `tfsdk:"id"`
	AccountID     types.String                            `tfsdk:"account_id"`
	Status        types.String                            `tfsdk:"status"`
	IntegrationID types.String                            `tfsdk:"integration_id"`
	Region        types.String                            `tfsdk:"region"`
	ProviderType  types.String                            `tfsdk:"provider_type"`
	Applications  []BackupAndDrApplicationDataSourceModel `tfsdk:"applications"`
}
