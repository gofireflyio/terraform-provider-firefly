# firefly_project (Resource)

Manages a Firefly project with hierarchical organization, variables, and metadata.

## Example Usage

```terraform
resource "firefly_project" "example" {
  name        = "Production Infrastructure"
  description = "Main production project"
  labels      = ["production", "critical"]
  
  # Optional parent project for hierarchy
  parent_id = firefly_project.parent.id
  
  # Cron pattern for scheduled executions
  cron_execution_pattern = "0 2 * * 1"
  
  # Project variables
  variables {
    key         = "ENVIRONMENT"
    value       = "production"
    sensitivity = "string"
    destination = "env"
  }
  
  variables {
    key         = "AWS_REGION" 
    value       = "us-west-2"
    sensitivity = "string"
    destination = "iac"
  }
  
  variables {
    key         = "DB_PASSWORD"
    value       = var.db_password
    sensitivity = "secret"
    destination = "env"
  }
}
```

## Schema

### Required

- `name` (String) - The name of the project

### Optional

- `description` (String) - The description of the project
- `labels` (List of String) - Labels to assign to the project
- `cron_execution_pattern` (String) - Cron pattern for scheduled executions
- `parent_id` (String) - ID of the parent project (for nested projects)
- `variables` (Block List) - Variables associated with the project

### Read-Only

- `id` (String) - The unique identifier of the project
- `account_id` (String) - ID of the account the project belongs to
- `members_count` (Number) - Number of members assigned to the project
- `workspace_count` (Number) - Number of workspaces in the project

### Nested Schema for `variables`

#### Required

- `key` (String) - The variable key
- `value` (String, Sensitive) - The variable value

#### Optional

- `sensitivity` (String) - The sensitivity of the variable. Valid values: `string`, `secret`. Defaults to `string`
- `destination` (String) - The destination of the variable. Valid values: `env`, `iac`. Defaults to `env`

## Import

Projects can be imported using their ID:

```shell
terraform import firefly_project.example project-id-here
```