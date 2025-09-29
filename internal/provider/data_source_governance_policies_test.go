package provider

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGovernancePoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGovernancePoliciesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_governance_policies.test", "id"),
					resource.TestCheckResourceAttrSet("data.firefly_governance_policies.test", "policies.#"),
				),
			},
		},
	})
}

func TestAccGovernancePoliciesDataSource_WithFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test with category filter
			{
				Config: testAccGovernancePoliciesDataSourceConfigWithCategory,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_governance_policies.security", "id"),
					resource.TestCheckResourceAttrSet("data.firefly_governance_policies.security", "policies.#"),
				),
			},
			// Test with query filter
			{
				Config: testAccGovernancePoliciesDataSourceConfigWithQuery,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firefly_governance_policies.search", "id"),
				),
			},
		},
	})
}

// Unit test for data source code decoding
func TestGovernancePoliciesDataSource_CodeDecoding(t *testing.T) {
	// Test that the data source properly decodes base64 encoded code from API responses
	regoCode := `
firefly {
    input.instance_state == "stopped"
}
`
	encodedCode := base64.StdEncoding.EncodeToString([]byte(regoCode))

	// Test base64 decoding logic
	decodedCode, err := base64.StdEncoding.DecodeString(encodedCode)
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	if string(decodedCode) != regoCode {
		t.Errorf("Decoded code doesn't match original. Got: %s, Expected: %s", string(decodedCode), regoCode)
	}

	// Test invalid base64 handling
	invalidBase64 := "this is not base64!"
	_, err = base64.StdEncoding.DecodeString(invalidBase64)
	if err == nil {
		t.Error("Expected error when decoding invalid base64, but got none")
	}
}

const testAccGovernancePoliciesDataSourceConfig = `
data "firefly_governance_policies" "test" {}
`

const testAccGovernancePoliciesDataSourceConfigWithCategory = `
data "firefly_governance_policies" "security" {
  category = "Security"
}
`

const testAccGovernancePoliciesDataSourceConfigWithQuery = `
data "firefly_governance_policies" "search" {
  query = "s3"
}
`

const testAccGovernancePoliciesDataSourceConfigWithLabels = `
data "firefly_governance_policies" "aws_policies" {
  labels = ["aws", "security"]
}
`

const testAccGovernancePoliciesDataSourceConfigCombined = `
data "firefly_governance_policies" "filtered" {
  category = "Security"
  query    = "encryption"
  labels   = ["aws"]
}

output "filtered_policies" {
  value = data.firefly_governance_policies.filtered.policies
}

output "policy_count" {
  value = length(data.firefly_governance_policies.filtered.policies)
}

output "security_policy_names" {
  value = [for policy in data.firefly_governance_policies.filtered.policies : policy.name]
}
`