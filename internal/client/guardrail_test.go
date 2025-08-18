package client

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestGuardrailService_CreateGuardrail(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock create guardrail
	mockServer.AddHandler("/v2/guardrails", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var createReq GuardrailRule
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if createReq.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		// Create response with generated ID
		createResp := CreateGuardrailResponse{
			RuleID:         "generated-rule-id",
			NotificationID: "notification-123",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(createResp)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	thresholdAmount := 500.0
	guardrail := &GuardrailRule{
		Name:      "Test Cost Guardrail",
		Type:      "cost",
		IsEnabled: true,
		Severity:  2,
		Scope: &GuardrailScope{
			Workspaces: &IncludeExcludeWildcard{
				Include: []string{"production-*"},
			},
		},
		Criteria: &GuardrailCriteria{
			Cost: &CostCriteria{
				ThresholdAmount: &thresholdAmount,
			},
		},
	}

	response, err := client.Guardrails.CreateGuardrail(guardrail)
	if err != nil {
		t.Fatalf("CreateGuardrail failed: %v", err)
	}

	if response.RuleID != "generated-rule-id" {
		t.Errorf("Expected rule ID 'generated-rule-id', got '%s'", response.RuleID)
	}

	if response.NotificationID != "notification-123" {
		t.Errorf("Expected notification ID 'notification-123', got '%s'", response.NotificationID)
	}
}

func TestGuardrailService_UpdateGuardrail(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock update guardrail
	mockServer.AddHandler("/v2/guardrails/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ruleID := r.URL.Path[len("/v2/guardrails/"):]
		if ruleID != "test-rule-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var updateReq GuardrailRule
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Return success response
		updateResp := UpdateGuardrailResponse{
			ID:        ruleID,
			Name:      updateReq.Name,
			Enabled:   updateReq.IsEnabled,
			UpdatedAt: "2023-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updateResp)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	guardrail := &GuardrailRule{
		Name:      "Updated Guardrail",
		IsEnabled: false,
		Severity:  1,
	}

	response, err := client.Guardrails.UpdateGuardrail("test-rule-id", guardrail)
	if err != nil {
		t.Fatalf("UpdateGuardrail failed: %v", err)
	}

	if response.ID != "test-rule-id" {
		t.Errorf("Expected rule ID 'test-rule-id', got '%s'", response.ID)
	}

	if response.Name != "Updated Guardrail" {
		t.Errorf("Expected name 'Updated Guardrail', got '%s'", response.Name)
	}

	if response.Enabled != false {
		t.Errorf("Expected enabled to be false, got %v", response.Enabled)
	}
}

func TestGuardrailService_DeleteGuardrail(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock delete guardrail
	mockServer.AddHandler("/v2/guardrails/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ruleID := r.URL.Path[len("/v2/guardrails/"):]
		if ruleID != "test-rule-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Return success response
		deleteResp := DeleteGuardrailResponse{
			Status:  200,
			Message: "Guardrail deleted successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deleteResp)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	response, err := client.Guardrails.DeleteGuardrail("test-rule-id")
	if err != nil {
		t.Fatalf("DeleteGuardrail failed: %v", err)
	}

	if response.Status != 200 {
		t.Errorf("Expected status 200, got %d", response.Status)
	}

	if response.Message != "Guardrail deleted successfully" {
		t.Errorf("Expected success message, got '%s'", response.Message)
	}
}