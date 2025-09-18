package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GovernanceInsightService provides access to the governance insight-related API methods
type GovernanceInsightService struct {
	client *Client
}

// GovernanceInsight represents a Firefly governance insight (policy)
type GovernanceInsight struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Code        string   `json:"code"`
	Type        []string `json:"type"`
	ProviderIDs []string `json:"providerIds"`
	Labels      []string `json:"labels,omitempty"`
	Severity    int      `json:"severity,omitempty"`
	Category    string   `json:"category"`
	Frameworks  []string `json:"frameworks,omitempty"`
	CreatedAt   string   `json:"createdAt,omitempty"`
	UpdatedAt   string   `json:"updatedAt,omitempty"`
	CreatedBy   string   `json:"createdBy,omitempty"`
	AccountID   string   `json:"accountId,omitempty"`
}

// GovernanceInsightFilters represents filters for listing insights
type GovernanceInsightFilters struct {
	Labels     []string `json:"labels,omitempty"`
	Frameworks []string `json:"frameworks,omitempty"`
	Severity   []int    `json:"severity,omitempty"`
}

// ListGovernanceInsightsRequest represents the request body for listing insights
type ListGovernanceInsightsRequest struct {
	Filters     *GovernanceInsightFilters `json:"filters,omitempty"`
	SearchValue string                     `json:"searchValue,omitempty"`
	Fields      []string                   `json:"fields,omitempty"`
}

// ListGovernanceInsights retrieves governance insights from Firefly
func (s *GovernanceInsightService) ListGovernanceInsights(request *ListGovernanceInsightsRequest, page, pageSize int) ([]GovernanceInsight, error) {
	// Construct query parameters
	queryParams := url.Values{}
	queryParams.Add("page", strconv.Itoa(page))
	queryParams.Add("page_size", strconv.Itoa(pageSize))

	// Add fields if specified
	if request != nil && len(request.Fields) > 0 {
		queryParams.Add("fields", strings.Join(request.Fields, ","))
	}

	// Add filters if specified
	if request != nil && request.Filters != nil {
		if len(request.Filters.Labels) > 0 {
			queryParams.Add("labels", strings.Join(request.Filters.Labels, ","))
		}
		if len(request.Filters.Frameworks) > 0 {
			for _, framework := range request.Filters.Frameworks {
				queryParams.Add("frameworks", framework)
			}
		}
		if len(request.Filters.Severity) > 0 {
			for _, sev := range request.Filters.Severity {
				queryParams.Add("severity", strconv.Itoa(sev))
			}
		}
	}

	// Create the request
	req, err := s.client.newRequest(http.MethodPost, fmt.Sprintf("/v2/governance/insights?%s", queryParams.Encode()), nil)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list governance insights: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var insights []GovernanceInsight
	if err := json.NewDecoder(resp.Body).Decode(&insights); err != nil {
		return nil, fmt.Errorf("error parsing governance insight list response: %s", err)
	}

	return insights, nil
}

// CreateGovernanceInsight creates a new governance insight
func (s *GovernanceInsightService) CreateGovernanceInsight(insight *GovernanceInsight) (*GovernanceInsight, error) {
	// Create the request - note the correct endpoint path
	req, err := s.client.newRequest(http.MethodPost, "/v2/governance/insights/create", insight)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create governance insight: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var createdInsight GovernanceInsight
	if err := json.NewDecoder(resp.Body).Decode(&createdInsight); err != nil {
		return nil, fmt.Errorf("error parsing create governance insight response: %s", err)
	}

	return &createdInsight, nil
}

// GetGovernanceInsight retrieves a governance insight by ID
func (s *GovernanceInsightService) GetGovernanceInsight(insightID string) (*GovernanceInsight, error) {
	// We'll use the list API with a filter since there's no direct get endpoint
	request := &ListGovernanceInsightsRequest{
		Fields: []string{"id", "name", "description", "code", "type", "providerIds", "labels", "severity", "category", "frameworks"},
	}

	// Make the request
	insights, err := s.ListGovernanceInsights(request, 0, 100)
	if err != nil {
		return nil, err
	}

	// Find the insight with the matching ID
	for _, insight := range insights {
		if insight.ID == insightID {
			return &insight, nil
		}
	}

	return nil, fmt.Errorf("governance insight with ID %s not found", insightID)
}

// UpdateGovernanceInsight updates an existing governance insight
func (s *GovernanceInsightService) UpdateGovernanceInsight(insightID string, insight *GovernanceInsight) (*GovernanceInsight, error) {
	// Create the request
	req, err := s.client.newRequest(http.MethodPut, fmt.Sprintf("/v2/governance/insights/%s", insightID), insight)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update governance insight: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var updatedInsight GovernanceInsight
	if err := json.NewDecoder(resp.Body).Decode(&updatedInsight); err != nil {
		return nil, fmt.Errorf("error parsing update governance insight response: %s", err)
	}

	return &updatedInsight, nil
}

// DeleteGovernanceInsight deletes a governance insight by ID
func (s *GovernanceInsightService) DeleteGovernanceInsight(insightID string) error {
	// Create the request
	req, err := s.client.newRequest(http.MethodDelete, fmt.Sprintf("/v2/governance/insights/%s", insightID), nil)
	if err != nil {
		return err
	}

	// Execute the request
	resp, err := s.client.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete governance insight: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	return nil
}