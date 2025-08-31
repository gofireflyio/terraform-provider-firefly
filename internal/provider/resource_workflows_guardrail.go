package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &guardrailResource{}
	_ resource.ResourceWithConfigure   = &guardrailResource{}
	_ resource.ResourceWithImportState = &guardrailResource{}
)

// NewGuardrailResource is a helper function to simplify the provider implementation
func NewGuardrailResource() resource.Resource {
	return &guardrailResource{}
}

// guardrailResource is the resource implementation
type guardrailResource struct {
	client *client.Client
}

// Metadata returns the resource type name
func (r *guardrailResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows_guardrail"
}

// Configure adds the provider configured client to the resource
func (r *guardrailResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates a new guardrail rule
func (r *guardrailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan GuardrailResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new guardrail
	guardrail, err := r.planToAPIGuardrail(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Guardrail",
			fmt.Sprintf("Could not convert plan to API guardrail: %s", err),
		)
		return
	}

	// Create in the API
	tflog.Debug(ctx, "Creating guardrail", map[string]interface{}{
		"name": guardrail.Name,
		"type": guardrail.Type,
	})

	createResp, err := r.client.Guardrails.CreateGuardrail(guardrail)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Guardrail",
			fmt.Sprintf("Could not create guardrail: %s", err),
		)
		return
	}

	// Set ID from response
	plan.ID = types.StringValue(createResp.RuleID)
	
	// Set notification ID from response or empty string if none
	if createResp.NotificationID != "" {
		plan.NotificationID = types.StringValue(createResp.NotificationID)
	} else {
		plan.NotificationID = types.StringValue("")
	}

	// Fetch the created guardrail to get computed properties only
	createdGuardrail, err := r.client.Guardrails.GetGuardrail(createResp.RuleID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Guardrail",
			fmt.Sprintf("Could not read created guardrail: %s", err),
		)
		return
	}

	// Update only computed fields from the API response
	// Don't overwrite user-provided scope with API defaults
	if createdGuardrail.CreatedAt != "" {
		plan.CreatedAt = types.StringValue(createdGuardrail.CreatedAt)
	}
	if createdGuardrail.UpdatedAt != "" {
		plan.UpdatedAt = types.StringValue(createdGuardrail.UpdatedAt)
	}

	// Set state to fully populated plan
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data
func (r *guardrailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state GuardrailResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get guardrail from API
	guardrail, err := r.client.Guardrails.GetGuardrail(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Guardrail",
			fmt.Sprintf("Could not read guardrail ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Update only the basic fields - preserve original scope structure
	state.ID = types.StringValue(guardrail.ID)
	state.Name = types.StringValue(guardrail.Name)
	state.Type = types.StringValue(guardrail.Type)
	state.IsEnabled = types.BoolValue(guardrail.IsEnabled)
	state.Severity = types.Int64Value(int64(guardrail.Severity))

	if guardrail.NotificationID != "" {
		state.NotificationID = types.StringValue(guardrail.NotificationID)
	}

	if guardrail.CreatedAt != "" {
		state.CreatedAt = types.StringValue(guardrail.CreatedAt)
	}

	if guardrail.UpdatedAt != "" {
		state.UpdatedAt = types.StringValue(guardrail.UpdatedAt)
	}

	// Only update scope fields that were explicitly configured in the original state
	if state.Scope != nil && guardrail.Scope != nil {
		// Update user-provided scope values based on what they originally configured
		if state.Scope.Workspaces != nil && guardrail.Scope.Workspaces != nil {
			if !state.Scope.Workspaces.Include.IsNull() && guardrail.Scope.Workspaces.Include != nil {
				state.Scope.Workspaces.Include = types.ListValueMust(types.StringType, listToValues(guardrail.Scope.Workspaces.Include))
			}
			if !state.Scope.Workspaces.Exclude.IsNull() && guardrail.Scope.Workspaces.Exclude != nil {
				state.Scope.Workspaces.Exclude = types.ListValueMust(types.StringType, listToValues(guardrail.Scope.Workspaces.Exclude))
			}
		}
		
		if state.Scope.Repositories != nil && guardrail.Scope.Repositories != nil {
			if !state.Scope.Repositories.Include.IsNull() && guardrail.Scope.Repositories.Include != nil {
				state.Scope.Repositories.Include = types.ListValueMust(types.StringType, listToValues(guardrail.Scope.Repositories.Include))
			}
			if !state.Scope.Repositories.Exclude.IsNull() && guardrail.Scope.Repositories.Exclude != nil {
				state.Scope.Repositories.Exclude = types.ListValueMust(types.StringType, listToValues(guardrail.Scope.Repositories.Exclude))
			}
		}
		
		if state.Scope.Branches != nil && guardrail.Scope.Branches != nil {
			if !state.Scope.Branches.Include.IsNull() && guardrail.Scope.Branches.Include != nil {
				state.Scope.Branches.Include = types.ListValueMust(types.StringType, listToValues(guardrail.Scope.Branches.Include))
			}
			if !state.Scope.Branches.Exclude.IsNull() && guardrail.Scope.Branches.Exclude != nil {
				state.Scope.Branches.Exclude = types.ListValueMust(types.StringType, listToValues(guardrail.Scope.Branches.Exclude))
			}
		}
		
		if state.Scope.Labels != nil && guardrail.Scope.Labels != nil {
			if !state.Scope.Labels.Include.IsNull() && guardrail.Scope.Labels.Include != nil {
				state.Scope.Labels.Include = types.ListValueMust(types.StringType, listToValues(guardrail.Scope.Labels.Include))
			}
			if !state.Scope.Labels.Exclude.IsNull() && guardrail.Scope.Labels.Exclude != nil {
				state.Scope.Labels.Exclude = types.ListValueMust(types.StringType, listToValues(guardrail.Scope.Labels.Exclude))
			}
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState imports a resource by ID
func (r *guardrailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
