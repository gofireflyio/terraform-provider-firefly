package client

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestBackupAndDrService_Create(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock create policy
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

		// Validate required fields
		if createReq.PolicyName == "" {
			http.Error(w, "PolicyName is required", http.StatusBadRequest)
			return
		}

		// Create response
		createResp := PolicyResponse{
			PolicyID:        "policy-123",
			AccountID:       "test-account",
			PolicyName:      createReq.PolicyName,
			IntegrationID:   createReq.IntegrationID,
			Region:          createReq.Region,
			ProviderType:    createReq.ProviderType,
			Schedule:        createReq.Schedule,
			Description:     createReq.Description,
			Scope:           createReq.Scope,
			NotificationID:  createReq.NotificationID,
			VCS:             createReq.VCS,
			BackupOnSave:    createReq.BackupOnSave,
			Status:          "Active",
			SnapshotsCount:  0,
			CreatedAt:       "2025-01-01T00:00:00Z",
			UpdatedAt:       "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
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

	policy := &PolicyCreateRequest{
		PolicyName:    "Test Daily Backup",
		IntegrationID: "int-123",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Schedule: ScheduleConfig{
			Frequency: "Daily",
			Hour:      2,
			Minute:    30,
		},
		Description:  "Test backup policy",
		BackupOnSave: true,
	}

	response, err := client.BackupAndDr.Create(policy)
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

	if response.Schedule.Frequency != "Daily" {
		t.Errorf("Expected frequency 'Daily', got '%s'", response.Schedule.Frequency)
	}
}

func TestBackupAndDrService_CreateWithWeeklySchedule(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock create policy
	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		var createReq PolicyCreateRequest
		json.NewDecoder(r.Body).Decode(&createReq)

		createResp := PolicyResponse{
			PolicyID:       "policy-456",
			AccountID:      "test-account",
			PolicyName:     createReq.PolicyName,
			IntegrationID:  createReq.IntegrationID,
			Region:         createReq.Region,
			ProviderType:   createReq.ProviderType,
			Schedule:       createReq.Schedule,
			BackupOnSave:   createReq.BackupOnSave,
			Status:         "Active",
			SnapshotsCount: 0,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
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

	policy := &PolicyCreateRequest{
		PolicyName:    "Test Weekly Backup",
		IntegrationID: "int-123",
		Region:        "us-west-2",
		ProviderType:  "aws",
		Schedule: ScheduleConfig{
			Frequency:  "Weekly",
			DaysOfWeek: []string{"Sunday", "Wednesday"},
			Hour:       1,
			Minute:     0,
		},
		BackupOnSave: true,
	}

	response, err := client.BackupAndDr.Create(policy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if len(response.Schedule.DaysOfWeek) != 2 {
		t.Errorf("Expected 2 days of week, got %d", len(response.Schedule.DaysOfWeek))
	}

	if response.Schedule.DaysOfWeek[0] != "Sunday" {
		t.Errorf("Expected first day 'Sunday', got '%s'", response.Schedule.DaysOfWeek[0])
	}
}

func TestBackupAndDrService_CreateWithMonthlySchedule(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock create policy
	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		var createReq PolicyCreateRequest
		json.NewDecoder(r.Body).Decode(&createReq)

		createResp := PolicyResponse{
			PolicyID:       "policy-789",
			AccountID:      "test-account",
			PolicyName:     createReq.PolicyName,
			IntegrationID:  createReq.IntegrationID,
			Region:         createReq.Region,
			ProviderType:   createReq.ProviderType,
			Schedule:       createReq.Schedule,
			BackupOnSave:   createReq.BackupOnSave,
			Status:         "Active",
			SnapshotsCount: 0,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
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

	policy := &PolicyCreateRequest{
		PolicyName:    "Test Monthly Backup",
		IntegrationID: "int-123",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Schedule: ScheduleConfig{
			Frequency:           "Monthly",
			MonthlyScheduleType: "specific_weekday",
			WeekdayOrdinal:      "First",
			WeekdayName:         "Sunday",
			Hour:                3,
			Minute:              0,
		},
		BackupOnSave: true,
	}

	response, err := client.BackupAndDr.Create(policy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if response.Schedule.MonthlyScheduleType != "specific_weekday" {
		t.Errorf("Expected monthly_schedule_type 'specific_weekday', got '%s'", response.Schedule.MonthlyScheduleType)
	}

	if response.Schedule.WeekdayOrdinal != "First" {
		t.Errorf("Expected weekday_ordinal 'First', got '%s'", response.Schedule.WeekdayOrdinal)
	}
}

func TestBackupAndDrService_CreateWithScope(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock create policy
	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		var createReq PolicyCreateRequest
		json.NewDecoder(r.Body).Decode(&createReq)

		createResp := PolicyResponse{
			PolicyID:       "policy-scope",
			AccountID:      "test-account",
			PolicyName:     createReq.PolicyName,
			IntegrationID:  createReq.IntegrationID,
			Region:         createReq.Region,
			ProviderType:   createReq.ProviderType,
			Schedule:       createReq.Schedule,
			Scope:          createReq.Scope,
			BackupOnSave:   createReq.BackupOnSave,
			Status:         "Active",
			SnapshotsCount: 0,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
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

	policy := &PolicyCreateRequest{
		PolicyName:    "Test Scoped Backup",
		IntegrationID: "int-123",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Schedule: ScheduleConfig{
			Frequency: "Daily",
			Hour:      2,
			Minute:    0,
		},
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

	response, err := client.BackupAndDr.Create(policy)
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

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock create policy
	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		var createReq PolicyCreateRequest
		json.NewDecoder(r.Body).Decode(&createReq)

		createResp := PolicyResponse{
			PolicyID:       "policy-vcs",
			AccountID:      "test-account",
			PolicyName:     createReq.PolicyName,
			IntegrationID:  createReq.IntegrationID,
			Region:         createReq.Region,
			ProviderType:   createReq.ProviderType,
			Schedule:       createReq.Schedule,
			VCS:            createReq.VCS,
			BackupOnSave:   createReq.BackupOnSave,
			Status:         "Active",
			SnapshotsCount: 0,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
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

	policy := &PolicyCreateRequest{
		PolicyName:    "Test VCS Backup",
		IntegrationID: "int-123",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Schedule: ScheduleConfig{
			Frequency: "Daily",
			Hour:      2,
			Minute:    0,
		},
		VCS: &VCSConfig{
			ProjectID:        "project-456",
			VCSIntegrationID: "github-integration-789",
			RepoID:           "backup-repo-123",
		},
		BackupOnSave: true,
	}

	response, err := client.BackupAndDr.Create(policy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if response.VCS == nil {
		t.Fatal("Expected VCS config, got nil")
	}

	if response.VCS.ProjectID != "project-456" {
		t.Errorf("Expected VCS project ID 'project-456', got '%s'", response.VCS.ProjectID)
	}

	if response.VCS.VCSIntegrationID != "github-integration-789" {
		t.Errorf("Expected VCS integration ID 'github-integration-789', got '%s'", response.VCS.VCSIntegrationID)
	}
}

func TestBackupAndDrService_Get(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock get policy
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
			PolicyID:      "policy-123",
			AccountID:     "test-account",
			PolicyName:    "Test Policy",
			IntegrationID: "int-123",
			Region:        "us-east-1",
			ProviderType:  "aws",
			Schedule: ScheduleConfig{
				Frequency: "Daily",
				Hour:      2,
				Minute:    30,
			},
			Status:         "Active",
			SnapshotsCount: 5,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-01T10:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(policy)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	response, err := client.BackupAndDr.Get("policy-123")
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

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock update policy
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

		var updateReq PolicyCreateRequest
		json.NewDecoder(r.Body).Decode(&updateReq)

		updateResp := PolicyResponse{
			PolicyID:       policyID,
			AccountID:      "test-account",
			PolicyName:     updateReq.PolicyName,
			IntegrationID:  updateReq.IntegrationID,
			Region:         updateReq.Region,
			ProviderType:   updateReq.ProviderType,
			Schedule:       updateReq.Schedule,
			Description:    updateReq.Description,
			BackupOnSave:   updateReq.BackupOnSave,
			Status:         "Active",
			SnapshotsCount: 10,
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-02T00:00:00Z",
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

	policy := &PolicyCreateRequest{
		PolicyName:    "Updated Policy",
		IntegrationID: "int-123",
		Region:        "us-east-1",
		ProviderType:  "aws",
		Schedule: ScheduleConfig{
			Frequency:  "Weekly",
			DaysOfWeek: []string{"Monday", "Friday"},
			Hour:       3,
			Minute:     0,
		},
		Description:  "Updated description",
		BackupOnSave: true,
	}

	updatePolicy := ConvertCreateToUpdate(policy)
	response, err := client.BackupAndDr.Update("policy-123", updatePolicy)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if response.PolicyName != "Updated Policy" {
		t.Errorf("Expected name 'Updated Policy', got '%s'", response.PolicyName)
	}

	if response.Schedule.Frequency != "Weekly" {
		t.Errorf("Expected frequency 'Weekly', got '%s'", response.Schedule.Frequency)
	}

	if response.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", response.Description)
	}
}

func TestBackupAndDrService_Delete(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock delete policy
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

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.BackupAndDr.Delete("policy-123")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestBackupAndDrService_List(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock list policies
	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := PolicyListResponse{
			Data: []PolicyResponse{
				{
					PolicyID:      "policy-1",
					AccountID:     "test-account",
					PolicyName:    "Daily Backup",
					IntegrationID: "int-123",
					Region:        "us-east-1",
					ProviderType:  "aws",
					Schedule: ScheduleConfig{
						Frequency: "Daily",
						Hour:      2,
						Minute:    0,
					},
					Status:         "Active",
					SnapshotsCount: 5,
				},
				{
					PolicyID:      "policy-2",
					AccountID:     "test-account",
					PolicyName:    "Weekly Backup",
					IntegrationID: "int-456",
					Region:        "us-west-2",
					ProviderType:  "aws",
					Schedule: ScheduleConfig{
						Frequency:  "Weekly",
						DaysOfWeek: []string{"Sunday"},
						Hour:       1,
						Minute:     0,
					},
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

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	response, err := client.BackupAndDr.List(nil)
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

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock list policies with query parameters
	mockServer.AddHandler("/v2/backup-and-dr/policies", func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		status := r.URL.Query().Get("status")
		region := r.URL.Query().Get("region")

		if status != "Active" || region != "us-east-1" {
			http.Error(w, "Invalid filters", http.StatusBadRequest)
			return
		}

		response := PolicyListResponse{
			Data: []PolicyResponse{
				{
					PolicyID:      "policy-filtered",
					AccountID:     "test-account",
					PolicyName:    "Filtered Policy",
					IntegrationID: "int-123",
					Region:        "us-east-1",
					ProviderType:  "aws",
					Schedule: ScheduleConfig{
						Frequency: "Daily",
						Hour:      2,
						Minute:    0,
					},
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

	client, err := NewClient(Config{
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

	response, err := client.BackupAndDr.List(filters)
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

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock error responses
	mockServer.AddHandler("/v2/backup-and-dr/policies/not-found", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Policy not found", http.StatusNotFound)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test Get with not found
	_, err = client.BackupAndDr.Get("not-found")
	if err == nil {
		t.Error("Expected error for not found policy, got nil")
	}

	// Test Delete with not found
	err = client.BackupAndDr.Delete("not-found")
	if err == nil {
		t.Error("Expected error for deleting not found policy, got nil")
	}
}
