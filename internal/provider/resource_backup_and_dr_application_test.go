package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBackupAndDRApplicationResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with daily schedule
			{
				Config: testAccBackupAndDRApplicationConfig("tf-test-backup-policy", "Test backup policy", "Daily"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "policy_name", "tf-test-backup-policy"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "description", "Test backup policy"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "status", "Active"),
					resource.TestCheckResourceAttrSet("firefly_backup_and_dr_application.test", "id"),
					resource.TestCheckResourceAttrSet("firefly_backup_and_dr_application.test", "created_at"),
				),
			},
			// ImportState
			{
				ResourceName:            "firefly_backup_and_dr_application.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_on_save"},
			},
			// Update name and description
			{
				Config: testAccBackupAndDRApplicationConfig("tf-test-backup-policy-updated", "Updated description", "Daily"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "policy_name", "tf-test-backup-policy-updated"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "description", "Updated description"),
				),
			},
			// Delete automatically occurs
		},
	})
}

func TestAccBackupAndDRApplicationResource_WeeklySchedule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDRApplicationWeeklyConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.weekly", "policy_name", "tf-test-weekly-backup"),
					resource.TestCheckResourceAttrSet("firefly_backup_and_dr_application.weekly", "id"),
				),
			},
		},
	})
}

func TestAccBackupAndDRApplicationResource_WithScope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDRApplicationScopeConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "policy_name", "tf-test-scoped-backup"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "scope.#", "1"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "scope.0.type", "tags"),
					resource.TestCheckResourceAttrSet("firefly_backup_and_dr_application.scoped", "id"),
				),
			},
		},
	})
}

func TestAccBackupAndDRApplicationResource_StatusToggle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create as Active
			{
				Config: testAccBackupAndDRApplicationStatusConfig("Active"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.status_test", "status", "Active"),
				),
			},
			// Toggle to Inactive
			{
				Config: testAccBackupAndDRApplicationStatusConfig("Inactive"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.status_test", "status", "Inactive"),
				),
			},
		},
	})
}

func testAccBackupAndDRApplicationConfig(name, description, frequency string) string {
	return fmt.Sprintf(`
resource "firefly_backup_and_dr_application" "test" {
  policy_name    = %[1]q
  description    = %[2]q
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = %[3]q
    hour      = 10
    minute    = 30
  }
}
`, name, description, frequency)
}

func testAccBackupAndDRApplicationWeeklyConfig() string {
	return `
resource "firefly_backup_and_dr_application" "weekly" {
  policy_name    = "tf-test-weekly-backup"
  description    = "Weekly backup test"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency    = "Weekly"
    hour         = 14
    minute       = 0
    days_of_week = ["Monday", "Friday"]
  }
}
`
}

func testAccBackupAndDRApplicationScopeConfig() string {
	return `
resource "firefly_backup_and_dr_application" "scoped" {
  policy_name    = "tf-test-scoped-backup"
  description    = "Scoped backup test"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 0
  }

  scope {
    type  = "tags"
    value = ["env:prod", "backup:true"]
  }
}
`
}

func testAccBackupAndDRApplicationStatusConfig(status string) string {
	return fmt.Sprintf(`
resource "firefly_backup_and_dr_application" "status_test" {
  policy_name    = "tf-test-status-toggle"
  description    = "Status toggle test"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"
  status         = %[1]q

  schedule {
    frequency = "Daily"
    hour      = 6
    minute    = 0
  }
}
`, status)
}
