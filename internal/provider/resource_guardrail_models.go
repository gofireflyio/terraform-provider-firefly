package provider

import (
	"context"
	"terraform-provider-firefly/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IncludeExcludeWildcardModel represents a pattern for including and excluding items
type IncludeExcludeWildcardModel struct {
	Include types.List `tfsdk:"include"`
	Exclude types.List `tfsdk:"exclude"`
}

// GuardrailScopeModel defines the scope of a guardrail rule
type GuardrailScopeModel struct {
	Workspaces   *IncludeExcludeWildcardModel `tfsdk:"workspaces"`
	Repositories *IncludeExcludeWildcardModel `tfsdk:"repositories"`
	Branches     *IncludeExcludeWildcardModel `tfsdk:"branches"`
	Labels       *IncludeExcludeWildcardModel `tfsdk:"labels"`
}

// CostCriteriaModel defines criteria for cost-based guardrails
type CostCriteriaModel struct {
	ThresholdAmount    types.Float64 `tfsdk:"threshold_amount"`
	ThresholdPercentage types.Float64 `tfsdk:"threshold_percentage"`
}

// PolicyCriteriaModel defines criteria for policy-based guardrails
type PolicyCriteriaModel struct {
	Severity types.String                `tfsdk:"severity"`
	Policies *IncludeExcludeWildcardModel `tfsdk:"policies"`
}

// ResourceCriteriaModel defines criteria for resource-based guardrails
type ResourceCriteriaModel struct {
	Actions          types.List                 `tfsdk:"actions"`
	Regions          *IncludeExcludeWildcardModel `tfsdk:"regions"`
	AssetTypes       *IncludeExcludeWildcardModel `tfsdk:"asset_types"`
	SpecificResources types.List                 `tfsdk:"specific_resources"`
}

// TagCriteriaModel defines criteria for tag-based guardrails
type TagCriteriaModel struct {
	TagEnforcementMode types.String  `tfsdk:"tag_enforcement_mode"`
	RequiredTags       types.List    `tfsdk:"required_tags"`
	RequiredValues     types.Map     `tfsdk:"required_values"`
}

// GuardrailCriteriaModel defines the criteria for a guardrail rule
type GuardrailCriteriaModel struct {
	Cost     *CostCriteriaModel     `tfsdk:"cost"`
	Policy   *PolicyCriteriaModel   `tfsdk:"policy"`
	Resource *ResourceCriteriaModel `tfsdk:"resource"`
	Tag      *TagCriteriaModel      `tfsdk:"tag"`
}

// GuardrailResourceModel describes the guardrail resource data model
type GuardrailResourceModel struct {
	ID            types.String           `tfsdk:"id"`
	Name          types.String           `tfsdk:"name"`
	Type          types.String           `tfsdk:"type"`
	Scope         *GuardrailScopeModel   `tfsdk:"scope"`
	Criteria      *GuardrailCriteriaModel `tfsdk:"criteria"`
	IsEnabled     types.Bool             `tfsdk:"is_enabled"`
	CreatedAt     types.String           `tfsdk:"created_at"`
	UpdatedAt     types.String           `tfsdk:"updated_at"`
	NotificationID types.String           `tfsdk:"notification_id"`
	Severity      types.Int64            `tfsdk:"severity"`
}

// planToAPIGuardrail converts the Terraform plan to a client.GuardrailRule
func (r *guardrailResource) planToAPIGuardrail(ctx context.Context, plan GuardrailResourceModel) (*client.GuardrailRule, error) {
	guardrail := &client.GuardrailRule{
		Name:      plan.Name.ValueString(),
		Type:      plan.Type.ValueString(),
		IsEnabled: plan.IsEnabled.ValueBool(),
		Severity:  int(plan.Severity.ValueInt64()),
		CreatedBy: "terraform-provider", // Auto-populate required field
	}

	if !plan.NotificationID.IsNull() {
		guardrail.NotificationID = plan.NotificationID.ValueString()
	}

	// Convert scope if it exists
	if plan.Scope != nil {
		guardrail.Scope = &client.GuardrailScope{}

		// Workspaces
		if plan.Scope.Workspaces != nil {
			guardrail.Scope.Workspaces = &client.IncludeExcludeWildcard{}

			if !plan.Scope.Workspaces.Include.IsNull() {
				var include []string
				plan.Scope.Workspaces.Include.ElementsAs(ctx, &include, false)
				guardrail.Scope.Workspaces.Include = include
			}

			if !plan.Scope.Workspaces.Exclude.IsNull() {
				var exclude []string
				plan.Scope.Workspaces.Exclude.ElementsAs(ctx, &exclude, false)
				guardrail.Scope.Workspaces.Exclude = exclude
			}
		}

		// Repositories
		if plan.Scope.Repositories != nil {
			guardrail.Scope.Repositories = &client.IncludeExcludeWildcard{}

			if !plan.Scope.Repositories.Include.IsNull() {
				var include []string
				plan.Scope.Repositories.Include.ElementsAs(ctx, &include, false)
				guardrail.Scope.Repositories.Include = include
			}

			if !plan.Scope.Repositories.Exclude.IsNull() {
				var exclude []string
				plan.Scope.Repositories.Exclude.ElementsAs(ctx, &exclude, false)
				guardrail.Scope.Repositories.Exclude = exclude
			}
		}

		// Branches
		if plan.Scope.Branches != nil {
			guardrail.Scope.Branches = &client.IncludeExcludeWildcard{}

			if !plan.Scope.Branches.Include.IsNull() {
				var include []string
				plan.Scope.Branches.Include.ElementsAs(ctx, &include, false)
				guardrail.Scope.Branches.Include = include
			}

			if !plan.Scope.Branches.Exclude.IsNull() {
				var exclude []string
				plan.Scope.Branches.Exclude.ElementsAs(ctx, &exclude, false)
				guardrail.Scope.Branches.Exclude = exclude
			}
		}

		// Labels
		if plan.Scope.Labels != nil {
			guardrail.Scope.Labels = &client.IncludeExcludeWildcard{}

			if !plan.Scope.Labels.Include.IsNull() {
				var include []string
				plan.Scope.Labels.Include.ElementsAs(ctx, &include, false)
				guardrail.Scope.Labels.Include = include
			}

			if !plan.Scope.Labels.Exclude.IsNull() {
				var exclude []string
				plan.Scope.Labels.Exclude.ElementsAs(ctx, &exclude, false)
				guardrail.Scope.Labels.Exclude = exclude
			}
		}
	}

	// Convert criteria if it exists
	if plan.Criteria != nil {
		guardrail.Criteria = &client.GuardrailCriteria{}

		// Cost criteria
		if plan.Criteria.Cost != nil {
			guardrail.Criteria.Cost = &client.CostCriteria{}

			if !plan.Criteria.Cost.ThresholdAmount.IsNull() {
				thresholdAmount := plan.Criteria.Cost.ThresholdAmount.ValueFloat64()
				guardrail.Criteria.Cost.ThresholdAmount = &thresholdAmount
			}

			if !plan.Criteria.Cost.ThresholdPercentage.IsNull() {
				thresholdPercentage := plan.Criteria.Cost.ThresholdPercentage.ValueFloat64()
				guardrail.Criteria.Cost.ThresholdPercentage = &thresholdPercentage
			}
		}

		// Policy criteria
		if plan.Criteria.Policy != nil {
			guardrail.Criteria.Policy = &client.PolicyCriteria{}

			if !plan.Criteria.Policy.Severity.IsNull() {
				guardrail.Criteria.Policy.Severity = plan.Criteria.Policy.Severity.ValueString()
			}

			if plan.Criteria.Policy.Policies != nil {
				guardrail.Criteria.Policy.Policies = &client.IncludeExcludeWildcard{}

				if !plan.Criteria.Policy.Policies.Include.IsNull() {
					var include []string
					plan.Criteria.Policy.Policies.Include.ElementsAs(ctx, &include, false)
					guardrail.Criteria.Policy.Policies.Include = include
				}

				if !plan.Criteria.Policy.Policies.Exclude.IsNull() {
					var exclude []string
					plan.Criteria.Policy.Policies.Exclude.ElementsAs(ctx, &exclude, false)
					guardrail.Criteria.Policy.Policies.Exclude = exclude
				}
			}
		}

		// Resource criteria
		if plan.Criteria.Resource != nil {
			guardrail.Criteria.Resource = &client.ResourceCriteria{}

			if !plan.Criteria.Resource.Actions.IsNull() {
				var actions []string
				plan.Criteria.Resource.Actions.ElementsAs(ctx, &actions, false)
				guardrail.Criteria.Resource.Actions = actions
			}

			if !plan.Criteria.Resource.SpecificResources.IsNull() {
				var specificResources []string
				plan.Criteria.Resource.SpecificResources.ElementsAs(ctx, &specificResources, false)
				guardrail.Criteria.Resource.SpecificResources = specificResources
			}

			// Regions - always provide default if not specified
			if plan.Criteria.Resource.Regions != nil {
				guardrail.Criteria.Resource.Regions = &client.IncludeExcludeWildcard{}

				if !plan.Criteria.Resource.Regions.Include.IsNull() {
					var include []string
					plan.Criteria.Resource.Regions.Include.ElementsAs(ctx, &include, false)
					guardrail.Criteria.Resource.Regions.Include = include
				}

				if !plan.Criteria.Resource.Regions.Exclude.IsNull() {
					var exclude []string
					plan.Criteria.Resource.Regions.Exclude.ElementsAs(ctx, &exclude, false)
					guardrail.Criteria.Resource.Regions.Exclude = exclude
				}
			} else {
				// Provide default regions if not specified (required by API)
				guardrail.Criteria.Resource.Regions = &client.IncludeExcludeWildcard{
					Include: []string{"*"},
				}
			}

			// Asset types - always provide default if not specified
			if plan.Criteria.Resource.AssetTypes != nil {
				guardrail.Criteria.Resource.AssetTypes = &client.IncludeExcludeWildcard{}

				if !plan.Criteria.Resource.AssetTypes.Include.IsNull() {
					var include []string
					plan.Criteria.Resource.AssetTypes.Include.ElementsAs(ctx, &include, false)
					guardrail.Criteria.Resource.AssetTypes.Include = include
				}

				if !plan.Criteria.Resource.AssetTypes.Exclude.IsNull() {
					var exclude []string
					plan.Criteria.Resource.AssetTypes.Exclude.ElementsAs(ctx, &exclude, false)
					guardrail.Criteria.Resource.AssetTypes.Exclude = exclude
				}
			} else {
				// Provide default asset types if not specified (may be required by API)
				guardrail.Criteria.Resource.AssetTypes = &client.IncludeExcludeWildcard{
					Include: []string{"*"},
				}
			}
		}

		// Tag criteria
		if plan.Criteria.Tag != nil {
			guardrail.Criteria.Tag = &client.TagCriteria{}

			if !plan.Criteria.Tag.TagEnforcementMode.IsNull() {
				guardrail.Criteria.Tag.TagEnforcementMode = plan.Criteria.Tag.TagEnforcementMode.ValueString()
			}

			if !plan.Criteria.Tag.RequiredTags.IsNull() {
				var requiredTags []string
				plan.Criteria.Tag.RequiredTags.ElementsAs(ctx, &requiredTags, false)
				guardrail.Criteria.Tag.RequiredTags = requiredTags
			}

			if !plan.Criteria.Tag.RequiredValues.IsNull() {
				requiredValues := make(map[string][]string)
				for k, v := range plan.Criteria.Tag.RequiredValues.Elements() {
					if strVal, ok := v.(types.String); ok {
						// Store single values as single-element arrays
						requiredValues[k] = []string{strVal.ValueString()}
					}
				}
				guardrail.Criteria.Tag.RequiredValues = requiredValues
			}
		}
	}

	return guardrail, nil
}

// apiGuardrailToPlan converts a client.GuardrailRule to a GuardrailResourceModel
func (r *guardrailResource) apiGuardrailToPlan(ctx context.Context, apiGuardrail client.GuardrailRule, plan *GuardrailResourceModel) error {
	plan.ID = types.StringValue(apiGuardrail.ID)
	plan.Name = types.StringValue(apiGuardrail.Name)
	plan.Type = types.StringValue(apiGuardrail.Type)
	plan.IsEnabled = types.BoolValue(apiGuardrail.IsEnabled)
	plan.Severity = types.Int64Value(int64(apiGuardrail.Severity))

	if apiGuardrail.NotificationID != "" {
		plan.NotificationID = types.StringValue(apiGuardrail.NotificationID)
	}

	if apiGuardrail.CreatedAt != "" {
		plan.CreatedAt = types.StringValue(apiGuardrail.CreatedAt)
	}

	if apiGuardrail.UpdatedAt != "" {
		plan.UpdatedAt = types.StringValue(apiGuardrail.UpdatedAt)
	}

	// Convert scope if it exists
	if apiGuardrail.Scope != nil {
		if plan.Scope == nil {
			plan.Scope = &GuardrailScopeModel{}
		}

		// Workspaces
		if apiGuardrail.Scope.Workspaces != nil {
			if plan.Scope.Workspaces == nil {
				plan.Scope.Workspaces = &IncludeExcludeWildcardModel{}
			}

			if apiGuardrail.Scope.Workspaces.Include != nil {
				plan.Scope.Workspaces.Include = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Scope.Workspaces.Include))
			}

			if apiGuardrail.Scope.Workspaces.Exclude != nil {
				plan.Scope.Workspaces.Exclude = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Scope.Workspaces.Exclude))
			}
		}

		// Repositories
		if apiGuardrail.Scope.Repositories != nil {
			if plan.Scope.Repositories == nil {
				plan.Scope.Repositories = &IncludeExcludeWildcardModel{}
			}

			if apiGuardrail.Scope.Repositories.Include != nil {
				plan.Scope.Repositories.Include = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Scope.Repositories.Include))
			}

			if apiGuardrail.Scope.Repositories.Exclude != nil {
				plan.Scope.Repositories.Exclude = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Scope.Repositories.Exclude))
			}
		}

		// Branches
		if apiGuardrail.Scope.Branches != nil {
			if plan.Scope.Branches == nil {
				plan.Scope.Branches = &IncludeExcludeWildcardModel{}
			}

			if apiGuardrail.Scope.Branches.Include != nil {
				plan.Scope.Branches.Include = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Scope.Branches.Include))
			}

			if apiGuardrail.Scope.Branches.Exclude != nil {
				plan.Scope.Branches.Exclude = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Scope.Branches.Exclude))
			}
		}

		// Labels
		if apiGuardrail.Scope.Labels != nil {
			if plan.Scope.Labels == nil {
				plan.Scope.Labels = &IncludeExcludeWildcardModel{}
			}

			if apiGuardrail.Scope.Labels.Include != nil {
				plan.Scope.Labels.Include = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Scope.Labels.Include))
			}

			if apiGuardrail.Scope.Labels.Exclude != nil {
				plan.Scope.Labels.Exclude = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Scope.Labels.Exclude))
			}
		}
	}

	// Convert criteria if it exists
	if apiGuardrail.Criteria != nil {
		if plan.Criteria == nil {
			plan.Criteria = &GuardrailCriteriaModel{}
		}

		// Cost criteria
		if apiGuardrail.Criteria.Cost != nil {
			if plan.Criteria.Cost == nil {
				plan.Criteria.Cost = &CostCriteriaModel{}
			}

			if apiGuardrail.Criteria.Cost.ThresholdAmount != nil {
				plan.Criteria.Cost.ThresholdAmount = types.Float64Value(*apiGuardrail.Criteria.Cost.ThresholdAmount)
			}

			if apiGuardrail.Criteria.Cost.ThresholdPercentage != nil {
				plan.Criteria.Cost.ThresholdPercentage = types.Float64Value(*apiGuardrail.Criteria.Cost.ThresholdPercentage)
			}
		}

		// Policy criteria
		if apiGuardrail.Criteria.Policy != nil {
			if plan.Criteria.Policy == nil {
				plan.Criteria.Policy = &PolicyCriteriaModel{}
			}

			if apiGuardrail.Criteria.Policy.Severity != "" {
				plan.Criteria.Policy.Severity = types.StringValue(apiGuardrail.Criteria.Policy.Severity)
			}

			if apiGuardrail.Criteria.Policy.Policies != nil {
				if plan.Criteria.Policy.Policies == nil {
					plan.Criteria.Policy.Policies = &IncludeExcludeWildcardModel{}
				}

				if apiGuardrail.Criteria.Policy.Policies.Include != nil {
					plan.Criteria.Policy.Policies.Include = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Criteria.Policy.Policies.Include))
				}

				if apiGuardrail.Criteria.Policy.Policies.Exclude != nil {
					plan.Criteria.Policy.Policies.Exclude = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Criteria.Policy.Policies.Exclude))
				}
			}
		}

		// Resource criteria
		if apiGuardrail.Criteria.Resource != nil {
			if plan.Criteria.Resource == nil {
				plan.Criteria.Resource = &ResourceCriteriaModel{}
			}

			if apiGuardrail.Criteria.Resource.Actions != nil {
				plan.Criteria.Resource.Actions = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Criteria.Resource.Actions))
			}

			if apiGuardrail.Criteria.Resource.SpecificResources != nil {
				plan.Criteria.Resource.SpecificResources = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Criteria.Resource.SpecificResources))
			}

			// Regions
			if apiGuardrail.Criteria.Resource.Regions != nil {
				if plan.Criteria.Resource.Regions == nil {
					plan.Criteria.Resource.Regions = &IncludeExcludeWildcardModel{}
				}

				if apiGuardrail.Criteria.Resource.Regions.Include != nil {
					plan.Criteria.Resource.Regions.Include = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Criteria.Resource.Regions.Include))
				}

				if apiGuardrail.Criteria.Resource.Regions.Exclude != nil {
					plan.Criteria.Resource.Regions.Exclude = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Criteria.Resource.Regions.Exclude))
				}
			}

			// Asset types
			if apiGuardrail.Criteria.Resource.AssetTypes != nil {
				if plan.Criteria.Resource.AssetTypes == nil {
					plan.Criteria.Resource.AssetTypes = &IncludeExcludeWildcardModel{}
				}

				if apiGuardrail.Criteria.Resource.AssetTypes.Include != nil {
					plan.Criteria.Resource.AssetTypes.Include = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Criteria.Resource.AssetTypes.Include))
				}

				if apiGuardrail.Criteria.Resource.AssetTypes.Exclude != nil {
					plan.Criteria.Resource.AssetTypes.Exclude = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Criteria.Resource.AssetTypes.Exclude))
				}
			}
		}

		// Tag criteria
		if apiGuardrail.Criteria.Tag != nil {
			if plan.Criteria.Tag == nil {
				plan.Criteria.Tag = &TagCriteriaModel{}
			}

			if apiGuardrail.Criteria.Tag.TagEnforcementMode != "" {
				plan.Criteria.Tag.TagEnforcementMode = types.StringValue(apiGuardrail.Criteria.Tag.TagEnforcementMode)
			}

			if apiGuardrail.Criteria.Tag.RequiredTags != nil {
				plan.Criteria.Tag.RequiredTags = types.ListValueMust(types.StringType, listToValues(apiGuardrail.Criteria.Tag.RequiredTags))
			}

			if apiGuardrail.Criteria.Tag.RequiredValues != nil {
				elements := make(map[string]attr.Value)
				for k, v := range apiGuardrail.Criteria.Tag.RequiredValues {
					elements[k] = types.ListValueMust(types.StringType, listToValues(v))
				}
				plan.Criteria.Tag.RequiredValues = types.MapValueMust(types.ListType{ElemType: types.StringType}, elements)
			}
		}
	}

	return nil
}

// Helper function to convert a string slice to a slice of types.String
func listToValues(list []string) []attr.Value {
	values := make([]attr.Value, len(list))
	for i, v := range list {
		values[i] = types.StringValue(v)
	}
	return values
}
