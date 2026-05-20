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

# Test: Backup & DR Application with 24-hour backup frequency
resource "firefly_backup_and_dr_application" "daily_equivalent" {
  account_id       = "66169d5af4992fc0bab04510"
  application_name = "terraform-test-24h"
  integration_id   = "692ec8acce65b3dc46cfceb5"
  region           = "eu-west-1"
  provider_type    = "aws"
  frequency        = 24

  scope {
    type  = "tags"
    value = ["Demo: BackupDR"]
  }
}

# Test: Backup & DR Application with 8-hour backup frequency
resource "firefly_backup_and_dr_application" "frequent_backup" {
  account_id       = "66169d5af4992fc0bab04510"
  application_name = "terraform-test-8h"
  integration_id   = "692ec8acce65b3dc46cfceb5"
  region           = "eu-west-1"
  provider_type    = "aws"
  description      = "Frequent backup test"
  frequency        = 8

  scope {
    type  = "tags"
    value = ["Demo: BackupDR"]
  }
}
