# firefly_workspace_runs (Data Source)

Fetches the list of runs for a Firefly workspace with optional filtering capabilities.

## Example Usage

```terraform
# Get all runs for a workspace
data "firefly_workspace_runs" "all" {
  workspace_id = "workspace-123"
}

# Search runs by name or other criteria
data "firefly_workspace_runs" "search" {
  workspace_id = "workspace-123"
  search_value = "deploy"
}

# Filter runs by specific criteria
data "firefly_workspace_runs" "filtered" {
  workspace_id = "workspace-123"
  
  filters {
    status     = ["success", "running"]
    branch     = ["main", "develop"]
    vcs_type   = ["github"]
  }
}

# Access run information
output "latest_successful_run" {
  value = [
    for run in data.firefly_workspace_runs.filtered.runs :
    run if run.status == "success"
  ][0]
}
```

## Schema

### Required

- `workspace_id` (String) - ID of the workspace to fetch runs for

### Optional

- `search_value` (String) - Search value to filter runs by name or other text fields
- `filters` (Block) - Filters for workspace runs (see [below for nested schema](#nestedblock--filters))

### Read-Only

- `runs` (List of Object) - List of workspace runs matching the criteria (see [below for nested schema](#nestedatt--runs))

<a id="nestedblock--filters"></a>
### Nested Schema for `filters`

#### Optional

- `run_id` (List of String) - Filter by run ID
- `run_name` (List of String) - Filter by run name
- `status` (List of String) - Filter by run status
- `branch` (List of String) - Filter by branch
- `commit_id` (List of String) - Filter by commit ID
- `ci_tool` (List of String) - Filter by CI tool
- `vcs_type` (List of String) - Filter by version control system type
- `repository` (List of String) - Filter by repository

<a id="nestedatt--runs"></a>
### Nested Schema for `runs`

#### Read-Only

- `id` (String) - Unique identifier of the run
- `workspace_id` (String) - ID of the workspace
- `workspace_name` (String) - Name of the workspace
- `run_id` (String) - Run ID
- `run_name` (String) - Name of the run
- `status` (String) - Status of the run
- `branch` (String) - Branch associated with the run
- `commit_id` (String) - Commit ID
- `commit_url` (String) - URL to the commit
- `runner_type` (String) - Type of runner used
- `build_id` (String) - Build ID
- `build_url` (String) - URL to the build
- `build_name` (String) - Name of the build
- `vcs_type` (String) - Version control system type
- `repo` (String) - Repository name
- `repo_url` (String) - Repository URL
- `title` (String) - Title of the run
- `created_at` (String) - Timestamp when the run was created
- `updated_at` (String) - Timestamp when the run was last updated