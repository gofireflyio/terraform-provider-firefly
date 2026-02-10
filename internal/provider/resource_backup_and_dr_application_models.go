package provider

import (
	"context"
	"fmt"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BackupAndDRApplicationResourceModel struct {
	ID               types.String  `tfsdk:"id"`
	AccountID        types.String  `tfsdk:"account_id"`
	PolicyName       types.String  `tfsdk:"policy_name"`
	Description      types.String  `tfsdk:"description"`
	IntegrationID    types.String  `tfsdk:"integration_id"`
	Region           types.String  `tfsdk:"region"`
	ProviderType     types.String  `tfsdk:"provider_type"`
	Schedule         *ScheduleModel `tfsdk:"schedule"`
	Scope            []ScopeModel  `tfsdk:"scope"`
	NotificationID   types.String  `tfsdk:"notification_id"`
	BackupOnSave     types.Bool    `tfsdk:"backup_on_save"`
	Status           types.String  `tfsdk:"status"`
	SnapshotsCount   types.Int64   `tfsdk:"snapshots_count"`
	LastBackupTime   types.String  `tfsdk:"last_backup_time"`
	LastBackupStatus types.String  `tfsdk:"last_backup_status"`
	NextBackupTime   types.String  `tfsdk:"next_backup_time"`
	VCS              *VCSModel     `tfsdk:"vcs"`
	CreatedAt        types.String  `tfsdk:"created_at"`
	UpdatedAt        types.String  `tfsdk:"updated_at"`
}

type ScheduleModel struct {
	Frequency           types.String `tfsdk:"frequency"`
	Hour                types.Int64  `tfsdk:"hour"`
	Minute              types.Int64  `tfsdk:"minute"`
	DaysOfWeek          types.List   `tfsdk:"days_of_week"`
	MonthlyScheduleType types.String `tfsdk:"monthly_schedule_type"`
	DayOfMonth          types.Int64  `tfsdk:"day_of_month"`
	WeekdayOrdinal      types.String `tfsdk:"weekday_ordinal"`
	WeekdayName         types.String `tfsdk:"weekday_name"`
	CronExpression      types.String `tfsdk:"cron_expression"`
}

type ScopeModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.List   `tfsdk:"value"`
}

type VCSModel struct {
	ProjectID        types.String `tfsdk:"project_id"`
	VCSIntegrationID types.String `tfsdk:"vcs_integration_id"`
	RepoID           types.String `tfsdk:"repo_id"`
}

func mapModelToCreateRequest(model *BackupAndDRApplicationResourceModel) (*client.BackupPolicyCreateRequest, error) {
	req := &client.BackupPolicyCreateRequest{
		PolicyName:    model.PolicyName.ValueString(),
		IntegrationID: model.IntegrationID.ValueString(),
		Region:        model.Region.ValueString(),
		ProviderType:  model.ProviderType.ValueString(),
	}

	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		req.Description = model.Description.ValueString()
	}

	if !model.NotificationID.IsNull() && !model.NotificationID.IsUnknown() {
		req.NotificationID = model.NotificationID.ValueString()
	}

	if !model.BackupOnSave.IsNull() && !model.BackupOnSave.IsUnknown() {
		v := model.BackupOnSave.ValueBool()
		req.BackupOnSave = &v
	}

	if model.Schedule != nil {
		schedule, err := mapScheduleModelToConfig(model.Schedule)
		if err != nil {
			return nil, fmt.Errorf("error converting schedule: %w", err)
		}
		req.Schedule = schedule
	}

	if len(model.Scope) > 0 {
		scopes, err := mapScopeModelsToConfigs(model.Scope)
		if err != nil {
			return nil, fmt.Errorf("error converting scope: %w", err)
		}
		req.Scope = scopes
	}

	if model.VCS != nil && isVCSPopulated(model.VCS) {
		req.VCS = &client.VCSConfig{
			ProjectID:        model.VCS.ProjectID.ValueString(),
			VCSIntegrationID: model.VCS.VCSIntegrationID.ValueString(),
			RepoID:           model.VCS.RepoID.ValueString(),
		}
	}

	return req, nil
}

