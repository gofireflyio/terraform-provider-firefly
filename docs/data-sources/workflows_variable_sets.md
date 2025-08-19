# firefly_workflows_variable_sets (Data Source)

Fetches a list of Firefly variable sets with optional filtering.

## Example Usage

```terraform
# Get all variable sets
data "firefly_workflows_variable_sets" "all" {}

# Search for specific variable sets
data "firefly_workflows_variable_sets" "aws_configs" {
  search_query = "aws"
}

# Use variable set data
resource "firefly_workflows_runners_workspace" "app" {
  name = "production-app"
  consumed_variable_sets = [for vs in data.firefly_workflows_variable_sets.aws_configs.variable_sets : vs.id]
  # ... other configuration
}

# Output variable set names
output "variable_set_names" {
  value = [for vs in data.firefly_workflows_variable_sets.all.variable_sets : vs.name]
}
```

## Schema

### Optional

- `search_query` (String) - Search query to filter variable sets by name or description

### Read-Only

- `variable_sets` (List of Object) - List of variable sets (see [below for nested schema](#nestedatt--variable_sets))

<a id="nestedatt--variable_sets"></a>
### Nested Schema for `variable_sets`

Read-Only:

- `id` (String) - The unique identifier of the variable set
- `name` (String) - The name of the variable set
- `description` (String) - The description of the variable set
- `labels` (List of String) - Labels assigned to the variable set
- `parents` (List of String) - Parent variable set IDs
- `version` (Number) - Version number of the variable set