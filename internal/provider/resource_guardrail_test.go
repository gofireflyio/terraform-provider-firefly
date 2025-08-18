package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGuardrailResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGuardrailResourceConfig("test-guardrail", "cost", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_guardrail.test", "name", "test-guardrail"),
					resource.TestCheckResourceAttr("firefly_guardrail.test", "type", "cost"),
					resource.TestCheckResourceAttr("firefly_guardrail.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("firefly_guardrail.test", "severity", "2"),
					resource.TestCheckResourceAttrSet("firefly_guardrail.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firefly_guardrail.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccGuardrailResourceConfig("test-guardrail-updated", "cost", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_guardrail.test", "name", "test-guardrail-updated"),
					resource.TestCheckResourceAttr("firefly_guardrail.test", "severity", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGuardrailResource_costThreshold(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardrailResourceCostConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_guardrail.cost", "name", "Cost Threshold Guardrail"),
					resource.TestCheckResourceAttr("firefly_guardrail.cost", "type", "cost"),
					resource.TestCheckResourceAttr("firefly_guardrail.cost", "criteria.0.cost.0.threshold_amount", "1000"),
				),
			},
		},
	})
}

func TestAccGuardrailResource_tagPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardrailResourceTagConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_guardrail.tag", "name", "Tag Policy Guardrail"),
					resource.TestCheckResourceAttr("firefly_guardrail.tag", "type", "tag"),
					resource.TestCheckResourceAttr("firefly_guardrail.tag", "criteria.0.tag.0.tag_enforcement_mode", "requiredTags"),
					resource.TestCheckResourceAttr("firefly_guardrail.tag", "criteria.0.tag.0.required_tags.#", "2"),
					resource.TestCheckResourceAttr("firefly_guardrail.tag", "criteria.0.tag.0.required_tags.0", "Environment"),
					resource.TestCheckResourceAttr("firefly_guardrail.tag", "criteria.0.tag.0.required_tags.1", "Owner"),
				),
			},
		},
	})
}

func TestAccGuardrailResource_withScope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardrailResourceWithScopeConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_guardrail.scoped", "name", "Scoped Guardrail"),
					resource.TestCheckResourceAttr("firefly_guardrail.scoped", "scope.0.workspaces.0.include.#", "1"),
					resource.TestCheckResourceAttr("firefly_guardrail.scoped", "scope.0.workspaces.0.include.0", "production-*"),
					resource.TestCheckResourceAttr("firefly_guardrail.scoped", "scope.0.labels.0.include.#", "2"),
					resource.TestCheckResourceAttr("firefly_guardrail.scoped", "scope.0.labels.0.include.0", "critical"),
					resource.TestCheckResourceAttr("firefly_guardrail.scoped", "scope.0.labels.0.include.1", "production"),
				),
			},
		},
	})
}

func testAccGuardrailResourceConfig(name, guardrailType string, severity int) string {
	return fmt.Sprintf(`
resource "firefly_guardrail" "test" {
  name       = %[1]q
  type       = %[2]q
  is_enabled = true
  severity   = %[3]d
}
`, name, guardrailType, severity)
}

func testAccGuardrailResourceCostConfig() string {
	return `
resource "firefly_guardrail" "cost" {
  name       = "Cost Threshold Guardrail"
  type       = "cost"
  is_enabled = true
  severity   = 2
  
  criteria {
    cost {
      threshold_amount = 1000
    }
  }
}
`
}

func testAccGuardrailResourceTagConfig() string {
	return `
resource "firefly_guardrail" "tag" {
  name       = "Tag Policy Guardrail"
  type       = "tag"
  is_enabled = true
  severity   = 3
  
  criteria {
    tag {
      tag_enforcement_mode = "requiredTags"
      required_tags       = ["Environment", "Owner"]
    }
  }
}
`
}

func testAccGuardrailResourceWithScopeConfig() string {
	return `
resource "firefly_guardrail" "scoped" {
  name       = "Scoped Guardrail"
  type       = "cost"
  is_enabled = true
  severity   = 1
  
  scope {
    workspaces {
      include = ["production-*"]
    }
    
    labels {
      include = ["critical", "production"]
    }
  }
  
  criteria {
    cost {
      threshold_amount = 500
    }
  }
}
`
}