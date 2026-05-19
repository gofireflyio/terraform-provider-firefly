package provider

import (
	"context"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringValueOrNull returns a types.String with a value if the input is non-empty,
// otherwise returns a null types.String.
func StringValueOrNull(s string) types.String {
	if s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

// mapModelToAPIRequest converts the Terraform model to an API request
func mapModelToAPIRequest(ctx context.Context, model *BackupAndDrApplicationResourceModel) (*client.PolicyCreateRequest, error) {
	request := &client.PolicyCreateRequest{
		PolicyName:    model.ApplicationName.ValueString(),
		IntegrationID: model.IntegrationID.ValueString(),
		Region:        model.Region.ValueString(),
		ProviderType:  model.ProviderType.ValueString(),
		BackupOnSave:  true, // Default value
	}

	if !model.Frequency.IsNull() {
		request.Frequency = int(model.Frequency.ValueInt64())
	}

	if !model.Description.IsNull() {
		request.Description = model.Description.ValueString()
	}

	if !model.BackupOnSave.IsNull() {
		request.BackupOnSave = model.BackupOnSave.ValueBool()
	}

	if !model.NotificationID.IsNull() {
		request.NotificationID = model.NotificationID.ValueString()
	}

	if !model.RestoreInstructions.IsNull() {
		request.RestoreInstructions = model.RestoreInstructions.ValueString()
	}

	if !model.TargetAccount.IsNull() {
		request.TargetAccount = model.TargetAccount.ValueString()
	}

	if !model.TargetRegion.IsNull() {
		request.TargetRegion = model.TargetRegion.ValueString()
	}

	if !model.AutoCreatePR.IsNull() {
		request.AutoCreatePR = model.AutoCreatePR.ValueBool()
	}

	if !model.ResilienceEnabled.IsNull() {
		request.ResilienceEnabled = model.ResilienceEnabled.ValueBool()
	}

	// Build scope array
	if len(model.Scope) > 0 {
		scopes := make([]client.ScopeConfig, len(model.Scope))
		for i, scopeModel := range model.Scope {
			var values []string
			scopeModel.Value.ElementsAs(ctx, &values, false)
			scopes[i] = client.ScopeConfig{
				Type:  scopeModel.Type.ValueString(),
				Value: values,
			}
		}
		request.Scope = scopes
	}

	// Build VCS config
	if model.VCS != nil {
		vcs := &client.VCSConfig{}

		if !model.VCS.VCSIntegrationID.IsNull() {
			vcs.VCSIntegrationID = model.VCS.VCSIntegrationID.ValueString()
		}

		if !model.VCS.RepoID.IsNull() {
			vcs.RepoID = model.VCS.RepoID.ValueString()
		}

		request.VCS = vcs
	}

	return request, nil
}

// mapAPIResponseToModel converts an API response to the Terraform model
// This function updates ALL fields from the API response, treating the API as the source of truth
func mapAPIResponseToModel(response *client.PolicyResponse, model *BackupAndDrApplicationResourceModel) error {
	// Update ID and account ID
	model.ID = types.StringValue(response.PolicyID)
	model.AccountID = types.StringValue(response.AccountID)

	// Update all user-provided fields from API (API is source of truth)
	model.ApplicationName = types.StringValue(response.PolicyName)
	model.IntegrationID = types.StringValue(response.IntegrationID)
	model.Region = types.StringValue(response.Region)
	model.ProviderType = types.StringValue(response.ProviderType)

	model.Description = StringValueOrNull(response.Description)
	model.NotificationID = StringValueOrNull(response.NotificationID)
	model.RestoreInstructions = StringValueOrNull(response.RestoreInstructions)

	model.Frequency = types.Int64Value(int64(response.Frequency))

	model.TargetAccount = StringValueOrNull(response.TargetAccount)
	model.TargetRegion = StringValueOrNull(response.TargetRegion)
	model.AutoCreatePR = types.BoolValue(response.AutoCreatePR)
	model.ResilienceEnabled = types.BoolValue(response.ResilienceEnabled)

	// Update scope array
	if len(response.Scope) > 0 {
		scopes := make([]ScopeModel, len(response.Scope))
		for i, scopeResp := range response.Scope {
			valuesList := make([]attr.Value, len(scopeResp.Value))
			for j, val := range scopeResp.Value {
				valuesList[j] = types.StringValue(val)
			}
			scopes[i] = ScopeModel{
				Type:  types.StringValue(scopeResp.Type),
				Value: types.ListValueMust(types.StringType, valuesList),
			}
		}
		model.Scope = scopes
	} else {
		model.Scope = nil
	}

	// Update VCS config
	if response.VCS != nil {
		vcsModel := &VCSModel{}
		vcsModel.VCSIntegrationID = StringValueOrNull(response.VCS.VCSIntegrationID)
		vcsModel.RepoID = StringValueOrNull(response.VCS.RepoID)
		model.VCS = vcsModel
	} else {
		model.VCS = nil
	}

	// Update computed fields
	model.Status = types.StringValue(response.Status)
	model.SnapshotsCount = types.Int64Value(int64(response.SnapshotsCount))

	model.LastBackupSnapshotID = StringValueOrNull(response.LastBackupSnapshotID)
	model.LastBackupTime = StringValueOrNull(response.LastBackupTime)
	model.LastBackupStatus = StringValueOrNull(response.LastBackupStatus)
	model.NextBackupTime = StringValueOrNull(response.NextBackupTime)

	model.CreatedAt = types.StringValue(response.CreatedAt)
	model.UpdatedAt = types.StringValue(response.UpdatedAt)

	return nil
}
