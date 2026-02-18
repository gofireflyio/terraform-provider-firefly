package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBackupAndDrApplicationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with Daily schedule
			{
				Config: testAccBackupAndDrApplicationResourceConfig("test-daily-backup", "Daily"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "policy_name", "test-daily-backup"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "description", "Test backup policy"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "region", "us-east-1"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "provider_type", "aws"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "schedule.frequency", "Daily"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "schedule.hour", "2"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "schedule.minute", "30"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "backup_on_save", "true"),
					resource.TestCheckResourceAttrSet("firefly_backup_and_dr_application.test", "id"),
					resource.TestCheckResourceAttrSet("firefly_backup_and_dr_application.test", "status"),
					resource.TestCheckResourceAttrSet("firefly_backup_and_dr_application.test", "created_at"),
					resource.TestCheckResourceAttrSet("firefly_backup_and_dr_application.test", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firefly_backup_and_dr_application.test",
				ImportState:       true,
				ImportStateIdFunc: testAccBackupAndDrApplicationImportStateIdFunc("firefly_backup_and_dr_application.test"),
				ImportStateVerify: true,
			},
			// Update to Weekly schedule
			{
				Config: testAccBackupAndDrApplicationResourceConfig("test-daily-backup-updated", "Weekly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "policy_name", "test-daily-backup-updated"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "schedule.frequency", "Weekly"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "schedule.days_of_week.#", "2"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "schedule.days_of_week.0", "Monday"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.test", "schedule.days_of_week.1", "Friday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccBackupAndDrApplicationResource_WithScope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDrApplicationResourceConfigWithScope(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "policy_name", "scoped-backup"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "scope.#", "2"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "scope.0.type", "tags"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "scope.0.value.#", "2"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "scope.0.value.0", "Environment:Production"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "scope.0.value.1", "Backup:Required"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "scope.1.type", "asset_types"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.scoped", "scope.1.value.#", "2"),
				),
			},
		},
	})
}

func TestAccBackupAndDrApplicationResource_MonthlySchedule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Monthly with specific_day
			{
				Config: testAccBackupAndDrApplicationResourceConfigMonthlySpecificDay(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.monthly", "schedule.frequency", "Monthly"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.monthly", "schedule.monthly_schedule_type", "specific_day"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.monthly", "schedule.day_of_month", "15"),
				),
			},
			// Update to specific_weekday
			{
				Config: testAccBackupAndDrApplicationResourceConfigMonthlySpecificWeekday(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.monthly", "schedule.frequency", "Monthly"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.monthly", "schedule.monthly_schedule_type", "specific_weekday"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.monthly", "schedule.weekday_ordinal", "First"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.monthly", "schedule.weekday_name", "Sunday"),
				),
			},
		},
	})
}

func TestAccBackupAndDrApplicationResource_WithVCS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDrApplicationResourceConfigWithVCS(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.vcs", "policy_name", "vcs-backup"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.vcs", "vcs.project_id", "project-123"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.vcs", "vcs.vcs_integration_id", "github-456"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.vcs", "vcs.repo_id", "repo-789"),
				),
			},
		},
	})
}

func TestAccBackupAndDrApplicationResource_RestoreInstructions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupAndDrApplicationResourceConfigWithRestoreInstructions(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.instructions", "restore_instructions", "Step 1: Access console\nStep 2: Select snapshot\nStep 3: Restore"),
					resource.TestCheckResourceAttr("firefly_backup_and_dr_application.instructions", "notification_id", "slack-notification-123"),
				),
			},
		},
	})
}

// Helper function for import state ID
func testAccBackupAndDrApplicationImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		accountID := rs.Primary.Attributes["account_id"]
		policyID := rs.Primary.ID

		if accountID == "" || policyID == "" {
			return "", fmt.Errorf("account_id or id is not set")
		}

		return fmt.Sprintf("%s:%s", accountID, policyID), nil
	}
}

// Test configuration functions

func testAccBackupAndDrApplicationResourceConfig(name string, frequency string) string {
	if frequency == "Weekly" {
		return fmt.Sprintf(`
resource "firefly_backup_and_dr_application" "test" {
  account_id     = "test-account-id"
  policy_name    = %[1]q
  description    = "Test backup policy"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency    = "Weekly"
    days_of_week = ["Monday", "Friday"]
    hour         = 2
    minute       = 30
  }

  backup_on_save = true
}
`, name)
	}

	return fmt.Sprintf(`
resource "firefly_backup_and_dr_application" "test" {
  account_id     = "test-account-id"
  policy_name    = %[1]q
  description    = "Test backup policy"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 30
  }

  backup_on_save = true
}
`, name)
}

func testAccBackupAndDrApplicationResourceConfigWithScope() string {
	return `
resource "firefly_backup_and_dr_application" "scoped" {
  account_id     = "test-account-id"
  policy_name    = "scoped-backup"
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
    value = ["Environment:Production", "Backup:Required"]
  }

  scope {
    type  = "asset_types"
    value = ["aws_instance", "aws_db_instance"]
  }

  backup_on_save = true
}
`
}

func testAccBackupAndDrApplicationResourceConfigMonthlySpecificDay() string {
	return `
resource "firefly_backup_and_dr_application" "monthly" {
  account_id     = "test-account-id"
  policy_name    = "monthly-backup"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency             = "Monthly"
    monthly_schedule_type = "specific_day"
    day_of_month          = 15
    hour                  = 3
    minute                = 0
  }

  backup_on_save = true
}
`
}

func testAccBackupAndDrApplicationResourceConfigMonthlySpecificWeekday() string {
	return `
resource "firefly_backup_and_dr_application" "monthly" {
  account_id     = "test-account-id"
  policy_name    = "monthly-backup"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency             = "Monthly"
    monthly_schedule_type = "specific_weekday"
    weekday_ordinal       = "First"
    weekday_name          = "Sunday"
    hour                  = 3
    minute                = 0
  }

  backup_on_save = true
}
`
}

func testAccBackupAndDrApplicationResourceConfigWithVCS() string {
	return `
resource "firefly_backup_and_dr_application" "vcs" {
  account_id     = "test-account-id"
  policy_name    = "vcs-backup"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 0
  }

  vcs {
    project_id         = "project-123"
    vcs_integration_id = "github-456"
    repo_id            = "repo-789"
  }

  backup_on_save = true
}
`
}

func testAccBackupAndDrApplicationResourceConfigWithRestoreInstructions() string {
	return `
resource "firefly_backup_and_dr_application" "instructions" {
  account_id     = "test-account-id"
  policy_name    = "documented-backup"
  integration_id = "test-integration-id"
  region         = "us-east-1"
  provider_type  = "aws"

  schedule {
    frequency = "Daily"
    hour      = 2
    minute    = 0
  }

  restore_instructions = "Step 1: Access console\nStep 2: Select snapshot\nStep 3: Restore"
  notification_id      = "slack-notification-123"
  backup_on_save       = true
}
`
}
