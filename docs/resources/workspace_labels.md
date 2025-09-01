# firefly_workspace_labels (Resource)

Manages labels for a Firefly workspace.

## Example Usage

```terraform
resource "firefly_workspace_labels" "example" {
  workspace_id = "workspace-123"
  labels       = ["production", "critical", "database"]
}

# Using with a workspace data source
data "firefly_workspaces" "example" {
  name = "my-workspace"
}

resource "firefly_workspace_labels" "from_data" {
  workspace_id = data.firefly_workspaces.example.workspaces[0].id
  labels       = ["managed-by-terraform"]
}
```

## Schema

### Required

- `workspace_id` (String) - The ID of the workspace to manage labels for. Changing this forces a new resource to be created.
- `labels` (List of String) - List of labels to assign to the workspace

### Read-Only

- `id` (String) - Unique identifier of the workspace
- `workspace_name` (String) - The name of the workspace
- `updated_at` (String) - Timestamp when the labels were last updated

## Import

Workspace labels can be imported using the workspace ID:

```shell
terraform import firefly_workspace_labels.example workspace-id-here
```

## Notes

- This resource manages the complete set of labels for a workspace. Any labels not included in the configuration will be removed.
- Changing the `workspace_id` will force the creation of a new resource.