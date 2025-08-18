package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig("test-project"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project.test", "name", "test-project"),
					resource.TestCheckResourceAttr("firefly_project.test", "description", "Test project description"),
					resource.TestCheckResourceAttrSet("firefly_project.test", "id"),
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
				Config: testAccProjectResourceConfig("test-project-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project.test", "name", "test-project-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "firefly_project" "test" {
  name        = %[1]q
  description = "Test project description"
  labels      = ["test", "terraform"]
}
`, name)
}