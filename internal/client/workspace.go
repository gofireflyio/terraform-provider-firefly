package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// WorkspaceService provides access to the workspace-related API methods
type WorkspaceService struct {
	client *Client
}

// Workspace represents a Firefly workspace
type Workspace struct {
	ID                string   `json:"id,omitempty"`
	AccountID         string   `json:"accountId,omitempty"`
	WorkspaceID       string   `json:"workspaceId,omitempty"`
	WorkspaceName     string   `json:"workspaceName,omitempty"`
	Repo              string   `json:"repo,omitempty"`
	RepoURL           string   `json:"repoUrl,omitempty"`
	VCSType           string   `json:"vcsType,omitempty"`
	RunnerType        string   `json:"runnerType,omitempty"`
	LastRunStatus     string   `json:"lastRunStatus,omitempty"`
	LastApplyTime     string   `json:"lastApplyTime,omitempty"`
	LastPlanTime      string   `json:"lastPlanTime,omitempty"`
	LastRunTime       string   `json:"lastRunTime,omitempty"`
	IACType           string   `json:"iacType,omitempty"`
	IACTypeVersion    string   `json:"iacTypeVersion,omitempty"`
	Labels            []string `json:"labels,omitempty"`
	RunsCount         int      `json:"runsCount,omitempty"`
	IsWorkflowManaged bool     `json:"isWorkflowManaged,omitempty"`
	CreatedAt         string   `json:"createdAt,omitempty"`
	UpdatedAt         string   `json:"updatedAt,omitempty"`
}

// WorkspaceFilters represents filters for listing workspaces
type WorkspaceFilters struct {
	WorkspaceName    []string `json:"workspaceName,omitempty"`
	Repositories     []string `json:"repositories,omitempty"`
	CITool           []string `json:"ciTool,omitempty"`
	Labels           []string `json:"labels,omitempty"`
	Status           []string `json:"status,omitempty"`
	IsManagedWorkflow *bool    `json:"isManagedWorkflow,omitempty"`
	VCSType          []string `json:"vcsType,omitempty"`
}

// ListWorkspacesRequest represents the request body for listing workspaces
type ListWorkspacesRequest struct {
	Filters      *WorkspaceFilters `json:"filters,omitempty"`
	SearchValue  string            `json:"searchValue,omitempty"`
	Projection   []string          `json:"projection,omitempty"`
}

// WorkspaceRun represents a Firefly workspace run
type WorkspaceRun struct {
	ID            string `json:"id,omitempty"`
	WorkspaceID   string `json:"workspaceId,omitempty"`
	WorkspaceName string `json:"workspaceName,omitempty"`
	RunID         string `json:"runId,omitempty"`
	RunName       string `json:"runName,omitempty"`
	Status        string `json:"status,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	UpdatedAt     string `json:"updatedAt,omitempty"`
}

// WorkspaceRunFilters represents filters for listing workspace runs
type WorkspaceRunFilters struct {
	RunID         []string `json:"runId,omitempty"`
	RunName       []string `json:"runName,omitempty"`
	Status        []string `json:"status,omitempty"`
	Branch        []string `json:"branch,omitempty"`
	CommitID      []string `json:"commitId,omitempty"`
	CITool        []string `json:"ciTool,omitempty"`
	VCSType       []string `json:"vcsType,omitempty"`
	Repository    []string `json:"repository,omitempty"`
}

// ListWorkspaceRunsRequest represents the request body for listing workspace runs
type ListWorkspaceRunsRequest struct {
	Filters     *WorkspaceRunFilters `json:"filters,omitempty"`
	SearchValue string               `json:"searchValue,omitempty"`
	Projection  []string             `json:"projection,omitempty"`
}

// UpdateWorkspaceLabelsRequest represents the request body for updating workspace labels
type UpdateWorkspaceLabelsRequest struct {
	Labels []string `json:"labels"`
}

// UpdateWorkspaceLabelsResponse represents the response from updating workspace labels
type UpdateWorkspaceLabelsResponse struct {
	ID            string   `json:"id"`
	WorkspaceName string   `json:"workspaceName"`
	Labels        []string `json:"labels"`
	UpdatedAt     string   `json:"updatedAt"`
}

// DeleteWorkspaceResponse represents the response from deleting a workspace
type DeleteWorkspaceResponse struct {
	Status int `json:"status"`
	Data   struct {
		Message string `json:"message"`
	} `json:"data"`
}

// ListWorkspaces retrieves workspaces from Firefly
func (s *WorkspaceService) ListWorkspaces(request *ListWorkspacesRequest, page, pageSize int) ([]Workspace, error) {
	// Construct query parameters
	queryParams := url.Values{}
	queryParams.Add("page", strconv.Itoa(page))
	queryParams.Add("pageSize", strconv.Itoa(pageSize))
	
	// Create the request
	req, err := s.client.newRequest(http.MethodPost, fmt.Sprintf("/workspaces/search?%s", queryParams.Encode()), request)
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
		return nil, fmt.Errorf("failed to list workspaces: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}
	
	// Parse the response
	var workspaces []Workspace
	if err := json.NewDecoder(resp.Body).Decode(&workspaces); err != nil {
		return nil, fmt.Errorf("error parsing workspace list response: %s", err)
	}
	
	return workspaces, nil
}

// DeleteWorkspace deletes a workspace by ID
func (s *WorkspaceService) DeleteWorkspace(workspaceID string) (*DeleteWorkspaceResponse, error) {
	// Create the request
	req, err := s.client.newRequest(http.MethodDelete, fmt.Sprintf("/workspaces/%s", workspaceID), nil)
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
		return nil, fmt.Errorf("failed to delete workspace: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}
	
	// Parse the response
	var deleteResp DeleteWorkspaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
		return nil, fmt.Errorf("error parsing delete workspace response: %s", err)
	}
	
	return &deleteResp, nil
}

// UpdateWorkspaceLabels updates the labels for a workspace
func (s *WorkspaceService) UpdateWorkspaceLabels(workspaceID string, labels []string) (*UpdateWorkspaceLabelsResponse, error) {
	// Create the request body
	reqBody := UpdateWorkspaceLabelsRequest{
		Labels: labels,
	}
	
	// Create the request
	req, err := s.client.newRequest(http.MethodPut, fmt.Sprintf("/workspaces/%s/labels", workspaceID), reqBody)
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
		return nil, fmt.Errorf("failed to update workspace labels: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}
	
	// Parse the response
	var updateResp UpdateWorkspaceLabelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return nil, fmt.Errorf("error parsing update workspace labels response: %s", err)
	}
	
	return &updateResp, nil
}

// ListWorkspaceRuns retrieves runs for a workspace
func (s *WorkspaceService) ListWorkspaceRuns(workspaceID string, request *ListWorkspaceRunsRequest, page, pageSize int) ([]WorkspaceRun, error) {
	// Construct query parameters
	queryParams := url.Values{}
	queryParams.Add("page", strconv.Itoa(page))
	queryParams.Add("pageSize", strconv.Itoa(pageSize))
	
	// Create the request
	req, err := s.client.newRequest(http.MethodPost, fmt.Sprintf("/workspaces/%s/runs/search?%s", workspaceID, queryParams.Encode()), request)
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
		return nil, fmt.Errorf("failed to list workspace runs: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}
	
	// Parse the response
	var runs []WorkspaceRun
	if err := json.NewDecoder(resp.Body).Decode(&runs); err != nil {
		return nil, fmt.Errorf("error parsing workspace runs response: %s", err)
	}
	
	return runs, nil
}
