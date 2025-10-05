package provider

import (
	"context"
	"fmt"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &GovernancePolicyResource{}
var _ resource.ResourceWithImportState = &GovernancePolicyResource{}

// NewGovernancePolicyResource creates a new governance policy resource
func NewGovernancePolicyResource() resource.Resource {
	return &GovernancePolicyResource{}
}

// GovernancePolicyResource defines the resource implementation
type GovernancePolicyResource struct {
	client *client.Client
}

func (r *GovernancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_governance_policy"
}

func (r *GovernancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Firefly governance policy (custom policy rule)",
		
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the governance policy",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the governance policy",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the governance policy",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"code": schema.StringAttribute{
				MarkdownDescription: "The Rego code for the policy rule (can be base64 encoded)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of resource types this policy applies to (e.g., 'aws_cloudwatch_event_target')",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"provider_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of provider IDs this policy applies to (e.g., 'aws_all', specific account IDs)",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"labels": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of labels for categorizing the policy",
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: "The severity level of the policy (trace, info, low, medium, high, critical)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("low"),
				Validators: []validator.String{
					stringvalidator.OneOf("trace", "info", "low", "medium", "high", "critical"),
				},
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The category of the policy (e.g., 'Misconfiguration')",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"frameworks": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of compliance frameworks this policy relates to (e.g., 'SOC2', 'ISO27001')",
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
		},
	}
}

func (r *GovernancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider is not configured
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

func (r *GovernancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GovernancePolicyResourceModel
	
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	// Convert model to API request
	policy, err := mapModelToGovernancePolicy(&data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating governance policy",
			fmt.Sprintf("Could not convert model to API request: %s", err),
		)
		return
	}
	
	tflog.Debug(ctx, "Creating governance policy", map[string]interface{}{
		"name": policy.Name,
	})
	
	// Create the policy
	createdPolicy, err := r.client.GovernancePolicies.Create(policy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating governance policy",
			fmt.Sprintf("Could not create governance policy: %s", err),
		)
		return
	}
	
	tflog.Debug(ctx, "Created governance policy", map[string]interface{}{
		"id":   createdPolicy.ID,
		"name": createdPolicy.Name,
	})
	
	// Map response to model
	err = mapGovernancePolicyToModel(createdPolicy, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating governance policy",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}
	
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GovernancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GovernancePolicyResourceModel
	
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	tflog.Debug(ctx, "Reading governance policy", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	
	// Get the policy
	policy, err := r.client.GovernancePolicies.Get(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading governance policy",
			fmt.Sprintf("Could not read governance policy: %s", err),
		)
		return
	}
	
	// Map response to model
	err = mapGovernancePolicyToModel(policy, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading governance policy",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}
	
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GovernancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GovernancePolicyResourceModel
	
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	// Convert model to API request
	policy, err := mapModelToGovernancePolicy(&data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating governance policy",
			fmt.Sprintf("Could not convert model to API request: %s", err),
		)
		return
	}
	
	tflog.Debug(ctx, "Updating governance policy", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": policy.Name,
	})
	
	// Update the policy
	updatedPolicy, err := r.client.GovernancePolicies.Update(data.ID.ValueString(), policy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating governance policy",
			fmt.Sprintf("Could not update governance policy: %s", err),
		)
		return
	}
	
	// Preserve the ID if the API response doesn't include it
	// This ensures consistent state management even when the API has inconsistent response formats
	if updatedPolicy.ID == "" {
		updatedPolicy.ID = data.ID.ValueString()
	}
	
	tflog.Debug(ctx, "Updated governance policy", map[string]interface{}{
		"id":   updatedPolicy.ID,
		"name": updatedPolicy.Name,
	})
	
	// Map response to model
	err = mapGovernancePolicyToModel(updatedPolicy, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating governance policy",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}
	
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GovernancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GovernancePolicyResourceModel
	
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	tflog.Debug(ctx, "Deleting governance policy", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	
	// Delete the policy
	err := r.client.GovernancePolicies.Delete(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting governance policy",
			fmt.Sprintf("Could not delete governance policy: %s", err),
		)
		return
	}
	
	tflog.Debug(ctx, "Deleted governance policy", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (r *GovernancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the import ID
	policyID, err := parseGovernancePolicyImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing governance policy",
			fmt.Sprintf("Invalid import ID format: %s", err),
		)
		return
	}
	
	// Set the ID in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), policyID)...)
}