package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVariableSetsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableSetsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_variable_sets.all", "variable_sets.#"),
				),
			},
		},
	})
}

func TestAccVariableSetsDataSource_withSearch(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableSetsDataSourceWithSearchConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.firefly_variable_sets.filtered", "search_query", "aws"),
					resource.TestCheckResourceAttrSet("data.firefly_variable_sets.filtered", "variable_sets.#"),
				),
			},
		},
	})
}

func TestAccVariableSetsDataSource_withCreatedVariableSet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableSetsDataSourceWithCreatedVariableSetConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_variable_set.test", "name", "datasource-test-varset"),
					resource.TestCheckResourceAttrSet("data.firefly_variable_sets.search", "variable_sets.#"),
					// Check that variable sets have expected attributes
					resource.TestCheckResourceAttrSet("data.firefly_variable_sets.search", "variable_sets.0.id"),
					resource.TestCheckResourceAttrSet("data.firefly_variable_sets.search", "variable_sets.0.name"),
					resource.TestCheckResourceAttrSet("data.firefly_variable_sets.search", "variable_sets.0.version"),
				),
			},
		},
	})
}

func testAccVariableSetsDataSourceConfig() string {
	return `
data "firefly_variable_sets" "all" {}
`
}

func testAccVariableSetsDataSourceWithSearchConfig() string {
	return `
data "firefly_variable_sets" "filtered" {
  search_query = "aws"
}
`
}

func testAccVariableSetsDataSourceWithCreatedVariableSetConfig() string {
	return `
resource "firefly_variable_set" "test" {
  name        = "datasource-test-varset"
  description = "Variable set created for data source testing"
  labels      = ["datasource", "test"]
  
  variables {
    key         = "TEST_VAR"
    value       = "test-value"
    sensitivity = "string"
    destination = "env"
  }
}

data "firefly_variable_sets" "search" {
  search_query = "datasource"
  depends_on   = [firefly_variable_set.test]
}
`
}