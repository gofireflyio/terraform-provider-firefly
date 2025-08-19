resource "firefly_workflows_project" "example" {
  name        = "my-project"
  description = "Example project with members"
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

resource "firefly_project_membership" "viewer" {
  project_id = firefly_workflows_project.example.id
  user_id    = "user789"
  email      = "viewer@example.com"
  role       = "viewer"
}