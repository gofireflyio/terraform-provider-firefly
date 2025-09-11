package provider

import (
	"context"
	"fmt"
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
	
	model.Code = types.StringValue(policy.Code)
	
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
	policy := &client.GovernancePolicy{
		Name: model.Name.ValueString(),
		Code: model.Code.ValueString(),
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

// parseGovernancePolicyImportID parses the import ID for a governance policy
func parseGovernancePolicyImportID(id string) (string, error) {
	// For governance policies, the import ID is just the policy ID
	id = strings.TrimSpace(id)
	if id == "" {
		return "", fmt.Errorf("invalid import ID: empty string")
	}
	return id, nil
}