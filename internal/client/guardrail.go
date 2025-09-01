package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// SeverityToString converts integer severity to string representation
// Flexible = 1, Strict = 2, Warning = 3
func SeverityToString(severity int) string {
	switch severity {
	case 1:
		return "Flexible"
	case 2:
		return "Strict"
	case 3:
		return "Warning"
	default:
		return "Unknown"
	}
}

// SeverityToInt converts string severity to integer representation
// Flexible = 1, Strict = 2, Warning = 3
func SeverityToInt(severity string) int {
	switch severity {
	case "Flexible":
		return 1
	case "Strict":
		return 2
	case "Warning":
		return 3
	default:
		return 0 // Unknown/invalid severity
	}
}

// GuardrailService provides access to the guardrail-related API methods
type GuardrailService struct {
	client *Client
}

// IncludeExcludeWildcard represents a pattern for including and excluding items
type IncludeExcludeWildcard struct {
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
}

// GuardrailScope defines the scope of a guardrail rule
type GuardrailScope struct {
	Workspaces   *IncludeExcludeWildcard `json:"workspaces,omitempty"`
	Repositories *IncludeExcludeWildcard `json:"repositories,omitempty"`
	Branches     *IncludeExcludeWildcard `json:"branches,omitempty"`
	Labels       *IncludeExcludeWildcard `json:"labels,omitempty"`
}

// CostCriteria defines criteria for cost-based guardrails
type CostCriteria struct {
	ThresholdAmount     *float64 `json:"thresholdAmount,omitempty"`
	ThresholdPercentage *float64 `json:"thresholdPercentage,omitempty"`
}

// PolicyCriteria defines criteria for policy-based guardrails
type PolicyCriteria struct {
	Severity string                  `json:"severity,omitempty"`
	Policies *IncludeExcludeWildcard `json:"policies,omitempty"`
}

// ResourceCriteria defines criteria for resource-based guardrails
type ResourceCriteria struct {
	Actions           []string                `json:"actions,omitempty"`
	Regions           *IncludeExcludeWildcard `json:"regions,omitempty"`
	AssetTypes        *IncludeExcludeWildcard `json:"assetTypes,omitempty"`
	SpecificResources []string                `json:"specificResources,omitempty"`
}

// TagCriteria defines criteria for tag-based guardrails
type TagCriteria struct {
	TagEnforcementMode string              `json:"tagEnforcementMode,omitempty"`
	RequiredTags       []string            `json:"requiredTags,omitempty"`
	RequiredValues     map[string][]string `json:"requiredValues,omitempty"`
}

// GuardrailCriteria defines the criteria for a guardrail rule
type GuardrailCriteria struct {
	Cost     *CostCriteria     `json:"cost,omitempty"`
	Policy   *PolicyCriteria   `json:"policy,omitempty"`
	Resource *ResourceCriteria `json:"resource,omitempty"`
	Tag      *TagCriteria      `json:"tag,omitempty"`
}

// GuardrailRule represents a Firefly guardrail rule
type GuardrailRule struct {
	ID             string             `json:"id,omitempty"`
	AccountID      string             `json:"accountId,omitempty"`
	CreatedBy      string             `json:"createdBy,omitempty"`
	Name           string             `json:"name"`
	Type           string             `json:"type"`
	Scope          *GuardrailScope    `json:"scope,omitempty"`
	Criteria       *GuardrailCriteria `json:"criteria,omitempty"`
	IsEnabled      bool               `json:"isEnabled"`
	CreatedAt      string             `json:"createdAt,omitempty"`
	UpdatedAt      string             `json:"updatedAt,omitempty"`
	NotificationID string             `json:"notificationId,omitempty"`
	Severity       int                `json:"severity"`
}

// GuardrailFilters represents filters for listing guardrails
type GuardrailFilters struct {
	CreatedBy    []string `json:"createdBy,omitempty"`
	Type         []string `json:"type,omitempty"`
	Labels       []string `json:"labels,omitempty"`
	Repositories []string `json:"repositories,omitempty"`
	Workspaces   []string `json:"workspaces,omitempty"`
	Branches     []string `json:"branches,omitempty"`
}

// ListGuardrailsRequest represents the request body for listing guardrails
type ListGuardrailsRequest struct {
	Filters     *GuardrailFilters `json:"filters,omitempty"`
	SearchValue string            `json:"searchValue,omitempty"`
	Projection  []string          `json:"projection,omitempty"`
}

