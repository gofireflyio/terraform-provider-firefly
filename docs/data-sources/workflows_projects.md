# firefly_workflows_projects (Data Source)

Fetches a list of Firefly projects with optional filtering.

## Example Usage

```terraform
# Get all projects
data "firefly_workflows_projects" "all" {}

# Search for specific projects
data "firefly_workflows_projects" "production" {
  search_query = "production"
}

# Use project data
resource "firefly_workflows_runners_workspace" "app" {
  name       = "new-workspace"
  project_id = data.firefly_workflows_projects.production.projects[0].id
  # ... other configuration
}

# Output project information
output "project_names" {
  value = [for project in data.firefly_workflows_projects.all.projects : project.name]
}
```

## Schema

### Optional

- `search_query` (String) - Search query to filter projects by name or description

### Read-Only

- `projects` (List of Object) - List of projects (see [below for nested schema](#nestedatt--projects))

<a id="nestedatt--projects"></a>
### Nested Schema for `projects`

Read-Only:

- `id` (String) - The unique identifier of the project
- `name` (String) - The name of the project
- `description` (String) - The description of the project
- `labels` (List of String) - Labels assigned to the project
- `cron_execution_pattern` (String) - Cron pattern for scheduled executions
- `parent_id` (String) - ID of the parent project
- `account_id` (String) - ID of the account the project belongs to
- `members_count` (Number) - Number of members assigned to the project
- `workspace_count` (Number) - Number of workspaces in the project