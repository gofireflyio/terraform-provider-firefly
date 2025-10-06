# firefly_governance_policies (Data Source)

Retrieves a list of governance policies that match the specified criteria.

## Example Usage

```terraform
# Get all governance policies
data "firefly_governance_policies" "all" {}

# Get policies by category
data "firefly_governance_policies" "security_policies" {
  category = "Security"
}

# Get policies with specific labels
data "firefly_governance_policies" "aws_policies" {
  labels = ["aws", "security"]
}

# Search for policies by name or description
data "firefly_governance_policies" "s3_policies" {
  query = "s3"
}

# Combined search criteria
data "firefly_governance_policies" "strict_aws_policies" {
  category = "Security"
  labels   = ["aws"]
  query    = "encryption"
}

# Use the results to reference specific policies
output "policy_names" {
  value = [for policy in data.firefly_governance_policies.all.policies : policy.name]
}

output "security_policy_ids" {
  value = [for policy in data.firefly_governance_policies.security_policies.policies : policy.id]
}
```

## Schema

### Optional

- `query` (String) - Search query to filter policies by name or description
- `labels` (List of String) - Filter policies that have all of the specified labels
- `category` (String) - Filter policies by category (e.g., `Security`, `Governance`, `Misconfiguration`)

### Read-Only

- `id` (String) - The data source identifier
- `policies` (List of Object) - List of governance policies matching the criteria (see [below for nested schema](#nestedatt--policies))

<a id="nestedatt--policies"></a>
### Nested Schema for `policies`

#### Read-Only

- `id` (String) - The unique identifier of the governance policy
- `name` (String) - The name of the governance policy
- `description` (String) - The description of the governance policy
- `code` (String) - The Rego code for the policy rule
- `type` (List of String) - List of resource types this policy applies to
- `provider_ids` (List of String) - List of provider IDs this policy applies to
- `labels` (List of String) - List of labels associated with the policy
- `severity` (String) - The severity level of the policy (`trace`, `info`, `low`, `medium`, `high`, `critical`)
- `category` (String) - The category of the policy
- `frameworks` (List of String) - List of compliance frameworks this policy relates to

## Usage Examples

### Finding Policies for Specific Resource Types
```terraform
# Get all policies that apply to S3 buckets
data "firefly_governance_policies" "s3_policies" {
  query = "s3_bucket"
}

locals {
  s3_policies = [
    for policy in data.firefly_governance_policies.s3_policies.policies :
    policy if contains(policy.type, "aws_s3_bucket")
  ]
}
```

### Filtering by Compliance Framework
```terraform
# Get all SOC2 compliance policies
data "firefly_governance_policies" "soc2_policies" {}

locals {
  soc2_policies = [
    for policy in data.firefly_governance_policies.soc2_policies.policies :
    policy if contains(policy.frameworks, "SOC2")
  ]
}
```

### Creating Reports
```terraform
# Generate a report of all governance policies
data "firefly_governance_policies" "all" {}

output "governance_policy_report" {
  value = {
    total_policies = length(data.firefly_governance_policies.all.policies)
    by_category = {
      for category in distinct([for p in data.firefly_governance_policies.all.policies : p.category]) :
      category => length([for p in data.firefly_governance_policies.all.policies : p if p.category == category])
    }
    by_severity = {
      for severity in distinct([for p in data.firefly_governance_policies.all.policies : p.severity]) :
      severity => length([for p in data.firefly_governance_policies.all.policies : p if p.severity == severity])
    }
  }
}