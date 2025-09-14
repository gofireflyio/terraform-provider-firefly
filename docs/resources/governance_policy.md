# firefly_governance_policy (Resource)

Manages a Firefly governance policy for custom infrastructure compliance rules using Rego code.

## Example Usage

```terraform
# Basic Rego policy for AWS CloudWatch Events
resource "firefly_governance_policy" "cloudwatch_rule" {
  name        = "CloudWatch Event Rule Policy"
  description = "Ensures CloudWatch event rules have proper configuration"
  
  code = <<-EOT
    package firefly
    
    import rego.v1
    
    default allow := false
    
    allow if {
        input.resource_type == "aws_cloudwatch_event_rule"
        input.configuration.state == "ENABLED"
        input.configuration.name != ""
    }
    
    deny[msg] if {
        input.resource_type == "aws_cloudwatch_event_rule"
        input.configuration.state != "ENABLED"
        msg := "CloudWatch event rule must be enabled"
    }
  EOT
  
  type         = ["aws_cloudwatch_event_rule"]
  provider_ids = ["aws_all"]
  severity     = "warning"
  category     = "Misconfiguration"
  labels       = ["aws", "cloudwatch", "monitoring"]
  frameworks   = ["SOC2"]
}

# S3 bucket security policy
resource "firefly_governance_policy" "s3_encryption" {
  name        = "S3 Bucket Encryption Policy"
  description = "Enforces encryption on S3 buckets"
  
  code = <<-EOT
    package firefly
    
    import rego.v1
    
    default allow := false
    
    allow if {
        input.resource_type == "aws_s3_bucket"
        input.configuration.server_side_encryption_configuration
    }
    
    deny[msg] if {
        input.resource_type == "aws_s3_bucket"
        not input.configuration.server_side_encryption_configuration
        msg := "S3 bucket must have server-side encryption enabled"
    }
  EOT
  
  type         = ["aws_s3_bucket"]
  provider_ids = ["aws_all"]
  severity     = "strict"
  category     = "Security"
  labels       = ["aws", "s3", "encryption", "security"]
  frameworks   = ["SOC2", "ISO27001", "PCI-DSS"]
}

# Multi-resource policy for specific AWS account
resource "firefly_governance_policy" "production_tagging" {
  name        = "Production Resource Tagging"
  description = "Ensures production resources have required tags"
  
  code = <<-EOT
    package firefly
    
    import rego.v1
    
    required_tags := ["Environment", "Owner", "CostCenter"]
    
    default allow := false
    
    allow if {
        input.resource_type in ["aws_instance", "aws_db_instance", "aws_s3_bucket"]
        tags := object.get(input.configuration, "tags", {})
        every tag in required_tags {
            tags[tag]
        }
    }
    
    deny[msg] if {
        input.resource_type in ["aws_instance", "aws_db_instance", "aws_s3_bucket"]
        tags := object.get(input.configuration, "tags", {})
        some tag in required_tags
        not tags[tag]
        msg := sprintf("Resource missing required tag: %s", [tag])
    }
  EOT
  
  type         = ["aws_instance", "aws_db_instance", "aws_s3_bucket"]
  provider_ids = ["123456789012"]  # Specific AWS account ID
  severity     = "flexible"
  category     = "Governance"
  labels       = ["aws", "tagging", "governance"]
  frameworks   = ["SOC2"]
}

# Base64 encoded policy (alternative format)
resource "firefly_governance_policy" "encoded_policy" {
  name        = "Base64 Encoded Policy"
  description = "Example using base64 encoded Rego code"
  
  # This is the same Rego code as above, but base64 encoded
  code = base64encode(<<-EOT
    package firefly
    
    import rego.v1
    
    default allow := false
    
    allow if {
        input.resource_type == "aws_instance"
        input.configuration.instance_type
    }
  EOT
  )
  
  type         = ["aws_instance"]
  provider_ids = ["aws_all"]
  severity     = "warning"
  category     = "Configuration"
}
```

## Schema

### Required

- `name` (String) - The name of the governance policy
- `code` (String) - The Rego code for the policy rule (can be base64 encoded)
- `type` (List of String) - List of resource types this policy applies to (e.g., `aws_cloudwatch_event_target`, `aws_s3_bucket`)
- `provider_ids` (List of String) - List of provider IDs this policy applies to (e.g., `aws_all`, specific account IDs like `123456789012`)

### Optional

- `description` (String) - The description of the governance policy. Defaults to empty string.
- `labels` (List of String) - List of labels for categorizing the policy. Defaults to empty list.
- `severity` (String) - The severity level of the policy. Valid values: `flexible`, `strict`, `warning`. Defaults to `warning`.
- `category` (String) - The category of the policy (e.g., `Misconfiguration`, `Security`, `Governance`). Defaults to empty string.
- `frameworks` (List of String) - List of compliance frameworks this policy relates to (e.g., `SOC2`, `ISO27001`, `PCI-DSS`). Defaults to empty list.

### Read-Only

- `id` (String) - The unique identifier of the governance policy

## Rego Policy Guidelines

### Policy Structure
Your Rego policies should follow this structure:

```rego
package firefly

import rego.v1

# Default decision
default allow := false

# Allow conditions
allow if {
    input.resource_type == "aws_s3_bucket"
    # Your allow conditions here
}

# Deny conditions with messages
deny[msg] if {
    input.resource_type == "aws_s3_bucket"
    # Your deny conditions here
    msg := "Your violation message here"
}
```

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
1. Always include `import rego.v1` for modern Rego syntax
2. Use meaningful violation messages in `deny` rules
3. Test your Rego code before deploying policies
4. Use appropriate severity levels:
   - `strict`: Blocks deployments
   - `flexible`: Allows override with justification
   - `warning`: Shows warnings but doesn't block

## Import

Governance policies can be imported using their ID:

```shell
terraform import firefly_governance_policy.example policy-id-here
```