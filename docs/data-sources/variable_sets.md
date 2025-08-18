# firefly_variable_sets (Data Source)

Fetches a list of Firefly variable sets with optional filtering.

## Example Usage

```terraform
# Get all variable sets
data "firefly_variable_sets" "all" {}

# Filter variable sets by search query
data "firefly_variable_sets" "aws" {
  search_query = "aws"
}

# Use variable set in workspace
resource "firefly_runners_workspace" "app" {
  name = "production-app"
  consumed_variable_sets = [
    data.firefly_variable_sets.aws.variable_sets[0].id
  ]
  # ... other configuration
}

# Output variable set information
output "aws_variable_sets" {
  value = [
    for vs in data.firefly_variable_sets.aws.variable_sets : {
      name        = vs.name
      id          = vs.id
      description = vs.description
    }
  ]
}
```

## Schema

### Optional

- `search_query` (String) - Optional search query to filter variable sets

### Read-Only

- `variable_sets` (List of Object) - List of variable sets

### Nested Schema for `variable_sets`

- `id` (String) - The unique identifier of the variable set
- `name` (String) - The name of the variable set
- `description` (String) - The description of the variable set
- `labels` (List of String) - Labels assigned to the variable set
- `parents` (List of String) - Parent variable set IDs
- `version` (Number) - Version number of the variable set