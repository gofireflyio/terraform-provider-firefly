# firefly_backup_and_dr_applications (Data Source)

Retrieves a list of Firefly Backup & DR applications that match the specified criteria.

## Example Usage

```terraform
# Get all backup applications for an account
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

# Get only active backup applications
data "firefly_backup_and_dr_applications" "active" {
  account_id = "66169d5af4992fc0bab04510"
  status     = "Active"
}

# Get backup applications for a specific cloud integration
data "firefly_backup_and_dr_applications" "aws_integration" {
  account_id     = "66169d5af4992fc0bab04510"
  integration_id = "692ec8acce65b3dc46cfceb5"
}

# Combined filters: Active AWS policies in us-east-1
data "firefly_backup_and_dr_applications" "filtered" {
  account_id     = "66169d5af4992fc0bab04510"
  status         = "Active"
  provider_type  = "aws"
  region         = "us-east-1"
  integration_id = "692ec8acce65b3dc46cfceb5"
}

# Use results in outputs
output "all_application_names" {
  value = [for p in data.firefly_backup_and_dr_applications.all.applications : p.application_name]
}

output "active_application_ids" {
  value = [for p in data.firefly_backup_and_dr_applications.active.applications : p.application_id]
}

# Find policies with failed last backup
locals {
  failed_backups = [
    for p in data.firefly_backup_and_dr_applications.all.applications :
    { name = p.application_name, id = p.application_id, last_backup_time = p.last_backup_time }
    if p.last_backup_status == "Failed"
  ]
}

# Group by frequency
output "frequency_distribution" {
  value = {
    for freq in distinct([for p in data.firefly_backup_and_dr_applications.all.applications : tostring(p.frequency)]) :
    freq => length([for p in data.firefly_backup_and_dr_applications.all.applications : p if tostring(p.frequency) == freq])
  }
}
```

## Schema

### Required

- `account_id` (String) - The account ID to list backup applications for

### Optional

- `status` (String) - Filter applications by status. Valid values: `Active`, `Inactive`
- `integration_id` (String) - Filter applications by cloud integration ID
- `region` (String) - Filter applications by cloud region (e.g., `us-east-1`, `eu-west-1`)
- `provider_type` (String) - Filter applications by cloud provider type (e.g., `aws`, `azure`, `gcp`)

### Read-Only

- `id` (String) - The data source identifier
- `applications` (List of Object) - List of backup applications matching the criteria (see [below for nested schema](#nestedatt--applications))

<a id="nestedatt--applications"></a>
### Nested Schema for `applications`

#### Read-Only

- `application_id` (String) - The unique identifier of the backup application
- `account_id` (String) - The account ID
- `application_name` (String) - The name of the backup application
- `integration_id` (String) - The integration ID for cloud provider credentials
- `region` (String) - The cloud region where backups are stored
- `provider_type` (String) - The cloud provider type
- `description` (String) - Description of the backup application
- `frequency` (Number) - Hours between scheduled backups (4, 8, 16, or 24)
- `notification_id` (String) - Notification channel ID for backup alerts
- `restore_instructions` (String) - Instructions for restoring from backups
- `target_account` (String) - Target account/integration ID where the restore should land
- `target_region` (String) - Target region where the restore should land
- `auto_create_pr` (Boolean) - Whether the restore flow automatically opens a VCS pull request
- `resilience_enabled` (Boolean) - Whether DR scheduling is enabled
- `status` (String) - Current status of the application (`Active` or `Inactive`)
- `snapshots_count` (Number) - Number of snapshots created by this application
- `last_backup_snapshot_id` (String) - ID of the most recent backup snapshot
- `last_backup_time` (String) - Timestamp of the last backup
- `last_backup_status` (String) - Status of the last backup (e.g., `Success`, `Failed`, `InProgress`)
- `next_backup_time` (String) - Timestamp of the next scheduled backup
- `created_at` (String) - Timestamp when the application was created
- `updated_at` (String) - Timestamp when the application was last updated

## Filter Behavior

- **No Filters**: Returns all backup applications for the specified account
- **Single Filter**: Returns only policies matching that specific criterion
- **Multiple Filters**: Returns policies matching ALL specified criteria (AND logic)
- **Empty Results**: If no policies match the criteria, the `applications` list will be empty (not an error)

## Notes

- All timestamp fields are in ISO 8601 format.
- The `snapshots_count` includes all snapshots created by the application, both successful and failed.
