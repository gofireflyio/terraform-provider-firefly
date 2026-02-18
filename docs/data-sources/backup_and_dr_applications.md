# firefly_backup_and_dr_applications (Data Source)

Retrieves a list of Firefly Backup & DR applications (backup policies) that match the specified criteria.

## Example Usage

```terraform
# Get all backup policies for an account
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

# Get only active backup policies
data "firefly_backup_and_dr_applications" "active" {
  account_id = "66169d5af4992fc0bab04510"
  status     = "Active"
}

# Get inactive backup policies
data "firefly_backup_and_dr_applications" "inactive" {
  account_id = "66169d5af4992fc0bab04510"
  status     = "Inactive"
}

# Get backup policies for a specific cloud integration
data "firefly_backup_and_dr_applications" "aws_integration" {
  account_id     = "66169d5af4992fc0bab04510"
  integration_id = "692ec8acce65b3dc46cfceb5"
}

# Get backup policies for a specific region
data "firefly_backup_and_dr_applications" "us_east" {
  account_id = "66169d5af4992fc0bab04510"
  region     = "us-east-1"
}

# Get backup policies for a specific cloud provider
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
output "all_policy_names" {
  value = [for policy in data.firefly_backup_and_dr_applications.all.policies : policy.policy_name]
}

output "active_policy_ids" {
  value = [for policy in data.firefly_backup_and_dr_applications.active.policies : policy.policy_id]
}

output "policy_count_by_region" {
  value = {
    for region in distinct([for p in data.firefly_backup_and_dr_applications.all.policies : p.region]) :
    region => length([for p in data.firefly_backup_and_dr_applications.all.policies : p if p.region == region])
  }
}

# Reference specific policies for other resources
locals {
  daily_backup_policies = [
    for policy in data.firefly_backup_and_dr_applications.all.policies :
    policy if policy.schedule_frequency == "Daily"
  ]

  production_backups = [
    for policy in data.firefly_backup_and_dr_applications.all.policies :
    policy if can(regex("production", lower(policy.policy_name)))
  ]
}

# Find policies with recent backups
locals {
  recently_backed_up = [
    for policy in data.firefly_backup_and_dr_applications.all.policies :
    policy if policy.last_backup_status == "Success"
  ]
}
```

## Schema

### Required

- `account_id` (String) - The account ID to list backup policies for

### Optional

- `status` (String) - Filter policies by status. Valid values: `Active`, `Inactive`
- `integration_id` (String) - Filter policies by cloud integration ID
- `region` (String) - Filter policies by cloud region (e.g., `us-east-1`, `eu-west-1`)
- `provider_type` (String) - Filter policies by cloud provider type (e.g., `aws`, `azure`, `gcp`)

### Read-Only

