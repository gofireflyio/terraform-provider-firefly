# firefly_backup_and_dr_application (Resource)

Manages a Firefly Backup & DR application (backup policy) for automated infrastructure backup and disaster recovery.

## Example Usage

```terraform
# Basic one-time backup with tag-based scope
resource "firefly_backup_and_dr_application" "simple_backup" {
  account_id     = "66169d5af4992fc0bab04510"
  policy_name    = "Simple One-time Backup"
  integration_id = "692ec8acce65b3dc46cfceb5"
  region         = "us-east-1"
  provider_type  = "aws"
  description    = "One-time backup of production resources"

  schedule {
    frequency = "One-time"
  }

  scope {
    type  = "tags"
    value = ["Environment:Production"]
  }
}

# Daily backup with specific time and multiple tag filters
resource "firefly_backup_and_dr_application" "daily_backup" {
  account_id            = "66169d5af4992fc0bab04510"
  policy_name           = "Daily Production Backup"
  integration_id        = "692ec8acce65b3dc46cfceb5"
  region                = "us-east-1"
  provider_type         = "aws"
  description           = "Daily backup of all production resources"
  notification_id       = "notification-channel-id"
  restore_instructions  = "Contact DevOps team for restore procedures"
  backup_on_save        = true

  schedule {
    frequency = "Daily"
    hour      = 2   # 2 AM UTC
    minute    = 30  # 2:30 AM UTC
  }

  scope {
    type  = "tags"
    value = ["Environment:Production", "Backup:Required"]
  }
}

# Weekly backup on specific days with asset type filtering
resource "firefly_backup_and_dr_application" "weekly_backup" {
  account_id     = "66169d5af4992fc0bab04510"
  policy_name    = "Weekly Database Backup"
  integration_id = "692ec8acce65b3dc46cfceb5"
  region         = "eu-west-1"
  provider_type  = "aws"
  description    = "Weekly backup of database resources"

  schedule {
    frequency    = "Weekly"
    days_of_week = ["Monday", "Friday"]
    hour         = 3
    minute       = 0
  }

  scope {
    type  = "asset_types"
    value = ["aws_db_instance", "aws_rds_cluster"]
  }
}

# Monthly backup on specific day with resource group scope
resource "firefly_backup_and_dr_application" "monthly_specific_day" {
  account_id     = "66169d5af4992fc0bab04510"
  policy_name    = "Monthly First-Day Backup"
  integration_id = "692ec8acce65b3dc46cfceb5"
  region         = "us-west-2"
  provider_type  = "aws"
  description    = "Monthly backup on the 1st of each month"

  schedule {
    frequency              = "Monthly"
    monthly_schedule_type  = "specific_day"
    day_of_month           = 1
    hour                   = 0
    minute                 = 0
  }

  scope {
    type  = "resource_group"
    value = ["production-rg", "critical-rg"]
  }
}

# Monthly backup on specific weekday (e.g., last Friday)
resource "firefly_backup_and_dr_application" "monthly_last_friday" {
  account_id     = "66169d5af4992fc0bab04510"
  policy_name    = "Last Friday Monthly Backup"
  integration_id = "692ec8acce65b3dc46cfceb5"
  region         = "us-east-1"
  provider_type  = "aws"
  description    = "Monthly backup on the last Friday of each month"

  schedule {
    frequency              = "Monthly"
    monthly_schedule_type  = "specific_weekday"
    weekday_ordinal        = "Last"
    weekday_name           = "Friday"
    hour                   = 23
    minute                 = 0
  }

  scope {
    type  = "tags"
    value = ["CriticalData:True"]
  }
}

# Monthly backup on last day of month
resource "firefly_backup_and_dr_application" "monthly_last_day" {
  account_id     = "66169d5af4992fc0bab04510"
  policy_name    = "End of Month Backup"
  integration_id = "692ec8acce65b3dc46cfceb5"
  region         = "us-east-1"
  provider_type  = "aws"
  description    = "Backup on the last day of each month"

  schedule {
    frequency              = "Monthly"
    monthly_schedule_type  = "last_day"
    hour                   = 23
    minute                 = 59
  }

  scope {
    type  = "tags"
    value = ["MonthlyBackup:Required"]
  }
}

# Backup with VCS integration for artifact storage
resource "firefly_backup_and_dr_application" "backup_with_vcs" {
  account_id     = "66169d5af4992fc0bab04510"
  policy_name    = "Backup to Git Repository"
  integration_id = "692ec8acce65b3dc46cfceb5"
  region         = "us-east-1"
  provider_type  = "aws"
  description    = "Backup stored in git repository"

  schedule {
    frequency = "Daily"
    hour      = 1
    minute    = 0
  }

  scope {
    type  = "tags"
    value = ["IaC:Terraform"]
  }

  vcs {
    project_id         = "project-123"
    vcs_integration_id = "vcs-integration-456"
    repo_id            = "repo-789"
  }
}

# Backup with multiple scope types
resource "firefly_backup_and_dr_application" "multi_scope_backup" {
  account_id     = "66169d5af4992fc0bab04510"
  policy_name    = "Multi-Scope Backup"
  integration_id = "692ec8acce65b3dc46cfceb5"
  region         = "us-east-1"
  provider_type  = "aws"
  description    = "Backup using multiple scope criteria"

  schedule {
    frequency = "Daily"
    hour      = 4
    minute    = 0
  }

  # Backup resources matching tags
  scope {
    type  = "tags"
    value = ["Environment:Production"]
  }

  # AND specific resource types
  scope {
    type  = "asset_types"
    value = ["aws_instance", "aws_ebs_volume"]
  }

  # AND from specific resource groups
  scope {
    type  = "resource_group"
    value = ["prod-app-rg"]
  }
}

# Backup with specific resource IDs
resource "firefly_backup_and_dr_application" "selected_resources_backup" {
  account_id     = "66169d5af4992fc0bab04510"
  policy_name    = "Specific Resources Backup"
  integration_id = "692ec8acce65b3dc46cfceb5"
  region         = "us-east-1"
  provider_type  = "aws"
  description    = "Backup of hand-picked critical resources"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 0
  }

  scope {
    type  = "selected_resources"
    value = [
      "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
      "arn:aws:s3:::my-critical-bucket",
      "arn:aws:rds:us-east-1:123456789012:db:production-db"
    ]
  }
}
```

