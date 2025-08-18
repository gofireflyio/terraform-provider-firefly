package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVariableSetResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccVariableSetResourceConfig("test-varset", "Test variable set description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_variable_set.test", "name", "test-varset"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "description", "Test variable set description"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "labels.0", "test"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "labels.1", "terraform"),
					resource.TestCheckResourceAttrSet("firefly_variable_set.test", "id"),
					resource.TestCheckResourceAttrSet("firefly_variable_set.test", "version"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firefly_variable_set.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccVariableSetResourceConfig("test-varset-updated", "Updated variable set description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_variable_set.test", "name", "test-varset-updated"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "description", "Updated variable set description"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVariableSetResource_withVariables(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableSetResourceWithVariablesConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_variable_set.test", "name", "test-varset-with-vars"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "variables.#", "3"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "variables.0.key", "AWS_REGION"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "variables.0.value", "us-west-2"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "variables.0.sensitivity", "string"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "variables.0.destination", "env"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "variables.1.key", "AWS_ACCESS_KEY_ID"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "variables.1.sensitivity", "secret"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "variables.2.key", "TF_VAR_region"),
					resource.TestCheckResourceAttr("firefly_variable_set.test", "variables.2.destination", "iac"),
				),
			},
		},
	})
}

func TestAccVariableSetResource_withInheritance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableSetResourceWithInheritanceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_variable_set.parent", "name", "parent-varset"),
					resource.TestCheckResourceAttr("firefly_variable_set.parent", "variables.#", "1"),
					resource.TestCheckResourceAttr("firefly_variable_set.child", "name", "child-varset"),
					resource.TestCheckResourceAttr("firefly_variable_set.child", "parents.#", "1"),
					resource.TestCheckResourceAttrPair("firefly_variable_set.child", "parents.0", "firefly_variable_set.parent", "id"),
				),
			},
		},
	})
}

func testAccVariableSetResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "firefly_variable_set" "test" {
  name        = %[1]q
  description = %[2]q
  labels      = ["test", "terraform"]
}
`, name, description)
}

func testAccVariableSetResourceWithVariablesConfig() string {
	return `
resource "firefly_variable_set" "test" {
  name        = "test-varset-with-vars"
  description = "Test variable set with variables"
  labels      = ["test", "variables"]
  
  variables {
    key         = "AWS_REGION"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "env"
  }
  
  variables {
    key         = "AWS_ACCESS_KEY_ID"
    value       = "test-access-key"
    sensitivity = "secret"
    destination = "env"
  }
  
  variables {
    key         = "TF_VAR_region"
    value       = "us-west-2"
    sensitivity = "string"
    destination = "iac"
  }
}
`
}

func testAccVariableSetResourceWithInheritanceConfig() string {
	return `
resource "firefly_variable_set" "parent" {
  name        = "parent-varset"
  description = "Parent variable set"
  labels      = ["parent", "base"]
  
  variables {
    key         = "COMPANY_NAME"
    value       = "ACME Corp"
    sensitivity = "string"
    destination = "env"
  }
}

resource "firefly_variable_set" "child" {
  name        = "child-varset"
  description = "Child variable set inheriting from parent"
  labels      = ["child", "derived"]
  parents     = [firefly_variable_set.parent.id]
  
  variables {
    key         = "SERVICE_NAME"
    value       = "api-service"
    sensitivity = "string"
    destination = "env"
  }
}
`
}