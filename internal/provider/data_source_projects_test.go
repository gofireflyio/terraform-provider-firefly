package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_projects.all", "projects.#"),
				),
			},
		},
	})
}

func TestAccProjectsDataSource_withSearch(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectsDataSourceWithSearchConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.firefly_projects.filtered", "search_query", "test"),
					resource.TestCheckResourceAttrSet("data.firefly_projects.filtered", "projects.#"),
				),
			},
		},
	})
}

func TestAccProjectsDataSource_withCreatedProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectsDataSourceWithCreatedProjectConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_project.test", "name", "datasource-test-project"),
					resource.TestCheckResourceAttrSet("data.firefly_projects.search", "projects.#"),
					// The created project should be findable in the data source
					resource.TestCheckResourceAttrSet("data.firefly_projects.search", "projects.0.id"),
					resource.TestCheckResourceAttrSet("data.firefly_projects.search", "projects.0.name"),
					resource.TestCheckResourceAttrSet("data.firefly_projects.search", "projects.0.account_id"),
				),
			},
		},
	})
}

func testAccProjectsDataSourceConfig() string {
	return `
data "firefly_projects" "all" {}
`
}

func testAccProjectsDataSourceWithSearchConfig() string {
	return `
data "firefly_projects" "filtered" {
  search_query = "test"
}
`
}

func testAccProjectsDataSourceWithCreatedProjectConfig() string {
	return `
resource "firefly_project" "test" {
  name        = "datasource-test-project"
  description = "Project created for data source testing"
  labels      = ["datasource", "test"]
}

data "firefly_projects" "search" {
  search_query = "datasource"
  depends_on   = [firefly_project.test]
}
`
}