## Schema

### Required

- `account_id` (String) - The account ID for the backup policy. **Note**: Changing this forces replacement of the resource.
- `policy_name` (String) - The name of the backup policy (max 100 characters)
- `integration_id` (String) - The integration ID for cloud provider credentials
- `region` (String) - The cloud region where backups will be stored
- `provider_type` (String) - The cloud provider type (max 50 characters, e.g., `aws`, `azure`, `gcp`)
- `schedule` (Block) - Backup schedule configuration (see [below for nested schema](#nestedblock--schedule))
- `scope` (Block List) - Resource scope configurations for backup targeting (see [below for nested schema](#nestedblock--scope))

### Optional

- `description` (String) - Description of the backup policy (max 500 characters)
- `notification_id` (String) - Notification channel ID for backup alerts
- `restore_instructions` (String) - Instructions for restoring from backups (max 2000 characters)
- `backup_on_save` (Boolean) - Whether to trigger a backup immediately on policy creation/update. Defaults to `true`.
- `vcs` (Block) - VCS integration configuration for backup artifacts (see [below for nested schema](#nestedblock--vcs))

### Read-Only

- `id` (String) - The unique identifier of the backup policy
- `status` (String) - Current status of the policy (`Active` or `Inactive`)
- `snapshots_count` (Number) - Number of snapshots created by this policy
- `last_backup_snapshot_id` (String) - ID of the most recent backup snapshot
- `last_backup_time` (String) - Timestamp of the last backup
- `last_backup_status` (String) - Status of the last backup
- `next_backup_time` (String) - Timestamp of the next scheduled backup
- `created_at` (String) - Timestamp when the policy was created
- `updated_at` (String) - Timestamp when the policy was last updated

<a id="nestedblock--schedule"></a>
### Nested Schema for `schedule`

#### Required

- `frequency` (String) - Backup frequency. Valid values: `One-time`, `Daily`, `Weekly`, `Monthly`

#### Optional

- `hour` (Number) - Hour of day for backup (0-23). Used for Daily, Weekly, and Monthly schedules.
- `minute` (Number) - Minute of hour for backup (0-59). Used for Daily, Weekly, and Monthly schedules.
- `days_of_week` (List of String) - Days of week for Weekly schedule (e.g., `["Monday", "Friday"]`). **Required** when `frequency` is `Weekly`.
- `monthly_schedule_type` (String) - Type of monthly schedule. Valid values: `specific_day`, `specific_weekday`, `last_day`. **Required** when `frequency` is `Monthly`.
- `day_of_month` (Number) - Day of month for `specific_day` monthly schedule (1-31). **Required** when `monthly_schedule_type` is `specific_day`.
- `weekday_ordinal` (String) - Weekday ordinal for `specific_weekday` monthly schedule. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`. **Required** when `monthly_schedule_type` is `specific_weekday`.
- `weekday_name` (String) - Weekday name for `specific_weekday` monthly schedule (e.g., `Sunday`, `Monday`, etc.). **Required** when `monthly_schedule_type` is `specific_weekday`.
- `cron_expression` (String) - Cron expression as alternative to explicit schedule parameters

<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

#### Required

- `type` (String) - Scope type. Valid values: `tags`, `resource_group`, `asset_types`, `selected_resources`
- `value` (List of String) - List of values for this scope type. Format depends on type:
  - `tags`: Tag key-value pairs in format `"key:value"` (e.g., `["Environment:Production", "Team:DevOps"]`)
  - `resource_group`: Resource group names (e.g., `["prod-rg", "staging-rg"]`)
  - `asset_types`: Terraform resource types (e.g., `["aws_instance", "aws_db_instance"]`)
  - `selected_resources`: Specific resource ARNs or IDs

**Note**: Multiple scope blocks can be defined. Resources must match ALL scope criteria (AND logic).

<a id="nestedblock--vcs"></a>
### Nested Schema for `vcs`

#### Optional

- `project_id` (String) - Project ID for VCS integration
- `vcs_integration_id` (String) - VCS integration ID for version control system
- `repo_id` (String) - Repository ID for storing backup artifacts

## Schedule Configuration Guidelines

### One-time Backup
For immediate, single-execution backups:
```terraform
schedule {
  frequency = "One-time"
}
```

### Daily Backup
For daily recurring backups:
```terraform
schedule {
  frequency = "Daily"
  hour      = 2   # Optional: 2 AM UTC
  minute    = 30  # Optional: 2:30 AM UTC
}
```
If `hour` and `minute` are not specified, the backup runs at midnight UTC.

### Weekly Backup
For weekly recurring backups on specific days:
```terraform
schedule {
  frequency    = "Weekly"
  days_of_week = ["Monday", "Wednesday", "Friday"]  # Required
  hour         = 3   # Optional: 3 AM UTC
  minute       = 0   # Optional: on the hour
}
```

### Monthly Backup - Specific Day
For monthly backups on a specific day of the month:
```terraform
schedule {
  frequency              = "Monthly"
  monthly_schedule_type  = "specific_day"
  day_of_month           = 15  # Required: 15th of each month
  hour                   = 1   # Optional: 1 AM UTC
  minute                 = 0   # Optional: on the hour
}
```

### Monthly Backup - Specific Weekday
For monthly backups on a specific weekday occurrence (e.g., "First Monday", "Last Friday"):
```terraform
schedule {
  frequency              = "Monthly"
  monthly_schedule_type  = "specific_weekday"
  weekday_ordinal        = "First"    # Required: First/Second/Third/Fourth/Last
  weekday_name           = "Monday"   # Required: day of week
  hour                   = 2          # Optional: 2 AM UTC
  minute                 = 0          # Optional: on the hour
}
```

### Monthly Backup - Last Day
For monthly backups on the last day of each month:
```terraform
schedule {
  frequency              = "Monthly"
  monthly_schedule_type  = "last_day"
  hour                   = 23  # Optional: 11 PM UTC
  minute                 = 59  # Optional: 11:59 PM UTC
}
```

## Scope Targeting Best Practices

### Tag-Based Scoping
Target resources by their tags. Use key:value format:
```terraform
scope {
  type  = "tags"
  value = ["Environment:Production", "Backup:Critical", "Team:Platform"]
}
```

### Resource Group Scoping
Target all resources within specific resource groups:
```terraform
scope {
  type  = "resource_group"
  value = ["production-rg", "database-rg"]
}
```

### Asset Type Scoping
Target specific Terraform resource types:
```terraform
scope {
  type  = "asset_types"
  value = ["aws_instance", "aws_db_instance", "aws_ebs_volume", "aws_s3_bucket"]
}
```

### Selected Resources Scoping
Target specific resources by their unique identifiers:
```terraform
scope {
  type  = "selected_resources"
  value = [
    "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
    "arn:aws:s3:::critical-data-bucket"
  ]
}
```

### Combining Multiple Scopes
Define multiple scope blocks to create AND conditions. Resources must match ALL criteria:
```terraform
# Backup only Production EC2 instances from specific resource group
scope {
  type  = "tags"
  value = ["Environment:Production"]
}

scope {
  type  = "asset_types"
  value = ["aws_instance"]
}

scope {
  type  = "resource_group"
  value = ["web-tier-rg"]
}
```

## Important Notes

- **Immediate Backup**: By default, `backup_on_save` is `true`, which triggers an immediate backup when the policy is created or updated. Set to `false` to disable this behavior.
- **Account ID Changes**: Changing the `account_id` attribute forces replacement (destroy and recreate) of the resource.
- **Computed Fields**: All timestamp and status fields are read-only and automatically updated by Firefly.
- **Schedule Validation**: The provider validates schedule parameters based on frequency type. For example, `days_of_week` is required for Weekly frequency but invalid for Daily.
- **VCS Integration**: When VCS is configured, backup artifacts are stored in the specified git repository in addition to cloud storage.

## Import

Backup & DR applications can be imported using the format `account_id:policy_id`:

```shell
terraform import firefly_backup_and_dr_application.example 66169d5af4992fc0bab04510:policy-id-here
```
