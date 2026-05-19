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

// Unit tests for model to API conversion

func TestMapModelToAPIRequest_WithFrequency(t *testing.T) {
	ctx := context.Background()
	model := &BackupAndDrApplicationResourceModel{
		AccountID:       types.StringValue("test-account"),
		ApplicationName: types.StringValue("Test Policy"),
		IntegrationID:   types.StringValue("int-123"),
		Region:          types.StringValue("us-east-1"),
		ProviderType:    types.StringValue("aws"),
		Description:     types.StringValue("Test description"),
		BackupOnSave:    types.BoolValue(true),
		Frequency:       types.Int64Value(24),
	}

	request, err := mapModelToAPIRequest(ctx, model)
	if err != nil {
		t.Fatalf("mapModelToAPIRequest failed: %v", err)
	}

	if request.PolicyName != "Test Policy" {
		t.Errorf("Expected PolicyName 'Test Policy', got '%s'", request.PolicyName)
	}

	if request.Frequency != 24 {
		t.Errorf("Expected Frequency 24, got %d", request.Frequency)
	}

	if !request.BackupOnSave {
		t.Error("Expected BackupOnSave true, got false")
	}
}

func TestMapModelToAPIRequest_FrequencyNull(t *testing.T) {
	ctx := context.Background()
	model := &BackupAndDrApplicationResourceModel{
		AccountID:       types.StringValue("test-account"),
		ApplicationName: types.StringValue("No Frequency Policy"),
		IntegrationID:   types.StringValue("int-123"),
		Region:          types.StringValue("us-east-1"),
		ProviderType:    types.StringValue("aws"),
		Frequency:       types.Int64Null(),
	}

	request, err := mapModelToAPIRequest(ctx, model)
	if err != nil {
		t.Fatalf("mapModelToAPIRequest failed: %v", err)
	}

	if request.Frequency != 0 {
		t.Errorf("Expected Frequency 0 (not set), got %d", request.Frequency)
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
		AccountID:       types.StringValue("test-account"),
		ApplicationName: types.StringValue("Scoped Policy"),
		IntegrationID:   types.StringValue("int-123"),
		Region:          types.StringValue("us-east-1"),
		ProviderType:    types.StringValue("aws"),
		Frequency:       types.Int64Value(8),
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
		AccountID:       types.StringValue("test-account"),
		ApplicationName: types.StringValue("VCS Policy"),
		IntegrationID:   types.StringValue("int-123"),
		Region:          types.StringValue("us-east-1"),
		ProviderType:    types.StringValue("aws"),
		Frequency:       types.Int64Value(24),
		VCS: &VCSModel{
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

	if request.VCS.VCSIntegrationID != "github-456" {
		t.Errorf("Expected VCSIntegrationID 'github-456', got '%s'", request.VCS.VCSIntegrationID)
	}

	if request.VCS.RepoID != "repo-789" {
		t.Errorf("Expected RepoID 'repo-789', got '%s'", request.VCS.RepoID)
	}
}

func TestMapModelToAPIRequest_WithResilienceFields(t *testing.T) {
	ctx := context.Background()
	model := &BackupAndDrApplicationResourceModel{
		AccountID:         types.StringValue("test-account"),
		ApplicationName:   types.StringValue("Resilience Policy"),
		IntegrationID:     types.StringValue("int-123"),
		Region:            types.StringValue("us-east-1"),
		ProviderType:      types.StringValue("aws"),
		Frequency:         types.Int64Value(4),
		TargetAccount:     types.StringValue("target-int-456"),
		TargetRegion:      types.StringValue("eu-west-1"),
		AutoCreatePR:      types.BoolValue(true),
		ResilienceEnabled: types.BoolValue(true),
	}

	request, err := mapModelToAPIRequest(ctx, model)
	if err != nil {
		t.Fatalf("mapModelToAPIRequest failed: %v", err)
	}

	if request.TargetAccount != "target-int-456" {
		t.Errorf("Expected TargetAccount 'target-int-456', got '%s'", request.TargetAccount)
	}

	if request.TargetRegion != "eu-west-1" {
		t.Errorf("Expected TargetRegion 'eu-west-1', got '%s'", request.TargetRegion)
	}

	if !request.AutoCreatePR {
		t.Error("Expected AutoCreatePR true, got false")
	}

	if !request.ResilienceEnabled {
		t.Error("Expected ResilienceEnabled true, got false")
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
		Frequency:      24,
		Status:         "Active",
		SnapshotsCount: 5,
		CreatedAt:      "2025-01-01T00:00:00Z",
		UpdatedAt:      "2025-01-01T10:00:00Z",
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

	if model.ApplicationName.ValueString() != "Test Policy" {
		t.Errorf("Expected ApplicationName 'Test Policy', got '%s'", model.ApplicationName.ValueString())
	}

	if model.Status.ValueString() != "Active" {
		t.Errorf("Expected Status 'Active', got '%s'", model.Status.ValueString())
	}

	if model.SnapshotsCount.ValueInt64() != 5 {
		t.Errorf("Expected SnapshotsCount 5, got %d", model.SnapshotsCount.ValueInt64())
	}

	if model.Frequency.ValueInt64() != 24 {
		t.Errorf("Expected Frequency 24, got %d", model.Frequency.ValueInt64())
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
		Frequency:     8,
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
		Status:    "Active",
		CreatedAt: "2025-01-01T00:00:00Z",
		UpdatedAt: "2025-01-01T10:00:00Z",
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

func TestMapAPIResponseToModel_WithResilienceFields(t *testing.T) {
	response := &client.PolicyResponse{
		PolicyID:          "policy-123",
		AccountID:         "account-456",
		PolicyName:        "Resilience Policy",
		IntegrationID:     "int-789",
		Region:            "us-east-1",
		ProviderType:      "aws",
		Frequency:         4,
		TargetAccount:     "target-int-456",
		TargetRegion:      "eu-west-1",
		AutoCreatePR:      true,
		ResilienceEnabled: true,
		Status:            "Active",
		CreatedAt:         "2025-01-01T00:00:00Z",
		UpdatedAt:         "2025-01-01T10:00:00Z",
	}

	model := &BackupAndDrApplicationResourceModel{}
	err := mapAPIResponseToModel(response, model)
	if err != nil {
		t.Fatalf("mapAPIResponseToModel failed: %v", err)
	}

	if model.TargetAccount.ValueString() != "target-int-456" {
		t.Errorf("Expected TargetAccount 'target-int-456', got '%s'", model.TargetAccount.ValueString())
	}

	if model.TargetRegion.ValueString() != "eu-west-1" {
		t.Errorf("Expected TargetRegion 'eu-west-1', got '%s'", model.TargetRegion.ValueString())
	}

	if !model.AutoCreatePR.ValueBool() {
		t.Error("Expected AutoCreatePR true, got false")
	}

	if !model.ResilienceEnabled.ValueBool() {
		t.Error("Expected ResilienceEnabled true, got false")
	}
}

func TestMapAPIResponseToModel_NullOptionalFields(t *testing.T) {
	response := &client.PolicyResponse{
		PolicyID:      "policy-123",
		AccountID:     "account-456",
		PolicyName:    "Minimal Policy",
		IntegrationID: "int-789",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Frequency:     0,
		Status:        "Active",
		CreatedAt:     "2025-01-01T00:00:00Z",
		UpdatedAt:     "2025-01-01T10:00:00Z",
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

	if !model.TargetAccount.IsNull() {
		t.Error("Expected TargetAccount to be null")
	}

	if !model.TargetRegion.IsNull() {
		t.Error("Expected TargetRegion to be null")
	}

	if model.Scope != nil {
		t.Errorf("Expected Scope to be nil, got %v", model.Scope)
	}

	if model.VCS != nil {
		t.Errorf("Expected VCS to be nil, got %v", model.VCS)
	}
}
