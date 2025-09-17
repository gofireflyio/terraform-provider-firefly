package client

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestRunnersWorkspaceService_CreateRunnersWorkspace(t *testing.T) {
	tests := []struct {
		name            string
		request         CreateRunnersWorkspaceRequest
		expectVariables bool
	}{
		{
			name: "with variables",
			request: CreateRunnersWorkspaceRequest{
				RunnerType:    "github-actions",
				IacType:       "terraform",
				WorkspaceName: "test-workspace",
				Description:   "Test workspace",
				Labels:        []string{"test"},
				VcsID:         "vcs-123",
				Repo:          "test/repo",
				DefaultBranch: "main",
				VcsType:       "github",
				WorkDir:       "/",
				Variables: []Variable{
					{
						Key:         "TEST_VAR",
						Value:       "test-value",
						Sensitivity: SensitivityString,
						Destination: DestinationEnv,
					},
				},
				ConsumedVariableSets: []string{"varset-1"},
				Execution: ExecutionConfig{
					Triggers:         []string{"merge"},
					ApplyRule:        "manual",
					TerraformVersion: "1.5.0",
				},
			},
			expectVariables: true,
		},
		{
			name: "without variables (empty slice)",
			request: CreateRunnersWorkspaceRequest{
				RunnerType:           "github-actions",
				IacType:              "terraform",
				WorkspaceName:        "test-workspace",
				Description:          "Test workspace",
				Labels:               []string{"test"},
				VcsID:                "vcs-123",
				Repo:                 "test/repo",
				DefaultBranch:        "main",
				VcsType:              "github",
				WorkDir:              "/",
				Variables:            []Variable{}, // Empty slice
				ConsumedVariableSets: []string{"varset-1"},
				Execution: ExecutionConfig{
					Triggers:         []string{"merge"},
					ApplyRule:        "manual",
					TerraformVersion: "1.5.0",
				},
			},
			expectVariables: false,
		},
		{
			name: "without variables (nil slice)",
			request: CreateRunnersWorkspaceRequest{
				RunnerType:           "github-actions",
				IacType:              "terraform",
				WorkspaceName:        "test-workspace",
				Description:          "Test workspace",
				Labels:               []string{"test"},
				VcsID:                "vcs-123",
				Repo:                 "test/repo",
				DefaultBranch:        "main",
				VcsType:              "github",
				WorkDir:              "/",
				Variables:            nil, // Nil slice
				ConsumedVariableSets: []string{"varset-1"},
				Execution: ExecutionConfig{
					Triggers:         []string{"merge"},
					ApplyRule:        "manual",
					TerraformVersion: "1.5.0",
				},
			},
			expectVariables: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := NewMockServer()
			defer mockServer.Close()

			// Mock login
			mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
				authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
				json.NewEncoder(w).Encode(authResp)
			})

			// Mock create workspace
			mockServer.AddHandler("/v2/runners/workspaces", func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					return
				}

				var createReq CreateRunnersWorkspaceRequest
				if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				// Verify the request structure
				if createReq.WorkspaceName != tt.request.WorkspaceName {
					http.Error(w, "Invalid workspace name", http.StatusBadRequest)
					return
				}

				// Check if variables are present in the request based on expectation
				hasVariables := len(createReq.Variables) > 0
				if hasVariables != tt.expectVariables {
					http.Error(w, "Variables presence mismatch", http.StatusBadRequest)
					return
				}

				// Create response workspace
				workspace := RunnersWorkspace{
					ID:               "workspace-123",
					Name:             createReq.WorkspaceName,
					Description:      createReq.Description,
					AccountID:        "account-1",
					Repository:       createReq.Repo,
					WorkingDirectory: createReq.WorkDir,
					VcsIntegrationID: createReq.VcsID,
					Vcs:              createReq.VcsType,
					DefaultBranch:    createReq.DefaultBranch,
					Labels:           createReq.Labels,
					ProjectID:        "project-1",
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(workspace)
			})

			client, err := NewClient(Config{
				AccessKey: "test-access",
				SecretKey: "test-secret",
				APIURL:    mockServer.URL(),
			})

			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			workspace, err := client.RunnersWorkspaces.CreateRunnersWorkspace(tt.request)
			if err != nil {
				t.Fatalf("CreateRunnersWorkspace failed: %v", err)
			}

			if workspace.ID != "workspace-123" {
				t.Errorf("Expected workspace ID 'workspace-123', got '%s'", workspace.ID)
			}

			if workspace.Name != tt.request.WorkspaceName {
				t.Errorf("Expected workspace name '%s', got '%s'", tt.request.WorkspaceName, workspace.Name)
			}
		})
	}
}

