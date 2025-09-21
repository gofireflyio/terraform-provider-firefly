package provider

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
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
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "severity", "low"),
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
					resource.TestCheckResourceAttr("firefly_governance_policy.test", "severity", "medium"),
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
		severity = "medium"
	} else {
		severity = "low"
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

// Unit tests for governance policy operations
func TestMapModelToGovernancePolicy_SmartEncoding(t *testing.T) {
	// Test with plain text Rego code
	t.Run("Plain text Rego code", func(t *testing.T) {
		regoCode := `
firefly {
    input.instance_state == "stopped"
}
`
		model := &GovernancePolicyResourceModel{
			Name:        types.StringValue("Test Policy"),
			Code:        types.StringValue(regoCode),
			Type:        types.ListValueMust(types.StringType, []attr.Value{types.StringValue("EC2")}),
			ProviderIDs: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("aws")}),
			Severity:    types.StringValue("low"),
		}

		policy, err := mapModelToGovernancePolicy(model)
		if err != nil {
			t.Fatalf("mapModelToGovernancePolicy failed: %v", err)
		}

		// Verify that the code was base64 encoded
		decodedCode, err := base64.StdEncoding.DecodeString(policy.Code)
		if err != nil {
			t.Fatalf("Policy code is not valid base64: %v", err)
		}

		if string(decodedCode) != regoCode {
			t.Errorf("Decoded code doesn't match original. Got: %s, Expected: %s", string(decodedCode), regoCode)
		}

		// Verify severity conversion
		if policy.Severity != 3 {
			t.Errorf("Expected severity 3 (low), got %d", policy.Severity)
		}
	})

	// Test with already base64 encoded code
	t.Run("Base64 encoded code", func(t *testing.T) {
		regoCode := `
firefly {
    input.instance_state == "stopped"
}
`
		encodedCode := base64.StdEncoding.EncodeToString([]byte(regoCode))

		model := &GovernancePolicyResourceModel{
			Name:        types.StringValue("Test Policy"),
			Code:        types.StringValue(encodedCode),
			Type:        types.ListValueMust(types.StringType, []attr.Value{types.StringValue("EC2")}),
			ProviderIDs: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("aws")}),
			Severity:    types.StringValue("high"),
		}

		policy, err := mapModelToGovernancePolicy(model)
		if err != nil {
			t.Fatalf("mapModelToGovernancePolicy failed: %v", err)
		}

		// Verify that the code remains the same (already base64)
		if policy.Code != encodedCode {
			t.Errorf("Base64 code was modified. Got: %s, Expected: %s", policy.Code, encodedCode)
		}

		// Verify severity conversion
		if policy.Severity != 5 {
			t.Errorf("Expected severity 5 (high), got %d", policy.Severity)
		}
	})
}

func TestMapGovernancePolicyToModel(t *testing.T) {
	// Test with base64 encoded code from API
	t.Run("Base64 encoded API response", func(t *testing.T) {
		regoCode := `
firefly {
    input.instance_state == "stopped"
}
`
		encodedCode := base64.StdEncoding.EncodeToString([]byte(regoCode))

		policy := &client.GovernancePolicy{
			ID:          "test-policy-id",
			Name:        "Test Policy",
			Description: "Test description",
			Code:        encodedCode,
			Type:        []string{"EC2"},
			ProviderIDs: []string{"aws"},
			Labels:      []string{"security"},
			Severity:    4,
			Category:    "security",
			Frameworks:  []string{"SOC2"},
		}

		model := &GovernancePolicyResourceModel{}
		err := mapGovernancePolicyToModel(policy, model)
		if err != nil {
			t.Fatalf("mapGovernancePolicyToModel failed: %v", err)
		}

		// Verify that the code was decoded back to plain text
		if model.Code.ValueString() != regoCode {
			t.Errorf("Code was not decoded properly. Got: %s, Expected: %s", model.Code.ValueString(), regoCode)
		}

		// Verify severity conversion
		if model.Severity.ValueString() != "medium" {
			t.Errorf("Expected severity 'medium', got '%s'", model.Severity.ValueString())
		}

		// Verify other fields
		if model.ID.ValueString() != "test-policy-id" {
			t.Errorf("Expected ID 'test-policy-id', got '%s'", model.ID.ValueString())
		}

		if model.Name.ValueString() != "Test Policy" {
			t.Errorf("Expected name 'Test Policy', got '%s'", model.Name.ValueString())
		}
	})
}

func TestIsLikelyBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid base64 string",
			input:    "CgpmaXJlZmx5IHsKICAgIGlucHV0Lmluc3RhbmNlX3N0YXRlID09ICJzdG9wcGVkIgp9Cgo=",
			expected: true,
		},
		{
			name:     "Plain Rego code with package",
			input:    "package firefly\n\nfirefly { input.state == \"stopped\" }",
			expected: false,
		},
		{
			name:     "Plain Rego code with firefly keyword",
			input:    "firefly { input.state == \"stopped\" }",
			expected: false,
		},
		{
			name:     "Plain Rego code with braces",
			input:    "{ input.state == \"stopped\" }",
			expected: false,
		},
		{
			name:     "Invalid base64 characters",
			input:    "invalid@base64!string",
			expected: false,
		},
		{
			name:     "Wrong length for base64",
			input:    "abc",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: true, // Empty string is technically valid base64
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isLikelyBase64(test.input)
			if result != test.expected {
				t.Errorf("isLikelyBase64(%q) = %v, expected %v", test.input, result, test.expected)
			}
		})
	}
}