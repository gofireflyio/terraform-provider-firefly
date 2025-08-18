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
  triggers            = ["merge"]
  
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
    key         = "AWS_REGION"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "iac"
  }
}