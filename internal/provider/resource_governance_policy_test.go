package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "category", "Optimization"),
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "severity", "3"),
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
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccGovernancePolicyResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "firefly_governance_policy" "test" {
  name        = %[1]q
  description = "Test governance policy"
  code        = <<-EOT
    package firefly
    import rego.v1
    firefly {
      input.test == true
    }
  EOT
  type         = ["aws_instance"]
  provider_ids = ["aws_all"]
  severity     = 3
  category     = "Optimization"
  labels       = ["test"]
}
`, name)
}

func TestAccGovernancePolicyResource_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with all optional fields
			{
				Config: testAccGovernancePolicyResourceConfigComplete(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firefly_governance_policy.complete", "name", "complete-policy"),
					resource.TestCheckResourceAttr("firefly_governance_policy.complete", "description", "Complete test policy"),
					resource.TestCheckResourceAttr("firefly_governance_policy.complete", "category", "Security"),
					resource.TestCheckResourceAttr("firefly_governance_policy.complete", "severity", "4"),
					resource.TestCheckResourceAttr("firefly_governance_policy.complete", "labels.#", "2"),
					resource.TestCheckResourceAttr("firefly_governance_policy.complete", "frameworks.#", "2"),
				),
			},
		},
	})
}

func testAccGovernancePolicyResourceConfigComplete() string {
	return `
resource "firefly_governance_policy" "complete" {
  name        = "complete-policy"
  description = "Complete test policy"
  code        = <<-EOT
    package firefly
    import rego.v1
    firefly {
      input.encrypted == false
    }
  EOT
  type         = ["aws_ebs_volume", "aws_instance"]
  provider_ids = ["aws_all"]
  severity     = 4
  category     = "Security"
  labels       = ["test", "security"]
  frameworks   = ["SOC2", "HIPAA"]
}
`
}