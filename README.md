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

### Examples

**Creating a guardrail rule:**

```hcl
resource "firefly_guardrail" "cost_guardrail" {
  name      = "Cost Threshold Alert"
  type      = "cost"
  is_enabled = true
  severity  = 2
  
  scope {
    workspaces {
      include = ["prod-*"]
      exclude = ["prod-test"]
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

**Managing workspace labels:**

```hcl
resource "firefly_workspace_labels" "prod_labels" {
  workspace_id = "507f1f77bcf86cd799439011"
  labels       = ["production", "terraform", "critical"]
}
```

**Reading workspaces:**

```hcl
data "firefly_workspaces" "all" {
  filters {
    labels = ["production"]
  }
}

output "production_workspaces" {
  value = data.firefly_workspaces.all
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
