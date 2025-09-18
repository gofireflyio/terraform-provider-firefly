package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &governancePolicyResource{}
	_ resource.ResourceWithConfigure   = &governancePolicyResource{}
	_ resource.ResourceWithImportState = &governancePolicyResource{}
)

// NewGovernancePolicyResource is a helper function to simplify the provider implementation
func NewGovernancePolicyResource() resource.Resource {
	return &governancePolicyResource{}
}

// governancePolicyResource is the resource implementation
type governancePolicyResource struct {
	client *client.Client
}

// GovernancePolicyResourceModel describes the resource data model
type GovernancePolicyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Code        types.String `tfsdk:"code"`
	Type        types.List   `tfsdk:"type"`
	ProviderIDs types.List   `tfsdk:"provider_ids"`
	Labels      types.List   `tfsdk:"labels"`
	Severity    types.Int64  `tfsdk:"severity"`
	Category    types.String `tfsdk:"category"`
	Frameworks  types.List   `tfsdk:"frameworks"`
}

// Metadata returns the resource type name
func (r *governancePolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_governance_policy"
}

// Schema defines the schema for the resource
func (r *governancePolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Firefly governance policy (insight).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the governance policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the governance policy.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the governance policy.",
				Optional:    true,
			},
			"code": schema.StringAttribute{
				Description: "The Rego code for the policy. Can be base64 encoded.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type": schema.ListAttribute{
				Description: "List of resource types this policy applies to (e.g., ['aws_instance']).",
				Required:    true,
				ElementType: types.StringType,
			},
			"provider_ids": schema.ListAttribute{
				Description: "List of provider IDs this policy applies to (e.g., ['aws_all']).",
				Required:    true,
				ElementType: types.StringType,
			},
			"labels": schema.ListAttribute{
				Description: "List of labels associated with the policy.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"severity": schema.Int64Attribute{
				Description: "The severity level of the policy. Valid values: 0 (TRACE), 1 (INFO), 2 (LOW), 3 (MEDIUM), 4 (HIGH), 5 (CRITICAL). Defaults to 3 (MEDIUM).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3), // Default to MEDIUM
				Validators: []validator.Int64{
					int64validator.Between(0, 5),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"category": schema.StringAttribute{
				Description: "The category of the policy (e.g., 'Optimization', 'Misconfiguration', 'Security').",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Optimization", "Misconfiguration", "Security", "Compliance", "Cost"),
				},
			},
			"frameworks": schema.ListAttribute{
				Description: "List of compliance frameworks associated with the policy (e.g., ['SOC2', 'HIPAA']).",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *governancePolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates a new governance policy
func (r *governancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan GovernancePolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert plan to API model
	insight, err := r.planToAPIModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Governance Policy",
			fmt.Sprintf("Could not convert plan to API model: %s", err),
		)
		return
	}

	// Create in the API
	tflog.Debug(ctx, "Creating governance policy", map[string]interface{}{
		"name":     insight.Name,
		"category": insight.Category,
	})

	createdInsight, err := r.client.GovernanceInsights.CreateGovernanceInsight(insight)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Governance Policy",
			fmt.Sprintf("Could not create governance policy: %s", err),
		)
		return
	}

	// Update plan with computed values
	plan.ID = types.StringValue(createdInsight.ID)
	if createdInsight.Severity > 0 {
		plan.Severity = types.Int64Value(int64(createdInsight.Severity))
	}

	// Save to state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the state with the latest data
