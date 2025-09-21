package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapGovernancePolicyToModel maps API response to Terraform model
func mapGovernancePolicyToModel(policy *client.GovernancePolicy, model *GovernancePolicyResourceModel) error {
	model.ID = types.StringValue(policy.ID)
	model.Name = types.StringValue(policy.Name)
	
	if policy.Description != "" {
		model.Description = types.StringValue(policy.Description)
	} else {
		model.Description = types.StringNull()
	}
	
	// Decode the base64 encoded Rego code from the API
	decodedCode, err := base64.StdEncoding.DecodeString(policy.Code)
	if err != nil {
		// If decoding fails, assume it's already plain text (for backward compatibility)
		model.Code = types.StringValue(policy.Code)
	} else {
		model.Code = types.StringValue(string(decodedCode))
	}
	
	// Convert Type array to list
	typeList, diags := types.ListValueFrom(context.Background(), types.StringType, policy.Type)
	if diags.HasError() {
		return fmt.Errorf("error converting type list: %v", diags)
	}
	model.Type = typeList
	
	// Convert ProviderIDs array to list
	providerList, diags := types.ListValueFrom(context.Background(), types.StringType, policy.ProviderIDs)
	if diags.HasError() {
		return fmt.Errorf("error converting provider IDs list: %v", diags)
	}
	model.ProviderIDs = providerList
	
	// Convert Labels array to list
	if len(policy.Labels) > 0 {
		labelsList, diags := types.ListValueFrom(context.Background(), types.StringType, policy.Labels)
		if diags.HasError() {
			return fmt.Errorf("error converting labels list: %v", diags)
		}
		model.Labels = labelsList
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
	
	// Convert severity integer to string
	if policy.Severity > 0 {
		model.Severity = types.StringValue(client.SeverityToString(policy.Severity))
	} else {
		model.Severity = types.StringNull()
	}
	
	if policy.Category != "" {
		model.Category = types.StringValue(policy.Category)
	} else {
		model.Category = types.StringNull()
	}
	
	// Convert Frameworks array to list
	if len(policy.Frameworks) > 0 {
		frameworksList, diags := types.ListValueFrom(context.Background(), types.StringType, policy.Frameworks)
		if diags.HasError() {
			return fmt.Errorf("error converting frameworks list: %v", diags)
		}
		model.Frameworks = frameworksList
	} else {
		model.Frameworks = types.ListNull(types.StringType)
	}
	
	return nil
}

// mapModelToGovernancePolicy maps Terraform model to API request
func mapModelToGovernancePolicy(model *GovernancePolicyResourceModel) (*client.GovernancePolicy, error) {
	// Smart encoding: encode if text, validate if already base64
	codeValue := model.Code.ValueString()
	var encodedCode string
	
	// Check if the code is already base64 encoded
	if _, err := base64.StdEncoding.DecodeString(codeValue); err == nil && isLikelyBase64(codeValue) {
		// Already valid base64, use as-is
		encodedCode = codeValue
	} else {
		// Plain text Rego code, encode it
		encodedCode = base64.StdEncoding.EncodeToString([]byte(codeValue))
	}
	
	policy := &client.GovernancePolicy{
		Name: model.Name.ValueString(),
		Code: encodedCode,
	}
	
	if !model.ID.IsNull() && !model.ID.IsUnknown() {
		policy.ID = model.ID.ValueString()
	}
	
	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		policy.Description = model.Description.ValueString()
	}
	
	// Convert Type list to array
	var typeArray []string
	diags := model.Type.ElementsAs(context.Background(), &typeArray, false)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting type list: %v", diags)
	}
	policy.Type = typeArray
	
	// Convert ProviderIDs list to array
	var providerArray []string
	diags = model.ProviderIDs.ElementsAs(context.Background(), &providerArray, false)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting provider IDs list: %v", diags)
	}
	policy.ProviderIDs = providerArray
	
	// Convert Labels list to array
	if !model.Labels.IsNull() && !model.Labels.IsUnknown() {
		var labelsArray []string
		diags = model.Labels.ElementsAs(context.Background(), &labelsArray, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting labels list: %v", diags)
		}
		policy.Labels = labelsArray
	}
	
	// Convert severity string to integer
	if !model.Severity.IsNull() && !model.Severity.IsUnknown() {
		policy.Severity = client.SeverityToInt(model.Severity.ValueString())
	}
	
	if !model.Category.IsNull() && !model.Category.IsUnknown() {
		policy.Category = model.Category.ValueString()
	}
	
	// Convert Frameworks list to array
	if !model.Frameworks.IsNull() && !model.Frameworks.IsUnknown() {
		var frameworksArray []string
		diags = model.Frameworks.ElementsAs(context.Background(), &frameworksArray, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting frameworks list: %v", diags)
		}
		policy.Frameworks = frameworksArray
	}
	
	return policy, nil
}

// isLikelyBase64 checks if a string is likely to be base64 encoded
func isLikelyBase64(s string) bool {
	// Base64 strings should only contain valid base64 characters and have proper padding
	matched, _ := regexp.MatchString(`^[A-Za-z0-9+/]*={0,2}$`, s)
	if !matched {
		return false
	}
	
	// Base64 strings should be divisible by 4 (with padding)
	if len(s)%4 != 0 {
		return false
	}
	
	// Should not look like typical Rego code
	if strings.Contains(s, "package ") || strings.Contains(s, "firefly") || strings.Contains(s, "{") {
		return false
	}
	
	return true
}

// parseGovernancePolicyImportID parses the import ID for a governance policy
func parseGovernancePolicyImportID(id string) (string, error) {
	// For governance policies, the import ID is just the policy ID
	id = strings.TrimSpace(id)
	if id == "" {
		return "", fmt.Errorf("invalid import ID: empty string")
	}
	return id, nil
}