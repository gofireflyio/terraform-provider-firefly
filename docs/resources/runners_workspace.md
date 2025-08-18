# firefly_runners_workspace (Resource)

Manages a Firefly runners workspace for Terraform/OpenTofu infrastructure with VCS integration.

## Example Usage

```terraform
resource "firefly_runners_workspace" "example" {
  name        = "production-app"
  description = "Production application infrastructure"
  
  # VCS Configuration
  repository           = "myorg/infrastructure"
  vcs_integration_id   = "your-vcs-integration-id"
  vcs_type            = "github"
  default_branch      = "main"
  working_directory   = "environments/production"
  
  # Infrastructure Configuration
  iac_type            = "terraform"
  terraform_version   = "1.6.0"
  apply_rule          = "manual"
  triggers            = ["merge", "push"]
  
  # Organization
  labels              = ["production", "terraform"]
  project_id          = firefly_project.main.id
  consumed_variable_sets = [firefly_variable_set.aws_config.id]
  
  # Workspace Variables
  variables {
    key         = "ENVIRONMENT"
    value       = "production"
    sensitivity = "string"
    destination = "env"
  }
  
  variables {
    key         = "TF_VAR_region"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "iac"
  }
}
```

## Schema

### Required

- `name` (String) - The name of the workspace
- `repository` (String) - Repository URL or name
- `vcs_integration_id` (String) - VCS integration ID
- `vcs_type` (String) - VCS type (e.g., github, gitlab)
- `default_branch` (String) - Default branch for the workspace

### Optional

- `description` (String) - The description of the workspace
- `working_directory` (String) - Working directory within the repository
- `cron_execution_pattern` (String) - Cron pattern for scheduled executions
- `iac_type` (String) - Infrastructure as Code type. Valid values: `terraform`, `opentofu`. Defaults to `terraform`
- `terraform_version` (String) - Terraform version to use. Defaults to `1.5.7`
- `apply_rule` (String) - Apply rule. Valid values: `manual`, `auto`. Defaults to `manual`
- `triggers` (List of String) - List of triggers for the workspace
- `labels` (List of String) - Labels to assign to the workspace
- `consumed_variable_sets` (List of String) - List of variable set IDs that this workspace consumes
- `project_id` (String) - Project ID for workspace assignment
- `variables` (Block List) - Variables associated with the workspace

### Read-Only

- `id` (String) - The unique identifier of the workspace
- `account_id` (String) - Account ID that the workspace belongs to

### Nested Schema for `variables`

#### Required

- `key` (String) - The variable key
- `value` (String, Sensitive) - The variable value

#### Optional

- `sensitivity` (String) - The sensitivity of the variable. Valid values: `string`, `secret`. Defaults to `string`
- `destination` (String) - The destination of the variable. Valid values: `env`, `iac`. Defaults to `env`

## Import

Runners workspaces can be imported using their ID:

```shell
terraform import firefly_runners_workspace.example workspace-id-here
```