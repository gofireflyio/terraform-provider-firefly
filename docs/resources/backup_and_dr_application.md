# firefly_backup_and_dr_application (Resource)

Manages a Firefly Backup & DR application for automated infrastructure backup and disaster recovery.

## Example Usage

```terraform
# Basic backup with tag-based scope, every 24 hours
resource "firefly_backup_and_dr_application" "daily" {
  account_id       = "66169d5af4992fc0bab04510"
  application_name = "Daily Production Backup"
  integration_id   = "692ec8acce65b3dc46cfceb5"
  region           = "us-east-1"
  provider_type    = "aws"
  description      = "Daily backup of production resources"
  frequency        = 24

  scope {
    type  = "tags"
    value = ["Environment:Production", "Backup:Required"]
  }
}

# 8-hour backup with asset type filtering
resource "firefly_backup_and_dr_application" "frequent" {
  account_id       = "66169d5af4992fc0bab04510"
  application_name = "Frequent Database Backup"
  integration_id   = "692ec8acce65b3dc46cfceb5"
  region           = "eu-west-1"
  provider_type    = "aws"
  frequency        = 8

  scope {
    type  = "asset_types"
    value = ["aws_db_instance", "aws_rds_cluster"]
  }
}

# Backup with disaster recovery (resilience) enabled
resource "firefly_backup_and_dr_application" "with_dr" {
  account_id         = "66169d5af4992fc0bab04510"
  application_name   = "DR-Enabled Backup"
  integration_id     = "692ec8acce65b3dc46cfceb5"
  region             = "us-east-1"
  provider_type      = "aws"
  frequency          = 4
  target_account     = "target-integration-id"
  target_region      = "us-west-2"
  auto_create_pr     = true
  resilience_enabled = true

  scope {
    type  = "tags"
    value = ["Environment:Production"]
  }
}

# Backup with VCS integration for artifact storage
resource "firefly_backup_and_dr_application" "with_vcs" {
  account_id       = "66169d5af4992fc0bab04510"
  application_name = "Backup to Git Repository"
  integration_id   = "692ec8acce65b3dc46cfceb5"
  region           = "us-east-1"
  provider_type    = "aws"
  frequency        = 24

  scope {
    type  = "tags"
    value = ["IaC:Terraform"]
  }

  vcs {
    vcs_integration_id = "vcs-integration-456"
    repo_id            = "repo-789"
  }
}

# Backup with multiple scope types (AND logic)
resource "firefly_backup_and_dr_application" "multi_scope" {
  account_id       = "66169d5af4992fc0bab04510"
  application_name = "Multi-Scope Backup"
  integration_id   = "692ec8acce65b3dc46cfceb5"
  region           = "us-east-1"
  provider_type    = "aws"
  frequency        = 16

  scope {
    type  = "tags"
    value = ["Environment:Production"]
  }

  scope {
    type  = "asset_types"
    value = ["aws_instance", "aws_ebs_volume"]
  }
}

# Backup of specific resources by ARN
resource "firefly_backup_and_dr_application" "selected" {
  account_id       = "66169d5af4992fc0bab04510"
  application_name = "Specific Resources Backup"
  integration_id   = "692ec8acce65b3dc46cfceb5"
  region           = "us-east-1"
  provider_type    = "aws"
  frequency        = 24

  scope {
    type  = "selected_resources"
    value = [
      "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
      "arn:aws:s3:::my-critical-bucket",
    ]
  }
}
```

## Schema

### Required

- `account_id` (String) - The account ID for the backup application. **Note**: Changing this forces replacement of the resource.
- `application_name` (String) - The name of the backup application (max 100 characters)
- `integration_id` (String) - The integration ID for cloud provider credentials
- `region` (String) - The cloud region where backups will be stored
- `provider_type` (String) - The cloud provider type (max 50 characters, e.g., `aws`, `azure`, `gcp`)

### Optional

- `description` (String) - Description of the backup application (max 500 characters)
- `frequency` (Number) - Hours between scheduled backups. Valid values: `4`, `8`, `16`, `24`
- `notification_id` (String) - Notification channel ID for backup alerts
- `restore_instructions` (String) - Instructions for restoring from backups (max 2000 characters)
- `backup_on_save` (Boolean) - Whether to trigger a backup immediately on application creation/update. Defaults to `true`.
- `target_account` (String) - Target account/integration ID where the restore should land (used with `resilience_enabled`)
- `target_region` (String) - Target region where the restore should land (used with `resilience_enabled`)
- `auto_create_pr` (Boolean) - If `true`, the restore flow automatically opens a VCS pull request with the restored IaC
- `resilience_enabled` (Boolean) - When `true`, DR scheduling applies. Requires `target_account`, `target_region`, and `frequency` to be set.
- `scope` (Block List) - Resource scope configurations for backup targeting (see [below for nested schema](#nestedblock--scope))
- `vcs` (Block) - VCS integration configuration for backup artifacts (see [below for nested schema](#nestedblock--vcs))

### Read-Only

- `id` (String) - The unique identifier of the backup application
- `status` (String) - Current status of the application (`Active` or `Inactive`)
- `snapshots_count` (Number) - Number of snapshots created by this application
- `last_backup_snapshot_id` (String) - ID of the most recent backup snapshot
- `last_backup_time` (String) - Timestamp of the last backup
- `last_backup_status` (String) - Status of the last backup
- `next_backup_time` (String) - Timestamp of the next scheduled backup
- `created_at` (String) - Timestamp when the application was created
- `updated_at` (String) - Timestamp when the application was last updated

<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

#### Required

- `type` (String) - Scope type. Valid values: `tags`, `resource_group`, `asset_types`, `excluded_asset_types`, `selected_resources`, `excluded_resources`
- `value` (List of String) - List of values for this scope type. Format depends on type:
  - `tags`: Tag key-value pairs in format `"key:value"` (e.g., `["Environment:Production", "Team:DevOps"]`)
  - `resource_group`: Resource group names (e.g., `["prod-rg", "staging-rg"]`)
  - `asset_types`: Terraform resource types to include (e.g., `["aws_instance", "aws_db_instance"]`)
  - `excluded_asset_types`: Terraform resource types to exclude (e.g., `["aws_s3_bucket"]`)
  - `selected_resources`: Specific resource ARNs or IDs
  - `excluded_resources`: Specific resource ARNs or IDs to exclude from all scope filters

**Note**: Multiple scope blocks can be defined. Resources must match ALL scope criteria (AND logic).

<a id="nestedblock--vcs"></a>
### Nested Schema for `vcs`

#### Optional

- `vcs_integration_id` (String) - VCS integration ID
- `repo_id` (String) - Repository ID for storing backup artifacts

## Important Notes

- **Immediate Backup**: By default, `backup_on_save` is `true`, which triggers an immediate backup when the application is created or updated. Set to `false` to disable this behavior.
- **Account ID Changes**: Changing the `account_id` attribute forces replacement (destroy and recreate) of the resource.
- **Computed Fields**: All timestamp and status fields are read-only and automatically updated by Firefly.
- **Disaster Recovery**: When `resilience_enabled` is `true`, `target_account`, `target_region`, and `frequency` are required.
- **VCS Integration**: When VCS is configured, backup artifacts are stored in the specified git repository in addition to cloud storage.

## Import

Backup & DR applications can be imported using the format `account_id:application_id`:

```shell
terraform import firefly_backup_and_dr_application.example 66169d5af4992fc0bab04510:application-id-here
```