func TestRunnersWorkspaceService_GetRunnersWorkspace(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock get workspace
	mockServer.AddHandler("/v2/runners/workspaces/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		workspaceID := r.URL.Path[len("/v2/runners/workspaces/"):]
		if workspaceID != "test-workspace-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		workspace := RunnersWorkspace{
			ID:               "test-workspace-id",
			Name:             "Test Workspace",
			Description:      "A test workspace for unit tests",
			AccountID:        "account-1",
			Repository:       "test/repo",
			WorkingDirectory: "/",
			VcsIntegrationID: "vcs-123",
			Vcs:              "github",
			DefaultBranch:    "main",
			Labels:           []string{"test", "unit"},
			ProjectID:        "project-1",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workspace)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	workspace, err := client.RunnersWorkspaces.GetRunnersWorkspace("test-workspace-id")
	if err != nil {
		t.Fatalf("GetRunnersWorkspace failed: %v", err)
	}

	if workspace.ID != "test-workspace-id" {
		t.Errorf("Expected workspace ID 'test-workspace-id', got '%s'", workspace.ID)
	}

	if workspace.Name != "Test Workspace" {
		t.Errorf("Expected workspace name 'Test Workspace', got '%s'", workspace.Name)
	}

	if workspace.Vcs != "github" {
		t.Errorf("Expected VCS 'github', got '%s'", workspace.Vcs)
	}
}

func TestRunnersWorkspaceService_UpdateRunnersWorkspace(t *testing.T) {
	tests := []struct {
		name            string
		request         UpdateRunnersWorkspaceRequest
		expectVariables bool
	}{
		{
			name: "with variables",
			request: UpdateRunnersWorkspaceRequest{
				Name:        "Updated Workspace",
				Description: "Updated description",
				Labels:      []string{"updated", "test"},
				Variables: []Variable{
					{
						Key:         "UPDATED_VAR",
						Value:       "updated-value",
						Sensitivity: SensitivityString,
						Destination: DestinationEnv,
					},
				},
			},
			expectVariables: true,
		},
		{
			name: "without variables (empty slice)",
			request: UpdateRunnersWorkspaceRequest{
				Name:        "Updated Workspace",
				Description: "Updated description",
				Labels:      []string{"updated", "test"},
				Variables:   []Variable{}, // Empty slice
			},
			expectVariables: false,
		},
		{
			name: "without variables (nil slice)",
			request: UpdateRunnersWorkspaceRequest{
				Name:        "Updated Workspace",
				Description: "Updated description",
				Labels:      []string{"updated", "test"},
				Variables:   nil, // Nil slice
			},
			expectVariables: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := NewMockServer()
			defer mockServer.Close()

			// Mock login
			mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
				authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
				json.NewEncoder(w).Encode(authResp)
			})

			// Mock update workspace
			mockServer.AddHandler("/v2/runners/workspaces/", func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					return
				}

				workspaceID := r.URL.Path[len("/v2/runners/workspaces/"):]
				if workspaceID != "test-workspace-id" {
					http.Error(w, "Not found", http.StatusNotFound)
					return
				}

				var updateReq UpdateRunnersWorkspaceRequest
				if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				// Check if variables are present in the request based on expectation
				hasVariables := len(updateReq.Variables) > 0
				if hasVariables != tt.expectVariables {
					http.Error(w, "Variables presence mismatch", http.StatusBadRequest)
					return
				}

				// Return updated workspace
				workspace := RunnersWorkspace{
					ID:               workspaceID,
					Name:             updateReq.Name,
					Description:      updateReq.Description,
					AccountID:        "account-1",
					Repository:       "test/repo",
					WorkingDirectory: "/",
					VcsIntegrationID: "vcs-123",
					Vcs:              "github",
					DefaultBranch:    "main",
					Labels:           updateReq.Labels,
					ProjectID:        "project-1",
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(workspace)
			})

			client, err := NewClient(Config{
				AccessKey: "test-access",
				SecretKey: "test-secret",
				APIURL:    mockServer.URL(),
			})

			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			workspace, err := client.RunnersWorkspaces.UpdateRunnersWorkspace("test-workspace-id", tt.request)
			if err != nil {
				t.Fatalf("UpdateRunnersWorkspace failed: %v", err)
			}

			if workspace.Name != tt.request.Name {
				t.Errorf("Expected workspace name '%s', got '%s'", tt.request.Name, workspace.Name)
			}

			if workspace.Description != tt.request.Description {
				t.Errorf("Expected description '%s', got '%s'", tt.request.Description, workspace.Description)
			}
		})
	}
}

