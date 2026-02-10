# Daily backup policy for AWS production resources
resource "firefly_backup_and_dr_application" "production_daily" {
  policy_name    = "production-daily-backup"
  description    = "Daily backup of production AWS infrastructure"
  integration_id = "aws-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 0
  }

  scope {
    type  = "tags"
    value = ["env:production", "backup:enabled"]
  }
}

# Weekly backup with specific days
resource "firefly_backup_and_dr_application" "staging_weekly" {
  policy_name    = "staging-weekly-backup"
  description    = "Weekly backup of staging environment on weekends"
  integration_id = "aws-integration-id"
  region         = "us-west-2"
  provider_type  = "aws"

  schedule {
    frequency    = "Weekly"
    hour         = 4
    minute       = 0
    days_of_week = ["Saturday", "Sunday"]
  }

  scope {
    type  = "tags"
    value = ["env:staging"]
  }
}

# Monthly backup on a specific day
resource "firefly_backup_and_dr_application" "compliance_monthly" {
  policy_name    = "compliance-monthly-backup"
  description    = "Monthly compliance backup on the 1st of each month"
  integration_id = "aws-integration-id"
  region         = "eu-west-1"
  provider_type  = "aws"

  schedule {
    frequency           = "Monthly"
    hour                = 1
    minute              = 0
    monthly_schedule_type = "DayOfMonth"
    day_of_month        = 1
  }
}

# Monthly backup on a specific weekday (e.g., first Monday)
resource "firefly_backup_and_dr_application" "monthly_weekday" {
  policy_name    = "first-monday-backup"
  description    = "Backup on the first Monday of each month"
  integration_id = "aws-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency           = "Monthly"
    hour                = 3
    minute              = 30
    monthly_schedule_type = "WeekdayOfMonth"
    weekday_ordinal     = "First"
    weekday_name        = "Monday"
  }
}

# Cron-based backup schedule
resource "firefly_backup_and_dr_application" "custom_cron" {
  policy_name    = "custom-cron-backup"
  description    = "Backup every 6 hours during business days"
  integration_id = "aws-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency       = "Cron"
    cron_expression = "0 */6 * * 1-5"
  }
}

# Backup with VCS integration for storing state
resource "firefly_backup_and_dr_application" "with_vcs" {
  policy_name    = "vcs-backed-backup"
  description    = "Backup with VCS storage for audit trail"
  integration_id = "aws-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 0
    minute    = 0
  }

  vcs {
    project_id         = "project-id"
    vcs_integration_id = "github-integration-id"
    repo_id            = "repo-id"
  }
}

# Multiple scope filters
resource "firefly_backup_and_dr_application" "multi_scope" {
  policy_name    = "multi-scope-backup"
  description    = "Backup with multiple scope filters"
  integration_id = "aws-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 5
    minute    = 0
  }

  scope {
    type  = "tags"
    value = ["env:production", "tier:critical"]
  }

  scope {
    type  = "resourceTypes"
    value = ["aws_s3_bucket", "aws_dynamodb_table", "aws_rds_cluster"]
  }
}

# Inactive policy (paused)
resource "firefly_backup_and_dr_application" "paused_backup" {
  policy_name    = "paused-backup-policy"
  description    = "Backup policy currently paused for maintenance"
  integration_id = "aws-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"
  status         = "Inactive"

  schedule {
    frequency = "Daily"
    hour      = 12
    minute    = 0
  }
}

# Azure backup policy
resource "firefly_backup_and_dr_application" "azure_daily" {
  policy_name    = "azure-daily-backup"
  description    = "Daily backup of Azure resources"
  integration_id = "azure-integration-id"
  region         = "eastus"
  provider_type  = "azure"

  schedule {
    frequency = "Daily"
    hour      = 3
    minute    = 0
  }

  scope {
    type  = "tags"
    value = ["environment:production"]
  }
}

# Outputs
output "production_backup_id" {
  description = "ID of the production daily backup policy"
  value       = firefly_backup_and_dr_application.production_daily.id
}

output "production_backup_status" {
  description = "Status of the production backup policy"
  value       = firefly_backup_and_dr_application.production_daily.status
}

output "production_next_backup" {
  description = "Next scheduled backup time"
  value       = firefly_backup_and_dr_application.production_daily.next_backup_time
}

output "production_snapshots" {
  description = "Number of snapshots taken"
  value       = firefly_backup_and_dr_application.production_daily.snapshots_count
}
