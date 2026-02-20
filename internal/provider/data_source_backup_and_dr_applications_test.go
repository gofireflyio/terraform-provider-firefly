package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBackupAndDrApplicationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a resource first, then read it with data source
			{
				Config: testAccBackupAndDrApplicationsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the resource was created
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "policy_name", "data-source-test-policy"),
					// Check the data source can read it
					resource.TestCheckResourceAttrSet("data.firefly_backup_and_dr_applications.test", "id"),
					resource.TestCheckResourceAttrSet("data.firefly_backup_and_dr_applications.test", "policies.#"),
				),
			},
		},
	})
}

func TestAccBackupAndDrApplicationsDataSource_FilterByStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDrApplicationsDataSourceConfigFilterByStatus(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_backup_and_dr_applications.active", "id"),
					resource.TestCheckResourceAttr("data.firefly_backup_and_dr_applications.active", "status", "Active"),
				),
			},
		},
	})
}

func TestAccBackupAndDrApplicationsDataSource_FilterByRegion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDrApplicationsDataSourceConfigFilterByRegion(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_backup_and_dr_applications.regional", "id"),
					resource.TestCheckResourceAttr("data.firefly_backup_and_dr_applications.regional", "region", "us-east-1"),
					resource.TestCheckResourceAttr("data.firefly_backup_and_dr_applications.regional", "provider_type", "aws"),
				),
			},
		},
	})
}

func TestAccBackupAndDrApplicationsDataSource_MultipleFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDrApplicationsDataSourceConfigMultipleFilters(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_backup_and_dr_applications.filtered", "id"),
					resource.TestCheckResourceAttr("data.firefly_backup_and_dr_applications.filtered", "status", "Active"),
					resource.TestCheckResourceAttr("data.firefly_backup_and_dr_applications.filtered", "region", "us-west-2"),
					resource.TestCheckResourceAttr("data.firefly_backup_and_dr_applications.filtered", "provider_type", "aws"),
				),
			},
		},
	})
}

// Test configuration functions

func testAccBackupAndDrApplicationsDataSourceConfig() string {
	return `
resource "firefly_backup_and_dr_application" "test" {
  account_id     = "test-account-id"
  policy_name    = "data-source-test-policy"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 0
  }

  backup_on_save = true
}

data "firefly_backup_and_dr_applications" "test" {
  account_id = firefly_backup_and_dr_application.test.account_id

  depends_on = [firefly_backup_and_dr_application.test]
}
`
}

func testAccBackupAndDrApplicationsDataSourceConfigFilterByStatus() string {
	return `
resource "firefly_backup_and_dr_application" "active_test" {
  account_id     = "test-account-id"
  policy_name    = "active-policy"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 0
  }

  backup_on_save = true
}

data "firefly_backup_and_dr_applications" "active" {
  account_id = "test-account-id"
  status     = "Active"

  depends_on = [firefly_backup_and_dr_application.active_test]
}
`
}

func testAccBackupAndDrApplicationsDataSourceConfigFilterByRegion() string {
	return `
resource "firefly_backup_and_dr_application" "regional_test" {
  account_id     = "test-account-id"
  policy_name    = "regional-policy"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 0
  }

  backup_on_save = true
}

data "firefly_backup_and_dr_applications" "regional" {
  account_id    = "test-account-id"
  region        = "us-east-1"
  provider_type = "aws"

  depends_on = [firefly_backup_and_dr_application.regional_test]
}
`
}

func testAccBackupAndDrApplicationsDataSourceConfigMultipleFilters() string {
	return `
resource "firefly_backup_and_dr_application" "multi_filter_test" {
  account_id     = "test-account-id"
  policy_name    = "multi-filter-policy"
  integration_id = "test-integration-id"
  region         = "us-west-2"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 3
    minute    = 0
  }

  backup_on_save = true
}

data "firefly_backup_and_dr_applications" "filtered" {
  account_id    = "test-account-id"
  status        = "Active"
  region        = "us-west-2"
  provider_type = "aws"

  depends_on = [firefly_backup_and_dr_application.multi_filter_test]
}
`
}
