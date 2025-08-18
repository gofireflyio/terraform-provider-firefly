resource "firefly_variable_set" "example" {
  name        = "AWS Configuration"
  description = "Shared AWS configuration variables"
  labels      = ["aws", "shared", "production"]
  
  # Parent variable sets for inheritance
  parents = [firefly_variable_set.base_config.id]
  
  # Variables in the set
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
    key         = "TERRAFORM_BACKEND_BUCKET"
    value       = "my-terraform-state-bucket"
    sensitivity = "string"
    destination = "iac"
  }
}