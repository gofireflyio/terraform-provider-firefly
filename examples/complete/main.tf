terraform {
  required_providers {
    firefly = {
      source  = "gofireflyio/firefly"
      version = "~> 0.0.2"
    }
  }
}

provider "firefly" {
  access_key = var.firefly_access_key
  secret_key = var.firefly_secret_key
  api_url    = var.firefly_api_url
}

# Variables
variable "firefly_access_key" {
  description = "Firefly access key"
  type        = string
  sensitive   = true
}

variable "firefly_secret_key" {
  description = "Firefly secret key"
  type        = string
  sensitive   = true
}

variable "firefly_api_url" {
  description = "Firefly API URL"
  type        = string
  default     = "https://api.firefly.ai"
}

variable "aws_access_key" {
  description = "AWS access key"
  type        = string
  sensitive   = true
}

variable "aws_secret_key" {
  description = "AWS secret key"
  type        = string
  sensitive   = true
}

# Create a parent project
resource "firefly_workflows_project" "organization" {
  name        = "organization-infrastructure"
  description = "Top-level organization project"
  labels      = ["organization", "root"]
}

# Add team members to organization project (makes it visible in UI)
resource "firefly_project_membership" "org_admin" {
  project_id = firefly_workflows_project.organization.id
  user_id    = "admin-user-123"
  email      = "admin@company.com"
  role       = "admin"
}

# Create environment projects
resource "firefly_workflows_project" "production" {
  name        = "production-environment"
  description = "Production infrastructure project"
  labels      = ["production", "critical"]
  parent_id   = firefly_workflows_project.organization.id

  # Scheduled daily execution at 2 AM
  cron_execution_pattern = "0 2 * * *"

  variables {
    key         = "ENVIRONMENT"
    value       = "production"
    sensitivity = "string"
    destination = "env"
  }
}

# Add production team members
resource "firefly_project_membership" "prod_admin" {
  project_id = firefly_workflows_project.production.id
  user_id    = "prod-admin-456"
  email      = "prod-admin@company.com"
  role       = "admin"
}

resource "firefly_project_membership" "prod_member" {
  project_id = firefly_workflows_project.production.id
  user_id    = "prod-dev-789"
  email      = "prod-dev@company.com"
  role       = "member"
}

resource "firefly_workflows_project" "staging" {
  name        = "staging-environment"
  description = "Staging infrastructure project"
  labels      = ["staging", "test"]
  parent_id   = firefly_workflows_project.organization.id

  variables {
    key         = "ENVIRONMENT"
    value       = "staging"
    sensitivity = "string"
    destination = "env"
  }
}

# Add staging team members
resource "firefly_project_membership" "staging_lead" {
  project_id = firefly_workflows_project.staging.id
  user_id    = "staging-lead-321"
  email      = "staging-lead@company.com"
  role       = "admin"
}

# Create base variable set
resource "firefly_workflows_variable_set" "base_config" {
  name        = "Base Configuration"
  description = "Base configuration variables"
  labels      = ["base", "shared"]

  variables {
    key         = "COMPANY_NAME"
    value       = "ACME Corp"
    sensitivity = "string"
    destination = "env"
  }

  variables {
    key         = "TF_LOG_LEVEL"
    value       = "INFO"
    sensitivity = "string"
    destination = "env"
  }
}

# Create AWS variable set
resource "firefly_workflows_variable_set" "aws_config" {
  name        = "AWS Configuration"
  description = "AWS configuration variables"
  labels      = ["aws", "cloud", "shared"]
  parents     = [firefly_workflows_variable_set.base_config.id]

  variables {
    key         = "AWS_DEFAULT_REGION"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "env"
  }

  variables {
    key         = "AWS_ACCESS_KEY_ID"
    value       = var.aws_access_key
    sensitivity = "secret"
    destination = "env"
  }

  variables {
    key         = "AWS_SECRET_ACCESS_KEY"
    value       = var.aws_secret_key
    sensitivity = "secret"
    destination = "env"
  }

  variables {
    key         = "TF_VAR_region"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "iac"
  }
}

