# firefly_workflows_guardrail (Resource)

Manages a Firefly guardrail for governance and compliance of infrastructure changes.

## Example Usage

```terraform
# Cost threshold guardrail
resource "firefly_workflows_guardrail" "cost_limit" {
  name       = "Monthly Cost Threshold"
  type       = "cost"
  is_enabled = true
  severity   = 2
  
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
  severity   = 1
  
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
  severity   = 1
  
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
```

## Schema

### Required

- `name` (String) - The name of the guardrail
- `type` (String) - Type of guardrail. Valid values: `cost`, `policy`, `resource`
- `is_enabled` (Boolean) - Whether the guardrail is enabled
- `severity` (Number) - Severity level (1-3, where 1 is highest)
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

#### Optional

- `cost` (Block) - Cost criteria (see [below for nested schema](#nestedblock--criteria--cost))
- `policy` (Block) - Policy criteria (see [below for nested schema](#nestedblock--criteria--policy))
- `resource` (Block) - Resource criteria (see [below for nested schema](#nestedblock--criteria--resource))

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

## Import

Guardrails can be imported using their ID:

```shell
terraform import firefly_workflows_guardrail.example guardrail-id-here
```