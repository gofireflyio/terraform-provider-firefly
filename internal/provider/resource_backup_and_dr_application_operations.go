package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func backupAndDRApplicationSchema() schema.Schema {
	return schema.Schema{
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
				MarkdownDescription: "The account ID associated with the backup policy",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_name": schema.StringAttribute{
				MarkdownDescription: "The name of the backup policy",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the backup policy",
				Optional:            true,
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "The cloud integration ID to use for backups",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The cloud region for the backup policy",
				Required:            true,
			},
			"provider_type": schema.StringAttribute{
				MarkdownDescription: "The cloud provider type (e.g., 'aws', 'azure', 'gcp')",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"notification_id": schema.StringAttribute{
				MarkdownDescription: "The notification channel ID for backup alerts",
				Optional:            true,
			},
			"backup_on_save": schema.BoolAttribute{
				MarkdownDescription: "Whether to trigger a backup immediately when the policy is saved",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the backup policy (Active or Inactive)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Active"),
				Validators: []validator.String{
					stringvalidator.OneOf("Active", "Inactive"),
				},
			},
			"snapshots_count": schema.Int64Attribute{
				MarkdownDescription: "The number of snapshots created by this policy",
				Computed:            true,
			},
			"last_backup_time": schema.StringAttribute{
				MarkdownDescription: "The timestamp of the last backup",
				Computed:            true,
			},
			"last_backup_status": schema.StringAttribute{
				MarkdownDescription: "The status of the last backup",
				Computed:            true,
			},
			"next_backup_time": schema.StringAttribute{
				MarkdownDescription: "The timestamp of the next scheduled backup",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The creation timestamp of the backup policy",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The last update timestamp of the backup policy",
				Computed:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"schedule": schema.SingleNestedBlock{
				MarkdownDescription: "The backup schedule configuration",
				Attributes: map[string]schema.Attribute{
					"frequency": schema.StringAttribute{
						MarkdownDescription: "The backup frequency (e.g., 'Daily', 'Weekly', 'Monthly', 'Cron')",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("Daily", "Weekly", "Monthly", "Cron"),
						},
					},
					"hour": schema.Int64Attribute{
						MarkdownDescription: "The hour of the day to run the backup (0-23)",
						Optional:            true,
					},
					"minute": schema.Int64Attribute{
						MarkdownDescription: "The minute of the hour to run the backup (0-59)",
						Optional:            true,
					},
					"days_of_week": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "Days of the week to run the backup (for Weekly frequency)",
						Optional:            true,
					},
					"monthly_schedule_type": schema.StringAttribute{
						MarkdownDescription: "The type of monthly schedule ('DayOfMonth' or 'WeekdayOfMonth')",
						Optional:            true,
					},
					"day_of_month": schema.Int64Attribute{
						MarkdownDescription: "The day of the month to run the backup (1-31, for Monthly frequency with DayOfMonth type)",
						Optional:            true,
					},
					"weekday_ordinal": schema.StringAttribute{
						MarkdownDescription: "The ordinal of the weekday (e.g., 'First', 'Second', 'Last', for Monthly with WeekdayOfMonth type)",
						Optional:            true,
					},
					"weekday_name": schema.StringAttribute{
						MarkdownDescription: "The name of the weekday (e.g., 'Monday', for Monthly with WeekdayOfMonth type)",
						Optional:            true,
					},
					"cron_expression": schema.StringAttribute{
						MarkdownDescription: "A cron expression for custom scheduling (for Cron frequency). Auto-generated by the API when using other frequency types.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
			"scope": schema.ListNestedBlock{
				MarkdownDescription: "Scope filters to limit which resources are backed up",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "The scope filter type (e.g., 'tags', 'resourceTypes', 'resourceIds')",
							Required:            true,
						},
						"value": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "The scope filter values",
							Required:            true,
						},
					},
				},
			},
			"vcs": schema.SingleNestedBlock{
				MarkdownDescription: "VCS (Version Control System) configuration for storing backups. When provided, all three attributes (project_id, vcs_integration_id, repo_id) are required.",
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						MarkdownDescription: "The Firefly project ID",
						Optional:            true,
					},
					"vcs_integration_id": schema.StringAttribute{
						MarkdownDescription: "The VCS integration ID",
						Optional:            true,
					},
					"repo_id": schema.StringAttribute{
						MarkdownDescription: "The repository ID",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *backupAndDRApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BackupAndDRApplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state BackupAndDRApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate schedule is present
	if plan.Schedule == nil {
		resp.Diagnostics.AddError(
			"Missing schedule",
			"A schedule block is required for backup policies.",
		)
		return
	}

	updateReq, err := mapModelToUpdateRequest(&plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backup policy",
			fmt.Sprintf("Could not convert model to API request: %s", err),
		)
		return
	}

	policyID := state.ID.ValueString()

	tflog.Debug(ctx, "Updating backup policy", map[string]interface{}{
		"id":   policyID,
		"name": updateReq.PolicyName,
	})

	_, err = r.client.BackupAndDR.Update(policyID, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backup policy",
			fmt.Sprintf("Could not update backup policy: %s", err),
		)
		return
	}

	// Handle status change via separate endpoint
	planStatus := plan.Status.ValueString()
	stateStatus := state.Status.ValueString()
	if planStatus != stateStatus {
		tflog.Debug(ctx, "Updating backup policy status", map[string]interface{}{
			"id":     policyID,
			"from":   stateStatus,
			"to":     planStatus,
		})

		err = r.client.BackupAndDR.SetStatus(policyID, planStatus)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating backup policy status",
				fmt.Sprintf("Could not update backup policy status: %s", err),
			)
			return
		}
	}

	// Re-fetch to get computed fields
	policy, err := r.client.BackupAndDR.Get(policyID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup policy after update",
			fmt.Sprintf("Could not read backup policy: %s", err),
		)
		return
	}

	err = mapBackupPolicyToModel(policy, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backup policy",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Updated backup policy", map[string]interface{}{
		"id":   policyID,
		"name": plan.PolicyName.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupAndDRApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BackupAndDRApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := data.ID.ValueString()

	tflog.Debug(ctx, "Deleting backup policy", map[string]interface{}{
		"id": policyID,
	})

	err := r.client.BackupAndDR.Delete(policyID, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting backup policy",
			fmt.Sprintf("Could not delete backup policy: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Deleted backup policy", map[string]interface{}{
		"id": policyID,
	})
}