func TestRunnersWorkspaceService_DeleteRunnersWorkspace(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock delete workspace
	mockServer.AddHandler("/v2/runners/workspaces/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		workspaceID := r.URL.Path[len("/v2/runners/workspaces/"):]
		if workspaceID != "test-workspace-id" {
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

	err = client.RunnersWorkspaces.DeleteRunnersWorkspace("test-workspace-id")
	if err != nil {
		t.Fatalf("DeleteRunnersWorkspace failed: %v", err)
	}
}

func TestRunnersWorkspaceService_DestroyWorkspaceResources(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock destroy workspace resources
	mockServer.AddHandler("/v2/runners/workspaces/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check if the path ends with /tasks/destroy
		expectedSuffix := "/tasks/destroy"
		if len(r.URL.Path) < len(expectedSuffix) || r.URL.Path[len(r.URL.Path)-len(expectedSuffix):] != expectedSuffix {
			http.Error(w, "Invalid path", http.StatusNotFound)
			return
		}

		workspaceID := r.URL.Path[:len(r.URL.Path)-len(expectedSuffix)]
		workspaceID = workspaceID[len("/v2/runners/workspaces/"):]
		if workspaceID != "test-workspace-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var destroyReq RunTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&destroyReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if destroyReq.TaskType != "destroy" {
			http.Error(w, "Invalid task type", http.StatusBadRequest)
			return
		}

		taskResp := TaskResponse{
			TaskID: "task-123",
			Status: "initiated",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(taskResp)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	destroyReq := RunTaskRequest{
		TaskType: "destroy",
	}

	taskResp, err := client.RunnersWorkspaces.DestroyWorkspaceResources("test-workspace-id", destroyReq)
	if err != nil {
		t.Fatalf("DestroyWorkspaceResources failed: %v", err)
	}

	if taskResp.TaskID != "task-123" {
		t.Errorf("Expected task ID 'task-123', got '%s'", taskResp.TaskID)
	}

	if taskResp.Status != "initiated" {
		t.Errorf("Expected status 'initiated', got '%s'", taskResp.Status)
	}
}

func TestRunnersWorkspaceService_Error_NotFound(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock 404 response
	mockServer.AddHandler("/v2/runners/workspaces/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Workspace not found", http.StatusNotFound)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.RunnersWorkspaces.GetRunnersWorkspace("nonexistent-workspace")
	if err == nil {
		t.Error("Expected error for nonexistent workspace")
	}
}

func TestRunnersWorkspaceService_Error_InvalidRequest(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock bad request response
	mockServer.AddHandler("/v2/runners/workspaces", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.RunnersWorkspaces.CreateRunnersWorkspace(CreateRunnersWorkspaceRequest{})
	if err == nil {
		t.Error("Expected error for invalid request")
	}
}