# Create production-specific variable set
resource "firefly_workflows_variable_set" "production_config" {
  name        = "Production Configuration"
  description = "Production-specific variables"
  labels      = ["production", "config"]
  parents     = [firefly_workflows_variable_set.aws_config.id]

  variables {
    key         = "INSTANCE_TYPE"
    value       = "m5.large"
    sensitivity = "string"
    destination = "iac"
  }

  variables {
    key         = "MIN_CAPACITY"
    value       = "3"
    sensitivity = "string"
    destination = "iac"
  }
}

# Create production workspaces
resource "firefly_workflows_runners_workspace" "prod_app" {
  name        = "production-application"
  description = "Production application infrastructure"
  project_id  = firefly_workflows_project.production.id

  # VCS Configuration
  repository         = "myorg/app-infrastructure"
  vcs_integration_id = "github-integration-id"
  vcs_type           = "github"
  default_branch     = "main"
  working_directory  = "environments/production"

  # Infrastructure Configuration
  iac_type          = "terraform"
  terraform_version = "1.6.0"
  apply_rule        = "manual"
  triggers          = ["merge"]

  # Labels and variable sets
  labels = ["production", "application", "terraform"]
  consumed_variable_sets = [
    firefly_workflows_variable_set.production_config.id
  ]

  # Workspace-specific variables
  variables {
    key         = "APP_NAME"
    value       = "production-app"
    sensitivity = "string"
    destination = "env"
  }

  variables {
    key         = "TF_VAR_app_name"
    value       = "production-app"
    sensitivity = "string"
    destination = "iac"
  }
}

resource "firefly_workflows_runners_workspace" "prod_database" {
  name        = "production-database"
  description = "Production database infrastructure"
  project_id  = firefly_workflows_project.production.id

  # VCS Configuration
  repository         = "myorg/database-infrastructure"
  vcs_integration_id = "github-integration-id"
  vcs_type           = "github"
  default_branch     = "main"
  working_directory  = "environments/production"

  # Infrastructure Configuration
  iac_type          = "terraform"
  terraform_version = "1.6.0"
  apply_rule        = "manual"
  triggers          = ["merge"]

  # Labels and variable sets
  labels = ["production", "database", "terraform"]
  consumed_variable_sets = [
    firefly_workflows_variable_set.production_config.id
  ]

  # Database-specific variables
  variables {
    key         = "DB_INSTANCE_CLASS"
    value       = "db.r5.xlarge"
    sensitivity = "string"
    destination = "iac"
  }

  variables {
    key         = "DB_BACKUP_RETENTION"
    value       = "30"
    sensitivity = "string"
    destination = "iac"
  }
}

# Create staging workspace
resource "firefly_workflows_runners_workspace" "staging_app" {
  name        = "staging-application"
  description = "Staging application infrastructure"
  project_id  = firefly_workflows_project.staging.id

  # VCS Configuration
  repository         = "myorg/app-infrastructure"
  vcs_integration_id = "github-integration-id"
  vcs_type           = "github"
  default_branch     = "develop"
  working_directory  = "environments/staging"

  # Infrastructure Configuration
  iac_type          = "terraform"
  terraform_version = "1.6.0"
  apply_rule        = "auto" # Auto-apply for staging
  triggers          = ["merge", "push"]

  # Labels and variable sets
  labels = ["staging", "application", "terraform"]
  consumed_variable_sets = [
    firefly_workflows_variable_set.aws_config.id # Use base AWS config for staging
  ]

  # Staging-specific variables
  variables {
    key         = "TF_VAR_instance_type"
    value       = "t3.medium"
    sensitivity = "string"
    destination = "iac"
  }

  variables {
    key         = "TF_VAR_min_capacity"
    value       = "1"
    sensitivity = "string"
    destination = "iac"
  }
}

