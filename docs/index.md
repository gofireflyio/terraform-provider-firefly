# Firefly Provider

The Firefly provider allows you to manage infrastructure workflows and governance using Firefly's platform.

## Example Usage

```terraform
terraform {
  required_providers {
    firefly = {
      source = "gofireflyio/firefly"
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
resource "firefly_workflows_project" "main" {
  name        = "Production Infrastructure"
  description = "Main production project"
  labels      = ["production", "critical"]
}

# Create a variable set
resource "firefly_workflows_variable_set" "aws_config" {
  name        = "AWS Configuration"
  description = "Shared AWS configuration"
  labels      = ["aws", "shared"]
  
  variables {
    key         = "AWS_DEFAULT_REGION"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "env"
  }
}

# Create a runners workspace
resource "firefly_workflows_runners_workspace" "app" {
  name                = "production-app"
  description         = "Production application workspace"
  project_id          = firefly_workflows_project.main.id
  
  repository          = "myorg/infrastructure"
  vcs_integration_id  = "your-vcs-integration-id"
  vcs_type           = "github"
  iac_type           = "terraform"
  terraform_version  = "1.6.0"
  
  consumed_variable_sets = [firefly_workflows_variable_set.aws_config.id]
}
```

## Schema

### Required

- `access_key` (String, Sensitive) - Firefly access key for authentication
- `secret_key` (String, Sensitive) - Firefly secret key for authentication

### Optional

- `api_url` (String) - Firefly API URL. Defaults to `https://api.firefly.ai`

## Environment Variables

The provider can be configured using environment variables:

- `FIREFLY_ACCESS_KEY` - Sets the access key
- `FIREFLY_SECRET_KEY` - Sets the secret key  
- `FIREFLY_API_URL` - Sets the API URL

## Authentication

The Firefly provider uses access key and secret key authentication. You can obtain these credentials from your Firefly account settings.