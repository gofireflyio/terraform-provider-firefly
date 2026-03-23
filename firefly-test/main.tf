terraform {
  required_providers {
    firefly = {
      source = "registry.terraform.io/gofireflyio/firefly"
    }
  }
}

provider "firefly" {
  access_key = var.access_key
  secret_key = var.secret_key
  api_url    = "https://ext-api-env2.dev.firefly.ai"
}

variable "access_key" {
  type      = string
  sensitive = true
}

variable "secret_key" {
  type      = string
  sensitive = true
}

# Test: Backup & DR Application with Monthly last_day schedule
# This reproduces the cron_expression inconsistency bug fix
resource "firefly_backup_and_dr_application" "monthly_last_day" {
  account_id     = "66169d5af4992fc0bab04510"
  application_name = "terraform-test-monthly-last-day"
  integration_id = "692ec8acce65b3dc46cfceb5"
  region         = "eu-west-1"
  provider_type  = "aws"

  schedule {
    frequency             = "Monthly"
    monthly_schedule_type = "last_day"
  }

  scope {
    type  = "tags"
    value = ["terraform-test:true"]
  }
}

# Test: Backup & DR Application with Daily schedule
resource "firefly_backup_and_dr_application" "daily_backup" {
  account_id       = "66169d5af4992fc0bab04510"
  application_name = "terraform-test-daily"
  integration_id   = "692ec8acce65b3dc46cfceb5"
  region           = "eu-west-1"
  provider_type    = "aws"
  description      = "Daily backup test"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 30
  }

  scope {
    type  = "tags"
    value = ["Demo: BackupDR"]
  }
}

