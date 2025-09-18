resource "firefly_governance_policy" "kms_key_rotation" {
  name        = "KMS Key Rotation Required"
  description = "Ensure KMS keys have rotation enabled"

  code = <<-EOT
    package firefly

    firefly {
      match
    }

    match {
      input.enable_key_rotation == false
      input.is_enabled == true
      not pending_deletion
    }

    pending_deletion {
      input.resource_status == "PendingDeletion"
    }
  EOT

  type         = ["aws_kms_key"]
  provider_ids = ["aws_all"]
  severity     = 4
  category     = "Security"
  labels       = ["kms", "encryption", "compliance"]
  frameworks   = ["SOC2", "HIPAA"]
}

resource "firefly_governance_policy" "s3_secure_transport" {
  name        = "S3 Secure Transport Required"
  description = "Ensure S3 buckets enforce secure transport"

  code = <<-EOT
    package firefly

    firefly {
      match
    }

    match {
      policy := json.unmarshal(input.policy)
      stmt := policy.Statement[_]
      not secure_transport_enforced(stmt)
    }

    secure_transport_enforced(stmt) {
      stmt.Condition.StringEquals["aws:SecureTransport"] == "true"
    }
  EOT

  type         = ["aws_s3_bucket"]
  provider_ids = ["aws_all"]
  severity     = 4
  category     = "Security"
  labels       = ["s3", "transport", "security"]
}

resource "firefly_governance_policy" "lb_header_validation" {
  name        = "Load Balancer Invalid Header Handling"
  description = "Ensure load balancers drop invalid headers"

  code = <<-EOT
    package firefly

    firefly {
      not lb_drops_invalid_headers
    }

    lb_drops_invalid_headers {
      input.drop_invalid_header_fields == true
    }
  EOT

  type         = ["aws_lb", "aws_alb"]
  provider_ids = ["aws_all"]
  severity     = 3
  category     = "Security"
  labels       = ["load-balancer", "headers", "security"]
}

resource "firefly_governance_policy" "stopped_ec2_instances" {
  name        = "Stopped EC2 Instances"
  description = "Identify stopped EC2 instances for cost optimization"

  code = <<-EOT
    package firefly

    firefly {
      input.instance_state == "stopped"
    }
  EOT

  type         = ["aws_instance"]
  provider_ids = ["aws_all"]
  severity     = 1
  category     = "Optimization"
  labels       = ["ec2", "cost-optimization"]
}