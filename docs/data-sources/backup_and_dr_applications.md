# firefly_backup_and_dr_applications (Data Source)

Retrieves a list of Firefly Backup & DR applications (backup applications) that match the specified criteria.

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

# Get inactive backup applications
data "firefly_backup_and_dr_applications" "inactive" {
  account_id = "66169d5af4992fc0bab04510"
  status     = "Inactive"
}

# Get backup applications for a specific cloud integration
data "firefly_backup_and_dr_applications" "aws_integration" {
  account_id     = "66169d5af4992fc0bab04510"
  integration_id = "692ec8acce65b3dc46cfceb5"
}

# Get backup applications for a specific region
data "firefly_backup_and_dr_applications" "us_east" {
  account_id = "66169d5af4992fc0bab04510"
  region     = "us-east-1"
}

# Get backup applications for a specific cloud provider
data "firefly_backup_and_dr_applications" "aws_apps" {
  account_id    = "66169d5af4992fc0bab04510"
  provider_type = "aws"
}

# Combined filters: Active AWS applications in us-east-1
data "firefly_backup_and_dr_applications" "filtered" {
  account_id     = "66169d5af4992fc0bab04510"
  status         = "Active"
  provider_type  = "aws"
  region         = "us-east-1"
  integration_id = "692ec8acce65b3dc46cfceb5"
}

# Use results in outputs
output "all_application_names" {
  value = [for app in data.firefly_backup_and_dr_applications.all.applications : app.application_name]
}

output "active_application_ids" {
  value = [for app in data.firefly_backup_and_dr_applications.active.applications : app.application_id]
}

output "application_count_by_region" {
  value = {
    for region in distinct([for app in data.firefly_backup_and_dr_applications.all.applications : app.region]) :
    region => length([for app in data.firefly_backup_and_dr_applications.all.applications : app if app.region == region])
  }
}

# Reference specific applications for other resources
locals {
  frequent_backup_apps = [
    for app in data.firefly_backup_and_dr_applications.all.applications :
    app if app.frequency <= 8
  ]

  production_backups = [
    for app in data.firefly_backup_and_dr_applications.all.applications :
    app if can(regex("production", lower(app.application_name)))
  ]
}

