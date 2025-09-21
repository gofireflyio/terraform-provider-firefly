# Example governance policy for S3 bucket encryption
resource "firefly_governance_policy" "s3_encryption" {
  name        = "S3 Bucket Encryption Policy"
  description = "Enforces that all S3 buckets have server-side encryption enabled"

  # Rego policy code
  code = <<-EOT
firefly {
    input.resource_type == "aws_s3_bucket"
    input.configuration.server_side_encryption_configuration
}
  EOT

  type         = ["aws_s3_bucket"]
  provider_ids = ["aws_all"]
  severity     = "high"
  category     = "Security"
  labels       = ["aws", "s3", "encryption", "security"]
  frameworks   = ["SOC2", "ISO27001", "PCI-DSS"]
}

# Example governance policy for CloudWatch Event Rules
resource "firefly_governance_policy" "cloudwatch_events" {
  name        = "CloudWatch Event Rules Policy"
  description = "Ensures CloudWatch event rules are properly configured and enabled"

  code = <<-EOT
firefly {
    input.resource_type == "aws_cloudwatch_event_rule"
    input.configuration.name != ""
    input.configuration.state == "ENABLED"
}
  EOT

  type         = ["aws_cloudwatch_event_rule"]
  provider_ids = ["aws_all"]
  severity     = "low"
  category     = "Misconfiguration"
  labels       = ["aws", "cloudwatch", "events", "monitoring"]
  frameworks   = ["SOC2"]
}

# Example governance policy for required resource tagging
resource "firefly_governance_policy" "required_tags" {
  name        = "Required Resource Tags"
  description = "Ensures EC2 instances and RDS instances have required tags"

  code = <<-EOT
firefly {
    input.resource_type in ["aws_instance", "aws_db_instance"]
    input.configuration.tags.Environment
    input.configuration.tags.Owner
    input.configuration.tags.CostCenter
    input.configuration.tags.Project
}
  EOT

  type         = ["aws_instance", "aws_db_instance"]
  provider_ids = ["123456789012"] # Replace with your AWS account ID
  severity     = "medium"
  category     = "Governance"
  labels       = ["aws", "tagging", "governance", "ec2", "rds"]
  frameworks   = ["SOC2"]
}

# Example with base64 encoded Rego code (the provider automatically detects and handles both formats)
resource "firefly_governance_policy" "base64_example" {
  name        = "Base64 Encoded Policy Example"
  description = "Example showing base64 encoded Rego code support"

  # This is the same Rego code as above, but base64 encoded
  # The provider automatically detects this is base64 and uses it as-is
  code = "CgpmaXJlZmx5IHsKICAgIGlucHV0Lmluc3RhbmNlX3N0YXRlID09ICJzdG9wcGVkIgp9Cgo="

  type         = ["aws_instance"]
  provider_ids = ["aws_all"]
  severity     = "critical"
  category     = "Security"
  labels       = ["base64", "encoding", "example"]
  frameworks   = ["ISO27001"]
}

# Example showing all available severity levels
resource "firefly_governance_policy" "severity_example" {
  name        = "Severity Levels Example"
  description = "Example showing all available severity levels: trace, info, low, medium, high, critical"

  code = <<-EOT
    firefly {
      input.instance_state == "stopped"
    }
  EOT

  type         = ["aws_instance"]
  provider_ids = ["aws_all"]
  severity     = "trace"  # Available: trace, info, low, medium, high, critical
  category     = "Example"
  labels       = ["severity", "example"]
}

# Output policy IDs for reference
output "s3_encryption_policy_id" {
  description = "ID of the S3 encryption policy"
  value       = firefly_governance_policy.s3_encryption.id
}

output "cloudwatch_events_policy_id" {
  description = "ID of the CloudWatch events policy"
  value       = firefly_governance_policy.cloudwatch_events.id
}

output "required_tags_policy_id" {
  description = "ID of the required tags policy"
  value       = firefly_governance_policy.required_tags.id
}

output "base64_example_policy_id" {
  description = "ID of the base64 example policy"
  value       = firefly_governance_policy.base64_example.id
}

output "severity_example_policy_id" {
  description = "ID of the severity example policy"
  value       = firefly_governance_policy.severity_example.id
}