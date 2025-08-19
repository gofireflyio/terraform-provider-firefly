package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccProjectMembershipResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectMembershipResourceConfig("test-user@example.com", "user123", "member"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project_membership.test", "email", "test-user@example.com"),
					resource.TestCheckResourceAttr("firefly_project_membership.test", "user_id", "user123"),
					resource.TestCheckResourceAttr("firefly_project_membership.test", "role", "member"),
					resource.TestCheckResourceAttrSet("firefly_project_membership.test", "project_id"),
					resource.TestCheckResourceAttrSet("firefly_project_membership.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firefly_project_membership.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccProjectMembershipImportStateIdFunc("firefly_project_membership.test"),
			},
			// Update and Read testing
			{
				Config: testAccProjectMembershipResourceConfig("test-user@example.com", "user123", "admin"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project_membership.test", "email", "test-user@example.com"),
					resource.TestCheckResourceAttr("firefly_project_membership.test", "user_id", "user123"),
					resource.TestCheckResourceAttr("firefly_project_membership.test", "role", "admin"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectMembershipResourceConfig(email, userID, role string) string {
	return fmt.Sprintf(`
resource "firefly_workflows_project" "test" {
  name        = "test-project-membership"
  description = "Test project for membership testing"
}

resource "firefly_project_membership" "test" {
  project_id = firefly_workflows_project.test.id
  user_id    = "%s"
  email      = "%s"
  role       = "%s"
}
`, userID, email, role)
}

func testAccProjectMembershipImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		userID := rs.Primary.Attributes["user_id"]

		if projectID == "" || userID == "" {
			return "", fmt.Errorf("project_id or user_id is not set")
		}

		return fmt.Sprintf("%s:%s", projectID, userID), nil
	}
}