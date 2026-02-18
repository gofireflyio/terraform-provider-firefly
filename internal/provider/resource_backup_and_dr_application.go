package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &BackupAndDrApplicationResource{}
var _ resource.ResourceWithImportState = &BackupAndDrApplicationResource{}

// NewBackupAndDrApplicationResource creates a new backup and DR application resource
func NewBackupAndDrApplicationResource() resource.Resource {
	return &BackupAndDrApplicationResource{}
}

// BackupAndDrApplicationResource defines the resource implementation
type BackupAndDrApplicationResource struct {
	client *client.Client
}

func (r *BackupAndDrApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_and_dr_application"
}

func (r *BackupAndDrApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Firefly Backup & DR application (backup policy)",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the backup policy",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				MarkdownDescription: "The account ID for the backup policy",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"policy_name": schema.StringAttribute{
				MarkdownDescription: "The name of the backup policy (max 100 characters)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "The integration ID for cloud provider credentials",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The cloud region where backups will be stored",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"provider_type": schema.StringAttribute{
				MarkdownDescription: "The cloud provider type (max 50 characters)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 50),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the backup policy (max 500 characters)",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
			},
			"notification_id": schema.StringAttribute{
				MarkdownDescription: "Notification channel ID for backup alerts",
				Optional:            true,
			},
			"restore_instructions": schema.StringAttribute{
				MarkdownDescription: "Instructions for restoring from backups (max 2000 characters)",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(2000),
				},
			},
			"backup_on_save": schema.BoolAttribute{
				MarkdownDescription: "Whether to trigger a backup immediately on policy creation/update",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			// Computed fields
			"status": schema.StringAttribute{
				MarkdownDescription: "Current status of the policy (Active/Inactive)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"snapshots_count": schema.Int64Attribute{
				MarkdownDescription: "Number of snapshots created by this policy",
				Computed:            true,
			},
			"last_backup_snapshot_id": schema.StringAttribute{
				MarkdownDescription: "ID of the most recent backup snapshot",
				Computed:            true,
			},
			"last_backup_time": schema.StringAttribute{
				MarkdownDescription: "Timestamp of the last backup",
				Computed:            true,
			},
			"last_backup_status": schema.StringAttribute{
				MarkdownDescription: "Status of the last backup",
				Computed:            true,
			},
			"next_backup_time": schema.StringAttribute{
				MarkdownDescription: "Timestamp of the next scheduled backup",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the policy was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the policy was last updated",
				Computed:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"schedule": schema.SingleNestedBlock{
				MarkdownDescription: "Backup schedule configuration",
				Attributes: map[string]schema.Attribute{
					"frequency": schema.StringAttribute{
						MarkdownDescription: "Backup frequency (One-time, Daily, Weekly, Monthly)",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("One-time", "Daily", "Weekly", "Monthly"),
						},
					},
					"hour": schema.Int64Attribute{
						MarkdownDescription: "Hour of day for backup (0-23)",
						Optional:            true,
						Validators: []validator.Int64{
							int64validator.Between(0, 23),
						},
					},
					"minute": schema.Int64Attribute{
						MarkdownDescription: "Minute of hour for backup (0-59)",
						Optional:            true,
						Validators: []validator.Int64{
							int64validator.Between(0, 59),
						},
					},
					"days_of_week": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "Days of week for Weekly schedule (e.g., ['Monday', 'Friday'])",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},
					},
					"monthly_schedule_type": schema.StringAttribute{
						MarkdownDescription: "Type of monthly schedule (specific_day, specific_weekday, last_day)",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("specific_day", "specific_weekday", "last_day"),
						},
					},
					"day_of_month": schema.Int64Attribute{
						MarkdownDescription: "Day of month for specific_day monthly schedule (1-31)",
						Optional:            true,
						Validators: []validator.Int64{
							int64validator.Between(1, 31),
						},
					},
					"weekday_ordinal": schema.StringAttribute{
						MarkdownDescription: "Weekday ordinal for specific_weekday monthly schedule (First, Second, Third, Fourth, Last)",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("First", "Second", "Third", "Fourth", "Last"),
						},
					},
					"weekday_name": schema.StringAttribute{
						MarkdownDescription: "Weekday name for specific_weekday monthly schedule (e.g., 'Sunday')",
						Optional:            true,
					},
					"cron_expression": schema.StringAttribute{
						MarkdownDescription: "Cron expression as alternative to explicit schedule",
						Optional:            true,
					},
				},
			},
			"scope": schema.ListNestedBlock{
				MarkdownDescription: "Resource scope configurations for backup targeting",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Scope type (tags, resource_group, asset_types, selected_resources)",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("tags", "resource_group", "asset_types", "selected_resources"),
							},
						},
						"value": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of values for this scope type",
							Required:            true,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
						},
					},
				},
			},
			"vcs": schema.SingleNestedBlock{
				MarkdownDescription: "VCS integration configuration for backup artifacts",
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						MarkdownDescription: "Project ID for VCS integration",
						Optional:            true,
					},
					"vcs_integration_id": schema.StringAttribute{
						MarkdownDescription: "VCS integration ID",
						Optional:            true,
					},
					"repo_id": schema.StringAttribute{
						MarkdownDescription: "Repository ID for storing backups",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *BackupAndDrApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BackupAndDrApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BackupAndDrApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert model to API request
	request, err := mapModelToAPIRequest(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup policy",
			fmt.Sprintf("Could not convert model to API request: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Creating backup policy", map[string]interface{}{
		"account_id":  data.AccountID.ValueString(),
		"policy_name": request.PolicyName,
	})

	// Create the policy
	createdPolicy, err := r.client.BackupAndDr.Create(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup policy",
			fmt.Sprintf("Could not create backup policy: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Created backup policy", map[string]interface{}{
		"policy_id":   createdPolicy.PolicyID,
		"policy_name": createdPolicy.PolicyName,
		"status":      createdPolicy.Status,
	})

	// Map response to model
	err = mapAPIResponseToModel(createdPolicy, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup policy",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupAndDrApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BackupAndDrApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policyID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading backup policy", map[string]interface{}{
		"account_id": data.AccountID.ValueString(),
		"policy_id":  policyID,
	})

	// Get the policy
	policy, err := r.client.BackupAndDr.Get(policyID)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "policy not found") ||
			strings.Contains(errorMsg, "not found") ||
			strings.Contains(errorMsg, "status 404") {
			tflog.Info(ctx, "Backup policy not found, removing from state", map[string]interface{}{
				"account_id": data.AccountID.ValueString(),
				"policy_id":  policyID,
			})
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading backup policy",
			fmt.Sprintf("Could not read backup policy: %s", err),
		)
		return
	}

	// Map response to model
	err = mapAPIResponseToModel(policy, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup policy",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupAndDrApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BackupAndDrApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert model to API request
	request, err := mapModelToAPIRequest(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backup policy",
			fmt.Sprintf("Could not convert model to API request: %s", err),
		)
		return
	}

	policyID := data.ID.ValueString()

	tflog.Debug(ctx, "Updating backup policy", map[string]interface{}{
		"account_id":  data.AccountID.ValueString(),
		"policy_id":   policyID,
		"policy_name": request.PolicyName,
	})

	// Convert create request to update request
	updateRequest := client.ConvertCreateToUpdate(request)

	// Update the policy
	updatedPolicy, err := r.client.BackupAndDr.Update(policyID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backup policy",
			fmt.Sprintf("Could not update backup policy: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Updated backup policy", map[string]interface{}{
		"policy_id":   updatedPolicy.PolicyID,
		"policy_name": updatedPolicy.PolicyName,
		"status":      updatedPolicy.Status,
	})

	// Map response to model
	err = mapAPIResponseToModel(updatedPolicy, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backup policy",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupAndDrApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BackupAndDrApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policyID := data.ID.ValueString()

	tflog.Debug(ctx, "Deleting backup policy", map[string]interface{}{
		"account_id": data.AccountID.ValueString(),
		"policy_id":  policyID,
	})

	err := r.client.BackupAndDr.Delete(policyID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting backup policy",
			fmt.Sprintf("Could not delete backup policy: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Deleted backup policy", map[string]interface{}{
		"account_id": data.AccountID.ValueString(),
		"policy_id":  policyID,
	})
}

func (r *BackupAndDrApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the import ID in format "account_id:policy_id"
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Error importing backup policy",
			fmt.Sprintf("Invalid import ID format. Expected 'account_id:policy_id', got: %s", req.ID),
		)
		return
	}

	accountID := parts[0]
	policyID := parts[1]

	// Set both IDs in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_id"), accountID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), policyID)...)
}
