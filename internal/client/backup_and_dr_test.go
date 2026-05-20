package client

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

const testAccountID = "test-account"

func TestBackupAndDrService_Create(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var createReq PolicyCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if createReq.PolicyName == "" {
			http.Error(w, "PolicyName is required", http.StatusBadRequest)
			return
		}

		createResp := PolicyResponse{
			PolicyID:        "policy-123",
			AccountID:       testAccountID,
			PolicyName:      createReq.PolicyName,
			IntegrationID:   createReq.IntegrationID,
			Region:          createReq.Region,
			ProviderType:    createReq.ProviderType,
			Frequency:       createReq.Frequency,
			Description:     createReq.Description,
			Scope:           createReq.Scope,
			NotificationID:  createReq.NotificationID,
			VCS:             createReq.VCS,
			Status:          "Active",
			SnapshotsCount:  0,
			CreatedAt:       "2025-01-01T00:00:00Z",
			UpdatedAt:       "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createResp)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	policy := &PolicyCreateRequest{
		PolicyName:    "Test Daily Backup",
		IntegrationID: "int-123",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Frequency:     24,
		Description:   "Test backup policy",
		BackupOnSave:  true,
	}

	response, err := c.BackupAndDr.Create(policy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if response.PolicyID != "policy-123" {
		t.Errorf("Expected policy ID 'policy-123', got '%s'", response.PolicyID)
	}

	if response.PolicyName != "Test Daily Backup" {
		t.Errorf("Expected name 'Test Daily Backup', got '%s'", response.PolicyName)
	}

	if response.Status != "Active" {
		t.Errorf("Expected status 'Active', got '%s'", response.Status)
	}

	if response.Frequency != 24 {
		t.Errorf("Expected frequency 24, got %d", response.Frequency)
	}
}

func TestBackupAndDrService_CreateWithScope(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		var createReq PolicyCreateRequest
		json.NewDecoder(r.Body).Decode(&createReq)

		createResp := PolicyResponse{
			PolicyID:       "policy-scope",
			AccountID:      testAccountID,
			PolicyName:     createReq.PolicyName,
			IntegrationID:  createReq.IntegrationID,
			Region:         createReq.Region,
			ProviderType:   createReq.ProviderType,
			Frequency:      createReq.Frequency,
			Scope:          createReq.Scope,
			Status:         "Active",
			SnapshotsCount: 0,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createResp)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	policy := &PolicyCreateRequest{
		PolicyName:    "Test Scoped Backup",
		IntegrationID: "int-123",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Frequency:     8,
		Scope: []ScopeConfig{
			{
				Type:  "tags",
				Value: []string{"Environment:Production", "Backup:Required"},
			},
			{
				Type:  "asset_types",
				Value: []string{"aws_instance", "aws_db_instance"},
			},
		},
		BackupOnSave: true,
	}

	response, err := c.BackupAndDr.Create(policy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if len(response.Scope) != 2 {
		t.Errorf("Expected 2 scopes, got %d", len(response.Scope))
	}

	if response.Scope[0].Type != "tags" {
		t.Errorf("Expected first scope type 'tags', got '%s'", response.Scope[0].Type)
	}

	if len(response.Scope[0].Value) != 2 {
		t.Errorf("Expected 2 tag values, got %d", len(response.Scope[0].Value))
	}
}

func TestBackupAndDrService_CreateWithVCS(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		var createReq PolicyCreateRequest
		json.NewDecoder(r.Body).Decode(&createReq)

		createResp := PolicyResponse{
			PolicyID:       "policy-vcs",
			AccountID:      testAccountID,
			PolicyName:     createReq.PolicyName,
			IntegrationID:  createReq.IntegrationID,
			Region:         createReq.Region,
			ProviderType:   createReq.ProviderType,
			Frequency:      createReq.Frequency,
			VCS:            createReq.VCS,
			Status:         "Active",
			SnapshotsCount: 0,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createResp)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	policy := &PolicyCreateRequest{
		PolicyName:    "Test VCS Backup",
		IntegrationID: "int-123",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Frequency:     24,
		VCS: &VCSConfig{
			VCSIntegrationID: "github-integration-789",
			RepoID:           "backup-repo-123",
		},
		BackupOnSave: true,
	}

	response, err := c.BackupAndDr.Create(policy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if response.VCS == nil {
		t.Fatal("Expected VCS config, got nil")
	}

	if response.VCS.VCSIntegrationID != "github-integration-789" {
		t.Errorf("Expected VCS integration ID 'github-integration-789', got '%s'", response.VCS.VCSIntegrationID)
	}

	if response.VCS.RepoID != "backup-repo-123" {
		t.Errorf("Expected VCS repo ID 'backup-repo-123', got '%s'", response.VCS.RepoID)
	}
}

func TestBackupAndDrService_CreateWithResilienceFields(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		var createReq PolicyCreateRequest
		json.NewDecoder(r.Body).Decode(&createReq)

		createResp := PolicyResponse{
			PolicyID:          "policy-resilience",
			AccountID:         testAccountID,
			PolicyName:        createReq.PolicyName,
			IntegrationID:     createReq.IntegrationID,
			Region:            createReq.Region,
			ProviderType:      createReq.ProviderType,
			Frequency:         createReq.Frequency,
			TargetAccount:     createReq.TargetAccount,
			TargetRegion:      createReq.TargetRegion,
			AutoCreatePR:      createReq.AutoCreatePR,
			ResilienceEnabled: createReq.ResilienceEnabled,
			Status:            "Active",
			SnapshotsCount:    0,
			CreatedAt:         "2025-01-01T00:00:00Z",
			UpdatedAt:         "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createResp)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	policy := &PolicyCreateRequest{
		PolicyName:        "Test Resilience Backup",
		IntegrationID:     "int-123",
		Region:            "us-east-1",
		ProviderType:      "aws",
		Frequency:         4,
		TargetAccount:     "target-int-456",
		TargetRegion:      "eu-west-1",
		AutoCreatePR:      true,
		ResilienceEnabled: true,
	}

	response, err := c.BackupAndDr.Create(policy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if response.TargetAccount != "target-int-456" {
		t.Errorf("Expected TargetAccount 'target-int-456', got '%s'", response.TargetAccount)
	}

	if response.TargetRegion != "eu-west-1" {
		t.Errorf("Expected TargetRegion 'eu-west-1', got '%s'", response.TargetRegion)
	}

	if !response.AutoCreatePR {
		t.Error("Expected AutoCreatePR true, got false")
	}

	if !response.ResilienceEnabled {
		t.Error("Expected ResilienceEnabled true, got false")
	}
}

func TestBackupAndDrService_Get(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		policyID := r.URL.Path[len("/v2/backup-and-dr/policies/"):]
		if policyID != "policy-123" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		policy := PolicyResponse{
			PolicyID:       "policy-123",
			AccountID:      testAccountID,
			PolicyName:     "Test Policy",
			IntegrationID:  "int-123",
			Region:         "us-east-1",
			ProviderType:   "aws",
			Frequency:      24,
			Status:         "Active",
			SnapshotsCount: 5,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-01T10:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(policy)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	response, err := c.BackupAndDr.Get("policy-123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if response.PolicyID != "policy-123" {
		t.Errorf("Expected policy ID 'policy-123', got '%s'", response.PolicyID)
	}

	if response.PolicyName != "Test Policy" {
		t.Errorf("Expected name 'Test Policy', got '%s'", response.PolicyName)
	}

	if response.SnapshotsCount != 5 {
		t.Errorf("Expected snapshots count 5, got %d", response.SnapshotsCount)
	}
}

func TestBackupAndDrService_Update(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		policyID := r.URL.Path[len("/v2/backup-and-dr/policies/"):]
		if policyID != "policy-123" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var updateReq PolicyUpdateRequest
		json.NewDecoder(r.Body).Decode(&updateReq)

		freq := 0
		if updateReq.Frequency != nil {
			freq = *updateReq.Frequency
		}

		updateResp := PolicyResponse{
			PolicyID:       policyID,
			AccountID:      testAccountID,
			PolicyName:     *updateReq.PolicyName,
			IntegrationID:  *updateReq.IntegrationID,
			Region:         *updateReq.Region,
			ProviderType:   *updateReq.ProviderType,
			Frequency:      freq,
			Status:         "Active",
			SnapshotsCount: 10,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-02T00:00:00Z",
		}

		if updateReq.Description != nil {
			updateResp.Description = *updateReq.Description
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updateResp)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	policy := &PolicyCreateRequest{
		PolicyName:    "Updated Policy",
		IntegrationID: "int-123",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Frequency:     8,
		Description:   "Updated description",
		BackupOnSave:  true,
	}

	updatePolicy := ConvertCreateToUpdate(policy)
	response, err := c.BackupAndDr.Update("policy-123", updatePolicy)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if response.PolicyName != "Updated Policy" {
		t.Errorf("Expected name 'Updated Policy', got '%s'", response.PolicyName)
	}

	if response.Frequency != 8 {
		t.Errorf("Expected frequency 8, got %d", response.Frequency)
	}

	if response.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", response.Description)
	}
}

func TestBackupAndDrService_Delete(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		policyID := r.URL.Path[len("/v2/backup-and-dr/policies/"):]
		if policyID != "policy-123" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.BackupAndDr.Delete("policy-123")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestBackupAndDrService_List(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := PolicyListResponse{
			Data: []PolicyResponse{
				{
					PolicyID:       "policy-1",
					AccountID:      testAccountID,
					PolicyName:     "Daily Backup",
					IntegrationID:  "int-123",
					Region:         "us-east-1",
					ProviderType:   "aws",
					Frequency:      24,
					Status:         "Active",
					SnapshotsCount: 5,
				},
				{
					PolicyID:       "policy-2",
					AccountID:      testAccountID,
					PolicyName:     "Frequent Backup",
					IntegrationID:  "int-456",
					Region:         "us-west-2",
					ProviderType:   "aws",
					Frequency:      4,
					Status:         "Inactive",
					SnapshotsCount: 0,
				},
			},
			Pagination: Pagination{
				Page:     1,
				PageSize: 25,
				Total:    2,
				HasNext:  false,
				HasPrev:  false,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	response, err := c.BackupAndDr.List(nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 policies, got %d", len(response.Data))
	}

	if response.Data[0].PolicyID != "policy-1" {
		t.Errorf("Expected first policy ID 'policy-1', got '%s'", response.Data[0].PolicyID)
	}

	if response.Data[1].Status != "Inactive" {
		t.Errorf("Expected second policy status 'Inactive', got '%s'", response.Data[1].Status)
	}
}

func TestBackupAndDrService_ListWithFilters(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		region := r.URL.Query().Get("region")

		if status != "Active" || region != "us-east-1" {
			http.Error(w, "Invalid filters", http.StatusBadRequest)
			return
		}

		response := PolicyListResponse{
			Data: []PolicyResponse{
				{
					PolicyID:       "policy-filtered",
					AccountID:      testAccountID,
					PolicyName:     "Filtered Policy",
					IntegrationID:  "int-123",
					Region:         "us-east-1",
					ProviderType:   "aws",
					Frequency:      24,
					Status:         "Active",
					SnapshotsCount: 3,
				},
			},
			Pagination: Pagination{
				Page:     1,
				PageSize: 25,
				Total:    1,
				HasNext:  false,
				HasPrev:  false,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	filters := &PolicyListFilters{
		Status: "Active",
		Region: "us-east-1",
	}

	response, err := c.BackupAndDr.List(filters)
	if err != nil {
		t.Fatalf("List with filters failed: %v", err)
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 policy, got %d", len(response.Data))
	}

	if response.Data[0].Status != "Active" {
		t.Errorf("Expected status 'Active', got '%s'", response.Data[0].Status)
	}

	if response.Data[0].Region != "us-east-1" {
		t.Errorf("Expected region 'us-east-1', got '%s'", response.Data[0].Region)
	}
}

func TestBackupAndDrService_ErrorHandling(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	mockServer.AddHandler("/v2/backup-and-dr/policies/not-found", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Policy not found", http.StatusNotFound)
	})

	c, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = c.BackupAndDr.Get("not-found")
	if err == nil {
		t.Error("Expected error for not found policy, got nil")
	}

	err = c.BackupAndDr.Delete("not-found")
	if err == nil {
		t.Error("Expected error for deleting not found policy, got nil")
	}
}
