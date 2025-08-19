# firefly_workflows_variable_set (Data Source)

Fetches a single Firefly variable set by ID.

## Example Usage

```terraform
data "firefly_workflows_variable_set" "aws_config" {
  id = "existing-variable-set-id"
}

# Use the variable set in a workspace
resource "firefly_workflows_runners_workspace" "app" {
  name = "production-app"
  consumed_variable_sets = [data.firefly_workflows_variable_set.aws_config.id]
  # ... other configuration
}

# Output variable set information
output "variable_set_info" {
  value = {
    name        = data.firefly_workflows_variable_set.aws_config.name
    description = data.firefly_workflows_variable_set.aws_config.description
    version     = data.firefly_workflows_variable_set.aws_config.version
    labels      = data.firefly_workflows_variable_set.aws_config.labels
  }
}
```

## Schema

### Required

- `id` (String) - The unique identifier of the variable set

### Read-Only

- `name` (String) - The name of the variable set
- `description` (String) - The description of the variable set
- `labels` (List of String) - Labels assigned to the variable set
- `parents` (List of String) - Parent variable set IDs
- `version` (Number) - Version number of the variable set