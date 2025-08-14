
# Simple cost guardrail with ALL required fields
resource "firefly_guardrail" "simple_cost" {
  name        = "Simple Cost Guardrail Test"
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
      threshold_amount = 1000.0
    }
  }
}

# Resource guardrail with regions
resource "firefly_guardrail" "simple_resource" {
  name        = "Simple Resource Guardrail"
  type        = "resource"
  is_enabled  = true
  severity    = 1
  
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
    resource {
      actions = ["delete"]
      specific_resources = ["aws_instance"]
      
      regions {
        include = ["*"]
      }
      
      asset_types {
        include = ["*"]
      }
    }
  }
}

output "guardrail_results" {
  value = {
    cost_id = firefly_guardrail.simple_cost.id
    resource_id = firefly_guardrail.simple_resource.id
  }
}