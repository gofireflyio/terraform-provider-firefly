package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Update updates an existing guardrail rule
func (r *guardrailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan GuardrailResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request from plan
	guardrail, err := r.planToAPIGuardrail(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Guardrail",
			fmt.Sprintf("Could not convert plan to API guardrail: %s", err),
		)
		return
	}

	// Update in the API
	_, err = r.client.Guardrails.UpdateGuardrail(plan.ID.ValueString(), guardrail)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Guardrail",
			fmt.Sprintf("Could not update guardrail ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Get updated guardrail from API
	updatedGuardrail, err := r.client.Guardrails.GetGuardrail(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Guardrail",
			fmt.Sprintf("Could not read updated guardrail ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Update computed fields from the API response
	// Don't overwrite user-provided scope with API defaults
	if updatedGuardrail.CreatedAt != "" {
		plan.CreatedAt = types.StringValue(updatedGuardrail.CreatedAt)
	}
	if updatedGuardrail.UpdatedAt != "" {
		plan.UpdatedAt = types.StringValue(updatedGuardrail.UpdatedAt)
	}

	// Always set notification_id to a known value
	if updatedGuardrail.NotificationID != "" {
		plan.NotificationID = types.StringValue(updatedGuardrail.NotificationID)
	} else {
		// If no notification_id from API, set to empty string (known value)
		plan.NotificationID = types.StringValue("")
	}

	// Update scope fields that were provided in the plan
	if plan.Scope != nil && updatedGuardrail.Scope != nil {
		if plan.Scope.Workspaces != nil && updatedGuardrail.Scope.Workspaces != nil {
			if !plan.Scope.Workspaces.Include.IsNull() && updatedGuardrail.Scope.Workspaces.Include != nil {
				plan.Scope.Workspaces.Include = types.ListValueMust(types.StringType, listToValues(updatedGuardrail.Scope.Workspaces.Include))
			}
			if !plan.Scope.Workspaces.Exclude.IsNull() && updatedGuardrail.Scope.Workspaces.Exclude != nil {
				plan.Scope.Workspaces.Exclude = types.ListValueMust(types.StringType, listToValues(updatedGuardrail.Scope.Workspaces.Exclude))
			}
		}

		if plan.Scope.Repositories != nil && updatedGuardrail.Scope.Repositories != nil {
			if !plan.Scope.Repositories.Include.IsNull() && updatedGuardrail.Scope.Repositories.Include != nil {
				plan.Scope.Repositories.Include = types.ListValueMust(types.StringType, listToValues(updatedGuardrail.Scope.Repositories.Include))
			}
			if !plan.Scope.Repositories.Exclude.IsNull() && updatedGuardrail.Scope.Repositories.Exclude != nil {
				plan.Scope.Repositories.Exclude = types.ListValueMust(types.StringType, listToValues(updatedGuardrail.Scope.Repositories.Exclude))
			}
		}

		if plan.Scope.Branches != nil && updatedGuardrail.Scope.Branches != nil {
			if !plan.Scope.Branches.Include.IsNull() && updatedGuardrail.Scope.Branches.Include != nil {
				plan.Scope.Branches.Include = types.ListValueMust(types.StringType, listToValues(updatedGuardrail.Scope.Branches.Include))
			}
			if !plan.Scope.Branches.Exclude.IsNull() && updatedGuardrail.Scope.Branches.Exclude != nil {
				plan.Scope.Branches.Exclude = types.ListValueMust(types.StringType, listToValues(updatedGuardrail.Scope.Branches.Exclude))
			}
		}

		if plan.Scope.Labels != nil && updatedGuardrail.Scope.Labels != nil {
			if !plan.Scope.Labels.Include.IsNull() && updatedGuardrail.Scope.Labels.Include != nil {
				plan.Scope.Labels.Include = types.ListValueMust(types.StringType, listToValues(updatedGuardrail.Scope.Labels.Include))
			}
			if !plan.Scope.Labels.Exclude.IsNull() && updatedGuardrail.Scope.Labels.Exclude != nil {
				plan.Scope.Labels.Exclude = types.ListValueMust(types.StringType, listToValues(updatedGuardrail.Scope.Labels.Exclude))
			}
		}
	}

	// Set state to fully populated plan
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes an existing guardrail rule
func (r *guardrailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state GuardrailResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete guardrail
	_, err := r.client.Guardrails.DeleteGuardrail(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Guardrail",
			fmt.Sprintf("Could not delete guardrail ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

// Schema defines the schema for the resource
func (r *guardrailResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Firefly guardrail rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the guardrail rule",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the guardrail rule",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the guardrail rule (policy, cost, resource, or tag)",
				Required:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the guardrail rule is enabled",
				Required:    true,
				Computed:    false,
			},
			"notification_id": schema.StringAttribute{
				Description: "ID of the associated notification",
				Optional:    true,
				Computed:    true,
			},
			"severity": schema.StringAttribute{
				Description: "Severity level of the guardrail rule. Allowed values: Flexible, Strict, Warning.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Flexible", "Strict", "Warning"),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the guardrail rule was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the guardrail rule was last updated",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"scope": schema.SingleNestedBlock{
				Description: "Scope of the guardrail rule",
				Blocks: map[string]schema.Block{
					"workspaces": schema.SingleNestedBlock{
						Description: "Workspace patterns to include or exclude",
						Attributes: map[string]schema.Attribute{
							"include": schema.ListAttribute{
								Description: "List of patterns to include",
								ElementType: types.StringType,
								Optional:    true,
							},
							"exclude": schema.ListAttribute{
								Description: "List of patterns to exclude",
								ElementType: types.StringType,
								Optional:    true,
							},
						},
					},
					"repositories": schema.SingleNestedBlock{
						Description: "Repository patterns to include or exclude",
						Attributes: map[string]schema.Attribute{
							"include": schema.ListAttribute{
								Description: "List of patterns to include",
								ElementType: types.StringType,
								Optional:    true,
							},
							"exclude": schema.ListAttribute{
								Description: "List of patterns to exclude",
								ElementType: types.StringType,
								Optional:    true,
							},
						},
					},
					"branches": schema.SingleNestedBlock{
						Description: "Branch patterns to include or exclude",
						Attributes: map[string]schema.Attribute{
							"include": schema.ListAttribute{
								Description: "List of patterns to include",
								ElementType: types.StringType,
								Optional:    true,
							},
							"exclude": schema.ListAttribute{
								Description: "List of patterns to exclude",
								ElementType: types.StringType,
								Optional:    true,
							},
						},
					},
					"labels": schema.SingleNestedBlock{
						Description: "Label patterns to include or exclude",
						Attributes: map[string]schema.Attribute{
							"include": schema.ListAttribute{
								Description: "List of patterns to include",
								ElementType: types.StringType,
								Optional:    true,
							},
							"exclude": schema.ListAttribute{
								Description: "List of patterns to exclude",
								ElementType: types.StringType,
								Optional:    true,
							},
						},
					},
				},
			},
			"criteria": schema.SingleNestedBlock{
				Description: "Criteria for the guardrail rule",
				Blocks: map[string]schema.Block{
					"cost": schema.SingleNestedBlock{
						Description: "Cost criteria for the guardrail rule",
						Attributes: map[string]schema.Attribute{
							"threshold_amount": schema.Float64Attribute{
								Description: "Absolute threshold amount (in USD) for cost criteria",
								Optional:    true,
							},
							"threshold_percentage": schema.Float64Attribute{
								Description: "Percentage threshold for cost criteria",
								Optional:    true,
							},
						},
					},
					"policy": schema.SingleNestedBlock{
						Description: "Policy criteria for the guardrail rule",
						Attributes: map[string]schema.Attribute{
							"severity": schema.StringAttribute{
								Description: "Severity level for policy criteria",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"policies": schema.SingleNestedBlock{
								Description: "Policy patterns to include or exclude",
								Attributes: map[string]schema.Attribute{
									"include": schema.ListAttribute{
										Description: "List of patterns to include",
										ElementType: types.StringType,
										Optional:    true,
									},
									"exclude": schema.ListAttribute{
										Description: "List of patterns to exclude",
										ElementType: types.StringType,
										Optional:    true,
									},
								},
							},
						},
					},
					"resource": schema.SingleNestedBlock{
						Description: "Resource criteria for the guardrail rule",
						Attributes: map[string]schema.Attribute{
							"actions": schema.ListAttribute{
								Description: "List of actions to filter by",
								ElementType: types.StringType,
								Optional:    true,
							},
							"specific_resources": schema.ListAttribute{
								Description: "List of specific resources to filter by",
								ElementType: types.StringType,
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"regions": schema.SingleNestedBlock{
								Description: "Region patterns to include or exclude",
								Attributes: map[string]schema.Attribute{
									"include": schema.ListAttribute{
										Description: "List of patterns to include",
										ElementType: types.StringType,
										Optional:    true,
									},
									"exclude": schema.ListAttribute{
										Description: "List of patterns to exclude",
										ElementType: types.StringType,
										Optional:    true,
									},
								},
							},
							"asset_types": schema.SingleNestedBlock{
								Description: "Asset type patterns to include or exclude",
								Attributes: map[string]schema.Attribute{
									"include": schema.ListAttribute{
										Description: "List of patterns to include",
										ElementType: types.StringType,
										Optional:    true,
									},
									"exclude": schema.ListAttribute{
										Description: "List of patterns to exclude",
										ElementType: types.StringType,
										Optional:    true,
									},
								},
							},
						},
					},
					"tag": schema.SingleNestedBlock{
						Description: "Tag criteria for the guardrail rule",
						Attributes: map[string]schema.Attribute{
							"tag_enforcement_mode": schema.StringAttribute{
								Description: "Mode of tag enforcement (requiredTags, anyTags, requiredValues)",
								Optional:    true,
							},
							"required_tags": schema.ListAttribute{
								Description: "List of required tags",
								ElementType: types.StringType,
								Optional:    true,
							},
							"required_values": schema.MapAttribute{
								Description: "Map of tag keys to required values",
								ElementType: types.ListType{
									ElemType: types.StringType,
								},
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}
