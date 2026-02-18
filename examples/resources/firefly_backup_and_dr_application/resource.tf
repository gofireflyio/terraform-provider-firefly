# Daily backup with tags scope
resource "firefly_backup_and_dr_application" "daily_production" {
  account_id     = "your-account-id"
  policy_name    = "Daily Production Backup"
  description    = "Backs up all production resources daily at 2 AM"
  integration_id = "aws-integration-123"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 0
  }

  scope {
    type  = "tags"
    value = ["Environment:Production", "Backup:Required"]
  }

  backup_on_save = true
}

# Weekly backup with multiple scopes and VCS integration
resource "firefly_backup_and_dr_application" "weekly_multi_scope" {
  account_id     = "your-account-id"
  policy_name    = "Weekly Multi-Scope Backup"
  integration_id = "aws-integration-123"
  region         = "us-west-2"
  provider_type  = "aws"

  schedule {
    frequency    = "Weekly"
    days_of_week = ["Sunday", "Wednesday"]
    hour         = 1
    minute       = 0
  }

  scope {
    type  = "asset_types"
    value = ["aws_instance", "aws_db_instance"]
  }

  scope {
    type  = "resource_group"
    value = ["production-rg"]
  }

  vcs {
    project_id         = "project-456"
    vcs_integration_id = "github-integration-789"
    repo_id            = "backup-repo-123"
  }

  notification_id = "slack-channel-123"
  backup_on_save  = true
}

# Monthly backup with specific weekday schedule
resource "firefly_backup_and_dr_application" "monthly_archive" {
  account_id     = "your-account-id"
  policy_name    = "Monthly Archive Backup"
  integration_id = "aws-integration-123"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency             = "Monthly"
    monthly_schedule_type = "specific_weekday"
    weekday_ordinal       = "First"
    weekday_name          = "Sunday"
    hour                  = 3
    minute                = 0
  }

  scope {
    type = "selected_resources"
    value = [
      "arn:aws:s3:::important-bucket",
      "arn:aws:rds:us-east-1:123456789012:db:prod-db"
    ]
  }

  restore_instructions = <<-EOT
    To restore from this backup:
    1. Access AWS Console
    2. Navigate to backup vault in us-east-1
    3. Select snapshot and choose restore option
  EOT
}

# Monthly backup with specific day of month
resource "firefly_backup_and_dr_application" "monthly_first_day" {
  account_id     = "your-account-id"
  policy_name    = "Monthly First Day Backup"
  integration_id = "aws-integration-123"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency             = "Monthly"
    monthly_schedule_type = "specific_day"
    day_of_month          = 1
    hour                  = 0
    minute                = 0
  }

  scope {
    type  = "tags"
    value = ["CriticalData:True"]
  }
}
