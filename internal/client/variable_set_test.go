package client

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestVariableSetService_ListVariableSets(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock list variable sets
	mockServer.AddHandler("/v2/runners/variables/variable-sets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		variableSets := []VariableSet{
			{
				ID:          "varset-1",
				Name:        "AWS Configuration",
				Description: "AWS related variables",
				Labels:      []string{"aws", "cloud"},
				Parents:     []string{},
				Version:     1,
			},
			{
				ID:          "varset-2",
				Name:        "Database Configuration", 
				Description: "Database connection variables",
				Labels:      []string{"database", "config"},
				Parents:     []string{"varset-1"},
				Version:     2,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(variableSets)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	variableSets, err := client.VariableSets.ListVariableSets(10, 0, "")
	if err != nil {
		t.Fatalf("ListVariableSets failed: %v", err)
	}

	if len(variableSets) != 2 {
		t.Errorf("Expected 2 variable sets, got %d", len(variableSets))
	}

	if variableSets[0].Name != "AWS Configuration" {
		t.Errorf("Expected first variable set name 'AWS Configuration', got '%s'", variableSets[0].Name)
	}

	if len(variableSets[1].Parents) != 1 || variableSets[1].Parents[0] != "varset-1" {
		t.Errorf("Expected second variable set to have parent 'varset-1', got %v", variableSets[1].Parents)
	}
}

func TestVariableSetService_GetVariableSet(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock get variable set
	mockServer.AddHandler("/v2/runners/variables/variable-sets/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		varsetID := r.URL.Path[len("/v2/runners/variables/variable-sets/"):]
		if varsetID != "test-varset-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		variableSet := VariableSet{
			ID:          "test-varset-id",
			Name:        "Test Variable Set",
			Description: "A test variable set for unit tests",
			Labels:      []string{"test", "unit"},
			Parents:     []string{"parent-varset"},
			Version:     3,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(variableSet)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	variableSet, err := client.VariableSets.GetVariableSet("test-varset-id")
	if err != nil {
		t.Fatalf("GetVariableSet failed: %v", err)
	}

	if variableSet.ID != "test-varset-id" {
		t.Errorf("Expected variable set ID 'test-varset-id', got '%s'", variableSet.ID)
	}

	if variableSet.Version != 3 {
		t.Errorf("Expected version 3, got %d", variableSet.Version)
	}

	if len(variableSet.Parents) != 1 || variableSet.Parents[0] != "parent-varset" {
		t.Errorf("Expected parent 'parent-varset', got %v", variableSet.Parents)
	}
}

func TestVariableSetService_CreateVariableSet(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock create variable set
	mockServer.AddHandler("/v2/runners/variables/variable-sets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var createReq CreateVariableSetRequest
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate variables structure
		if len(createReq.Variables) > 0 {
			for _, variable := range createReq.Variables {
				if variable.Key == "" {
					http.Error(w, "Variable key cannot be empty", http.StatusBadRequest)
					return
				}
			}
		}

		// Create response with generated ID
		createResp := CreateVariableSetResponse{
			VariableSetID: "generated-varset-id",
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

	createReq := CreateVariableSetRequest{
		Name:        "New Variable Set",
		Description: "A newly created variable set",
		Labels:      []string{"new", "test"},
		Parents:     []string{"parent-varset"},
		Variables: []Variable{
			{
				Key:         "TEST_VAR",
				Value:       "test-value",
				Sensitivity: SensitivityString,
				Destination: DestinationEnv,
			},
			{
				Key:         "SECRET_VAR",
				Value:       "secret-value",
				Sensitivity: SensitivitySecret,
				Destination: DestinationEnv,
			},
		},
	}

	variableSet, err := client.VariableSets.CreateVariableSet(createReq)
	if err != nil {
		t.Fatalf("CreateVariableSet failed: %v", err)
	}

	if variableSet.VariableSetID != "generated-varset-id" {
		t.Errorf("Expected variable set ID 'generated-varset-id', got '%s'", variableSet.VariableSetID)
	}
}

func TestVariableSetService_UpdateVariableSet(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock update variable set
	mockServer.AddHandler("/v2/runners/variables/variable-sets/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		varsetID := r.URL.Path[len("/v2/runners/variables/variable-sets/"):]
		if varsetID != "test-varset-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var updateReq UpdateVariableSetRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Return updated variable set
		variableSet := VariableSet{
			ID:          varsetID,
			Name:        updateReq.Name,
			Description: updateReq.Description,
			Labels:      updateReq.Labels,
			Parents:     updateReq.Parents,
			Version:     2, // Version incremented
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(variableSet)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	updateReq := UpdateVariableSetRequest{
		Name:        "Updated Variable Set",
		Description: "Updated variable set description",
		Labels:      []string{"updated", "test"},
		Parents:     []string{"new-parent-varset"},
		Variables: []Variable{
			{
				Key:         "UPDATED_VAR",
				Value:       "updated-value",
				Sensitivity: SensitivityString,
				Destination: DestinationEnv,
			},
		},
	}

	variableSet, err := client.VariableSets.UpdateVariableSet("test-varset-id", updateReq)
	if err != nil {
		t.Fatalf("UpdateVariableSet failed: %v", err)
	}

	if variableSet.Name != updateReq.Name {
		t.Errorf("Expected name '%s', got '%s'", updateReq.Name, variableSet.Name)
	}

	if variableSet.Version != 2 {
		t.Errorf("Expected version 2, got %d", variableSet.Version)
	}

	if len(variableSet.Parents) != 1 || variableSet.Parents[0] != "new-parent-varset" {
		t.Errorf("Expected parents ['new-parent-varset'], got %v", variableSet.Parents)
	}
}

func TestVariableSetService_DeleteVariableSet(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock delete variable set
	mockServer.AddHandler("/v2/runners/variables/variable-sets/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		varsetID := r.URL.Path[len("/v2/runners/variables/variable-sets/"):]
		if varsetID != "test-varset-id" {
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

	err = client.VariableSets.DeleteVariableSet("test-varset-id")
	if err != nil {
		t.Fatalf("DeleteVariableSet failed: %v", err)
	}
}