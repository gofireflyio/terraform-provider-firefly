# Terraform Provider for Firefly

[![Go Report Card](https://goreportcard.com/badge/github.com/gofireflyio/terraform-provider-firefly)](https://goreportcard.com/report/github.com/gofireflyio/terraform-provider-firefly)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

A comprehensive Terraform Provider for managing [Firefly](https://gofirefly.io) infrastructure resources. This provider enables Infrastructure as Code (IaC) management of projects, workspaces, variable sets, guardrails, and more through Terraform.

## Features

- **Complete Resource Coverage** - Manage all major Firefly resources
- **Production Ready** - Comprehensive testing with 49 test scenarios
- **Modern Architecture** - Built with Terraform Plugin Framework v1.15.1
- **Professional Tooling** - Makefile, automated releases, and development tools
- **Extensive Documentation** - Complete examples and usage guides

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23.7 (for development)

## Installation

### Terraform Registry (Recommended)

```hcl
terraform {
  required_providers {
    firefly = {
      source  = "firefly/firefly"
      version = "~> 1.0"
    }
  }
}
```

### Building From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/gofireflyio/terraform-provider-firefly.git
   cd terraform-provider-firefly
   ```

2. Build the provider:
   ```bash
   make build
   # or: go install
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

## Testing

The provider includes comprehensive test coverage with 49 total tests:

### Unit Tests (21 tests)
```bash
# Run all unit tests
make test

# Run specific client tests
go test ./internal/client -v -run TestProject
```

### Acceptance Tests (28 scenarios)
Acceptance tests create real resources and require valid Firefly credentials.

```bash
# Set required environment variables
export FIREFLY_ACCESS_KEY="your-access-key"
export FIREFLY_SECRET_KEY="your-secret-key"
export TF_ACC=1

# Run all acceptance tests
make testacc

# Run specific resource tests
go test ./internal/provider -v -run TestAccProject
```

## Development

### Prerequisites
- [Go](https://golang.org/doc/install) >= 1.23.7
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- Valid Firefly API credentials (for acceptance tests)

### Setup
```bash
# Clone the repository
git clone https://github.com/gofireflyio/terraform-provider-firefly.git
cd terraform-provider-firefly

# Install dependencies
make deps

# Build the provider
make build

# Run tests
make test
```

### Development Commands
```bash
# Format code
make fmt

# Run linting
make vet

# Generate documentation
make docs

# Build for debugging
make debug

# Run specific tests
go test ./internal/client -v
go test ./internal/provider -v
```

### Testing Your Changes
1. Build the provider: `make build`
2. Run unit tests: `make test`
3. Run acceptance tests: `TF_ACC=1 make testacc` (requires credentials)
4. Test manually with debug mode: `make debug`

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Format code (`make fmt`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.