# Governance policies for infrastructure compliance
resource "firefly_governance_policy" "s3_encryption" {
  name        = "S3 Bucket Encryption Policy"
  description = "Enforces that all S3 buckets have server-side encryption enabled"

  code = <<-EOT
    package firefly
    
    import rego.v1
    
    default allow := false
    
    allow if {
        input.resource_type == "aws_s3_bucket"
        input.configuration.server_side_encryption_configuration
        count(input.configuration.server_side_encryption_configuration) > 0
    }
    
    deny[msg] if {
        input.resource_type == "aws_s3_bucket"
        not input.configuration.server_side_encryption_configuration
        msg := "S3 bucket must have server-side encryption enabled"
    }
  EOT

  type         = ["aws_s3_bucket"]
  provider_ids = ["aws_all"]
  severity     = "high"
  category     = "Security"
  labels       = ["aws", "s3", "encryption", "security"]
  frameworks   = ["SOC2", "ISO27001"]
}

resource "firefly_governance_policy" "required_tags" {
  name        = "Required Resource Tags"
  description = "Ensures production resources have required tags"

  code = <<-EOT
    package firefly
    
    import rego.v1
    
    required_tags := ["Environment", "Owner", "CostCenter"]
    
    default allow := false
    
    allow if {
        input.resource_type in ["aws_instance", "aws_db_instance"]
        tags := object.get(input.configuration, "tags", {})
        every tag in required_tags {
            tags[tag]
            tags[tag] != ""
        }
    }
    
    deny[msg] if {
        input.resource_type in ["aws_instance", "aws_db_instance"]
        tags := object.get(input.configuration, "tags", {})
        some tag in required_tags
        not tags[tag]
        msg := sprintf("Resource missing required tag: %s", [tag])
    }
  EOT

  type         = ["aws_instance", "aws_db_instance"]
  provider_ids = ["aws_all"]
  severity     = "medium"
  category     = "Governance"
  labels       = ["aws", "tagging", "governance"]
  frameworks   = ["SOC2"]
}

# Data sources for existing resources
data "firefly_workflows_projects" "all_projects" {
  search_query = "infrastructure"
}

data "firefly_workflows_variable_sets" "shared_sets" {
  search_query = "shared"
}

data "firefly_governance_policies" "security_policies" {
  category = "Security"
}

# Outputs
output "production_project_info" {
  description = "Production project information"
  value = {
    id              = firefly_workflows_project.production.id
    name            = firefly_workflows_project.production.name
    workspace_count = firefly_workflows_project.production.workspace_count
  }
}

output "workspace_ids" {
  description = "All workspace IDs"
  value = {
    prod_app      = firefly_workflows_runners_workspace.prod_app.id
    prod_database = firefly_workflows_runners_workspace.prod_database.id
    staging_app   = firefly_workflows_runners_workspace.staging_app.id
  }
}

output "variable_set_info" {
  description = "Variable set information"
  value = {
    aws_config = {
      id      = firefly_workflows_variable_set.aws_config.id
      version = firefly_workflows_variable_set.aws_config.version
    }
    production_config = {
      id      = firefly_workflows_variable_set.production_config.id
      version = firefly_workflows_variable_set.production_config.version
    }
  }
}

output "governance_policy_info" {
  description = "Governance policy information"
  value = {
    s3_encryption = {
      id       = firefly_governance_policy.s3_encryption.id
      name     = firefly_governance_policy.s3_encryption.name
      severity = firefly_governance_policy.s3_encryption.severity
    }
    required_tags = {
      id       = firefly_governance_policy.required_tags.id
      name     = firefly_governance_policy.required_tags.name
      severity = firefly_governance_policy.required_tags.severity
    }
  }
}

output "security_policies_summary" {
  description = "Summary of security policies"
  value = {
    count = length(data.firefly_governance_policies.security_policies.policies)
    names = [for policy in data.firefly_governance_policies.security_policies.policies : policy.name]
  }
}