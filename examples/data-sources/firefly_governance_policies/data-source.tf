# Get all governance policies
data "firefly_governance_policies" "all_policies" {}

# Get policies by category
data "firefly_governance_policies" "security_policies" {
  category = "Security"
}

# Get policies with specific labels
data "firefly_governance_policies" "aws_policies" {
  labels = ["aws", "security"]
}

# Search for policies by name/description
data "firefly_governance_policies" "s3_policies" {
  query = "s3"
}

# Get policies for a specific compliance framework
data "firefly_governance_policies" "soc2_policies" {
  labels = ["SOC2"]
}

# Filter policies by combining multiple criteria
data "firefly_governance_policies" "strict_aws_security" {
  category = "Security"
  labels   = ["aws"]
}

# Output examples showing how to use the data
output "all_policy_names" {
  description = "Names of all governance policies"
  value       = [for policy in data.firefly_governance_policies.all_policies.policies : policy.name]
}

output "security_policy_count" {
  description = "Number of security policies"
  value       = length(data.firefly_governance_policies.security_policies.policies)
}

output "aws_policy_details" {
  description = "Details of AWS-related policies"
  value = [
    for policy in data.firefly_governance_policies.aws_policies.policies : {
      name        = policy.name
      id          = policy.id
      severity    = policy.severity
      category    = policy.category
      types       = policy.type
      frameworks  = policy.frameworks
    }
  ]
}

output "s3_policies_summary" {
  description = "Summary of S3-related policies"
  value = {
    count = length(data.firefly_governance_policies.s3_policies.policies)
    names = [for policy in data.firefly_governance_policies.s3_policies.policies : policy.name]
    severities = distinct([for policy in data.firefly_governance_policies.s3_policies.policies : policy.severity])
  }
}

# Example: Create a governance report
locals {
  all_policies = data.firefly_governance_policies.all_policies.policies
  
  # Group policies by category
  policies_by_category = {
    for category in distinct([for p in local.all_policies : p.category if p.category != ""]) :
    category => [for p in local.all_policies : p if p.category == category]
  }
  
  # Group policies by severity
  policies_by_severity = {
    for severity in distinct([for p in local.all_policies : p.severity]) :
    severity => [for p in local.all_policies : p if p.severity == severity]
  }
}

output "governance_report" {
  description = "Comprehensive governance policy report"
  value = {
    total_policies = length(local.all_policies)
    
    by_category = {
      for category, policies in local.policies_by_category :
      category => {
        count = length(policies)
        policy_names = [for p in policies : p.name]
      }
    }
    
    by_severity = {
      for severity, policies in local.policies_by_severity :
      severity => {
        count = length(policies)
        categories = distinct([for p in policies : p.category if p.category != ""])
      }
    }
    
    frameworks_coverage = distinct(flatten([
      for policy in local.all_policies : policy.frameworks
    ]))
    
    resource_types_covered = distinct(flatten([
      for policy in local.all_policies : policy.type
    ]))
  }
}

# Example: Find policies that apply to specific resource types
locals {
  ec2_policies = [
    for policy in data.firefly_governance_policies.all_policies.policies :
    policy if contains(policy.type, "aws_instance")
  ]
  
  s3_policies = [
    for policy in data.firefly_governance_policies.all_policies.policies :
    policy if contains(policy.type, "aws_s3_bucket")
  ]
}

output "resource_specific_policies" {
  description = "Policies that apply to specific resource types"
  value = {
    ec2_instances = {
      count = length(local.ec2_policies)
      policies = [for p in local.ec2_policies : { name = p.name, severity = p.severity }]
    }
    s3_buckets = {
      count = length(local.s3_policies)
      policies = [for p in local.s3_policies : { name = p.name, severity = p.severity }]
    }
  }
}