func (r *governancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state GovernancePolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get governance policy from API
	insight, err := r.client.GovernanceInsights.GetGovernanceInsight(state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Governance Policy",
			fmt.Sprintf("Could not read governance policy ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Update state from API response
	r.apiToState(ctx, insight, &state)

	// Save updated state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates an existing governance policy
func (r *governancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan and state
	var plan GovernancePolicyResourceModel
	var state GovernancePolicyResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert plan to API model
	insight, err := r.planToAPIModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Governance Policy",
			fmt.Sprintf("Could not convert plan to API model: %s", err),
		)
		return
	}

	// Update in the API
	tflog.Debug(ctx, "Updating governance policy", map[string]interface{}{
		"id":   state.ID.ValueString(),
		"name": insight.Name,
	})

	updatedInsight, err := r.client.GovernanceInsights.UpdateGovernanceInsight(state.ID.ValueString(), insight)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Governance Policy",
			fmt.Sprintf("Could not update governance policy ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Update state with any computed values from response
	plan.ID = state.ID // Preserve the ID
	if updatedInsight.Severity > 0 {
		plan.Severity = types.Int64Value(int64(updatedInsight.Severity))
	}

	// Save to state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes an existing governance policy
func (r *governancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state GovernancePolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete from the API
	tflog.Debug(ctx, "Deleting governance policy", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := r.client.GovernanceInsights.DeleteGovernanceInsight(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Governance Policy",
			fmt.Sprintf("Could not delete governance policy ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

// ImportState imports an existing governance policy
func (r *governancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID directly as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to convert plan to API model
func (r *governancePolicyResource) planToAPIModel(ctx context.Context, plan GovernancePolicyResourceModel) (*client.GovernanceInsight, error) {
	insight := &client.GovernanceInsight{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Code:        plan.Code.ValueString(),
		Category:    plan.Category.ValueString(),
	}

	// Convert severity
	if !plan.Severity.IsNull() {
		insight.Severity = int(plan.Severity.ValueInt64())
	}

	// Convert type list
	if !plan.Type.IsNull() {
		var typeList []string
		diags := plan.Type.ElementsAs(ctx, &typeList, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting type list")
		}
		insight.Type = typeList
	}

	// Convert provider IDs list
	if !plan.ProviderIDs.IsNull() {
		var providerList []string
		diags := plan.ProviderIDs.ElementsAs(ctx, &providerList, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting provider IDs list")
		}
		insight.ProviderIDs = providerList
	}

	// Convert labels list
	if !plan.Labels.IsNull() {
		var labelList []string
		diags := plan.Labels.ElementsAs(ctx, &labelList, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting labels list")
		}
		insight.Labels = labelList
	}

	// Convert frameworks list
	if !plan.Frameworks.IsNull() {
		var frameworkList []string
		diags := plan.Frameworks.ElementsAs(ctx, &frameworkList, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting frameworks list")
		}
		insight.Frameworks = frameworkList
	}

	return insight, nil
}

// Helper function to convert API response to state
func (r *governancePolicyResource) apiToState(ctx context.Context, insight *client.GovernanceInsight, state *GovernancePolicyResourceModel) {
	state.ID = types.StringValue(insight.ID)
	state.Name = types.StringValue(insight.Name)
	state.Description = types.StringValue(insight.Description)
	state.Code = types.StringValue(insight.Code)
	state.Category = types.StringValue(insight.Category)

	if insight.Severity > 0 {
		state.Severity = types.Int64Value(int64(insight.Severity))
	}

	// Convert type array
	if len(insight.Type) > 0 {
		typeList, _ := types.ListValueFrom(ctx, types.StringType, insight.Type)
		state.Type = typeList
	}

	// Convert provider IDs array
	if len(insight.ProviderIDs) > 0 {
		providerList, _ := types.ListValueFrom(ctx, types.StringType, insight.ProviderIDs)
		state.ProviderIDs = providerList
	}

	// Convert labels array
	if len(insight.Labels) > 0 {
		labelList, _ := types.ListValueFrom(ctx, types.StringType, insight.Labels)
		state.Labels = labelList
	}

	// Convert frameworks array
	if len(insight.Frameworks) > 0 {
		frameworkList, _ := types.ListValueFrom(ctx, types.StringType, insight.Frameworks)
		state.Frameworks = frameworkList
	}
}