# firefly_workflows_project (Data Source)

Fetches a single Firefly project by ID or path (name).

## Example Usage

```terraform
# Fetch project by ID
data "firefly_workflows_project" "by_id" {
  id = "existing-project-id"
}

# Fetch project by path (name)  
data "firefly_workflows_project" "by_path" {
  path = "Production Infrastructure"
}

# Use the project data
resource "firefly_workflows_runners_workspace" "app" {
  name       = "new-app"
  project_id = data.firefly_workflows_project.by_path.id
  # ... other configuration
}

# Output project information
output "project_info" {
  value = {
    id              = data.firefly_workflows_project.by_path.id
    path            = data.firefly_workflows_project.by_path.path
    name            = data.firefly_workflows_project.by_path.name
    description     = data.firefly_workflows_project.by_path.description
    workspace_count = data.firefly_workflows_project.by_path.workspace_count
  }
}
```

## Schema

### Optional (exactly one required)

- `id` (String) - The unique identifier of the project
- `path` (String) - The path (name) of the project

### Read-Only

- `id` (String) - The unique identifier of the project (computed when using path)
- `path` (String) - The path (name) of the project (computed when using id)  
- `name` (String) - The name of the project
- `description` (String) - The description of the project
- `labels` (List of String) - Labels assigned to the project
- `cron_execution_pattern` (String) - Cron pattern for scheduled executions
- `parent_id` (String) - ID of the parent project
- `account_id` (String) - ID of the account the project belongs to
- `members_count` (Number) - Number of members assigned to the project
- `workspace_count` (Number) - Number of workspaces in the project