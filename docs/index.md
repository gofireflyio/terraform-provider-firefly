# Firefly Provider

The Firefly provider allows you to manage your Firefly infrastructure as code using Terraform.

## Example Usage

```terraform
terraform {
  required_providers {
    firefly = {
      source  = "firefly/firefly"
      version = "~> 1.0"
    }
  }
}

provider "firefly" {
  access_key = var.firefly_access_key
  secret_key = var.firefly_secret_key
  api_url    = "https://api.firefly.ai" # Optional
}

# Create a project
resource "firefly_project" "main" {
  name        = "Production Infrastructure"
  description = "Main production project"
  labels      = ["production", "critical"]
}

# Create a variable set
resource "firefly_variable_set" "aws_config" {
  name        = "AWS Configuration"
  description = "Shared AWS configuration"
  
  variables {
    key         = "AWS_REGION"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "env"
  }
}

# Create a runners workspace
resource "firefly_runners_workspace" "app" {
  name                 = "production-app"
  project_id           = firefly_project.main.id
  repository           = "myorg/infrastructure"
  vcs_integration_id   = "your-vcs-integration-id"
  vcs_type            = "github"
  default_branch      = "main"
  iac_type            = "terraform"
  consumed_variable_sets = [firefly_variable_set.aws_config.id]
}
```

## Authentication

The Firefly provider requires authentication using access and secret keys:

```terraform
provider "firefly" {
  access_key = var.firefly_access_key
  secret_key = var.firefly_secret_key
}
```

### Environment Variables

You can also configure authentication using environment variables:

- `FIREFLY_ACCESS_KEY` - Your Firefly access key
- `FIREFLY_SECRET_KEY` - Your Firefly secret key  
- `FIREFLY_API_URL` - API endpoint (optional, defaults to https://api.firefly.ai)

## Schema

### Required

- `access_key` (String, Sensitive) - The access key for API operations
- `secret_key` (String, Sensitive) - The secret key for API operations

### Optional

- `api_url` (String) - The URL of the Firefly API. Defaults to `https://api.firefly.ai`