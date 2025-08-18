package client

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestProjectService_ListProjects(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock list projects
	mockServer.AddHandler("/v2/runners/projects/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		projects := ProjectsListResponse{
			Data: []Project{
				{
					ID:           "project-1",
					Name:         "Test Project 1",
					Description:  "First test project",
					Labels:       []string{"test", "project1"},
					AccountID:    "account-1",
					MembersCount: 5,
				},
				{
					ID:           "project-2", 
					Name:         "Test Project 2",
					Description:  "Second test project",
					Labels:       []string{"test", "project2"},
					AccountID:    "account-1",
					MembersCount: 3,
				},
			},
			TotalCount: 2,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(projects)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	projects, err := client.Projects.ListProjects(10, 0, "")
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}

	if len(projects.Data) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects.Data))
	}

	if projects.Data[0].ID != "project-1" {
		t.Errorf("Expected first project ID 'project-1', got '%s'", projects.Data[0].ID)
	}

	if projects.Data[1].Name != "Test Project 2" {
		t.Errorf("Expected second project name 'Test Project 2', got '%s'", projects.Data[1].Name)
	}
}

func TestProjectService_GetProject(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock get project
	mockServer.AddHandler("/v2/runners/projects/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		projectID := r.URL.Path[len("/v2/runners/projects/"):]
		if projectID != "test-project-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		project := Project{
			ID:              "test-project-id",
			Name:            "Test Project",
			Description:     "A test project for unit tests",
			Labels:          []string{"test", "unit"},
			AccountID:       "account-1",
			MembersCount:    10,
			WorkspaceCount:  5,
			ParentID:        "parent-project-id",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(project)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	project, err := client.Projects.GetProject("test-project-id")
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}

	if project.ID != "test-project-id" {
		t.Errorf("Expected project ID 'test-project-id', got '%s'", project.ID)
	}

	if project.Name != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got '%s'", project.Name)
	}

	if project.WorkspaceCount != 5 {
		t.Errorf("Expected workspace count 5, got %d", project.WorkspaceCount)
	}
}

func TestProjectService_CreateProject(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock create project
	mockServer.AddHandler("/v2/runners/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var createReq CreateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Create response project with generated ID
		project := Project{
			ID:          "generated-project-id",
			Name:        createReq.Name,
			Description: createReq.Description,
			Labels:      createReq.Labels,
			ParentID:    createReq.ParentID,
			AccountID:   "account-1",
			MembersCount: 1,
			WorkspaceCount: 0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(project)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	createReq := CreateProjectRequest{
		Name:        "New Test Project",
		Description: "A newly created test project",
		Labels:      []string{"new", "test"},
		ParentID:    "parent-project",
	}

	project, err := client.Projects.CreateProject(createReq)
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	if project.ID != "generated-project-id" {
		t.Errorf("Expected project ID 'generated-project-id', got '%s'", project.ID)
	}

	if project.Name != createReq.Name {
		t.Errorf("Expected project name '%s', got '%s'", createReq.Name, project.Name)
	}

	if len(project.Labels) != 2 || project.Labels[0] != "new" || project.Labels[1] != "test" {
		t.Errorf("Expected labels [new, test], got %v", project.Labels)
	}
}

func TestProjectService_UpdateProject(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock update project
	mockServer.AddHandler("/v2/runners/projects/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		projectID := r.URL.Path[len("/v2/runners/projects/"):]
		if projectID != "test-project-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var updateReq UpdateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Return updated project
		project := Project{
			ID:          projectID,
			Name:        updateReq.Name,
			Description: updateReq.Description,
			Labels:      updateReq.Labels,
			AccountID:   "account-1",
			MembersCount: 5,
			WorkspaceCount: 3,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(project)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	updateReq := UpdateProjectRequest{
		Name:        "Updated Project Name",
		Description: "Updated project description",
		Labels:      []string{"updated", "test"},
	}

	project, err := client.Projects.UpdateProject("test-project-id", updateReq)
	if err != nil {
		t.Fatalf("UpdateProject failed: %v", err)
	}

	if project.Name != updateReq.Name {
		t.Errorf("Expected project name '%s', got '%s'", updateReq.Name, project.Name)
	}

	if project.Description != updateReq.Description {
		t.Errorf("Expected description '%s', got '%s'", updateReq.Description, project.Description)
	}
}

func TestProjectService_DeleteProject(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock delete project
	mockServer.AddHandler("/v2/runners/projects/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		projectID := r.URL.Path[len("/v2/runners/projects/"):]
		if projectID != "test-project-id" {
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

	err = client.Projects.DeleteProject("test-project-id")
	if err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}
}

func TestProjectService_Error_NotFound(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock 404 response
	mockServer.AddHandler("/v2/runners/projects/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Project not found", http.StatusNotFound)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret", 
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Projects.GetProject("nonexistent-project")
	if err == nil {
		t.Error("Expected error for nonexistent project")
	}
}