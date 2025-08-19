# firefly_workflows_project (Resource)

Manages a Firefly project for organizing and grouping infrastructure resources.

## Example Usage

```terraform
# Basic project
resource "firefly_workflows_project" "example" {
  name        = "Production Infrastructure"
  description = "Main production project"
  labels      = ["production", "critical"]
}

# Project with variables and scheduled execution
resource "firefly_workflows_project" "scheduled" {
  name        = "Scheduled Project"
  description = "Project with automatic scheduled runs"
  labels      = ["automation", "scheduled"]
  
  # Run daily at 2 AM
  cron_execution_pattern = "0 2 * * *"
  
  variables {
    key         = "ENVIRONMENT"
    value       = "production"
    sensitivity = "string"
    destination = "env"
  }
}

# Child project
resource "firefly_workflows_project" "child" {
  name        = "Child Project"
  description = "Child project inheriting from parent"
  parent_id   = firefly_workflows_project.example.id
  labels      = ["child", "development"]
}
```

## Schema

### Required

- `name` (String) - The name of the project

### Optional

- `description` (String) - The description of the project
- `labels` (List of String) - Labels to assign to the project
- `cron_execution_pattern` (String) - Cron pattern for scheduled executions
- `parent_id` (String) - ID of the parent project for hierarchical organization
- `variables` (Block Set) - Variables to define for the project (see [below for nested schema](#nestedblock--variables))

### Read-Only

- `id` (String) - The unique identifier of the project
- `account_id` (String) - ID of the account the project belongs to
- `members_count` (Number) - Number of members assigned to the project
- `workspace_count` (Number) - Number of workspaces in the project

<a id="nestedblock--variables"></a>
### Nested Schema for `variables`

#### Required

- `key` (String) - The variable key/name
- `value` (String) - The variable value
- `sensitivity` (String) - Variable sensitivity level. Valid values: `string`, `secret`
- `destination` (String) - Where the variable is used. Valid values: `env`, `iac`

## Import

Projects can be imported using their ID:

```shell
terraform import firefly_workflows_project.example project-id-here
```