# firefly_governance_policy (Resource)

Manages a Firefly governance policy for custom infrastructure compliance rules using Rego code.

## Example Usage

```terraform
# Basic Rego policy for AWS CloudWatch Events
resource "firefly_governance_policy" "cloudwatch_rule" {
  name        = "CloudWatch Event Rule Policy"
  description = "Ensures CloudWatch event rules have proper configuration"
  
  code = <<-EOT
firefly {
    input.resource_type == "aws_cloudwatch_event_rule"
    input.configuration.state == "ENABLED"
    input.configuration.name != ""
}
  EOT
  
  type         = ["aws_cloudwatch_event_rule"]
  provider_ids = ["aws_all"]
  severity     = "low"
  category     = "Misconfiguration"
  labels       = ["aws", "cloudwatch", "monitoring"]
  frameworks   = ["SOC2"]
}

# S3 bucket security policy
resource "firefly_governance_policy" "s3_encryption" {
  name        = "S3 Bucket Encryption Policy"
  description = "Enforces encryption on S3 buckets"
  
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

# Multi-resource policy for specific AWS account
resource "firefly_governance_policy" "production_tagging" {
  name        = "Production Resource Tagging"
  description = "Ensures production resources have required tags"
  
  code = <<-EOT
firefly {
    input.resource_type in ["aws_instance", "aws_db_instance", "aws_s3_bucket"]
    input.configuration.tags.Environment
    input.configuration.tags.Owner
    input.configuration.tags.CostCenter
}
  EOT
  
  type         = ["aws_instance", "aws_db_instance", "aws_s3_bucket"]
  provider_ids = ["123456789012"]  # Specific AWS account ID
  severity     = "medium"
  category     = "Governance"
  labels       = ["aws", "tagging", "governance"]
  frameworks   = ["SOC2"]
}

# Base64 encoded policy (alternative format)
resource "firefly_governance_policy" "encoded_policy" {
  name        = "Base64 Encoded Policy"
  description = "Example using base64 encoded Rego code - provider auto-detects format"
  
  # This is base64 encoded Rego code - the provider automatically detects and handles it
  # You can use either plain text or base64 encoded code
  code = "CgpmaXJlZmx5IHsKICAgIGlucHV0Lmluc3RhbmNlX3N0YXRlID09ICJzdG9wcGVkIgp9Cgo="
  
  type         = ["aws_instance"]
  provider_ids = ["aws_all"]
  severity     = "critical"
  category     = "Security"
}
```

## Schema

### Required

- `name` (String) - The name of the governance policy
- `code` (String) - The Rego code for the policy rule. Can be provided as plain text or base64 encoded - the provider automatically detects and handles both formats.
- `type` (List of String) - List of resource types this policy applies to (e.g., `aws_cloudwatch_event_target`, `aws_s3_bucket`)
- `provider_ids` (List of String) - List of provider IDs this policy applies to (e.g., `aws_all`, specific account IDs like `123456789012`)

### Optional

- `description` (String) - The description of the governance policy. Defaults to empty string.
- `labels` (List of String) - List of labels for categorizing the policy. Defaults to empty list.
- `severity` (String) - The severity level of the policy. Valid values: `trace`, `info`, `low`, `medium`, `high`, `critical`. Defaults to `low`.
- `category` (String) - The category of the policy (e.g., `Misconfiguration`, `Security`, `Governance`). Defaults to empty string.
- `frameworks` (List of String) - List of compliance frameworks this policy relates to (e.g., `SOC2`, `ISO27001`, `PCI-DSS`). Defaults to empty list.

### Read-Only

- `id` (String) - The unique identifier of the governance policy

## Rego Policy Guidelines

### Policy Structure
Your Rego policies should follow this simple structure:

```rego
firefly {
    input.resource_type == "aws_s3_bucket"
    # Your conditions here
}
```

**Important**: Do not include `package` declarations or `import` statements. The Firefly API expects simple rule definitions.

### Available Input Data
The `input` object contains:
- `resource_type` (String) - The Terraform resource type
- `configuration` (Object) - The resource configuration/attributes
- `action` (String) - The action being performed (create, update, delete)

### Provider IDs
- Use `aws_all` to apply to all AWS accounts
- Use specific account IDs (e.g., `123456789012`) for account-specific policies
- Multiple provider IDs can be specified in the list

### Best Practices
1. Use the simple `firefly { ... }` rule format (no package or import statements)
2. Test your Rego code before deploying policies
3. Use appropriate severity levels:
   - `trace`: Detailed diagnostic information
   - `info`: Informational messages
   - `low`: Minor issues that should be noted
   - `medium`: Issues that require attention
   - `high`: Serious issues that should be addressed promptly
   - `critical`: Critical issues that must be resolved immediately

## Import

Governance policies can be imported using their ID:

```shell
terraform import firefly_governance_policy.example policy-id-here
```