func mapModelToUpdateRequest(model *BackupAndDRApplicationResourceModel) (*client.BackupPolicyUpdateRequest, error) {
	req := &client.BackupPolicyUpdateRequest{
		PolicyName:    model.PolicyName.ValueString(),
		IntegrationID: model.IntegrationID.ValueString(),
		Region:        model.Region.ValueString(),
		ProviderType:  model.ProviderType.ValueString(),
	}

	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		req.Description = model.Description.ValueString()
	}

	if !model.NotificationID.IsNull() && !model.NotificationID.IsUnknown() {
		req.NotificationID = model.NotificationID.ValueString()
	}

	if !model.BackupOnSave.IsNull() && !model.BackupOnSave.IsUnknown() {
		v := model.BackupOnSave.ValueBool()
		req.BackupOnSave = &v
	}

	if model.Schedule != nil {
		schedule, err := mapScheduleModelToConfig(model.Schedule)
		if err != nil {
			return nil, fmt.Errorf("error converting schedule: %w", err)
		}
		req.Schedule = schedule
	}

	if len(model.Scope) > 0 {
		scopes, err := mapScopeModelsToConfigs(model.Scope)
		if err != nil {
			return nil, fmt.Errorf("error converting scope: %w", err)
		}
		req.Scope = scopes
	}

	if model.VCS != nil && isVCSPopulated(model.VCS) {
		req.VCS = &client.VCSConfig{
			ProjectID:        model.VCS.ProjectID.ValueString(),
			VCSIntegrationID: model.VCS.VCSIntegrationID.ValueString(),
			RepoID:           model.VCS.RepoID.ValueString(),
		}
	}

	return req, nil
}

func mapBackupPolicyToModel(policy *client.BackupPolicy, model *BackupAndDRApplicationResourceModel) error {
	model.ID = types.StringValue(policy.ID)
	model.PolicyName = types.StringValue(policy.PolicyName)
	model.IntegrationID = types.StringValue(policy.IntegrationID)
	model.Region = types.StringValue(policy.Region)
	model.ProviderType = types.StringValue(policy.ProviderType)

	if policy.AccountID != "" {
		model.AccountID = types.StringValue(policy.AccountID)
	} else {
		model.AccountID = types.StringNull()
	}

	if policy.Description != "" {
		model.Description = types.StringValue(policy.Description)
	} else {
		model.Description = types.StringNull()
	}

	if policy.NotificationID != "" {
		model.NotificationID = types.StringValue(policy.NotificationID)
	} else {
		model.NotificationID = types.StringNull()
	}

	if policy.BackupOnSave != nil {
		model.BackupOnSave = types.BoolValue(*policy.BackupOnSave)
	} else {
		model.BackupOnSave = types.BoolNull()
	}

	if policy.Status != "" {
		model.Status = types.StringValue(policy.Status)
	} else {
		model.Status = types.StringValue("Active")
	}

	model.SnapshotsCount = types.Int64Value(int64(policy.SnapshotsCount))

	if policy.LastBackupTime != "" {
		model.LastBackupTime = types.StringValue(policy.LastBackupTime)
	} else {
		model.LastBackupTime = types.StringNull()
	}

	if policy.LastBackupStatus != "" {
		model.LastBackupStatus = types.StringValue(policy.LastBackupStatus)
	} else {
		model.LastBackupStatus = types.StringNull()
	}

	if policy.NextBackupTime != "" {
		model.NextBackupTime = types.StringValue(policy.NextBackupTime)
	} else {
		model.NextBackupTime = types.StringNull()
	}

	if policy.CreatedAt != "" {
		model.CreatedAt = types.StringValue(policy.CreatedAt)
	} else {
		model.CreatedAt = types.StringNull()
	}

	if policy.UpdatedAt != "" {
		model.UpdatedAt = types.StringValue(policy.UpdatedAt)
	} else {
		model.UpdatedAt = types.StringNull()
	}

	// Map schedule
	if policy.Schedule != nil {
		model.Schedule = mapScheduleConfigToModel(policy.Schedule)
	} else {
		model.Schedule = nil
	}

	// Map scope
	if len(policy.Scope) > 0 {
		scopeModels, err := mapScopeConfigsToModels(policy.Scope)
		if err != nil {
			return fmt.Errorf("error converting scope from API: %w", err)
		}
		model.Scope = scopeModels
	} else {
		model.Scope = nil
	}

	// Map VCS
	if policy.VCS != nil {
		model.VCS = &VCSModel{
			ProjectID:        types.StringValue(policy.VCS.ProjectID),
			VCSIntegrationID: types.StringValue(policy.VCS.VCSIntegrationID),
			RepoID:           types.StringValue(policy.VCS.RepoID),
		}
	} else {
		model.VCS = nil
	}

	return nil
}

