# Get all backup & DR applications
data "firefly_backup_and_dr_applications" "all" {}

# Get only AWS backup policies
data "firefly_backup_and_dr_applications" "aws_policies" {
  provider_type = "aws"
}

# Get only active backup policies
data "firefly_backup_and_dr_applications" "active_policies" {
  status = "Active"
}

# Get active AWS backup policies
data "firefly_backup_and_dr_applications" "active_aws" {
  provider_type = "aws"
  status        = "Active"
}

# Output: list all backup policy names
output "all_backup_policy_names" {
  description = "Names of all backup policies"
  value       = [for app in data.firefly_backup_and_dr_applications.all.applications : app.policy_name]
}

# Output: count of active vs inactive policies
output "backup_policy_summary" {
  description = "Summary of backup policies"
  value = {
    total_policies      = length(data.firefly_backup_and_dr_applications.all.applications)
    active_aws_policies = length(data.firefly_backup_and_dr_applications.active_aws.applications)
  }
}

# Output: details of each AWS policy
output "aws_backup_details" {
  description = "Details of AWS backup policies"
  value = [
    for app in data.firefly_backup_and_dr_applications.aws_policies.applications : {
      id              = app.id
      name            = app.policy_name
      region          = app.region
      status          = app.status
      snapshots_count = app.snapshots_count
    }
  ]
}
