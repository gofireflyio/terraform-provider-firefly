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

func TestProjectService_ListProjectMembers(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock list project members
	mockServer.AddHandler("/v2/runners/projects/test-project/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		members := []Member{
			{
				UserID: "user1",
				Email:  "user1@example.com",
				Role:   "admin",
			},
			{
				UserID: "user2",
				Email:  "user2@example.com",
				Role:   "member",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(members)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	members, err := client.Projects.ListProjectMembers("test-project")
	if err != nil {
		t.Fatalf("ListProjectMembers failed: %v", err)
	}

	if len(members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members))
	}

	if members[0].UserID != "user1" || members[0].Role != "admin" {
		t.Errorf("Expected first member user1/admin, got %s/%s", members[0].UserID, members[0].Role)
	}
}

func TestProjectService_AddProjectMember(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock add project member
	mockServer.AddHandler("/v2/runners/projects/test-project/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var members []Member
		if err := json.NewDecoder(r.Body).Decode(&members); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if len(members) == 0 {
			http.Error(w, "No members provided", http.StatusBadRequest)
			return
		}

		// Return the added member
		addedMembers := []Member{members[0]}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(addedMembers)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	member := Member{
		UserID: "newuser",
		Email:  "newuser@example.com",
		Role:   "member",
	}

	addedMember, err := client.Projects.AddProjectMember("test-project", member)
	if err != nil {
		t.Fatalf("AddProjectMember failed: %v", err)
	}

	if addedMember.UserID != member.UserID {
		t.Errorf("Expected user ID %s, got %s", member.UserID, addedMember.UserID)
	}

	if addedMember.Role != member.Role {
		t.Errorf("Expected role %s, got %s", member.Role, addedMember.Role)
	}
}

func TestProjectService_GetProjectMember(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock list project members (used by GetProjectMember)
	mockServer.AddHandler("/v2/runners/projects/test-project/members", func(w http.ResponseWriter, r *http.Request) {
		members := []Member{
			{
				UserID: "user1",
				Email:  "user1@example.com",
				Role:   "admin",
			},
			{
				UserID: "user2",
				Email:  "user2@example.com",
				Role:   "member",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(members)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test getting existing member
	member, err := client.Projects.GetProjectMember("test-project", "user1")
	if err != nil {
		t.Fatalf("GetProjectMember failed: %v", err)
	}

	if member.UserID != "user1" {
		t.Errorf("Expected user ID user1, got %s", member.UserID)
	}

	if member.Role != "admin" {
		t.Errorf("Expected role admin, got %s", member.Role)
	}

	// Test getting non-existent member
	_, err = client.Projects.GetProjectMember("test-project", "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent member")
	}
}

func TestProjectService_RemoveProjectMember(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock remove project member
	mockServer.AddHandler("/v2/runners/projects/test-project/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var userIDs []string
		if err := json.NewDecoder(r.Body).Decode(&userIDs); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if len(userIDs) == 0 || userIDs[0] != "user1" {
			http.Error(w, "Invalid user", http.StatusBadRequest)
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

	err = client.Projects.RemoveProjectMember("test-project", "user1")
	if err != nil {
		t.Fatalf("RemoveProjectMember failed: %v", err)
	}
}