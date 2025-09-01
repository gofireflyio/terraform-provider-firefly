# Terraform Provider for Firefly

[![Go Report Card](https://goreportcard.com/badge/github.com/gofireflyio/terraform-provider-firefly)](https://goreportcard.com/report/github.com/gofireflyio/terraform-provider-firefly)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

A comprehensive Terraform Provider for managing [Firefly](https://gofirefly.io) infrastructure resources. This provider enables Infrastructure as Code (IaC) management of projects, workspaces, variable sets, guardrails, and more through Terraform.

## Features

- **Complete Resource Coverage** - Manage all major Firefly resources
- **Project Membership Management** - Assign users to projects with roles (admin/member) for UI visibility
- **Production Ready** - Comprehensive testing with 49 test scenarios
- **Modern Architecture** - Built with Terraform Plugin Framework v1.15.1
- **Professional Tooling** - Makefile, automated releases, and development tools
- **Extensive Documentation** - Complete examples and usage guides

## Status

🚧 **Provider Status**: This provider is fully functional but not yet published to the Terraform Registry.

✅ **What Works**:
- All resources and data sources are implemented and tested
- Full CRUD operations for all Firefly resources (create, read, update, delete)
- Workspace-project relationships and state management
- Variable sets and consumed variable sets handling
- Comprehensive test coverage (49 test scenarios)
- Complete documentation and examples
- **Recently Fixed**: Critical workspace read/delete issues resolved (Aug 2025)

🔄 **Next Steps**:
- Terraform Registry publication pending
- For now, use local development mode (instructions below)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23.7 (for development)

## Installation

> **⚠️ Important**: This provider is not yet published to the Terraform Registry. Use the local development method below.

### Local Development (Required)

Since the provider is not yet in the Terraform Registry, you need to build and run it locally:

1. **Clone and build the provider:**
   ```bash
   git clone https://github.com/gofireflyio/terraform-provider-firefly.git
   cd terraform-provider-firefly
   go build -o terraform-provider-firefly
   ```

2. **Run the provider in debug mode:**
   ```bash
   go run main.go -debug
   ```
   
   This will output something like:
   ```
   Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:
   
   TF_REATTACH_PROVIDERS='{"registry.terraform.io/firefly/firefly":{"Protocol":"grpc","ProtocolVersion":6,"Pid":12345,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin123"}}}'
   ```

3. **In another terminal, set the environment variable and run Terraform:**
   ```bash
   export TF_REATTACH_PROVIDERS='{"registry.terraform.io/firefly/firefly":{"Protocol":"grpc","ProtocolVersion":6,"Pid":12345,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin123"}}}'
   
   # Now you can use terraform normally
   terraform init
   terraform plan
   terraform apply
   ```

### Future Registry Installation

Once published to the Terraform Registry, you'll be able to use:

```hcl
terraform {
  required_providers {
    firefly = {
      source  = "gofireflyio/firefly"
      version = "~> 1.0"
    }
  }
}
```

## Using the Provider

### Provider Configuration

When running locally (current method), configure your Terraform files like this:

```hcl
terraform {
  required_providers {
    firefly = {
      source = "registry.terraform.io/firefly/firefly"  # For local debug mode
    }
  }
}

provider "firefly" {
  access_key = "your-firefly-access-key"
  secret_key = "your-firefly-secret-key"
  api_url    = "https://api.firefly.ai"  # Optional, defaults to https://api.firefly.ai
}
```

### Authentication

You can configure authentication in several ways:

1. **Direct configuration (not recommended for production):**
   ```hcl
   provider "firefly" {
     access_key = "your-access-key"
     secret_key = "your-secret-key"
   }
   ```

2. **Environment variables (recommended):**
   ```bash
   export FIREFLY_ACCESS_KEY="your-access-key"
   export FIREFLY_SECRET_KEY="your-secret-key"
   export FIREFLY_API_URL="https://api.firefly.ai"  # Optional
   ```
   
   ```hcl
   provider "firefly" {
     # Configuration will be read from environment variables
   }
   ```

3. **Terraform variables:**
   ```hcl
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
   
   provider "firefly" {
     access_key = var.firefly_access_key
     secret_key = var.firefly_secret_key
   }
   ```

## Resources

The provider supports the following resources:

- **`firefly_workflows_project`** - Manage Firefly projects with hierarchical organization
- **`firefly_project_membership`** - Manage project member assignments and roles (admin, member)
- **`firefly_workflows_runners_workspace`** - Manage Terraform/OpenTofu runner workspaces with VCS integration
- **`firefly_workflows_variable_set`** - Manage reusable variable sets with inheritance
- **`firefly_workflows_guardrail`** - Manage guardrail rules for cost, policy, and resource governance
- **`firefly_workspace_labels`** - Manage workspace label assignments

## Data Sources

The provider supports the following data sources:

- **`firefly_workflows_projects`** - List and filter projects
- **`firefly_workflows_project`** - Get a single project by ID
- **`firefly_workflows_variable_sets`** - List and filter variable sets
- **`firefly_workflows_variable_set`** - Get a single variable set by ID
- **`firefly_workspaces`** - List and filter CI workspaces
- **`firefly_workspace_runs`** - Get workspace run information
- **`firefly_workflows_guardrails`** - List and filter guardrail rules

### Examples

**Creating a complete infrastructure setup:**

```hcl
# Create a project
resource "firefly_workflows_project" "main" {
  name        = "production-infrastructure"
  description = "Main production project"
  labels      = ["production", "critical"]
  
  variables {
    key         = "AWS_REGION"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "env"
  }
}

# Add team members to the project (makes it visible in UI)
resource "firefly_project_membership" "admin" {
  project_id = firefly_workflows_project.main.id
  user_id    = "user123"
  email      = "admin@company.com"
  role       = "admin"
}

resource "firefly_project_membership" "developer" {
  project_id = firefly_workflows_project.main.id
  user_id    = "user456" 
  email      = "dev@company.com"
  role       = "member"
}

# Create a variable set for shared configuration
resource "firefly_workflows_variable_set" "aws_config" {
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
resource "firefly_workflows_runners_workspace" "app" {
  name                 = "production-app"
  description          = "Production application infrastructure"
  project_id           = firefly_workflows_project.main.id
  
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
  consumed_variable_sets = [firefly_workflows_variable_set.aws_config.id]
  
  variables {
    key         = "ENVIRONMENT"
    value       = "production"
    sensitivity = "string"
    destination = "env"
  }
}

# Create a guardrail rule
resource "firefly_workflows_guardrail" "cost_guardrail" {
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
data "firefly_workflows_projects" "existing" {
  search_query = "legacy"
}

# Get specific project details
data "firefly_workflows_project" "main" {
  id = "existing-project-id"
}

# Create workspace in existing project
resource "firefly_workflows_runners_workspace" "new_service" {
  name       = "new-service"
  project_id = data.firefly_workflows_project.main.id
  # ... other configuration
}

# Find shared variable sets
data "firefly_workflows_variable_sets" "shared" {
  search_query = "aws"
}

# Reference existing variable set
resource "firefly_workflows_runners_workspace" "with_shared_vars" {
  name = "service-with-shared-config"
  consumed_variable_sets = [data.firefly_workflows_variable_sets.shared.variable_sets[0].id]
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

### Running the Provider Locally

To test the provider locally with Terraform:

1. **Build the provider:**
   ```bash
   make build
   # or: go install
   ```

2. **Run in debug mode:**
   ```bash
   go run main.go -debug
   ```
   This will output a `TF_REATTACH_PROVIDERS` environment variable that you need to set.

3. **Set up your Terraform configuration:**
   ```bash
   # Copy the TF_REATTACH_PROVIDERS from the debug output and export it
   # Note: Use the exact output from step 2, it will look like this:
   export TF_REATTACH_PROVIDERS='{"registry.terraform.io/firefly/firefly":{"Protocol":"grpc","ProtocolVersion":6,"Pid":12345,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/.../T/plugin123456789"}}}'
   
   # Create a test Terraform configuration
   mkdir test-terraform && cd test-terraform
   ```

4. **Create a test configuration file (main.tf):**
   ```hcl
   terraform {
     required_providers {
       firefly = {
         source = "registry.terraform.io/firefly/firefly"  # Required for debug mode
       }
     }
   }
   
   provider "firefly" {
     access_key = "your-firefly-access-key"
     secret_key = "your-firefly-secret-key"
     api_url    = "https://api.firefly.ai"
   }
   
   # Test resource - note that project names must be alphanumeric with hyphens/underscores
   resource "firefly_workflows_project" "test" {
     name        = "local-test-project"
     description = "Testing locally built provider"
     labels      = ["test", "local"]
     
     variables {
       key         = "TEST_ENV" 
       value       = "local"
       sensitivity = "string"
       destination = "env"
     }
   }
   ```

5. **Run Terraform commands:**
   ```bash
   terraform init
   terraform plan
   terraform apply
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
10. After PR is merged, switch back to main and pull latest changes:
    ```bash
    git checkout main
    git pull origin main
    ```

### Branch Management
- **Always work on feature branches** - avoid pushing directly to main
- **Use descriptive branch names** - e.g., `feature/project-path-lookup`, `fix/email-schema`
- **Main branch is protected** - all changes must go through Pull Requests

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.
