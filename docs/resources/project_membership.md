# firefly_project_membership Resource

Manages membership of a user in a Firefly project. This resource allows you to add users to projects with specific roles.

## Example Usage

```terraform
resource "firefly_workflows_project" "example" {
  name        = "my-project"
  description = "Example project"
}

resource "firefly_project_membership" "admin" {
  project_id = firefly_workflows_project.example.id
  user_id    = "user123"
  email      = "admin@example.com"
  role       = "admin"
}

resource "firefly_project_membership" "member" {
  project_id = firefly_workflows_project.example.id
  user_id    = "user456"
  email      = "member@example.com"
  role       = "member"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new resource to be created.
* `user_id` - (Required) The ID of the user to add to the project. Changing this forces a new resource to be created.
* `email` - (Required) The email address of the user.
* `role` - (Required) The role of the user in the project (e.g., 'admin', 'member', 'viewer').

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for the membership in the format `project_id:user_id`.

## Import

Project memberships can be imported using the format `project_id:user_id`:

```bash
terraform import firefly_project_membership.example project-123:user-456
```

## Notes

- When a project membership is deleted from Terraform, the user will be removed from the project.
- Changing the `project_id` or `user_id` will force the creation of a new resource.
- Role updates are handled by removing and re-adding the user with the new role.