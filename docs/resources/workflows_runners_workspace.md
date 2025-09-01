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
  terraform_version = "1.6.0"  # Note: still use terraform_version even for OpenTofu
  apply_rule        = "auto"
  triggers          = ["push", "merge"]
  
  labels = ["opentofu", "staging"]
}
```

## Schema

### Required

- `name` (String) - The name of the workspace (cannot contain spaces)
- `repository` (String) - Repository URL or name
- `vcs_integration_id` (String) - ID of the VCS integration to use
- `vcs_type` (String) - Type of VCS (e.g., github, gitlab)
- `default_branch` (String) - Default branch for the workspace

### Optional

- `description` (String) - The description of the workspace. Defaults to empty string
- `working_directory` (String) - Working directory within the repository. Defaults to empty string
- `cron_execution_pattern` (String) - Cron pattern for scheduled executions. Defaults to empty string
- `iac_type` (String) - Infrastructure as Code type (terraform, opentofu). Defaults to `terraform`
- `terraform_version` (String) - Terraform version to use. Defaults to `1.5.7`
- `apply_rule` (String) - Apply rule (manual or auto). Defaults to `manual`
- `project_id` (String) - Project ID for workspace assignment. Defaults to empty string
- `triggers` (List of String) - List of triggers for the workspace
- `labels` (List of String) - Labels to assign to the workspace
- `consumed_variable_sets` (List of String) - List of variable set IDs that this workspace consumes
- `variables` (Block Set) - Variables associated with the workspace (see [below for nested schema](#nestedblock--variables))

### Read-Only

- `id` (String) - The unique identifier of the workspace
- `account_id` (String) - Account ID that the workspace belongs to

<a id="nestedblock--variables"></a>
### Nested Schema for `variables`

#### Required

- `key` (String) - The variable key/name
- `value` (String) - The variable value

#### Optional

- `sensitivity` (String) - The sensitivity of the variable (string or secret). Defaults to `string`
- `destination` (String) - The destination of the variable (env or iac). Defaults to `env`

## Import

Runners workspaces can be imported using their ID:

```shell
terraform import firefly_workflows_runners_workspace.example workspace-id-here
```