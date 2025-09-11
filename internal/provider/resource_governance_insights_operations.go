package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapGovernanceInsightToModel maps API response to Terraform model
func mapGovernanceInsightToModel(insight *client.GovernanceInsight, model *GovernanceInsightResourceModel) error {
	model.ID = types.StringValue(insight.ID)
	model.Name = types.StringValue(insight.Name)
	
	if insight.Description != "" {
		model.Description = types.StringValue(insight.Description)
	} else {
		model.Description = types.StringNull()
	}
	
	model.Code = types.StringValue(insight.Code)
	
	// Convert Type array to list
	typeList, diags := types.ListValueFrom(context.Background(), types.StringType, insight.Type)
	if diags.HasError() {
		return fmt.Errorf("error converting type list: %v", diags)
	}
	model.Type = typeList
	
	// Convert ProviderIDs array to list
	providerList, diags := types.ListValueFrom(context.Background(), types.StringType, insight.ProviderIDs)
	if diags.HasError() {
		return fmt.Errorf("error converting provider IDs list: %v", diags)
	}
	model.ProviderIDs = providerList
	
	// Convert Labels array to list
	if len(insight.Labels) > 0 {
		labelsList, diags := types.ListValueFrom(context.Background(), types.StringType, insight.Labels)
		if diags.HasError() {
			return fmt.Errorf("error converting labels list: %v", diags)
		}
		model.Labels = labelsList
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
	
	// Convert severity integer to string
	if insight.Severity > 0 {
		model.Severity = types.StringValue(client.SeverityToString(insight.Severity))
	} else {
		model.Severity = types.StringNull()
	}
	
	if insight.Category != "" {
		model.Category = types.StringValue(insight.Category)
	} else {
		model.Category = types.StringNull()
	}
	
	// Convert Frameworks array to list
	if len(insight.Frameworks) > 0 {
		frameworksList, diags := types.ListValueFrom(context.Background(), types.StringType, insight.Frameworks)
		if diags.HasError() {
			return fmt.Errorf("error converting frameworks list: %v", diags)
		}
		model.Frameworks = frameworksList
	} else {
		model.Frameworks = types.ListNull(types.StringType)
	}
	
	return nil
}

// mapModelToGovernanceInsight maps Terraform model to API request
func mapModelToGovernanceInsight(model *GovernanceInsightResourceModel) (*client.GovernanceInsight, error) {
	insight := &client.GovernanceInsight{
		Name: model.Name.ValueString(),
		Code: model.Code.ValueString(),
	}
	
	if !model.ID.IsNull() && !model.ID.IsUnknown() {
		insight.ID = model.ID.ValueString()
	}
	
	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		insight.Description = model.Description.ValueString()
	}
	
	// Convert Type list to array
	var typeArray []string
	diags := model.Type.ElementsAs(context.Background(), &typeArray, false)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting type list: %v", diags)
	}
	insight.Type = typeArray
	
	// Convert ProviderIDs list to array
	var providerArray []string
	diags = model.ProviderIDs.ElementsAs(context.Background(), &providerArray, false)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting provider IDs list: %v", diags)
	}
	insight.ProviderIDs = providerArray
	
	// Convert Labels list to array
	if !model.Labels.IsNull() && !model.Labels.IsUnknown() {
		var labelsArray []string
		diags = model.Labels.ElementsAs(context.Background(), &labelsArray, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting labels list: %v", diags)
		}
		insight.Labels = labelsArray
	}
	
	// Convert severity string to integer
	if !model.Severity.IsNull() && !model.Severity.IsUnknown() {
		insight.Severity = client.SeverityToInt(model.Severity.ValueString())
	}
	
	if !model.Category.IsNull() && !model.Category.IsUnknown() {
		insight.Category = model.Category.ValueString()
	}
	
	// Convert Frameworks list to array
	if !model.Frameworks.IsNull() && !model.Frameworks.IsUnknown() {
		var frameworksArray []string
		diags = model.Frameworks.ElementsAs(context.Background(), &frameworksArray, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting frameworks list: %v", diags)
		}
		insight.Frameworks = frameworksArray
	}
	
	return insight, nil
}

// parseGovernanceInsightImportID parses the import ID for a governance insight
func parseGovernanceInsightImportID(id string) (string, error) {
	// For governance insights, the import ID is just the insight ID
	id = strings.TrimSpace(id)
	if id == "" {
		return "", fmt.Errorf("invalid import ID: empty string")
	}
	return id, nil
}