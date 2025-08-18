package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVariableSetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableSetDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_variable_set.test", "name", "single-datasource-test-varset"),
					resource.TestCheckResourceAttrPair("data.firefly_variable_set.test", "id", "firefly_variable_set.test", "id"),
					resource.TestCheckResourceAttrPair("data.firefly_variable_set.test", "name", "firefly_variable_set.test", "name"),
					resource.TestCheckResourceAttrPair("data.firefly_variable_set.test", "description", "firefly_variable_set.test", "description"),
					resource.TestCheckResourceAttrPair("data.firefly_variable_set.test", "labels", "firefly_variable_set.test", "labels"),
					resource.TestCheckResourceAttrSet("data.firefly_variable_set.test", "version"),
				),
			},
		},
	})
}

func testAccVariableSetDataSourceConfig() string {
	return `
resource "firefly_variable_set" "test" {
  name        = "single-datasource-test-varset"
  description = "Single variable set for data source testing"
  labels      = ["single", "datasource", "test"]
  
  variables {
    key         = "TEST_VARIABLE"
    value       = "test-value"
    sensitivity = "string"
    destination = "env"
  }
}

data "firefly_variable_set" "test" {
  id = firefly_variable_set.test.id
}
`
}