// CreateGuardrailResponse represents the response from creating a guardrail
type CreateGuardrailResponse struct {
	RuleID         string `json:"ruleId"`
	NotificationID string `json:"notificationId"`
}

// UpdateGuardrailResponse represents the response from updating a guardrail
type UpdateGuardrailResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	UpdatedAt string `json:"updatedAt"`
}

// DeleteGuardrailResponse represents the response from deleting a guardrail
type DeleteGuardrailResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ListGuardrails retrieves guardrail rules from Firefly
func (s *GuardrailService) ListGuardrails(request *ListGuardrailsRequest, page, pageSize int) ([]GuardrailRule, error) {
	// Construct query parameters
	queryParams := url.Values{}
	queryParams.Add("page", strconv.Itoa(page))
	queryParams.Add("pageSize", strconv.Itoa(pageSize))

	// Create the request
	req, err := s.client.newRequest(http.MethodPost, fmt.Sprintf("/v2/guardrails/search?%s", queryParams.Encode()), request)
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
		return nil, fmt.Errorf("failed to list guardrails: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var guardrails []GuardrailRule
	if err := json.NewDecoder(resp.Body).Decode(&guardrails); err != nil {
		return nil, fmt.Errorf("error parsing guardrail list response: %s", err)
	}

	return guardrails, nil
}

// CreateGuardrail creates a new guardrail rule
func (s *GuardrailService) CreateGuardrail(guardrail *GuardrailRule) (*CreateGuardrailResponse, error) {
	// Create the request
	req, err := s.client.newRequest(http.MethodPost, "/v2/guardrails", guardrail)
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
		return nil, fmt.Errorf("failed to create guardrail: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response - API returns just the rule ID as a string
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading create guardrail response: %s", err)
	}

	// Try to unmarshal as JSON object first
	var createResp CreateGuardrailResponse
	if err := json.Unmarshal(bodyBytes, &createResp); err != nil {
		// If that fails, try to unmarshal as a simple string (which the API actually returns)
		var ruleID string
		if err := json.Unmarshal(bodyBytes, &ruleID); err != nil {
			return nil, fmt.Errorf("error parsing create guardrail response: %s", err)
		}
		// Create response object from string
		createResp = CreateGuardrailResponse{
			RuleID:         ruleID,
			NotificationID: "", // Not provided in string response
		}
	}

	return &createResp, nil
}

// GetGuardrail retrieves a guardrail rule by ID
func (s *GuardrailService) GetGuardrail(ruleID string) (*GuardrailRule, error) {
	// We'll use the list API with a filter since there's no direct get endpoint in the OpenAPI spec
	request := &ListGuardrailsRequest{
		Filters: &GuardrailFilters{},
		// No direct ID filter available, we'll filter client-side
	}

	// Make the request
	guardrails, err := s.ListGuardrails(request, 0, 100)
	if err != nil {
		return nil, err
	}

	// Find the guardrail with the matching ID
	for _, rule := range guardrails {
		if rule.ID == ruleID {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("guardrail rule with ID %s not found", ruleID)
}

// UpdateGuardrail updates an existing guardrail rule
func (s *GuardrailService) UpdateGuardrail(ruleID string, guardrail *GuardrailRule) (*UpdateGuardrailResponse, error) {
	// Create the request
	req, err := s.client.newRequest(http.MethodPatch, fmt.Sprintf("/v2/guardrails/%s", ruleID), guardrail)
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
		return nil, fmt.Errorf("failed to update guardrail: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var updateResp UpdateGuardrailResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return nil, fmt.Errorf("error parsing update guardrail response: %s", err)
	}

	return &updateResp, nil
}

// DeleteGuardrail deletes a guardrail rule by ID
func (s *GuardrailService) DeleteGuardrail(ruleID string) (*DeleteGuardrailResponse, error) {
	// Create the request
	req, err := s.client.newRequest(http.MethodDelete, fmt.Sprintf("/v2/guardrails/%s", ruleID), nil)
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
		return nil, fmt.Errorf("failed to delete guardrail: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var deleteResp DeleteGuardrailResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
		return nil, fmt.Errorf("error parsing delete guardrail response: %s", err)
	}

	return &deleteResp, nil
}
