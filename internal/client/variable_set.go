package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// VariableSetService handles communication with the variable sets related methods of the Firefly API
type VariableSetService struct {
	client *Client
}

// CreateVariableSetRequest represents a request to create a new variable set
type CreateVariableSetRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
	Parents     []string   `json:"parents,omitempty"`
	Variables   []Variable `json:"variables,omitempty"`
}

// UpdateVariableSetRequest represents a request to update a variable set
type UpdateVariableSetRequest struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
	Parents     []string   `json:"parents,omitempty"`
	Variables   []Variable `json:"variables,omitempty"`
}

// CreateVariableSetResponse represents the response from creating a variable set
type CreateVariableSetResponse struct {
	VariableSetID string `json:"variableSetId"`
}

// VariableSet represents a Firefly variable set
type VariableSet struct {
	ID          string     `json:"id"`
	Version     int        `json:"version"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Labels      []string   `json:"labels"`
	Parents     []string   `json:"parents"`
	Variables   []Variable `json:"variables"`
	Descendants []string   `json:"descendants"`
}

// UpsertVariableSetVariablesRequest represents a request to upsert variables in a variable set
type UpsertVariableSetVariablesRequest struct {
	Variables []Variable `json:"variables"`
}

// DeleteVariablesRequest represents a request to delete variables from a variable set
type DeleteVariablesRequest struct {
	VariableIDs []string `json:"variableIds"`
}

// CreateVariableSet creates a new variable set
func (s *VariableSetService) CreateVariableSet(req CreateVariableSetRequest) (*CreateVariableSetResponse, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodPost, "/v2/runners/variables/variable-sets", req)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-201 responses
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create variable set: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var createResp CreateVariableSetResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &createResp, nil
}

// GetVariableSet retrieves a variable set by ID
func (s *VariableSetService) GetVariableSet(id string) (*VariableSet, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodGet, fmt.Sprintf("/v2/runners/variables/variable-sets/%s", id), nil)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get variable set: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var variableSet VariableSet
	if err := json.NewDecoder(resp.Body).Decode(&variableSet); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &variableSet, nil
}

// UpdateVariableSet updates an existing variable set
func (s *VariableSetService) UpdateVariableSet(id string, req UpdateVariableSetRequest) (*VariableSet, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodPut, fmt.Sprintf("/v2/runners/variables/variable-sets/%s", id), req)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update variable set: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var variableSet VariableSet
	if err := json.NewDecoder(resp.Body).Decode(&variableSet); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &variableSet, nil
}

// DeleteVariableSet deletes a variable set
func (s *VariableSetService) DeleteVariableSet(id string) error {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodDelete, fmt.Sprintf("/v2/runners/variables/variable-sets/%s", id), nil)
	if err != nil {
		return err
	}

	// Execute the request
	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle non-204 responses
	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete variable set: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	return nil
}

// UpsertVariablesInSet creates or updates variables within a variable set
func (s *VariableSetService) UpsertVariablesInSet(id string, req UpsertVariableSetVariablesRequest) ([]Variable, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodPost, fmt.Sprintf("/v2/runners/variables/variable-sets/%s/variables", id), req)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to upsert variables in set: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var variables []Variable
	if err := json.NewDecoder(resp.Body).Decode(&variables); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return variables, nil
}

// DeleteVariablesFromSet deletes variables from a variable set
func (s *VariableSetService) DeleteVariablesFromSet(id string, req DeleteVariablesRequest) error {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodDelete, fmt.Sprintf("/v2/runners/variables/variable-sets/%s/variables", id), req)
	if err != nil {
		return err
	}

	// Execute the request
	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle non-204 responses
	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete variables from set: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	return nil
}

// ListVariableSets retrieves all variable sets with pagination support
func (s *VariableSetService) ListVariableSets(pageSize, offset int, searchQuery string) ([]VariableSet, error) {
	// Create the request URL with query parameters
	url := fmt.Sprintf("/v2/runners/variables/variable-sets?pageSize=%d&offset=%d", pageSize, offset)
	if searchQuery != "" {
		url += "&searchQuery=" + searchQuery
	}

	// Create the request
	httpReq, err := s.client.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list variable sets: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var variableSets []VariableSet
	if err := json.NewDecoder(resp.Body).Decode(&variableSets); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return variableSets, nil
}