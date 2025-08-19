# firefly_workflows_project (Data Source)

Fetches a single Firefly project by ID.

## Example Usage

```terraform
data "firefly_workflows_project" "main" {
  id = "existing-project-id"
}

# Use the project data
resource "firefly_workflows_runners_workspace" "app" {
  name       = "new-app"
  project_id = data.firefly_workflows_project.main.id
  # ... other configuration
}

# Output project information
output "project_info" {
  value = {
    name            = data.firefly_workflows_project.main.name
    description     = data.firefly_workflows_project.main.description
    workspace_count = data.firefly_workflows_project.main.workspace_count
  }
}
```

## Schema

### Required

- `id` (String) - The unique identifier of the project

### Read-Only

- `name` (String) - The name of the project
- `description` (String) - The description of the project
- `labels` (List of String) - Labels assigned to the project
- `cron_execution_pattern` (String) - Cron pattern for scheduled executions
- `parent_id` (String) - ID of the parent project
- `account_id` (String) - ID of the account the project belongs to
- `members_count` (Number) - Number of members assigned to the project
- `workspace_count` (Number) - Number of workspaces in the project