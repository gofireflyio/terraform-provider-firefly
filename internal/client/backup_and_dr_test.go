package client

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func setupBackupAndDRMockServer(t *testing.T) (*MockServer, *Client) {
	t.Helper()
	mockServer := NewMockServer()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{
			AccessToken: "test-token",
			ExpiresAt:   time.Now().Add(time.Hour).Unix(),
			TokenType:   "Bearer",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authResp)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return mockServer, client
}

func TestBackupAndDRCreate(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req BackupPolicyCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		resp := BackupPolicy{
			ID:            "policy-123",
			AccountID:     "acc-456",
			PolicyName:    req.PolicyName,
			Description:   req.Description,
			IntegrationID: req.IntegrationID,
			Region:        req.Region,
			ProviderType:  req.ProviderType,
			Schedule:      req.Schedule,
			Scope:         req.Scope,
			Status:        "Active",
			CreatedAt:     "2025-01-01T00:00:00Z",
			UpdatedAt:     "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	})

	hour := 10
	minute := 30
	req := &BackupPolicyCreateRequest{
		PolicyName:    "test-backup-policy",
		Description:   "Test backup policy",
		IntegrationID: "int-789",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Schedule: &ScheduleConfig{
			Frequency: "Daily",
			Hour:      &hour,
			Minute:    &minute,
		},
		Scope: []ScopeConfig{
			{Type: "tags", Value: []string{"env:prod"}},
		},
	}

	policy, err := client.BackupAndDR.Create(req)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if policy.ID != "policy-123" {
		t.Errorf("Expected ID 'policy-123', got '%s'", policy.ID)
	}
	if policy.PolicyName != "test-backup-policy" {
		t.Errorf("Expected name 'test-backup-policy', got '%s'", policy.PolicyName)
	}
	if policy.Status != "Active" {
		t.Errorf("Expected status 'Active', got '%s'", policy.Status)
	}
	if policy.Schedule == nil || policy.Schedule.Frequency != "Daily" {
		t.Error("Expected schedule with frequency 'Daily'")
	}
}

func TestBackupAndDRCreate_Error(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	req := &BackupPolicyCreateRequest{
		PolicyName:    "test",
		IntegrationID: "int-789",
		Region:        "us-east-1",
		ProviderType:  "aws",
	}

	_, err := client.BackupAndDR.Create(req)
	if err == nil {
		t.Error("Expected error but got none")
	}
}

func TestBackupAndDRGet(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies/policy-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := BackupPolicy{
			ID:            "policy-123",
			AccountID:     "acc-456",
			PolicyName:    "test-backup-policy",
			Description:   "Test backup policy",
			IntegrationID: "int-789",
			Region:        "us-east-1",
			ProviderType:  "aws",
			Status:        "Active",
			SnapshotsCount: 5,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	policy, err := client.BackupAndDR.Get("policy-123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if policy.ID != "policy-123" {
		t.Errorf("Expected ID 'policy-123', got '%s'", policy.ID)
	}
	if policy.PolicyName != "test-backup-policy" {
		t.Errorf("Expected name 'test-backup-policy', got '%s'", policy.PolicyName)
	}
	if policy.SnapshotsCount != 5 {
		t.Errorf("Expected 5 snapshots, got %d", policy.SnapshotsCount)
	}
}

func TestBackupAndDRGet_NotFound(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies/nonexistent", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	_, err := client.BackupAndDR.Get("nonexistent")
	if err == nil {
		t.Error("Expected error but got none")
	}
}

func TestBackupAndDRUpdate(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies/policy-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req BackupPolicyUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		resp := BackupPolicy{
			ID:            "policy-123",
			PolicyName:    req.PolicyName,
			Description:   req.Description,
			IntegrationID: req.IntegrationID,
			Region:        req.Region,
			ProviderType:  req.ProviderType,
			Schedule:      req.Schedule,
			Status:        "Active",
			UpdatedAt:     "2025-01-02T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	hour := 14
	minute := 0
	req := &BackupPolicyUpdateRequest{
		PolicyName:    "updated-policy",
		Description:   "Updated description",
		IntegrationID: "int-789",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Schedule: &ScheduleConfig{
			Frequency: "Weekly",
			Hour:      &hour,
			Minute:    &minute,
			DaysOfWeek: []string{"Monday", "Friday"},
		},
	}

	policy, err := client.BackupAndDR.Update("policy-123", req)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if policy.PolicyName != "updated-policy" {
		t.Errorf("Expected name 'updated-policy', got '%s'", policy.PolicyName)
	}
	if policy.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", policy.Description)
	}
}

func TestBackupAndDRUpdate_Error(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies/policy-123", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	req := &BackupPolicyUpdateRequest{
		PolicyName:    "updated-policy",
		IntegrationID: "int-789",
		Region:        "us-east-1",
		ProviderType:  "aws",
	}

	_, err := client.BackupAndDR.Update("policy-123", req)
	if err == nil {
		t.Error("Expected error but got none")
	}
}

func TestBackupAndDRDelete(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies/policy-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.BackupAndDR.Delete("policy-123", false)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestBackupAndDRDelete_Error(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies/policy-123", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	err := client.BackupAndDR.Delete("policy-123", false)
	if err == nil {
		t.Error("Expected error but got none")
	}
}

func TestBackupAndDRSetStatus(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies/policy-123/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if body["status"] != "Inactive" {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	err := client.BackupAndDR.SetStatus("policy-123", "Inactive")
	if err != nil {
		t.Fatalf("SetStatus failed: %v", err)
	}
}

func TestBackupAndDRSetStatus_Error(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies/policy-123/status", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	err := client.BackupAndDR.SetStatus("policy-123", "Invalid")
	if err == nil {
		t.Error("Expected error but got none")
	}
}

func TestBackupAndDRList(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := BackupPolicyListResponse{
			Policies: []BackupPolicy{
				{
					ID:            "policy-1",
					PolicyName:    "policy-one",
					IntegrationID: "int-1",
					Region:        "us-east-1",
					ProviderType:  "aws",
					Status:        "Active",
				},
				{
					ID:            "policy-2",
					PolicyName:    "policy-two",
					IntegrationID: "int-2",
					Region:        "eu-west-1",
					ProviderType:  "aws",
					Status:        "Inactive",
				},
			},
			Total:    2,
			Page:     1,
			PageSize: 50,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	result, err := client.BackupAndDR.List(1, 50)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(result.Policies) != 2 {
		t.Errorf("Expected 2 policies, got %d", len(result.Policies))
	}
	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}
	if result.Policies[0].PolicyName != "policy-one" {
		t.Errorf("Expected first policy name 'policy-one', got '%s'", result.Policies[0].PolicyName)
	}
}

func TestBackupAndDRList_Error(t *testing.T) {
	mockServer, client := setupBackupAndDRMockServer(t)
	defer mockServer.Close()

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	_, err := client.BackupAndDR.List(1, 50)
	if err == nil {
		t.Error("Expected error but got none")
	}
}
