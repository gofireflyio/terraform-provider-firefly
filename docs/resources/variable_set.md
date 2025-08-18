# firefly_variable_set (Resource)

Manages a Firefly variable set for reusable configuration with inheritance support.

## Example Usage

```terraform
resource "firefly_variable_set" "example" {
  name        = "AWS Configuration"
  description = "Shared AWS configuration variables"
  labels      = ["aws", "shared", "production"]
  
  # Parent variable sets for inheritance
  parents = [firefly_variable_set.base_config.id]
  
  # Variables in the set
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
  
  variables {
    key         = "TF_VAR_aws_region"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "iac"
  }
}

# Use the variable set in a workspace
resource "firefly_runners_workspace" "app" {
  name = "production-app"
  consumed_variable_sets = [firefly_variable_set.example.id]
  # ... other configuration
}
```

## Schema

### Required

- `name` (String) - The name of the variable set

### Optional

- `description` (String) - The description of the variable set
- `labels` (List of String) - Labels to assign to the variable set
- `parents` (List of String) - Parent variable set IDs for inheritance
- `variables` (Block List) - Variables in the variable set

### Read-Only

- `id` (String) - The unique identifier of the variable set
- `version` (Number) - Version number of the variable set

### Nested Schema for `variables`

#### Required

- `key` (String) - The variable key
- `value` (String, Sensitive) - The variable value

#### Optional

- `sensitivity` (String) - The sensitivity of the variable. Valid values: `string`, `secret`. Defaults to `string`
- `destination` (String) - The destination of the variable. Valid values: `env`, `iac`. Defaults to `env`

## Import

Variable sets can be imported using their ID:

```shell
terraform import firefly_variable_set.example variable-set-id-here
```