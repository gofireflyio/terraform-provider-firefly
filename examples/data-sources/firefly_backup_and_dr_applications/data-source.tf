# List all backup policies for an account
data "firefly_backup_and_dr_applications" "all" {
  account_id = "your-account-id"
}

# Filter by status - get only active policies
data "firefly_backup_and_dr_applications" "active" {
  account_id = "your-account-id"
  status     = "Active"
}

# Filter by region and provider
data "firefly_backup_and_dr_applications" "aws_east" {
  account_id    = "your-account-id"
  region        = "us-east-1"
  provider_type = "aws"
}

# Filter by integration ID
data "firefly_backup_and_dr_applications" "specific_integration" {
  account_id     = "your-account-id"
  integration_id = "aws-integration-123"
}

# Use in outputs to display policy names
output "active_application_names" {
  value       = [for p in data.firefly_backup_and_dr_applications.active.applications : p.application_name]
  description = "List of all active backup policy names"
}

# Use in outputs to display policy details
output "aws_east_policies" {
  value = [
    for p in data.firefly_backup_and_dr_applications.aws_east.applications : {
      name              = p.application_name
      status            = p.status
      schedule          = p.schedule_frequency
      snapshots_count   = p.snapshots_count
      last_backup_time  = p.last_backup_time
      last_backup_status = p.last_backup_status
    }
  ]
  description = "Details of backup policies in us-east-1"
}

# Find policies by status and use them in other resources
locals {
  active_application_ids = [for p in data.firefly_backup_and_dr_applications.active.applications : p.application_id]
}