- `id` (String) - The data source identifier
- `policies` (List of Object) - List of backup policies matching the criteria (see [below for nested schema](#nestedatt--policies))

<a id="nestedatt--policies"></a>
### Nested Schema for `policies`

#### Read-Only

- `policy_id` (String) - The unique identifier of the backup policy
- `account_id` (String) - The account ID
- `policy_name` (String) - The name of the backup policy
- `integration_id` (String) - The integration ID for cloud provider credentials
- `region` (String) - The cloud region where backups are stored
- `provider_type` (String) - The cloud provider type
- `description` (String) - Description of the backup policy
- `schedule_frequency` (String) - Backup schedule frequency (`One-time`, `Daily`, `Weekly`, `Monthly`)
- `notification_id` (String) - Notification channel ID for backup alerts
- `restore_instructions` (String) - Instructions for restoring from backups
- `backup_on_save` (Boolean) - Whether backups are triggered on policy save
- `status` (String) - Current status of the policy (`Active` or `Inactive`)
- `snapshots_count` (Number) - Number of snapshots created by this policy
- `last_backup_snapshot_id` (String) - ID of the most recent backup snapshot
- `last_backup_time` (String) - Timestamp of the last backup
- `last_backup_status` (String) - Status of the last backup (e.g., `Success`, `Failed`, `InProgress`)
- `next_backup_time` (String) - Timestamp of the next scheduled backup
- `created_at` (String) - Timestamp when the policy was created
- `updated_at` (String) - Timestamp when the policy was last updated

## Usage Examples

### Listing All Policies
```terraform
data "firefly_backup_and_dr_applications" "all" {
  account_id = "66169d5af4992fc0bab04510"
}

output "total_policies" {
  value = length(data.firefly_backup_and_dr_applications.all.policies)
}

output "policy_summary" {
  value = {
    for policy in data.firefly_backup_and_dr_applications.all.policies :
    policy.policy_name => {
      id        = policy.policy_id
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
  active_us_east_policy_ids = [
    for policy in data.firefly_backup_and_dr_applications.active_us_east.policies :
    policy.policy_id
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
    for policy in data.firefly_backup_and_dr_applications.all.policies :
    {
      name              = policy.policy_name
      id                = policy.policy_id
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
    for policy in data.firefly_backup_and_dr_applications.all.policies :
    "${policy.provider_type}/${policy.region}" => policy.policy_name...
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
    for freq in distinct([for p in data.firefly_backup_and_dr_applications.all.policies : p.schedule_frequency]) :
    freq => length([for p in data.firefly_backup_and_dr_applications.all.policies : p if p.schedule_frequency == freq])
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
    total_policies   = length(data.firefly_backup_and_dr_applications.integration_policies.policies)
    active_count     = length([for p in data.firefly_backup_and_dr_applications.integration_policies.policies : p if p.status == "Active"])
    total_snapshots  = sum([for p in data.firefly_backup_and_dr_applications.integration_policies.policies : p.snapshots_count])
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
    for policy in data.firefly_backup_and_dr_applications.active.policies :
    policy if policy.last_backup_time == "" || policy.last_backup_status != "Success"
  ]
}

output "backup_health_alert" {
  value = length(local.policies_without_recent_backups) > 0 ? {
    alert = "WARNING: ${length(local.policies_without_recent_backups)} active policies without successful recent backups"
    policies = [for p in local.policies_without_recent_backups : p.policy_name]
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
      total_policies      = length(data.firefly_backup_and_dr_applications.all.policies)
      active_policies     = length([for p in data.firefly_backup_and_dr_applications.all.policies : p if p.status == "Active"])
      inactive_policies   = length([for p in data.firefly_backup_and_dr_applications.all.policies : p if p.status == "Inactive"])
      total_snapshots     = sum([for p in data.firefly_backup_and_dr_applications.all.policies : p.snapshots_count])
    }
    by_frequency = {
      for freq in distinct([for p in data.firefly_backup_and_dr_applications.all.policies : p.schedule_frequency]) :
      freq => {
        count    = length([for p in data.firefly_backup_and_dr_applications.all.policies : p if p.schedule_frequency == freq])
        policies = [for p in data.firefly_backup_and_dr_applications.all.policies : p.policy_name if p.schedule_frequency == freq]
      }
    }
    by_region = {
      for region in distinct([for p in data.firefly_backup_and_dr_applications.all.policies : p.region]) :
      region => length([for p in data.firefly_backup_and_dr_applications.all.policies : p if p.region == region])
    }
    by_provider = {
      for provider in distinct([for p in data.firefly_backup_and_dr_applications.all.policies : p.provider_type]) :
      provider => length([for p in data.firefly_backup_and_dr_applications.all.policies : p if p.provider_type == provider])
    }
  }
}
```

## Common Use Cases

### Audit and Compliance Reporting
Use this data source to generate compliance reports showing all backup policies, their schedules, and last backup status for audit purposes.

### Backup Monitoring and Alerting
Query backup policies to identify failed backups, missing backups, or policies without recent successful snapshots for monitoring dashboards.

### Cost Analysis
Analyze backup policies by region and provider to understand backup storage distribution and optimize costs.

### Policy Discovery
Find existing backup policies before creating new ones to avoid duplication and ensure consistent backup coverage.

### Integration Validation
Verify that all cloud integrations have appropriate backup policies configured and that policies are actively backing up resources.

## Filter Behavior

- **No Filters**: Returns all backup policies for the specified account
- **Single Filter**: Returns only policies matching that specific criterion
- **Multiple Filters**: Returns policies matching ALL specified criteria (AND logic)
- **Empty Results**: If no policies match the criteria, the `policies` list will be empty (not an error)

## Notes

- The data source returns a simplified representation of policies. Nested blocks like `schedule`, `scope`, and `vcs` are not included in the output. Use the `firefly_backup_and_dr_application` resource for full policy details.
- The `schedule_frequency` field provides a simplified view of the schedule (frequency only, without specific times or days).
- All timestamp fields are in ISO 8601 format.
- The `snapshots_count` includes all snapshots created by the policy, both successful and failed.
