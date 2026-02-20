package provider

import (
	"context"
	"testing"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Unit tests for StringValueOrNull helper

func TestStringValueOrNull_EmptyString(t *testing.T) {
	result := StringValueOrNull("")
	if !result.IsNull() {
		t.Error("Expected null types.String for empty string, got value")
	}
}

func TestStringValueOrNull_NonEmptyString(t *testing.T) {
	result := StringValueOrNull("test-value")
	if result.IsNull() {
		t.Error("Expected non-null types.String for non-empty string, got null")
	}
	if result.ValueString() != "test-value" {
		t.Errorf("Expected value 'test-value', got '%s'", result.ValueString())
	}
}

func TestStringValueOrNull_WhitespaceString(t *testing.T) {
	result := StringValueOrNull("   ")
	if result.IsNull() {
		t.Error("Expected non-null types.String for whitespace string (not treated as empty), got null")
	}
	if result.ValueString() != "   " {
		t.Errorf("Expected value '   ', got '%s'", result.ValueString())
	}
}

func TestStringValueOrNull_SingleCharacter(t *testing.T) {
	result := StringValueOrNull("a")
	if result.IsNull() {
		t.Error("Expected non-null types.String for single character, got null")
	}
	if result.ValueString() != "a" {
		t.Errorf("Expected value 'a', got '%s'", result.ValueString())
	}
}

// Unit tests for schedule validation

func TestValidateScheduleConfig_Daily(t *testing.T) {
	schedule := &ScheduleModel{
		Frequency: types.StringValue("Daily"),
		Hour:      types.Int64Value(2),
		Minute:    types.Int64Value(30),
	}

	err := validateScheduleConfig(schedule)
	if err != nil {
		t.Errorf("Daily schedule validation failed: %v", err)
	}
}

func TestValidateScheduleConfig_Weekly_Valid(t *testing.T) {
	daysOfWeek := []attr.Value{
		types.StringValue("Monday"),
		types.StringValue("Friday"),
	}

	schedule := &ScheduleModel{
		Frequency:  types.StringValue("Weekly"),
		DaysOfWeek: types.ListValueMust(types.StringType, daysOfWeek),
		Hour:       types.Int64Value(2),
		Minute:     types.Int64Value(30),
	}

	err := validateScheduleConfig(schedule)
	if err != nil {
		t.Errorf("Weekly schedule validation failed: %v", err)
	}
}

func TestValidateScheduleConfig_Weekly_MissingDays(t *testing.T) {
	schedule := &ScheduleModel{
		Frequency:  types.StringValue("Weekly"),
		DaysOfWeek: types.ListNull(types.StringType),
		Hour:       types.Int64Value(2),
		Minute:     types.Int64Value(30),
	}

	err := validateScheduleConfig(schedule)
	if err == nil {
		t.Error("Expected error for Weekly schedule without days_of_week, got nil")
	}

	expectedMsg := "weekly schedule requires days_of_week to be specified"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestValidateScheduleConfig_Monthly_SpecificDay(t *testing.T) {
	schedule := &ScheduleModel{
		Frequency:           types.StringValue("Monthly"),
		MonthlyScheduleType: types.StringValue("specific_day"),
		DayOfMonth:          types.Int64Value(15),
		Hour:                types.Int64Value(3),
		Minute:              types.Int64Value(0),
	}

	err := validateScheduleConfig(schedule)
	if err != nil {
		t.Errorf("Monthly specific_day validation failed: %v", err)
	}
}

func TestValidateScheduleConfig_Monthly_SpecificDay_Invalid(t *testing.T) {
	schedule := &ScheduleModel{
		Frequency:           types.StringValue("Monthly"),
		MonthlyScheduleType: types.StringValue("specific_day"),
		DayOfMonth:          types.Int64Value(32),
		Hour:                types.Int64Value(3),
		Minute:              types.Int64Value(0),
	}

	err := validateScheduleConfig(schedule)
	if err == nil {
		t.Error("Expected error for day_of_month > 31, got nil")
	}
}

func TestValidateScheduleConfig_Monthly_SpecificWeekday(t *testing.T) {
	schedule := &ScheduleModel{
		Frequency:           types.StringValue("Monthly"),
		MonthlyScheduleType: types.StringValue("specific_weekday"),
		WeekdayOrdinal:      types.StringValue("First"),
		WeekdayName:         types.StringValue("Sunday"),
		Hour:                types.Int64Value(3),
		Minute:              types.Int64Value(0),
	}

	err := validateScheduleConfig(schedule)
	if err != nil {
		t.Errorf("Monthly specific_weekday validation failed: %v", err)
	}
}

func TestValidateScheduleConfig_Monthly_MissingType(t *testing.T) {
	schedule := &ScheduleModel{
		Frequency:           types.StringValue("Monthly"),
		MonthlyScheduleType: types.StringNull(),
		Hour:                types.Int64Value(3),
		Minute:              types.Int64Value(0),
	}

	err := validateScheduleConfig(schedule)
	if err == nil {
		t.Error("Expected error for Monthly without monthly_schedule_type, got nil")
	}
}

func TestValidateScheduleConfig_InvalidFrequency(t *testing.T) {
	schedule := &ScheduleModel{
		Frequency: types.StringValue("Yearly"),
	}

	err := validateScheduleConfig(schedule)
	if err == nil {
		t.Error("Expected error for invalid frequency, got nil")
	}
}

// Unit tests for model to API conversion

func TestMapModelToAPIRequest_Daily(t *testing.T) {
	ctx := context.Background()
	model := &BackupAndDrApplicationResourceModel{
		AccountID:     types.StringValue("test-account"),
		PolicyName:    types.StringValue("Test Policy"),
		IntegrationID: types.StringValue("int-123"),
		Region:        types.StringValue("us-east-1"),
		ProviderType:  types.StringValue("aws"),
		Description:   types.StringValue("Test description"),
		BackupOnSave:  types.BoolValue(true),
		Schedule: &ScheduleModel{
			Frequency: types.StringValue("Daily"),
			Hour:      types.Int64Value(2),
			Minute:    types.Int64Value(30),
		},
	}

	request, err := mapModelToAPIRequest(ctx, model)
	if err != nil {
		t.Fatalf("mapModelToAPIRequest failed: %v", err)
	}

	if request.PolicyName != "Test Policy" {
		t.Errorf("Expected PolicyName 'Test Policy', got '%s'", request.PolicyName)
	}

	if request.Schedule.Frequency != "Daily" {
		t.Errorf("Expected Frequency 'Daily', got '%s'", request.Schedule.Frequency)
	}

	if request.Schedule.Hour != 2 {
		t.Errorf("Expected Hour 2, got %d", request.Schedule.Hour)
	}

	if request.Schedule.Minute != 30 {
		t.Errorf("Expected Minute 30, got %d", request.Schedule.Minute)
	}

	if !request.BackupOnSave {
		t.Error("Expected BackupOnSave true, got false")
	}
}

func TestMapModelToAPIRequest_WithScope(t *testing.T) {
	ctx := context.Background()

	tagValues := []attr.Value{
		types.StringValue("Environment:Production"),
		types.StringValue("Backup:Required"),
	}

	assetValues := []attr.Value{
		types.StringValue("aws_instance"),
		types.StringValue("aws_db_instance"),
	}

	model := &BackupAndDrApplicationResourceModel{
		AccountID:     types.StringValue("test-account"),
		PolicyName:    types.StringValue("Scoped Policy"),
		IntegrationID: types.StringValue("int-123"),
		Region:        types.StringValue("us-east-1"),
		ProviderType:  types.StringValue("aws"),
		Schedule: &ScheduleModel{
			Frequency: types.StringValue("Daily"),
			Hour:      types.Int64Value(2),
			Minute:    types.Int64Value(0),
		},
		Scope: []ScopeModel{
			{
				Type:  types.StringValue("tags"),
				Value: types.ListValueMust(types.StringType, tagValues),
			},
			{
				Type:  types.StringValue("asset_types"),
				Value: types.ListValueMust(types.StringType, assetValues),
			},
		},
	}

	request, err := mapModelToAPIRequest(ctx, model)
	if err != nil {
		t.Fatalf("mapModelToAPIRequest failed: %v", err)
	}

	if len(request.Scope) != 2 {
		t.Fatalf("Expected 2 scope items, got %d", len(request.Scope))
	}

	if request.Scope[0].Type != "tags" {
		t.Errorf("Expected first scope type 'tags', got '%s'", request.Scope[0].Type)
	}

	if len(request.Scope[0].Value) != 2 {
		t.Errorf("Expected 2 tag values, got %d", len(request.Scope[0].Value))
	}

	if request.Scope[1].Type != "asset_types" {
		t.Errorf("Expected second scope type 'asset_types', got '%s'", request.Scope[1].Type)
	}
}

func TestMapModelToAPIRequest_WithVCS(t *testing.T) {
	ctx := context.Background()
	model := &BackupAndDrApplicationResourceModel{
		AccountID:     types.StringValue("test-account"),
		PolicyName:    types.StringValue("VCS Policy"),
		IntegrationID: types.StringValue("int-123"),
		Region:        types.StringValue("us-east-1"),
		ProviderType:  types.StringValue("aws"),
		Schedule: &ScheduleModel{
			Frequency: types.StringValue("Daily"),
			Hour:      types.Int64Value(2),
			Minute:    types.Int64Value(0),
		},
		VCS: &VCSModel{
			ProjectID:        types.StringValue("project-123"),
			VCSIntegrationID: types.StringValue("github-456"),
			RepoID:           types.StringValue("repo-789"),
		},
	}

	request, err := mapModelToAPIRequest(ctx, model)
	if err != nil {
		t.Fatalf("mapModelToAPIRequest failed: %v", err)
	}

	if request.VCS == nil {
		t.Fatal("Expected VCS config, got nil")
	}

	if request.VCS.ProjectID != "project-123" {
		t.Errorf("Expected ProjectID 'project-123', got '%s'", request.VCS.ProjectID)
	}

	if request.VCS.VCSIntegrationID != "github-456" {
		t.Errorf("Expected VCSIntegrationID 'github-456', got '%s'", request.VCS.VCSIntegrationID)
	}

	if request.VCS.RepoID != "repo-789" {
		t.Errorf("Expected RepoID 'repo-789', got '%s'", request.VCS.RepoID)
	}
}

// Unit tests for API to model conversion

func TestMapAPIResponseToModel(t *testing.T) {
	response := &client.PolicyResponse{
		PolicyID:       "policy-123",
		AccountID:      "account-456",
		PolicyName:     "Test Policy",
		IntegrationID:  "int-789",
		Region:         "us-east-1",
		ProviderType:   "aws",
		Description:    "Test description",
		BackupOnSave:   true,
		Status:         "Active",
		SnapshotsCount: 5,
		Schedule: client.ScheduleConfig{
			Frequency: "Daily",
			Hour:      2,
			Minute:    30,
		},
		CreatedAt: "2025-01-01T00:00:00Z",
		UpdatedAt: "2025-01-01T10:00:00Z",
	}

	model := &BackupAndDrApplicationResourceModel{}
	err := mapAPIResponseToModel(response, model)
	if err != nil {
		t.Fatalf("mapAPIResponseToModel failed: %v", err)
	}

	if model.ID.ValueString() != "policy-123" {
		t.Errorf("Expected ID 'policy-123', got '%s'", model.ID.ValueString())
	}

	if model.AccountID.ValueString() != "account-456" {
		t.Errorf("Expected AccountID 'account-456', got '%s'", model.AccountID.ValueString())
	}

	if model.PolicyName.ValueString() != "Test Policy" {
		t.Errorf("Expected PolicyName 'Test Policy', got '%s'", model.PolicyName.ValueString())
	}

	if model.Status.ValueString() != "Active" {
		t.Errorf("Expected Status 'Active', got '%s'", model.Status.ValueString())
	}

	if model.SnapshotsCount.ValueInt64() != 5 {
		t.Errorf("Expected SnapshotsCount 5, got %d", model.SnapshotsCount.ValueInt64())
	}

	if model.Schedule == nil {
		t.Fatal("Expected schedule, got nil")
	}

	if model.Schedule.Frequency.ValueString() != "Daily" {
		t.Errorf("Expected Frequency 'Daily', got '%s'", model.Schedule.Frequency.ValueString())
	}
}

func TestMapAPIResponseToModel_WithScope(t *testing.T) {
	response := &client.PolicyResponse{
		PolicyID:      "policy-123",
		AccountID:     "account-456",
		PolicyName:    "Scoped Policy",
		IntegrationID: "int-789",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Schedule: client.ScheduleConfig{
			Frequency: "Daily",
			Hour:      2,
			Minute:    0,
		},
		Scope: []client.ScopeConfig{
			{
				Type:  "tags",
				Value: []string{"Environment:Production", "Backup:Required"},
			},
			{
				Type:  "asset_types",
				Value: []string{"aws_instance"},
			},
		},
		Status:     "Active",
		CreatedAt:  "2025-01-01T00:00:00Z",
		UpdatedAt:  "2025-01-01T10:00:00Z",
		BackupOnSave: true,
	}

	model := &BackupAndDrApplicationResourceModel{}
	err := mapAPIResponseToModel(response, model)
	if err != nil {
		t.Fatalf("mapAPIResponseToModel failed: %v", err)
	}

	if len(model.Scope) != 2 {
		t.Fatalf("Expected 2 scope items, got %d", len(model.Scope))
	}

	if model.Scope[0].Type.ValueString() != "tags" {
		t.Errorf("Expected first scope type 'tags', got '%s'", model.Scope[0].Type.ValueString())
	}

	if model.Scope[1].Type.ValueString() != "asset_types" {
		t.Errorf("Expected second scope type 'asset_types', got '%s'", model.Scope[1].Type.ValueString())
	}
}

func TestMapAPIResponseToModel_NullOptionalFields(t *testing.T) {
	response := &client.PolicyResponse{
		PolicyID:       "policy-123",
		AccountID:      "account-456",
		PolicyName:     "Minimal Policy",
		IntegrationID:  "int-789",
		Region:         "us-east-1",
		ProviderType:   "aws",
		Schedule: client.ScheduleConfig{
			Frequency: "Daily",
		},
		Status:       "Active",
		CreatedAt:    "2025-01-01T00:00:00Z",
		UpdatedAt:    "2025-01-01T10:00:00Z",
		BackupOnSave: false,
		// All optional fields are empty
	}

	model := &BackupAndDrApplicationResourceModel{}
	err := mapAPIResponseToModel(response, model)
	if err != nil {
		t.Fatalf("mapAPIResponseToModel failed: %v", err)
	}

	if !model.Description.IsNull() {
		t.Error("Expected Description to be null")
	}

	if !model.NotificationID.IsNull() {
		t.Error("Expected NotificationID to be null")
	}

	if !model.LastBackupSnapshotID.IsNull() {
		t.Error("Expected LastBackupSnapshotID to be null")
	}

	if model.Scope != nil {
		t.Errorf("Expected Scope to be nil, got %v", model.Scope)
	}

	if model.VCS != nil {
		t.Errorf("Expected VCS to be nil, got %v", model.VCS)
	}
}
