# firefly_workflows_variable_set (Resource)

Manages a Firefly variable set for sharing configuration across multiple workspaces.

## Example Usage

```terraform
# Basic variable set
resource "firefly_workflows_variable_set" "aws_config" {
  name        = "AWS Configuration"
  description = "Shared AWS configuration variables"
  labels      = ["aws", "shared"]
  
  variables {
    key         = "AWS_DEFAULT_REGION"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "env"
  }
  
  variables {
    key         = "AWS_ACCESS_KEY_ID"
    value       = var.aws_access_key
    sensitivity = "secret"
    destination = "env"
  }
}

# Variable set with inheritance
resource "firefly_workflows_variable_set" "production_config" {
  name        = "Production Configuration"
  description = "Production-specific configuration"
  labels      = ["production", "config"]
  parents     = [firefly_workflows_variable_set.aws_config.id]
  
  variables {
    key         = "ENVIRONMENT"
    value       = "production"
    sensitivity = "string"
    destination = "env"
  }
}
```

## Schema

### Required

- `name` (String) - The name of the variable set

### Optional

- `description` (String) - The description of the variable set
- `labels` (List of String) - Labels to assign to the variable set
- `parents` (List of String) - List of parent variable set IDs for inheritance
- `variables` (Block Set) - Variables to define in the set (see [below for nested schema](#nestedblock--variables))

### Read-Only

- `id` (String) - The unique identifier of the variable set
- `version` (Number) - Version number of the variable set

<a id="nestedblock--variables"></a>
### Nested Schema for `variables`

#### Required

- `key` (String) - The variable key/name
- `value` (String) - The variable value
- `sensitivity` (String) - Variable sensitivity level. Valid values: `string`, `secret`
- `destination` (String) - Where the variable is used. Valid values: `env`, `iac`

## Import

Variable sets can be imported using their ID:

```shell
terraform import firefly_workflows_variable_set.example variable-set-id-here
```