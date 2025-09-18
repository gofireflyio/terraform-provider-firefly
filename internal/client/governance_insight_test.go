package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateGovernanceInsight(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle authentication endpoint
		if r.URL.Path == "/v2/login" {
			authResp := AuthResponse{
				AccessToken: "test-token",
				ExpiresAt:   1234567890,
				TokenType:   "Bearer",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(authResp)
			return
		}

		// Check the path
		if r.URL.Path != "/v2/governance/insights/create" {
			t.Errorf("Expected path /v2/governance/insights/create, got %s", r.URL.Path)
		}

		// Check the method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Parse request body
		var insight GovernanceInsight
		if err := json.NewDecoder(r.Body).Decode(&insight); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify severity is an integer
		if insight.Severity != 4 {
			t.Errorf("Expected severity 4, got %d", insight.Severity)
		}

		// Return a mock response
		response := GovernanceInsight{
			ID:          "test-id-123",
			Name:        insight.Name,
			Description: insight.Description,
			Code:        insight.Code,
			Type:        insight.Type,
			ProviderIDs: insight.ProviderIDs,
			Severity:    insight.Severity,
			Category:    insight.Category,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create a client with the test server URL
	config := Config{
		AccessKey: "test-key",
		SecretKey: "test-secret",
		APIURL:    server.URL,
	}
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create a governance insight
	insight := &GovernanceInsight{
		Name:        "Test Policy",
		Description: "Test description",
		Code:        "package firefly\nfirefly { true }",
		Type:        []string{"aws_instance"},
		ProviderIDs: []string{"aws_all"},
		Severity:    4, // Integer value for "HIGH"
		Category:    "Optimization",
	}

	result, err := client.GovernanceInsights.CreateGovernanceInsight(insight)
	if err != nil {
		t.Fatalf("Failed to create governance insight: %v", err)
	}

	// Verify the result
	if result.ID != "test-id-123" {
		t.Errorf("Expected ID test-id-123, got %s", result.ID)
	}
	if result.Severity != 4 {
		t.Errorf("Expected severity 4, got %d", result.Severity)
	}
}

func TestDeleteGovernanceInsight(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle authentication endpoint
		if r.URL.Path == "/v2/login" {
			authResp := AuthResponse{
				AccessToken: "test-token",
				ExpiresAt:   1234567890,
				TokenType:   "Bearer",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(authResp)
			return
		}

		// Check the path
		if r.URL.Path != "/v2/governance/insights/test-id-123" {
			t.Errorf("Expected path /v2/governance/insights/test-id-123, got %s", r.URL.Path)
		}

		// Check the method
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}

		// Return 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Create a client with the test server URL
	config := Config{
		AccessKey: "test-key",
		SecretKey: "test-secret",
		APIURL:    server.URL,
	}
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Delete a governance insight
	err = client.GovernanceInsights.DeleteGovernanceInsight("test-id-123")
	if err != nil {
		t.Fatalf("Failed to delete governance insight: %v", err)
	}
}