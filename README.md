# Terraform Provider for Firefly

This Terraform Provider allows you to manage your Firefly SaaS resources using Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

To use the provider, create a `terraform` block:

```hcl
terraform {
  required_providers {
    firefly = {
      source = "firefly/firefly"
      version = "~> 1.0"
    }
  }
}

provider "firefly" {
  access_key = "your-access-key"
  secret_key = "your-secret-key"
  api_url    = "https://api.firefly.ai"  # Optional, defaults to https://api.firefly.ai
}
```

## Resources

The provider supports the following resources:

- **`firefly_project`** - Manage Firefly projects with hierarchical organization
- **`firefly_runners_workspace`** - Manage Terraform/OpenTofu runner workspaces with VCS integration
- **`firefly_variable_set`** - Manage reusable variable sets with inheritance
- **`firefly_guardrail`** - Manage guardrail rules for cost, policy, and resource governance
- **`firefly_workspace_labels`** - Manage workspace label assignments

## Data Sources

The provider supports the following data sources:

- **`firefly_projects`** - List and filter projects
- **`firefly_project`** - Get a single project by ID
- **`firefly_variable_sets`** - List and filter variable sets
- **`firefly_variable_set`** - Get a single variable set by ID
- **`firefly_workspaces`** - List and filter CI workspaces
- **`firefly_workspace_runs`** - Get workspace run information
- **`firefly_guardrails`** - List and filter guardrail rules

### Examples

**Creating a complete infrastructure setup:**

```hcl
# Create a project
resource "firefly_project" "main" {
  name        = "Production Infrastructure"
  description = "Main production project"
  labels      = ["production", "critical"]
  
  variables {
    key         = "AWS_REGION"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "env"
  }
}

# Create a variable set for shared configuration
resource "firefly_variable_set" "aws_config" {
  name        = "AWS Configuration"
  description = "Shared AWS configuration variables"
  labels      = ["aws", "shared"]
  
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
}

# Create a runners workspace
resource "firefly_runners_workspace" "app" {
  name                 = "production-app"
  description          = "Production application infrastructure"
  project_id           = firefly_project.main.id
  
  repository           = "myorg/infrastructure"
  vcs_integration_id   = "your-vcs-integration-id"
  vcs_type            = "github"
  default_branch      = "main"
  working_directory   = "environments/production"
  
  iac_type            = "terraform"
  terraform_version   = "1.6.0"
  apply_rule          = "manual"
  triggers            = ["merge"]
  
  labels              = ["production", "terraform"]
  consumed_variable_sets = [firefly_variable_set.aws_config.id]
  
  variables {
    key         = "ENVIRONMENT"
    value       = "production"
    sensitivity = "string"
    destination = "env"
  }
}

# Create a guardrail rule
resource "firefly_guardrail" "cost_guardrail" {
  name      = "Production Cost Threshold"
  type      = "cost"
  is_enabled = true
  severity  = 2
  
  scope {
    workspaces {
      include = ["production-*"]
    }
    
    labels {
      include = ["production"]
    }
  }
  
  criteria {
    cost {
      threshold_amount = 1000
    }
  }
}
```

**Using data sources to reference existing infrastructure:**

```hcl
# Find existing projects
data "firefly_projects" "existing" {
  search_query = "legacy"
}

# Get specific project details
data "firefly_project" "main" {
  id = "existing-project-id"
}

# Create workspace in existing project
resource "firefly_runners_workspace" "new_service" {
  name       = "new-service"
  project_id = data.firefly_project.main.id
  # ... other configuration
}

# Find shared variable sets
data "firefly_variable_sets" "shared" {
  search_query = "aws"
}

# Reference existing variable set
resource "firefly_runners_workspace" "with_shared_vars" {
  name = "service-with-shared-config"
  consumed_variable_sets = [data.firefly_variable_sets.shared.variable_sets[0].id]
  # ... other configuration
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
