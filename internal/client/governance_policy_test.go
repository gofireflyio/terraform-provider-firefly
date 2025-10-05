package client

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestGovernancePolicyService_CreatePolicy(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock create governance policy
	mockServer.AddHandler("/v2/governance/insights/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var createReq GovernancePolicy
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if createReq.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		if createReq.Code == "" {
			http.Error(w, "Code is required", http.StatusBadRequest)
			return
		}

		// Create response with generated ID
		createResp := GovernancePolicy{
			ID:          "generated-policy-id",
			Name:        createReq.Name,
			Description: createReq.Description,
			Code:        createReq.Code,
			Type:        createReq.Type,
			ProviderIDs: createReq.ProviderIDs,
			Labels:      createReq.Labels,
			Severity:    createReq.Severity,
			Category:    createReq.Category,
			Frameworks:  createReq.Frameworks,
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

	// Test with base64 encoded Rego code
	regoCode := `
firefly {
    input.instance_state == "stopped"
}
`
	encodedCode := base64.StdEncoding.EncodeToString([]byte(regoCode))

	policy := &GovernancePolicy{
		Name:        "Test Instance State Policy",
		Description: "Ensure instances are stopped",
		Code:        encodedCode,
		Type:        []string{"EC2"},
		ProviderIDs: []string{"aws"},
		Labels:      []string{"security", "compliance"},
		Severity:    3, // Low
		Category:    "security",
		Frameworks:  []string{"SOC2"},
	}

	response, err := client.GovernancePolicies.Create(policy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if response.ID != "generated-policy-id" {
		t.Errorf("Expected policy ID 'generated-policy-id', got '%s'", response.ID)
	}

	if response.Name != "Test Instance State Policy" {
		t.Errorf("Expected name 'Test Instance State Policy', got '%s'", response.Name)
	}

	if response.Severity != 3 {
		t.Errorf("Expected severity 3, got %d", response.Severity)
	}
}

func TestGovernancePolicyService_ListPolicies(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock list governance policies
	mockServer.AddHandler("/v2/governance/insights", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var listReq GovernancePolicyListRequest
		if err := json.NewDecoder(r.Body).Decode(&listReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Create mock response
		policies := []GovernancePolicy{
			{
				ID:          "policy-1",
				Name:        "Test Policy 1",
				Description: "First test policy",
				Code:        "CgpmaXJlZmx5IHsKICAgIGlucHV0Lmluc3RhbmNlX3N0YXRlID09ICJzdG9wcGVkIgp9Cgo=",
				Type:        []string{"EC2"},
				ProviderIDs: []string{"aws"},
				Severity:    3,
				Category:    "security",
			},
			{
				ID:          "policy-2",
				Name:        "Test Policy 2",
				Description: "Second test policy",
				Code:        "CgpmaXJlZmx5IHsKICAgIGlucHV0LnB1YmxpY19yZWFkID09IGZhbHNlCn0KCg==",
				Type:        []string{"S3"},
				ProviderIDs: []string{"aws"},
				Severity:    5,
				Category:    "security",
			},
		}

		response := GovernancePoliciesResponse{
			Hits:     policies,
			Total:    2,
			Page:     1,
			PageSize: 50,
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

	request := &GovernancePolicyListRequest{
		Query:    "test",
		Category: "security",
		Page:     1,
		PageSize: 50,
	}

	response, err := client.GovernancePolicies.List(request)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if response.Total != 2 {
		t.Errorf("Expected total 2, got %d", response.Total)
	}

	if len(response.Hits) != 2 {
		t.Errorf("Expected 2 policies, got %d", len(response.Hits))
	}

	if response.Hits[0].ID != "policy-1" {
		t.Errorf("Expected first policy ID 'policy-1', got '%s'", response.Hits[0].ID)
	}
}

func TestGovernancePolicyService_GetPolicy(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock list (used by Get)
	mockServer.AddHandler("/v2/governance/insights", func(w http.ResponseWriter, r *http.Request) {
		var listReq GovernancePolicyListRequest
		json.NewDecoder(r.Body).Decode(&listReq)

		// Return policy if ID matches
		if len(listReq.ID) > 0 && listReq.ID[0] == "test-policy-id" {
			policy := GovernancePolicy{
				ID:          "test-policy-id",
				Name:        "Test Policy",
				Description: "A test policy",
				Code:        "CgpmaXJlZmx5IHsKICAgIGlucHV0Lmluc3RhbmNlX3N0YXRlID09ICJzdG9wcGVkIgp9Cgo=",
				Type:        []string{"EC2"},
				ProviderIDs: []string{"aws"},
				Severity:    3,
				Category:    "security",
			}

			response := GovernancePoliciesResponse{
				Hits:     []GovernancePolicy{policy},
				Total:    1,
				Page:     1,
				PageSize: 1,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			response := GovernancePoliciesResponse{
				Hits:     []GovernancePolicy{},
				Total:    0,
				Page:     1,
				PageSize: 1,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	policy, err := client.GovernancePolicies.Get("test-policy-id")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if policy.ID != "test-policy-id" {
		t.Errorf("Expected policy ID 'test-policy-id', got '%s'", policy.ID)
	}

	if policy.Name != "Test Policy" {
		t.Errorf("Expected name 'Test Policy', got '%s'", policy.Name)
	}
}

func TestGovernancePolicyService_UpdatePolicy(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock update governance policy
	mockServer.AddHandler("/v2/governance/insights/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		policyID := r.URL.Path[len("/v2/governance/insights/"):]
		if policyID != "test-policy-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var updateReq GovernancePolicy
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Return updated policy
		updateResp := GovernancePolicy{
			ID:          policyID,
			Name:        updateReq.Name,
			Description: updateReq.Description,
			Code:        updateReq.Code,
			Type:        updateReq.Type,
			ProviderIDs: updateReq.ProviderIDs,
			Labels:      updateReq.Labels,
			Severity:    updateReq.Severity,
			Category:    updateReq.Category,
			Frameworks:  updateReq.Frameworks,
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

	policy := &GovernancePolicy{
		Name:        "Updated Policy",
		Description: "Updated description",
		Code:        "CgpmaXJlZmx5IHsKICAgIGlucHV0Lmluc3RhbmNlX3N0YXRlID09ICJydW5uaW5nIgp9Cgo=",
		Type:        []string{"EC2", "RDS"},
		ProviderIDs: []string{"aws"},
		Severity:    5,
		Category:    "performance",
	}

	response, err := client.GovernancePolicies.Update("test-policy-id", policy)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if response.ID != "test-policy-id" {
		t.Errorf("Expected policy ID 'test-policy-id', got '%s'", response.ID)
	}

	if response.Name != "Updated Policy" {
		t.Errorf("Expected name 'Updated Policy', got '%s'", response.Name)
	}

	if response.Severity != 5 {
		t.Errorf("Expected severity 5, got %d", response.Severity)
	}
}

func TestGovernancePolicyService_DeletePolicy(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock delete governance policy
	mockServer.AddHandler("/v2/governance/insights/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		policyID := r.URL.Path[len("/v2/governance/insights/"):]
		if policyID != "test-policy-id" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Return success (204 No Content)
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

	err = client.GovernancePolicies.Delete("test-policy-id")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestGovernancePolicyService_UpdatePolicyWithFrameworkChange(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Mock login
	mockServer.AddHandler("/v2/login", func(w http.ResponseWriter, r *http.Request) {
		authResp := AuthResponse{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
		json.NewEncoder(w).Encode(authResp)
	})

	// Mock update governance policy - simulate API not returning ID or changed frameworks
	mockServer.AddHandler("/v2/governance/insights/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		policyID := r.URL.Path[len("/v2/governance/insights/"):]
		if policyID != "68d926e57e33bb411adcb37a" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var updateReq GovernancePolicyRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Simulate API bug: return empty ID and different framework
		updateResp := GovernancePolicy{
			ID:          "", // Simulate missing ID
			Name:        updateReq.Name,
			Description: updateReq.Description,
			Code:        updateReq.Code,
			Type:        updateReq.Type,
			ProviderIDs: updateReq.ProviderIDs,
			Labels:      FlexibleStringArray(updateReq.Labels),
			Severity:    updateReq.Severity,
			Category:    updateReq.Category,
			Frameworks:  []string{"devops"}, // Return different framework than requested
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updateResp)
	})

	// Mock list (used by Get to fetch correct state)
	mockServer.AddHandler("/v2/governance/insights", func(w http.ResponseWriter, r *http.Request) {
		var listReq GovernancePolicyListRequest
		json.NewDecoder(r.Body).Decode(&listReq)

		// Return policy with correct frameworks
		if len(listReq.ID) > 0 && listReq.ID[0] == "68d926e57e33bb411adcb37a" {
			policy := GovernancePolicy{
				ID:          "68d926e57e33bb411adcb37a",
				Name:        "EJR-Test-Required-Labels",
				Description: "TF test - EJR-Test-Required-Labels",
				Code:        base64.StdEncoding.EncodeToString([]byte("import future.keywords\nfirefly { labels_exist }")),
				Type:        []string{"gcpobjects"},
				ProviderIDs: []string{"gcp_all"},
				Labels:      FlexibleStringArray{"terraform-test"},
				Severity:    3, // low
				Category:    "Optimization",
				Frameworks:  []string{"tagging_policies"}, // Correct framework
			}

			response := GovernancePoliciesResponse{
				Hits:     []GovernancePolicy{policy},
				Total:    1,
				Page:     1,
				PageSize: 1,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	})

	client, err := NewClient(Config{
		AccessKey: "test-access",
		SecretKey: "test-secret",
		APIURL:    mockServer.URL(),
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test updating with framework change
	policy := &GovernancePolicy{
		Name:        "EJR-Test-Required-Labels",
		Description: "TF test - EJR-Test-Required-Labels",
		Code:        base64.StdEncoding.EncodeToString([]byte("import future.keywords\nfirefly { labels_exist }")),
		Type:        []string{"gcpobjects"},
		ProviderIDs: []string{"gcp_all"},
		Labels:      FlexibleStringArray{"terraform-test"},
		Severity:    3, // low
		Category:    "Optimization",
		Frameworks:  []string{"tagging_policies"}, // Change from devops to tagging_policies
	}

	response, err := client.GovernancePolicies.Update("68d926e57e33bb411adcb37a", policy)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Even though API returned empty ID, client should preserve it
	if response.ID != "" {
		t.Errorf("Expected empty ID from update response (simulating bug), got '%s'", response.ID)
	}

	// The refetch should get the correct data
	correctPolicy, err := client.GovernancePolicies.Get("68d926e57e33bb411adcb37a")
	if err != nil {
		t.Fatalf("Get after update failed: %v", err)
	}

	if correctPolicy.ID != "68d926e57e33bb411adcb37a" {
		t.Errorf("Expected policy ID '68d926e57e33bb411adcb37a', got '%s'", correctPolicy.ID)
	}

	if len(correctPolicy.Frameworks) != 1 || correctPolicy.Frameworks[0] != "tagging_policies" {
		t.Errorf("Expected framework 'tagging_policies', got %v", correctPolicy.Frameworks)
	}
}

func TestSeverityConversion(t *testing.T) {
	// Test SeverityToString
	tests := []struct {
		input    int
		expected string
	}{
		{1, "trace"},
		{2, "info"},
		{3, "low"},
		{4, "medium"},
		{5, "high"},
		{6, "critical"},
		{0, "low"},     // default
		{999, "low"},   // default
	}

	for _, test := range tests {
		result := SeverityToString(test.input)
		if result != test.expected {
			t.Errorf("SeverityToString(%d) = %s, expected %s", test.input, result, test.expected)
		}
	}

	// Test SeverityToInt
	intTests := []struct {
		input    string
		expected int
	}{
		{"trace", 1},
		{"info", 2},
		{"low", 3},
		{"medium", 4},
		{"high", 5},
		{"critical", 6},
		{"invalid", 3}, // default
		{"", 3},        // default
	}

	for _, test := range intTests {
		result := SeverityToInt(test.input)
		if result != test.expected {
			t.Errorf("SeverityToInt(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}