# firefly_workspaces (Data Source)

Fetches the list of Firefly workspaces with optional filtering capabilities.

## Example Usage

```terraform
# Get all workspaces
data "firefly_workspaces" "all" {}

# Search for workspaces by name
data "firefly_workspaces" "search" {
  search_value = "production"
}

# Filter workspaces by multiple criteria
data "firefly_workspaces" "filtered" {
  filters {
    workspace_name = ["prod-app", "staging-app"]
    labels         = ["production", "critical"]
    status         = ["active"]
    vcs_type      = ["github"]
  }
}

# Use the results
resource "firefly_workspace_labels" "example" {
  workspace_id = data.firefly_workspaces.filtered.workspaces[0].id
  labels       = ["managed-by-terraform"]
}
```

## Schema

### Optional

- `search_value` (String) - Search value to filter workspaces by name or other text fields
- `filters` (Block) - Filters for workspaces (see [below for nested schema](#nestedblock--filters))

### Read-Only

- `workspaces` (List of Object) - List of workspaces matching the criteria (see [below for nested schema](#nestedatt--workspaces))

<a id="nestedblock--filters"></a>
### Nested Schema for `filters`

#### Optional

- `workspace_name` (List of String) - Filter by workspace name
- `repositories` (List of String) - Filter by repositories
- `ci_tool` (List of String) - Filter by CI tool
- `labels` (List of String) - Filter by labels
- `status` (List of String) - Filter by status
- `is_managed_workflow` (Boolean) - Filter by whether the workflow is managed
- `vcs_type` (List of String) - Filter by version control system type

<a id="nestedatt--workspaces"></a>
### Nested Schema for `workspaces`

#### Read-Only

- `id` (String) - Unique identifier of the workspace
- `account_id` (String) - Account ID associated with the workspace
- `workspace_id` (String) - Workspace ID
- `workspace_name` (String) - Name of the workspace
- `repo` (String) - Repository associated with the workspace
- `repo_url` (String) - Repository URL
- `vcs_type` (String) - Version control system type
- `runner_type` (String) - CI/CD runner type
- `last_run_status` (String) - Status of the last run
- `last_apply_time` (String) - Timestamp of the last apply
- `last_plan_time` (String) - Timestamp of the last plan
- `last_run_time` (String) - Timestamp of the last run
- `iac_type` (String) - Type of Infrastructure as Code
- `iac_type_version` (String) - Version of the IaC tool
- `labels` (List of String) - List of labels associated with the workspace
- `runs_count` (Number) - Number of runs
- `is_workflow_managed` (Boolean) - Whether the workflow is managed
- `created_at` (String) - Timestamp when the workspace was created
- `updated_at` (String) - Timestamp when the workspace was last updated