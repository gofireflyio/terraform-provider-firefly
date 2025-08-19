# firefly_workflows_runners_workspace (Resource)

Manages a Firefly runners workspace for executing Terraform or OpenTofu configurations.

## Example Usage

```terraform
# Basic runners workspace
resource "firefly_workflows_runners_workspace" "example" {
  name                = "production-app"
  description         = "Production application infrastructure"
  project_id          = firefly_workflows_project.main.id
  
  # VCS Configuration
  repository          = "myorg/infrastructure"
  vcs_integration_id  = "github-integration-id"
  vcs_type           = "github"
  default_branch     = "main"
  working_directory  = "environments/production"
  
  # Infrastructure Configuration
  iac_type           = "terraform"
  terraform_version  = "1.6.0"
  apply_rule         = "manual"
  triggers           = ["merge"]
  
  labels = ["production", "terraform"]
  consumed_variable_sets = [firefly_workflows_variable_set.aws_config.id]
  
  variables {
    key         = "ENVIRONMENT"
    value       = "production"
    sensitivity = "string"
    destination = "env"
  }
}

# OpenTofu workspace
resource "firefly_workflows_runners_workspace" "opentofu_example" {
  name               = "opentofu-workspace"
  description        = "OpenTofu workspace example"
  project_id         = firefly_workflows_project.main.id
  
  repository         = "myorg/opentofu-configs"
  vcs_integration_id = "github-integration-id"
  vcs_type          = "github"
  default_branch    = "main"
  
  iac_type          = "opentofu"
  opentofu_version  = "1.6.0"
  apply_rule        = "auto"
  triggers          = ["push", "merge"]
  
  labels = ["opentofu", "staging"]
}
```

## Schema

### Required

- `name` (String) - The name of the workspace
- `project_id` (String) - ID of the project this workspace belongs to
- `repository` (String) - Repository name in format "org/repo"
- `vcs_integration_id` (String) - ID of the VCS integration to use
- `vcs_type` (String) - Type of VCS. Valid values: `github`, `gitlab`, `bitbucket`, `azure_devops`
- `iac_type` (String) - Infrastructure as Code type. Valid values: `terraform`, `opentofu`

### Optional

- `description` (String) - The description of the workspace
- `labels` (List of String) - Labels to assign to the workspace
- `default_branch` (String) - Default branch to use. Defaults to `main`
- `working_directory` (String) - Working directory within the repository
- `terraform_version` (String) - Terraform version to use (required if iac_type is "terraform")
- `opentofu_version` (String) - OpenTofu version to use (required if iac_type is "opentofu")
- `apply_rule` (String) - When to apply changes. Valid values: `manual`, `auto`. Defaults to `manual`
- `triggers` (List of String) - Events that trigger runs. Valid values: `push`, `merge`, `pull_request`
- `consumed_variable_sets` (List of String) - List of variable set IDs to consume
- `variables` (Block Set) - Workspace-specific variables (see [below for nested schema](#nestedblock--variables))

### Read-Only

- `id` (String) - The unique identifier of the workspace

<a id="nestedblock--variables"></a>
### Nested Schema for `variables`

#### Required

- `key` (String) - The variable key/name
- `value` (String) - The variable value
- `sensitivity` (String) - Variable sensitivity level. Valid values: `string`, `secret`
- `destination` (String) - Where the variable is used. Valid values: `env`, `iac`

## Import

Runners workspaces can be imported using their ID:

```shell
terraform import firefly_workflows_runners_workspace.example workspace-id-here
```