# Find applications with recent successful backups
locals {
  recently_backed_up = [
    for app in data.firefly_backup_and_dr_applications.all.applications :
    app if app.last_backup_status == "Success"
  ]
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

## Usage Examples

### Listing All Applications
```terraform
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

output "total_applications" {
  value = length(data.firefly_backup_and_dr_applications.all.applications)
}

output "application_summary" {
  value = {
    for app in data.firefly_backup_and_dr_applications.all.applications :
    app.application_name => {
      id        = app.application_id
      status    = app.status
      frequency = app.frequency
      region    = app.region
    }
  }
}
```

### Filtering Active Applications by Region
```terraform
data "firefly_backup_and_dr_applications" "active_us_east" {
  account_id = "66169d5af4992fc0bab04510"
  status     = "Active"
  region     = "us-east-1"
}

locals {
  active_us_east_application_ids = [
    for app in data.firefly_backup_and_dr_applications.active_us_east.applications :
    app.application_id
  ]
}
```

### Finding Applications with Backup Issues
```terraform
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

locals {
  failed_backups = [
    for app in data.firefly_backup_and_dr_applications.all.applications :
    {
      name             = app.application_name
      id               = app.application_id
      last_backup_time = app.last_backup_time
      status           = app.last_backup_status
    }
    if app.last_backup_status == "Failed"
  ]
}

output "failed_backup_applications" {
  value = local.failed_backups
}
```

### Grouping Applications by Provider and Region
```terraform
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

output "applications_by_provider_region" {
  value = {
    for app in data.firefly_backup_and_dr_applications.all.applications :
    "${app.provider_type}/${app.region}" => app.application_name...
  }
}
```

### Analyzing Backup Frequency Distribution
```terraform
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

output "frequency_distribution" {
  value = {
    for freq in distinct([for app in data.firefly_backup_and_dr_applications.all.applications : tostring(app.frequency)]) :
    freq => length([for app in data.firefly_backup_and_dr_applications.all.applications : app if tostring(app.frequency) == freq])
  }
}

# Example output:
# {
#   "4"  = 3
#   "8"  = 8
#   "16" = 5
#   "24" = 15
# }
```

### Finding Applications by Integration
```terraform
data "firefly_backup_and_dr_applications" "integration_apps" {
  account_id     = "66169d5af4992fc0bab04510"
  integration_id = "692ec8acce65b3dc46cfceb5"
}

output "integration_backup_summary" {
  value = {
    total_applications = length(data.firefly_backup_and_dr_applications.integration_apps.applications)
    active_count       = length([for app in data.firefly_backup_and_dr_applications.integration_apps.applications : app if app.status == "Active"])
    total_snapshots    = sum([for app in data.firefly_backup_and_dr_applications.integration_apps.applications : app.snapshots_count])
  }
}
```

### Monitoring Backup Health
```terraform
data "firefly_backup_and_dr_applications" "active" {
  account_id = "66169d5af4992fc0bab04510"
  status     = "Active"
}

locals {
  applications_without_recent_backups = [
    for app in data.firefly_backup_and_dr_applications.active.applications :
    app if app.last_backup_time == "" || app.last_backup_status != "Success"
  ]
}

output "backup_health_alert" {
  value = length(local.applications_without_recent_backups) > 0 ? {
    alert        = "WARNING: ${length(local.applications_without_recent_backups)} active applications without successful recent backups"
    applications = [for app in local.applications_without_recent_backups : app.application_name]
  } : {
    alert = "OK: All active applications have successful recent backups"
  }
}
```

### Creating Backup Inventory Report
```terraform
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

output "backup_inventory_report" {
  value = {
    summary = {
      total_applications  = length(data.firefly_backup_and_dr_applications.all.applications)
      active_applications = length([for app in data.firefly_backup_and_dr_applications.all.applications : app if app.status == "Active"])
      inactive_applications = length([for app in data.firefly_backup_and_dr_applications.all.applications : app if app.status == "Inactive"])
      total_snapshots     = sum([for app in data.firefly_backup_and_dr_applications.all.applications : app.snapshots_count])
    }
    by_frequency = {
      for freq in distinct([for app in data.firefly_backup_and_dr_applications.all.applications : tostring(app.frequency)]) :
      freq => {
        count        = length([for app in data.firefly_backup_and_dr_applications.all.applications : app if tostring(app.frequency) == freq])
        applications = [for app in data.firefly_backup_and_dr_applications.all.applications : app.application_name if tostring(app.frequency) == freq]
      }
    }
    by_region = {
      for region in distinct([for app in data.firefly_backup_and_dr_applications.all.applications : app.region]) :
      region => length([for app in data.firefly_backup_and_dr_applications.all.applications : app if app.region == region])
    }
    by_provider = {
      for provider in distinct([for app in data.firefly_backup_and_dr_applications.all.applications : app.provider_type]) :
      provider => length([for app in data.firefly_backup_and_dr_applications.all.applications : app if app.provider_type == provider])
    }
  }
}
```

## Common Use Cases

### Audit and Compliance Reporting
Use this data source to generate compliance reports showing all backup applications, their schedules, and last backup status for audit purposes.

### Backup Monitoring and Alerting
Query backup applications to identify failed backups, missing backups, or applications without recent successful snapshots for monitoring dashboards.

### Cost Analysis
Analyze backup applications by region and provider to understand backup storage distribution and optimize costs.

### Application Discovery
Find existing backup applications before creating new ones to avoid duplication and ensure consistent backup coverage.

### Integration Validation
Verify that all cloud integrations have appropriate backup applications configured and that applications are actively backing up resources.

## Filter Behavior

- **No Filters**: Returns all backup applications for the specified account
- **Single Filter**: Returns only applications matching that specific criterion
- **Multiple Filters**: Returns applications matching ALL specified criteria (AND logic)
- **Empty Results**: If no applications match the criteria, the `applications` list will be empty (not an error)

## Notes

- The data source returns a simplified representation of applications. Nested blocks like `scope` and `vcs` are not included in the output. Use the `firefly_backup_and_dr_application` resource for full application details.
- All timestamp fields are in ISO 8601 format.
- The `snapshots_count` includes all snapshots created by the application, both successful and failed.
