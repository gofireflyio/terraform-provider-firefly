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
data "firefly_backup_and_dr_applications" "aws_policies" {
  account_id    = "66169d5af4992fc0bab04510"
  provider_type = "aws"
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
  value = [for policy in data.firefly_backup_and_dr_applications.all.applications : policy.application_name]
}

output "active_application_ids" {
  value = [for policy in data.firefly_backup_and_dr_applications.active.applications : policy.application_id]
}

output "policy_count_by_region" {
  value = {
    for region in distinct([for p in data.firefly_backup_and_dr_applications.all.applications : p.region]) :
    region => length([for p in data.firefly_backup_and_dr_applications.all.applications : p if p.region == region])
  }
}

# Reference specific policies for other resources
locals {
  daily_backup_policies = [
    for policy in data.firefly_backup_and_dr_applications.all.applications :
    policy if policy.schedule_frequency == "Daily"
  ]

  production_backups = [
    for policy in data.firefly_backup_and_dr_applications.all.applications :
    policy if can(regex("production", lower(policy.application_name)))
  ]
}

# Find policies with recent backups
locals {
  recently_backed_up = [
    for policy in data.firefly_backup_and_dr_applications.all.applications :
    policy if policy.last_backup_status == "Success"
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
- `schedule_frequency` (String) - Backup schedule frequency (`One-time`, `Daily`, `Weekly`, `Monthly`)
- `notification_id` (String) - Notification channel ID for backup alerts
- `restore_instructions` (String) - Instructions for restoring from backups
- `backup_on_save` (Boolean) - Whether backups are triggered on application save
- `status` (String) - Current status of the application (`Active` or `Inactive`)
- `snapshots_count` (Number) - Number of snapshots created by this application
- `last_backup_snapshot_id` (String) - ID of the most recent backup snapshot
- `last_backup_time` (String) - Timestamp of the last backup
- `last_backup_status` (String) - Status of the last backup (e.g., `Success`, `Failed`, `InProgress`)
- `next_backup_time` (String) - Timestamp of the next scheduled backup
- `created_at` (String) - Timestamp when the application was created
- `updated_at` (String) - Timestamp when the application was last updated

## Usage Examples

### Listing All Policies
```terraform
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

output "total_policies" {
  value = length(data.firefly_backup_and_dr_applications.all.applications)
}

output "policy_summary" {
  value = {
    for policy in data.firefly_backup_and_dr_applications.all.applications :
    policy.application_name => {
      id        = policy.application_id
      status    = policy.status
      frequency = policy.schedule_frequency
      region    = policy.region
    }
  }
}
```

### Filtering Active Policies by Region
```terraform
data "firefly_backup_and_dr_applications" "active_us_east" {
  account_id = "66169d5af4992fc0bab04510"
  status     = "Active"
  region     = "us-east-1"
}

# Use the filtered policies
locals {
  active_us_east_application_ids = [
    for policy in data.firefly_backup_and_dr_applications.active_us_east.applications :
    policy.application_id
  ]
}
```

### Finding Policies with Backup Issues
```terraform
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

# Identify policies with failed last backup
locals {
  failed_backups = [
    for policy in data.firefly_backup_and_dr_applications.all.applications :
    {
      name              = policy.application_name
      id                = policy.application_id
      last_backup_time  = policy.last_backup_time
      status            = policy.last_backup_status
    }
    if policy.last_backup_status == "Failed"
  ]
}

output "failed_backup_policies" {
  value = local.failed_backups
}
```

### Grouping Policies by Provider and Region
```terraform
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

output "policies_by_provider_region" {
  value = {
    for policy in data.firefly_backup_and_dr_applications.all.applications :
    "${policy.provider_type}/${policy.region}" => policy.application_name...
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
    for freq in distinct([for p in data.firefly_backup_and_dr_applications.all.applications : p.schedule_frequency]) :
    freq => length([for p in data.firefly_backup_and_dr_applications.all.applications : p if p.schedule_frequency == freq])
  }
}

# Example output:
# {
#   "Daily" = 15
#   "Weekly" = 8
#   "Monthly" = 3
#   "One-time" = 2
# }
```

### Finding Policies by Integration
```terraform
# Get all policies for a specific cloud integration
data "firefly_backup_and_dr_applications" "integration_policies" {
  account_id     = "66169d5af4992fc0bab04510"
  integration_id = "692ec8acce65b3dc46cfceb5"
}

output "integration_backup_summary" {
  value = {
    total_policies   = length(data.firefly_backup_and_dr_applications.integration_policies.applications)
    active_count     = length([for p in data.firefly_backup_and_dr_applications.integration_policies.applications : p if p.status == "Active"])
    total_snapshots  = sum([for p in data.firefly_backup_and_dr_applications.integration_policies.applications : p.snapshots_count])
  }
}
```

### Monitoring Backup Health
```terraform
data "firefly_backup_and_dr_applications" "active" {
  account_id = "66169d5af4992fc0bab04510"
  status     = "Active"
}

# Check for policies without recent backups
locals {
  policies_without_recent_backups = [
    for policy in data.firefly_backup_and_dr_applications.active.applications :
    policy if policy.last_backup_time == "" || policy.last_backup_status != "Success"
  ]
}

output "backup_health_alert" {
  value = length(local.applications_without_recent_backups) > 0 ? {
    alert = "WARNING: ${length(local.applications_without_recent_backups)} active policies without successful recent backups"
    policies = [for p in local.applications_without_recent_backups : p.application_name]
  } : {
    alert = "OK: All active policies have successful recent backups"
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
      total_policies      = length(data.firefly_backup_and_dr_applications.all.applications)
      active_policies     = length([for p in data.firefly_backup_and_dr_applications.all.applications : p if p.status == "Active"])
      inactive_policies   = length([for p in data.firefly_backup_and_dr_applications.all.applications : p if p.status == "Inactive"])
      total_snapshots     = sum([for p in data.firefly_backup_and_dr_applications.all.applications : p.snapshots_count])
    }
    by_frequency = {
      for freq in distinct([for p in data.firefly_backup_and_dr_applications.all.applications : p.schedule_frequency]) :
      freq => {
        count    = length([for p in data.firefly_backup_and_dr_applications.all.applications : p if p.schedule_frequency == freq])
        policies = [for p in data.firefly_backup_and_dr_applications.all.applications : p.application_name if p.schedule_frequency == freq]
      }
    }
    by_region = {
      for region in distinct([for p in data.firefly_backup_and_dr_applications.all.applications : p.region]) :
      region => length([for p in data.firefly_backup_and_dr_applications.all.applications : p if p.region == region])
    }
    by_provider = {
      for provider in distinct([for p in data.firefly_backup_and_dr_applications.all.applications : p.provider_type]) :
      provider => length([for p in data.firefly_backup_and_dr_applications.all.applications : p if p.provider_type == provider])
    }
  }
}
```

## Common Use Cases

### Audit and Compliance Reporting
Use this data source to generate compliance reports showing all backup applications, their schedules, and last backup status for audit purposes.

### Backup Monitoring and Alerting
Query backup applications to identify failed backups, missing backups, or policies without recent successful snapshots for monitoring dashboards.

### Cost Analysis
Analyze backup applications by region and provider to understand backup storage distribution and optimize costs.

### Policy Discovery
Find existing backup applications before creating new ones to avoid duplication and ensure consistent backup coverage.

### Integration Validation
Verify that all cloud integrations have appropriate backup applications configured and that policies are actively backing up resources.

## Filter Behavior

- **No Filters**: Returns all backup applications for the specified account
- **Single Filter**: Returns only policies matching that specific criterion
- **Multiple Filters**: Returns policies matching ALL specified criteria (AND logic)
- **Empty Results**: If no policies match the criteria, the `applications` list will be empty (not an error)

## Notes

- The data source returns a simplified representation of policies. Nested blocks like `schedule`, `scope`, and `vcs` are not included in the output. Use the `firefly_backup_and_dr_application` resource for full policy details.
- The `schedule_frequency` field provides a simplified view of the schedule (frequency only, without specific times or days).
- All timestamp fields are in ISO 8601 format.
- The `snapshots_count` includes all snapshots created by the application, both successful and failed.
