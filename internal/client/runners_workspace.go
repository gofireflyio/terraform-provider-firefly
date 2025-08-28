package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RunnersWorkspaceService handles communication with the runners workspace related methods of the Firefly API
type RunnersWorkspaceService struct {
	client *Client
}

// IacProvisioner represents the IaC provisioner configuration
type IacProvisioner struct {
	Type    string `json:"type"`    // terraform or opentofu
	Version string `json:"version"`
}

// ExecutionConfig represents execution configuration
type ExecutionConfig struct {
	Triggers         []string `json:"triggers"`         // e.g., ["merge"]
	ApplyRule        string   `json:"applyRule"`        // manual or auto
	TerraformVersion string   `json:"terraformVersion"`
}

// CreateRunnersWorkspaceRequest represents a request to create a new runners workspace
type CreateRunnersWorkspaceRequest struct {
	RunnerType              string          `json:"runnerType"`
	IacType                 string          `json:"iacType"`
	WorkspaceName           string          `json:"workspaceName"`
	Description             string          `json:"description,omitempty"`
	Labels                  []string        `json:"labels,omitempty"`
	VcsID                   string          `json:"vcsId"`
	Repo                    string          `json:"repo"`
	DefaultBranch           string          `json:"defaultBranch"`
	VcsType                 string          `json:"vcsType"`
	WorkDir                 string          `json:"workDir"`
	Variables               []Variable      `json:"variables"`
	ConsumedVariableSets    []string        `json:"consumedVariableSets,omitempty"`
	Execution               ExecutionConfig `json:"execution"`
	Project                 *string         `json:"project,omitempty"`
	TerraformVariables      map[string]interface{} `json:"terraformVariables,omitempty"`
	TerraformSensitiveVariables map[string]interface{} `json:"terraformSensitiveVariables,omitempty"`
	ProvidersCredentials    map[string]interface{} `json:"providersCredentials,omitempty"`
	RunnerEnvironment       map[string]interface{} `json:"runnerEnvironment,omitempty"`
}

// UpdateRunnersWorkspaceRequest represents a request to update a runners workspace
type UpdateRunnersWorkspaceRequest struct {
	Name                    string          `json:"name,omitempty"`
	Description             string          `json:"description,omitempty"`
	Labels                  []string        `json:"labels,omitempty"`
	VcsIntegrationID        string          `json:"vcsIntegrationId,omitempty"`
	Repository              string          `json:"repository,omitempty"`
	DefaultBranch           string          `json:"defaultBranch,omitempty"`
	WorkingDirectory        string          `json:"workingDirectory,omitempty"`
	CronExecutionPattern    string          `json:"cronExecutionPattern,omitempty"`
	IacProvisioner          *IacProvisioner `json:"iacProvisioner,omitempty"`
	Variables               []Variable      `json:"variables,omitempty"`
	ConsumedVariableSets    []string        `json:"consumedVariableSets,omitempty"`
}

// RunnersWorkspace represents a Firefly runners workspace
type RunnersWorkspace struct {
	ID                   string          `json:"id"`
	Name                 string          `json:"name"`
	Description          string          `json:"description"`
	AccountID            string          `json:"accountId"`
	Repository           string          `json:"repository"`
	WorkingDirectory     string          `json:"workingDirectory"`
	VcsIntegrationID     string          `json:"vcsIntegrationId"`
	Vcs                  string          `json:"vcs"`
	DefaultBranch        string          `json:"defaultBranch"`
	CronExecutionPattern string          `json:"cronExecutionPattern"`
	IacProvisioner       *IacProvisioner `json:"iacProvisioner"`
	Labels               []string        `json:"labels"`
	ProjectID            string          `json:"projectId,omitempty"`  // Added missing project ID field
}

// TaskResponse represents the response from a task operation
type TaskResponse struct {
	TaskID string `json:"taskId"`
	Status string `json:"status"`
}

// RunTaskRequest represents a request to run a task on a workspace
type RunTaskRequest struct {
	TaskType string `json:"taskType"` // e.g., "destroy", "plan", "apply"
}

// CreateRunnersWorkspace creates a new runners workspace
func (s *RunnersWorkspaceService) CreateRunnersWorkspace(req CreateRunnersWorkspaceRequest) (*RunnersWorkspace, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodPost, "/v2/runners/workspaces", req)
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
		return nil, fmt.Errorf("failed to create runners workspace: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var workspace RunnersWorkspace
	if err := json.NewDecoder(resp.Body).Decode(&workspace); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &workspace, nil
}

// GetRunnersWorkspace retrieves a runners workspace by ID
func (s *RunnersWorkspaceService) GetRunnersWorkspace(id string) (*RunnersWorkspace, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodGet, fmt.Sprintf("/v2/runners/workspaces/%s", id), nil)
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
		return nil, fmt.Errorf("failed to get runners workspace: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var workspace RunnersWorkspace
	if err := json.NewDecoder(resp.Body).Decode(&workspace); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &workspace, nil
}

// UpdateRunnersWorkspace updates an existing runners workspace
func (s *RunnersWorkspaceService) UpdateRunnersWorkspace(id string, req UpdateRunnersWorkspaceRequest) (*RunnersWorkspace, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodPut, fmt.Sprintf("/v2/runners/workspaces/%s", id), req)
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
		return nil, fmt.Errorf("failed to update runners workspace: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var workspace RunnersWorkspace
	if err := json.NewDecoder(resp.Body).Decode(&workspace); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &workspace, nil
}

// DeleteRunnersWorkspace deletes a runners workspace
func (s *RunnersWorkspaceService) DeleteRunnersWorkspace(id string) error {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodDelete, fmt.Sprintf("/v2/runners/workspaces/%s", id), nil)
	if err != nil {
		return err
	}

	// Execute the request
	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle non-success responses (accept both 200 and 204)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete runners workspace: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	return nil
}

// DestroyWorkspaceResources initiates a destroy task to clean up cloud infrastructure resources
func (s *RunnersWorkspaceService) DestroyWorkspaceResources(id string, req RunTaskRequest) (*TaskResponse, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodPost, fmt.Sprintf("/v2/runners/workspaces/%s/tasks/destroy", id), req)
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
		return nil, fmt.Errorf("failed to destroy workspace resources: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var taskResp TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &taskResp, nil
}