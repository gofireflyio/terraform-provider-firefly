package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_workflows_project.test", "name", "single-datasource-test-project"),
					resource.TestCheckResourceAttrPair("data.firefly_workflows_project.test", "id", "firefly_workflows_project.test", "id"),
					resource.TestCheckResourceAttrPair("data.firefly_workflows_project.test", "name", "firefly_workflows_project.test", "name"),
					resource.TestCheckResourceAttrPair("data.firefly_workflows_project.test", "description", "firefly_workflows_project.test", "description"),
					resource.TestCheckResourceAttrPair("data.firefly_workflows_project.test", "labels", "firefly_workflows_project.test", "labels"),
					resource.TestCheckResourceAttrSet("data.firefly_workflows_project.test", "account_id"),
					resource.TestCheckResourceAttrSet("data.firefly_workflows_project.test", "members_count"),
					resource.TestCheckResourceAttrSet("data.firefly_workflows_project.test", "workspace_count"),
				),
			},
		},
	})
}

func testAccProjectDataSourceConfig() string {
	return `
resource "firefly_project" "test" {
  name        = "single-datasource-test-project"
  description = "Single project for data source testing"
  labels      = ["single", "datasource", "test"]
}

data "firefly_workflows_project" "test" {
  id = firefly_workflows_project.test.id
}
`
}