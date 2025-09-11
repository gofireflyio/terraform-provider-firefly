package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GovernancePolicyService provides access to the governance policy API methods
type GovernancePolicyService struct {
	client *Client
}

// GovernancePolicy represents a governance policy
type GovernancePolicy struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Code        string   `json:"code"`
	Type        []string `json:"type"`
	ProviderIDs []string `json:"providerIds"`
	Labels      []string `json:"labels,omitempty"`
	Severity    int      `json:"severity,omitempty"`
	Category    string   `json:"category,omitempty"`
	Frameworks  []string `json:"frameworks,omitempty"`
}

// GovernancePoliciesResponse represents the response from the policies list endpoint
type GovernancePoliciesResponse struct {
	Data     []GovernancePolicy `json:"data"`
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

// GovernancePolicyListRequest represents the request body for listing policies
type GovernancePolicyListRequest struct {
	Query                 string   `json:"query,omitempty"`
	Labels                []string `json:"labels,omitempty"`
	Frameworks            []string `json:"frameworks,omitempty"`
	Category              string   `json:"category,omitempty"`
	IsDefault             *bool    `json:"isDefault,omitempty"`
	OnlySubscribed        bool     `json:"onlySubscribed,omitempty"`
	OnlyProduction        bool     `json:"onlyProduction,omitempty"`
	OnlyMatchingAssets    bool     `json:"onlyMatchingAssets,omitempty"`
	OnlyEnabled           bool     `json:"onlyEnabled,omitempty"`
	OnlyAvailableProviders bool     `json:"onlyAvailableProviders,omitempty"`
	ShowExclusion         bool     `json:"showExclusion,omitempty"`
	Type                  []string `json:"type,omitempty"`
	Providers             []string `json:"providers,omitempty"`
	Integrations          []string `json:"integrations,omitempty"`
	Severity              []int    `json:"severity,omitempty"`
	ID                    []string `json:"id,omitempty"`
	Page                  int      `json:"page,omitempty"`
	PageSize              int      `json:"page_size,omitempty"`
	Sorting               []string `json:"sorting,omitempty"`
	ProvidersAccounts     []string `json:"providersAcoounts,omitempty"`
}

// List retrieves governance policies
func (s *GovernancePolicyService) List(request *GovernancePolicyListRequest) (*GovernancePoliciesResponse, error) {
	endpoint := "/governance/insights"
	
	// Set default values if not provided
	if request.Page == 0 {
		request.Page = 1
	}
	if request.PageSize == 0 {
		request.PageSize = 50
	}
	
	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}
	
	req, err := s.client.newRequest("POST", endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	
	resp, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var result GovernancePoliciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	return &result, nil
}

// Get retrieves a specific governance policy by ID
func (s *GovernancePolicyService) Get(id string) (*GovernancePolicy, error) {
	// Use the List endpoint with specific ID filter to get a single policy
	request := &GovernancePolicyListRequest{
		ID:       []string{id},
		PageSize: 1,
	}
	
	response, err := s.List(request)
	if err != nil {
		return nil, err
	}
	
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("policy not found: %s", id)
	}
	
	return &response.Data[0], nil
}

// Create creates a new governance policy
func (s *GovernancePolicyService) Create(policy *GovernancePolicy) (*GovernancePolicy, error) {
	endpoint := "/governance/insights/create"
	
	data, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("error marshaling policy: %w", err)
	}
	
	req, err := s.client.newRequest("POST", endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	
	resp, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var result GovernancePolicy
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	return &result, nil
}

// Update updates an existing governance policy
func (s *GovernancePolicyService) Update(id string, policy *GovernancePolicy) (*GovernancePolicy, error) {
	endpoint := fmt.Sprintf("/governance/insights/%s", url.PathEscape(id))
	
	data, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("error marshaling policy: %w", err)
	}
	
	req, err := s.client.newRequest("PUT", endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	
	resp, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var result GovernancePolicy
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	return &result, nil
}

// Delete deletes a governance policy
func (s *GovernancePolicyService) Delete(id string) error {
	endpoint := fmt.Sprintf("/governance/insights/%s", url.PathEscape(id))
	
	req, err := s.client.newRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	
	resp, err := s.client.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}