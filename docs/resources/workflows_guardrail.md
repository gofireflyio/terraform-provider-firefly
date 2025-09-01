# firefly_workflows_guardrail (Resource)

Manages a Firefly guardrail for governance and compliance of infrastructure changes.

> **Important**: Each guardrail can only have **one criteria type** (cost, policy, or resource). Multiple criteria types in a single guardrail are not supported.

## Example Usage

```terraform
# Cost threshold guardrail
resource "firefly_workflows_guardrail" "cost_limit" {
  name       = "Monthly Cost Threshold"
  type       = "cost"
  is_enabled = true
  severity   = "strict"
  
  scope {
    workspaces {
      include = ["production-*"]
      exclude = ["production-test"]
    }
  }
  
  criteria {
    cost {
      threshold_amount = 5000.0
    }
  }
}

# Policy enforcement guardrail
resource "firefly_workflows_guardrail" "security_policy" {
  name       = "Security Policy Enforcement"
  type       = "policy"
  is_enabled = true
  severity   = "flexible"
  
  scope {
    workspaces {
      include = ["*"]
    }
  }
  
  criteria {
    policy {
      severity = "high"
    }
  }
}

# Resource control guardrail
resource "firefly_workflows_guardrail" "resource_control" {
  name       = "Prevent Production Deletions"
  type       = "resource"
  is_enabled = true
  severity   = "flexible"
  
  scope {
    workspaces {
      include = ["prod-*"]
    }
    regions {
      include = ["us-west-2", "us-east-1"]
    }
  }
  
  criteria {
    resource {
      actions           = ["delete"]
      specific_resources = ["aws_instance", "aws_db_instance"]
    }
  }
}

# Invalid example - DO NOT DO THIS
# This will cause an error because multiple criteria types are specified
resource "firefly_workflows_guardrail" "invalid_multiple_criteria" {
  name       = "Invalid Example"
  type       = "cost"
  is_enabled = true
  severity   = "strict"
  
  scope {
    workspaces {
      include = ["*"]
    }
  }
  
  criteria {
    cost {
      threshold_amount = 1000.0
    }
    policy {  # ERROR: Cannot specify both cost AND policy criteria
      severity = "high"
    }
  }
}

# Tag enforcement guardrail
resource "firefly_workflows_guardrail" "tag_policy" {
  name       = "Required Tags Policy"
  type       = "tag"
  is_enabled = true
  severity   = "warning"
  
  scope {
    workspaces {
      include = ["production-*"]
    }
  }
  
  criteria {
    tag {
      tag_enforcement_mode = "requiredTags"
      required_tags       = ["Environment", "Owner", "CostCenter"]
    }
  }
}
```

## Schema

### Required

- `name` (String) - The name of the guardrail
- `type` (String) - Type of guardrail. Valid values: `cost`, `policy`, `resource`, `tag`
- `is_enabled` (Boolean) - Whether the guardrail is enabled
- `severity` (String) - Severity level. Valid values: `flexible`, `strict`, `warning`
- `scope` (Block) - Scope configuration (see [below for nested schema](#nestedblock--scope))
- `criteria` (Block) - Criteria configuration (see [below for nested schema](#nestedblock--criteria))

### Optional

- `description` (String) - The description of the guardrail

### Read-Only

- `id` (String) - The unique identifier of the guardrail

<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

#### Optional

- `workspaces` (Block) - Workspace scope (see [below for nested schema](#nestedblock--scope--workspaces))
- `regions` (Block) - Region scope (see [below for nested schema](#nestedblock--scope--regions))

<a id="nestedblock--scope--workspaces"></a>
### Nested Schema for `scope.workspaces`

#### Optional

- `include` (List of String) - Workspace patterns to include
- `exclude` (List of String) - Workspace patterns to exclude

<a id="nestedblock--scope--regions"></a>
### Nested Schema for `scope.regions`

#### Optional

- `include` (List of String) - Regions to include
- `exclude` (List of String) - Regions to exclude

<a id="nestedblock--criteria"></a>
### Nested Schema for `criteria`

> **Note**: Exactly one of the following criteria types must be specified per guardrail.

#### Optional

- `cost` (Block) - Cost criteria (see [below for nested schema](#nestedblock--criteria--cost))
- `policy` (Block) - Policy criteria (see [below for nested schema](#nestedblock--criteria--policy))
- `resource` (Block) - Resource criteria (see [below for nested schema](#nestedblock--criteria--resource))
- `tag` (Block) - Tag criteria (see [below for nested schema](#nestedblock--criteria--tag))

<a id="nestedblock--criteria--cost"></a>
### Nested Schema for `criteria.cost`

#### Required

- `threshold_amount` (Number) - Cost threshold amount

<a id="nestedblock--criteria--policy"></a>
### Nested Schema for `criteria.policy`

#### Required

- `severity` (String) - Policy severity level. Valid values: `low`, `medium`, `high`

<a id="nestedblock--criteria--resource"></a>
### Nested Schema for `criteria.resource`

#### Required

- `actions` (List of String) - Resource actions to monitor. Valid values: `create`, `update`, `delete`

#### Optional

- `specific_resources` (List of String) - Specific resource types to monitor

<a id="nestedblock--criteria--tag"></a>
### Nested Schema for `criteria.tag`

#### Required

- `tag_enforcement_mode` (String) - Tag enforcement mode. Valid values: `requiredTags`
- `required_tags` (List of String) - List of required tags that resources must have

## Import

Guardrails can be imported using their ID:

```shell
terraform import firefly_workflows_guardrail.example guardrail-id-here
```