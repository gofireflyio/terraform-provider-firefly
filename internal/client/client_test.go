package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockServer provides a test HTTP server for client tests
type MockServer struct {
	server *httptest.Server
	mux    *http.ServeMux
}

// NewMockServer creates a new mock server for testing
func NewMockServer() *MockServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	return &MockServer{
		server: server,
		mux:    mux,
	}
}

// Close shuts down the mock server
func (m *MockServer) Close() {
	m.server.Close()
}

// URL returns the mock server URL
func (m *MockServer) URL() string {
	return m.server.URL
}

// AddHandler adds a handler for a specific path
func (m *MockServer) AddHandler(path string, handler http.HandlerFunc) {
	m.mux.HandleFunc(path, handler)
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config",
			config: Config{
				AccessKey: "test-access-key",
				SecretKey: "test-secret-key",
				APIURL:    "https://api.example.com",
			},
			expectError: false,
		},
		{
			name: "missing access key",
			config: Config{
				SecretKey: "test-secret-key",
				APIURL:    "https://api.example.com",
			},
			expectError: true,
		},
		{
			name: "missing secret key",
			config: Config{
				AccessKey: "test-access-key",
				APIURL:    "https://api.example.com",
			},
			expectError: true,
		},
		{
			name: "invalid URL",
			config: Config{
				AccessKey: "test-access-key",
				SecretKey: "test-secret-key",
				APIURL:    "://invalid-url",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if client == nil {
				t.Error("expected client but got nil")
				return
			}
			
			// Verify services are initialized
			if client.Workspaces == nil {
				t.Error("Workspaces service not initialized")
			}
			if client.Guardrails == nil {
				t.Error("Guardrails service not initialized")
			}
			if client.Projects == nil {
				t.Error("Projects service not initialized")
			}
			if client.RunnersWorkspaces == nil {
				t.Error("RunnersWorkspaces service not initialized")
			}
			if client.VariableSets == nil {
				t.Error("VariableSets service not initialized")
			}
		})
	}
}

func TestAuthentication(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock successful login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var loginReq map[string]string
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if loginReq["accessKey"] != "test-access" || loginReq["secretKey"] != "test-secret" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

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

	// Test authentication
	err = client.ensureAuthenticated()
	if err != nil {
		t.Fatalf("Authentication failed: %v", err)
	}

	if client.authToken != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", client.authToken)
	}
}

func TestAuthenticationFailure(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock failed login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})

	client, err := NewClient(Config{
		AccessKey: "invalid-access",
		SecretKey: "invalid-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test authentication failure
	err = client.ensureAuthenticated()
	if err == nil {
		t.Error("Expected authentication to fail")
	}
}

func TestRequestWithAuthentication(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login endpoint
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{
			AccessToken: "test-token",
			ExpiresAt:   time.Now().Add(time.Hour).Unix(),
			TokenType:   "Bearer",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock API endpoint
	mockServer.AddHandler("/v2/test", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		fmt.Fprint(w, `{"success": true}`)
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Make authenticated request
	req, err := client.newRequest(http.MethodGet, "/v2/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.doRequest(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}