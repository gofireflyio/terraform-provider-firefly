package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig("test-project", "Test project description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project.test", "name", "test-project"),
					resource.TestCheckResourceAttr("firefly_project.test", "description", "Test project description"),
					resource.TestCheckResourceAttr("firefly_project.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("firefly_project.test", "labels.0", "test"),
					resource.TestCheckResourceAttr("firefly_project.test", "labels.1", "terraform"),
					resource.TestCheckResourceAttrSet("firefly_project.test", "id"),
					resource.TestCheckResourceAttrSet("firefly_project.test", "account_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firefly_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfig("test-project-updated", "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project.test", "name", "test-project-updated"),
					resource.TestCheckResourceAttr("firefly_project.test", "description", "Updated description"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProjectResource_withVariables(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceWithVariablesConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project.test", "name", "test-project-with-vars"),
					resource.TestCheckResourceAttr("firefly_project.test", "variables.#", "2"),
					resource.TestCheckResourceAttr("firefly_project.test", "variables.0.key", "ENVIRONMENT"),
					resource.TestCheckResourceAttr("firefly_project.test", "variables.0.value", "test"),
					resource.TestCheckResourceAttr("firefly_project.test", "variables.0.sensitivity", "string"),
					resource.TestCheckResourceAttr("firefly_project.test", "variables.0.destination", "env"),
					resource.TestCheckResourceAttr("firefly_project.test", "variables.1.key", "TF_LOG_LEVEL"),
					resource.TestCheckResourceAttr("firefly_project.test", "variables.1.sensitivity", "string"),
				),
			},
		},
	})
}

func TestAccProjectResource_withSchedule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceWithScheduleConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project.test", "name", "scheduled-project"),
					resource.TestCheckResourceAttr("firefly_project.test", "cron_execution_pattern", "0 2 * * *"),
				),
			},
		},
	})
}

func testAccProjectResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "firefly_project" "test" {
  name        = %[1]q
  description = %[2]q
  labels      = ["test", "terraform"]
}
`, name, description)
}

func testAccProjectResourceWithVariablesConfig() string {
	return `
resource "firefly_project" "test" {
  name        = "test-project-with-vars"
  description = "Test project with variables"
  labels      = ["test", "variables"]
  
  variables {
    key         = "ENVIRONMENT"
    value       = "test"
    sensitivity = "string"
    destination = "env"
  }
  
  variables {
    key         = "TF_LOG_LEVEL"
    value       = "INFO"
    sensitivity = "string"
    destination = "env"
  }
}
`
}

func testAccProjectResourceWithScheduleConfig() string {
	return `
resource "firefly_project" "test" {
  name                   = "scheduled-project"
  description            = "Project with scheduled execution"
  labels                 = ["scheduled", "test"]
  cron_execution_pattern = "0 2 * * *"
}
`
}