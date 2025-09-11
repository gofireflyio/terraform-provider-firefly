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
var _ resource.Resource = &GovernanceInsightResource{}
var _ resource.ResourceWithImportState = &GovernanceInsightResource{}

// NewGovernanceInsightResource creates a new governance insight resource
func NewGovernanceInsightResource() resource.Resource {
	return &GovernanceInsightResource{}
}

// GovernanceInsightResource defines the resource implementation
type GovernanceInsightResource struct {
	client *client.Client
}

func (r *GovernanceInsightResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_governance_insight"
}

func (r *GovernanceInsightResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Firefly governance insight (custom policy rule)",
		
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the governance insight",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the governance insight",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the governance insight",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"code": schema.StringAttribute{
				MarkdownDescription: "The Rego code for the insight rule (can be base64 encoded)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of resource types this insight applies to (e.g., 'aws_cloudwatch_event_target')",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"provider_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of provider IDs this insight applies to (e.g., 'aws_all', specific account IDs)",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"labels": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of labels for categorizing the insight",
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: "The severity level of the insight (flexible, strict, warning)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("warning"),
				Validators: []validator.String{
					stringvalidator.OneOf("flexible", "strict", "warning"),
				},
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The category of the insight (e.g., 'Misconfiguration')",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"frameworks": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of compliance frameworks this insight relates to (e.g., 'SOC2', 'ISO27001')",
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
		},
	}
}

func (r *GovernanceInsightResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GovernanceInsightResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GovernanceInsightResourceModel
	
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	// Convert model to API request
	insight, err := mapModelToGovernanceInsight(&data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating governance insight",
			fmt.Sprintf("Could not convert model to API request: %s", err),
		)
		return
	}
	
	tflog.Debug(ctx, "Creating governance insight", map[string]interface{}{
		"name": insight.Name,
	})
	
	// Create the insight
	createdInsight, err := r.client.GovernanceInsights.Create(insight)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating governance insight",
			fmt.Sprintf("Could not create governance insight: %s", err),
		)
		return
	}
	
	tflog.Debug(ctx, "Created governance insight", map[string]interface{}{
		"id":   createdInsight.ID,
		"name": createdInsight.Name,
	})
	
	// Map response to model
	err = mapGovernanceInsightToModel(createdInsight, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating governance insight",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}
	
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GovernanceInsightResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GovernanceInsightResourceModel
	
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	tflog.Debug(ctx, "Reading governance insight", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	
	// Get the insight
	insight, err := r.client.GovernanceInsights.Get(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading governance insight",
			fmt.Sprintf("Could not read governance insight: %s", err),
		)
		return
	}
	
	// Map response to model
	err = mapGovernanceInsightToModel(insight, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading governance insight",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}
	
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GovernanceInsightResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GovernanceInsightResourceModel
	
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	// Convert model to API request
	insight, err := mapModelToGovernanceInsight(&data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating governance insight",
			fmt.Sprintf("Could not convert model to API request: %s", err),
		)
		return
	}
	
	tflog.Debug(ctx, "Updating governance insight", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": insight.Name,
	})
	
	// Update the insight
	updatedInsight, err := r.client.GovernanceInsights.Update(data.ID.ValueString(), insight)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating governance insight",
			fmt.Sprintf("Could not update governance insight: %s", err),
		)
		return
	}
	
	tflog.Debug(ctx, "Updated governance insight", map[string]interface{}{
		"id":   updatedInsight.ID,
		"name": updatedInsight.Name,
	})
	
	// Map response to model
	err = mapGovernanceInsightToModel(updatedInsight, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating governance insight",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}
	
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GovernanceInsightResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GovernanceInsightResourceModel
	
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	tflog.Debug(ctx, "Deleting governance insight", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	
	// Delete the insight
	err := r.client.GovernanceInsights.Delete(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting governance insight",
			fmt.Sprintf("Could not delete governance insight: %s", err),
		)
		return
	}
	
	tflog.Debug(ctx, "Deleted governance insight", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (r *GovernanceInsightResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the import ID
	insightID, err := parseGovernanceInsightImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing governance insight",
			fmt.Sprintf("Invalid import ID format: %s", err),
		)
		return
	}
	
	// Set the ID in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), insightID)...)
}