package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &backupAndDRApplicationResource{}
var _ resource.ResourceWithImportState = &backupAndDRApplicationResource{}

func NewBackupAndDRApplicationResource() resource.Resource {
	return &backupAndDRApplicationResource{}
}

type backupAndDRApplicationResource struct {
	client *client.Client
}

func (r *backupAndDRApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_and_dr_application"
}

func (r *backupAndDRApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = backupAndDRApplicationSchema()
}

func (r *backupAndDRApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *backupAndDRApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BackupAndDRApplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate schedule is present
	if data.Schedule == nil {
		resp.Diagnostics.AddError(
			"Missing schedule",
			"A schedule block is required for backup policies.",
		)
		return
	}

	createReq, err := mapModelToCreateRequest(&data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup policy",
			fmt.Sprintf("Could not convert model to API request: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Creating backup policy", map[string]interface{}{
		"name": createReq.PolicyName,
	})

	createdPolicy, err := r.client.BackupAndDR.Create(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup policy",
			fmt.Sprintf("Could not create backup policy: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Created backup policy", map[string]interface{}{
		"id":   createdPolicy.ID,
		"name": createdPolicy.PolicyName,
	})

	// If desired status is Inactive, toggle it after creation
	desiredStatus := data.Status.ValueString()
	if desiredStatus == "Inactive" {
		err = r.client.BackupAndDR.SetStatus(createdPolicy.ID, "Inactive")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting backup policy status",
				fmt.Sprintf("Policy was created but status could not be set to Inactive: %s", err),
			)
			return
		}
	}

	// Re-fetch to get all computed fields
	policy, err := r.client.BackupAndDR.Get(createdPolicy.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup policy after creation",
			fmt.Sprintf("Could not read backup policy: %s", err),
		)
		return
	}

	err = mapBackupPolicyToModel(policy, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup policy",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *backupAndDRApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BackupAndDRApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading backup policy", map[string]interface{}{
		"id": policyID,
	})

	policy, err := r.client.BackupAndDR.Get(policyID)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "not found") ||
			strings.Contains(errorMsg, "status 404") {
			tflog.Info(ctx, "Backup policy not found, removing from state", map[string]interface{}{
				"id": policyID,
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

	err = mapBackupPolicyToModel(policy, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup policy",
			fmt.Sprintf("Could not map API response to model: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *backupAndDRApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
