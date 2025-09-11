package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGovernanceInsightResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGovernanceInsightResourceConfig("test-insight"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "name", "test-insight"),
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "description", "Test governance insight"),
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "category", "Misconfiguration"),
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "severity", "warning"),
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "type.#", "1"),
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "type.0", "aws_s3_bucket"),
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "provider_ids.#", "1"),
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "provider_ids.0", "aws_all"),
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "labels.#", "2"),
					resource.TestCheckResourceAttrSet("firefly_governance_insight.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firefly_governance_insight.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccGovernanceInsightResourceConfig("test-insight-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "name", "test-insight-updated"),
					resource.TestCheckResourceAttr("firefly_governance_insight.test", "severity", "strict"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGovernanceInsightResource_BasicRegoCode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test with basic Rego code
			{
				Config: testAccGovernanceInsightResourceConfigBasicRego(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_governance_insight.basic", "name", "basic-rego-insight"),
					resource.TestCheckResourceAttrSet("firefly_governance_insight.basic", "code"),
					resource.TestCheckResourceAttrSet("firefly_governance_insight.basic", "id"),
				),
			},
		},
	})
}

func testAccGovernanceInsightResourceConfig(name string) string {
	var severity string
	if name == "test-insight-updated" {
		severity = "strict"
	} else {
		severity = "warning"
	}
	
	return fmt.Sprintf(`
resource "firefly_governance_insight" "test" {
  name        = %[1]q
  description = "Test governance insight"
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

func testAccGovernanceInsightResourceConfigBasicRego() string {
	return `
resource "firefly_governance_insight" "basic" {
  name = "basic-rego-insight"
  
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