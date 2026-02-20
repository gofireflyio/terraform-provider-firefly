package provider

import (
	"context"
	"fmt"

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
	// Validate schedule first
	if model.Schedule == nil {
		return nil, fmt.Errorf("schedule is required")
	}

	if err := validateScheduleConfig(model.Schedule); err != nil {
		return nil, err
	}

	// Build schedule config
	schedule := client.ScheduleConfig{
		Frequency: model.Schedule.Frequency.ValueString(),
	}

	if !model.Schedule.Hour.IsNull() {
		schedule.Hour = int(model.Schedule.Hour.ValueInt64())
	}

	if !model.Schedule.Minute.IsNull() {
		schedule.Minute = int(model.Schedule.Minute.ValueInt64())
	}

	if !model.Schedule.DaysOfWeek.IsNull() {
		var daysOfWeek []string
		model.Schedule.DaysOfWeek.ElementsAs(ctx, &daysOfWeek, false)
		schedule.DaysOfWeek = daysOfWeek
	}

	if !model.Schedule.MonthlyScheduleType.IsNull() {
		schedule.MonthlyScheduleType = model.Schedule.MonthlyScheduleType.ValueString()
	}

	if !model.Schedule.DayOfMonth.IsNull() {
		schedule.DayOfMonth = int(model.Schedule.DayOfMonth.ValueInt64())
	}

	if !model.Schedule.WeekdayOrdinal.IsNull() {
		schedule.WeekdayOrdinal = model.Schedule.WeekdayOrdinal.ValueString()
	}

	if !model.Schedule.WeekdayName.IsNull() {
		schedule.WeekdayName = model.Schedule.WeekdayName.ValueString()
	}

	if !model.Schedule.CronExpression.IsNull() {
		schedule.CronExpression = model.Schedule.CronExpression.ValueString()
	}

	// Build base request
	request := &client.PolicyCreateRequest{
		PolicyName:    model.PolicyName.ValueString(),
		IntegrationID: model.IntegrationID.ValueString(),
		Region:        model.Region.ValueString(),
		ProviderType:  model.ProviderType.ValueString(),
		Schedule:      schedule,
		BackupOnSave:  true, // Default value
	}

	// Optional fields
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

		if !model.VCS.ProjectID.IsNull() {
			vcs.ProjectID = model.VCS.ProjectID.ValueString()
		}

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
	model.PolicyName = types.StringValue(response.PolicyName)
	model.IntegrationID = types.StringValue(response.IntegrationID)
	model.Region = types.StringValue(response.Region)
	model.ProviderType = types.StringValue(response.ProviderType)

	model.Description = StringValueOrNull(response.Description)

	// Note: backup_on_save is write-only (only in CREATE requests), not returned by API
	// We preserve the value from the plan/config instead of reading from API response

	model.NotificationID = StringValueOrNull(response.NotificationID)

	model.RestoreInstructions = StringValueOrNull(response.RestoreInstructions)

	// Update schedule (always present)
	scheduleModel := &ScheduleModel{
		Frequency: types.StringValue(response.Schedule.Frequency),
	}

	if response.Schedule.Hour != 0 || response.Schedule.Minute != 0 {
		scheduleModel.Hour = types.Int64Value(int64(response.Schedule.Hour))
		scheduleModel.Minute = types.Int64Value(int64(response.Schedule.Minute))
	} else {
		scheduleModel.Hour = types.Int64Null()
		scheduleModel.Minute = types.Int64Null()
	}

	if len(response.Schedule.DaysOfWeek) > 0 {
		daysValues := make([]attr.Value, len(response.Schedule.DaysOfWeek))
		for i, day := range response.Schedule.DaysOfWeek {
			daysValues[i] = types.StringValue(day)
		}
		scheduleModel.DaysOfWeek = types.ListValueMust(types.StringType, daysValues)
	} else {
		scheduleModel.DaysOfWeek = types.ListNull(types.StringType)
	}

	scheduleModel.MonthlyScheduleType = StringValueOrNull(response.Schedule.MonthlyScheduleType)

	if response.Schedule.DayOfMonth != 0 {
		scheduleModel.DayOfMonth = types.Int64Value(int64(response.Schedule.DayOfMonth))
	} else {
		scheduleModel.DayOfMonth = types.Int64Null()
	}

	scheduleModel.WeekdayOrdinal = StringValueOrNull(response.Schedule.WeekdayOrdinal)

	scheduleModel.WeekdayName = StringValueOrNull(response.Schedule.WeekdayName)

	scheduleModel.CronExpression = StringValueOrNull(response.Schedule.CronExpression)

	model.Schedule = scheduleModel

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

		vcsModel.ProjectID = StringValueOrNull(response.VCS.ProjectID)

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

// validateScheduleConfig validates the schedule configuration based on frequency
func validateScheduleConfig(schedule *ScheduleModel) error {
	frequency := schedule.Frequency.ValueString()

	switch frequency {
	case "One-time", "Daily":
		// No special validation needed
		return nil

	case "Weekly":
		// Weekly requires days_of_week
		if schedule.DaysOfWeek.IsNull() || len(schedule.DaysOfWeek.Elements()) == 0 {
			return fmt.Errorf("weekly schedule requires days_of_week to be specified")
		}
		return nil

	case "Monthly":
		// Monthly requires monthly_schedule_type
		if schedule.MonthlyScheduleType.IsNull() {
			return fmt.Errorf("monthly schedule requires monthly_schedule_type to be specified")
		}

		scheduleType := schedule.MonthlyScheduleType.ValueString()

		switch scheduleType {
		case "specific_day":
			// Requires day_of_month
			if schedule.DayOfMonth.IsNull() {
				return fmt.Errorf("monthly schedule with specific_day requires day_of_month to be specified")
			}
			dayOfMonth := schedule.DayOfMonth.ValueInt64()
			if dayOfMonth < 1 || dayOfMonth > 31 {
				return fmt.Errorf("day_of_month must be between 1 and 31, got %d", dayOfMonth)
			}
			return nil

		case "specific_weekday":
			// Requires weekday_ordinal and weekday_name
			if schedule.WeekdayOrdinal.IsNull() {
				return fmt.Errorf("monthly schedule with specific_weekday requires weekday_ordinal to be specified")
			}
			if schedule.WeekdayName.IsNull() {
				return fmt.Errorf("monthly schedule with specific_weekday requires weekday_name to be specified")
			}
			return nil

		case "last_day":
			// No additional requirements
			return nil

		default:
			return fmt.Errorf("invalid monthly_schedule_type: %s (must be specific_day, specific_weekday, or last_day)", scheduleType)
		}

	default:
		return fmt.Errorf("invalid frequency: %s (must be One-time, Daily, Weekly, or Monthly)", frequency)
	}
}
