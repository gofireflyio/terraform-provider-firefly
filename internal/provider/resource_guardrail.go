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
	resp.TypeName = req.ProviderTypeName + "_guardrail"
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

	// Fetch the created guardrail to get all of its properties
	createdGuardrail, err := r.client.Guardrails.GetGuardrail(createResp.RuleID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Guardrail",
			fmt.Sprintf("Could not read created guardrail: %s", err),
		)
		return
	}

	// Map created guardrail to plan
	err = r.apiGuardrailToPlan(ctx, *createdGuardrail, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting API Response",
			fmt.Sprintf("Could not convert API response to plan: %s", err),
		)
		return
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

	// Map response to state
	err = r.apiGuardrailToPlan(ctx, *guardrail, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting API Response",
			fmt.Sprintf("Could not convert API response to state: %s", err),
		)
		return
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
