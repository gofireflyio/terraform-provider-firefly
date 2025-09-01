terraform {
  required_providers {
    firefly = {
      source = "gofireflyio/firefly"
      version = "~> 0.0.1"
    }
  }
}

# Fetch project by path (name)
data "firefly_workflows_project" "by_path" {
  path = "Production Infrastructure"
}

# Output the fetched project details
output "project_details" {
  value = {
    id              = data.firefly_workflows_project.by_path.id
    path            = data.firefly_workflows_project.by_path.path
    name            = data.firefly_workflows_project.by_path.name
    description     = data.firefly_workflows_project.by_path.description
    workspace_count = data.firefly_workflows_project.by_path.workspace_count
    members_count   = data.firefly_workflows_project.by_path.members_count
  }
}