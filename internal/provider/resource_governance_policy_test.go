package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGovernancePolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGovernancePolicyResourceConfig("test-policy"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "name", "test-policy"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "description", "Test governance policy"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "category", "Misconfiguration"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "severity", "warning"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "type.#", "1"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "type.0", "aws_s3_bucket"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "provider_ids.#", "1"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "provider_ids.0", "aws_all"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "labels.#", "2"),
					resource.TestCheckResourceAttrSet("firefly_governance_policy.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firefly_governance_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccGovernancePolicyResourceConfig("test-policy-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "name", "test-policy-updated"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "severity", "strict"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGovernancePolicyResource_BasicRegoCode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test with basic Rego code
			{
				Config: testAccGovernancePolicyResourceConfigBasicRego(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_governance_policy.basic", "name", "basic-rego-policy"),
					resource.TestCheckResourceAttrSet("firefly_governance_policy.basic", "code"),
					resource.TestCheckResourceAttrSet("firefly_governance_policy.basic", "id"),
				),
			},
		},
	})
}

func testAccGovernancePolicyResourceConfig(name string) string {
	var severity string
	if name == "test-policy-updated" {
		severity = "strict"
	} else {
		severity = "warning"
	}
	
	return fmt.Sprintf(`
resource "firefly_governance_policy" "test" {
  name        = %[1]q
  description = "Test governance policy"
  category    = "Misconfiguration"
  severity    = %[2]q
  
  code = <<-EOT
    firefly {
      match
    }
    
    match {
      input.public_access_block == false
    }
  EOT
  
  type         = ["aws_s3_bucket"]
  provider_ids = ["aws_all"]
  labels       = ["security", "compliance"]
  frameworks   = ["SOC2", "ISO27001"]
}
`, name, severity)
}

func testAccGovernancePolicyResourceConfigBasicRego() string {
	return `
resource "firefly_governance_policy" "basic" {
  name = "basic-rego-policy"
  
  code = <<-EOT
    firefly {
      match
    }
    
    match {
      true
    }
  EOT
  
  type         = ["aws_cloudwatch_event_target"]
  provider_ids = ["aws_all"]
}
`
}