terraform {
  required_providers {
    firefly = {
      source = "registry.terraform.io/firefly/firefly"
      version = "1.0.0"
    }
  }
}

provider "firefly" {
  access_key = var.firefly_access_key
  secret_key = var.firefly_secret_key
  api_url    = "https://api.firefly.ai"
}

# Example cost guardrail
resource "firefly_guardrail" "cost_limit" {
  name        = "Monthly Cost Limit Alert"
  type        = "cost"
  is_enabled  = true
  severity    = 2
  
  scope {
    workspaces {
      include = ["*"]
    }
    repositories {
      include = ["*"]
    }
    branches {
      include = ["*"]
    }
    labels {
      include = ["*"]
    }
  }
  
  criteria {
    cost {
      threshold_amount = 5000.0
    }
  }
}

# Example policy guardrail
resource "firefly_guardrail" "policy_check" {
  name        = "Security Policy Enforcement"
  type        = "policy"
  is_enabled  = true
  severity    = 1
  
  scope {
    workspaces {
      include = ["prod-*"]
    }
    repositories {
      include = ["*"]
    }
    branches {
      include = ["*"]
    }
    labels {
      include = ["*"]
    }
  }
  
  criteria {
    policy {
      severity = "high"
    }
  }
}

# Example resource guardrail
resource "firefly_guardrail" "resource_control" {
  name        = "Prevent Production Deletions"
  type        = "resource"
  is_enabled  = true
  severity    = 1
  
  scope {
    workspaces {
      include = ["prod-*"]
    }
    
    repositories {
      include = ["*"]
    }
    
    branches {
      include = ["main", "master"]
    }
    
    labels {
      include = ["*"]
    }
  }
  
  criteria {
    resource {
      actions = ["delete"]
      specific_resources = ["aws_instance", "aws_db_instance"]
    }
  }
}

output "guardrail_ids" {
  value = {
    cost_limit     = firefly_guardrail.cost_limit.id
    policy_check   = firefly_guardrail.policy_check.id
    resource_control = firefly_guardrail.resource_control.id
  }
}