# firefly_projects (Data Source)

Fetches a list of Firefly projects with optional filtering.

## Example Usage

```terraform
# Get all projects
data "firefly_projects" "all" {}

# Filter projects by search query
data "firefly_projects" "production" {
  search_query = "production"
}

# Use project data
resource "firefly_runners_workspace" "app" {
  name       = "new-app"
  project_id = data.firefly_projects.production.projects[0].id
  # ... other configuration
}

# Output project information
output "production_projects" {
  value = [
    for project in data.firefly_projects.production.projects : {
      name = project.name
      id   = project.id
    }
  ]
}
```

## Schema

### Optional

- `search_query` (String) - Optional search query to filter projects

### Read-Only

- `projects` (List of Object) - List of projects

### Nested Schema for `projects`

- `id` (String) - The unique identifier of the project
- `name` (String) - The name of the project
- `description` (String) - The description of the project
- `labels` (List of String) - Labels assigned to the project
- `cron_execution_pattern` (String) - Cron pattern for scheduled executions
- `parent_id` (String) - ID of the parent project
- `account_id` (String) - ID of the account the project belongs to
- `members_count` (Number) - Number of members assigned to the project
- `workspace_count` (Number) - Number of workspaces in the project