func mapScheduleModelToConfig(m *ScheduleModel) (*client.ScheduleConfig, error) {
	config := &client.ScheduleConfig{
		Frequency: m.Frequency.ValueString(),
	}

	if !m.Hour.IsNull() && !m.Hour.IsUnknown() {
		v := int(m.Hour.ValueInt64())
		config.Hour = &v
	}

	if !m.Minute.IsNull() && !m.Minute.IsUnknown() {
		v := int(m.Minute.ValueInt64())
		config.Minute = &v
	}

	if !m.DaysOfWeek.IsNull() && !m.DaysOfWeek.IsUnknown() {
		var days []string
		diags := m.DaysOfWeek.ElementsAs(context.Background(), &days, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting days_of_week: %v", diags)
		}
		config.DaysOfWeek = days
	}

	if !m.MonthlyScheduleType.IsNull() && !m.MonthlyScheduleType.IsUnknown() {
		config.MonthlyScheduleType = m.MonthlyScheduleType.ValueString()
	}

	if !m.DayOfMonth.IsNull() && !m.DayOfMonth.IsUnknown() {
		v := int(m.DayOfMonth.ValueInt64())
		config.DayOfMonth = &v
	}

	if !m.WeekdayOrdinal.IsNull() && !m.WeekdayOrdinal.IsUnknown() {
		config.WeekdayOrdinal = m.WeekdayOrdinal.ValueString()
	}

	if !m.WeekdayName.IsNull() && !m.WeekdayName.IsUnknown() {
		config.WeekdayName = m.WeekdayName.ValueString()
	}

	if !m.CronExpression.IsNull() && !m.CronExpression.IsUnknown() {
		config.CronExpression = m.CronExpression.ValueString()
	}

	return config, nil
}

func mapScheduleConfigToModel(config *client.ScheduleConfig) *ScheduleModel {
	m := &ScheduleModel{
		Frequency: types.StringValue(config.Frequency),
	}

	if config.Hour != nil {
		m.Hour = types.Int64Value(int64(*config.Hour))
	} else {
		m.Hour = types.Int64Null()
	}

	if config.Minute != nil {
		m.Minute = types.Int64Value(int64(*config.Minute))
	} else {
		m.Minute = types.Int64Null()
	}

	if len(config.DaysOfWeek) > 0 {
		list, diags := types.ListValueFrom(context.Background(), types.StringType, config.DaysOfWeek)
		if !diags.HasError() {
			m.DaysOfWeek = list
		} else {
			m.DaysOfWeek = types.ListNull(types.StringType)
		}
	} else {
		m.DaysOfWeek = types.ListNull(types.StringType)
	}

	if config.MonthlyScheduleType != "" {
		m.MonthlyScheduleType = types.StringValue(config.MonthlyScheduleType)
	} else {
		m.MonthlyScheduleType = types.StringNull()
	}

	if config.DayOfMonth != nil {
		m.DayOfMonth = types.Int64Value(int64(*config.DayOfMonth))
	} else {
		m.DayOfMonth = types.Int64Null()
	}

	if config.WeekdayOrdinal != "" {
		m.WeekdayOrdinal = types.StringValue(config.WeekdayOrdinal)
	} else {
		m.WeekdayOrdinal = types.StringNull()
	}

	if config.WeekdayName != "" {
		m.WeekdayName = types.StringValue(config.WeekdayName)
	} else {
		m.WeekdayName = types.StringNull()
	}

	if config.CronExpression != "" {
		m.CronExpression = types.StringValue(config.CronExpression)
	} else {
		m.CronExpression = types.StringNull()
	}

	return m
}

func mapScopeModelsToConfigs(models []ScopeModel) ([]client.ScopeConfig, error) {
	configs := make([]client.ScopeConfig, len(models))
	for i, m := range models {
		var values []string
		diags := m.Value.ElementsAs(context.Background(), &values, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting scope value at index %d: %v", i, diags)
		}
		configs[i] = client.ScopeConfig{
			Type:  m.Type.ValueString(),
			Value: values,
		}
	}
	return configs, nil
}

func isVCSPopulated(vcs *VCSModel) bool {
	if vcs == nil {
		return false
	}
	return (!vcs.ProjectID.IsNull() && !vcs.ProjectID.IsUnknown()) ||
		(!vcs.VCSIntegrationID.IsNull() && !vcs.VCSIntegrationID.IsUnknown()) ||
		(!vcs.RepoID.IsNull() && !vcs.RepoID.IsUnknown())
}

func mapScopeConfigsToModels(configs []client.ScopeConfig) ([]ScopeModel, error) {
	models := make([]ScopeModel, len(configs))
	for i, c := range configs {
		list, diags := types.ListValueFrom(context.Background(), types.StringType, c.Value)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting scope value at index %d: %v", i, diags)
		}
		models[i] = ScopeModel{
			Type:  types.StringValue(c.Type),
			Value: list,
		}
	}
	return models, nil
}
