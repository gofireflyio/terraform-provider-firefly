package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ProjectService handles communication with the projects related methods of the Firefly API
type ProjectService struct {
	client *Client
}

// VariableSensitivity represents the sensitivity level of a variable
type VariableSensitivity string

const (
	SensitivityString VariableSensitivity = "string"
	SensitivitySecret VariableSensitivity = "secret"
)

// VariableDestination represents where the variable should be used
type VariableDestination string

const (
	DestinationEnv VariableDestination = "env"
	DestinationIAC VariableDestination = "iac"
)

// Variable represents a project variable
type Variable struct {
	Key         string               `json:"key"`
	Value       string               `json:"value"`
	Sensitivity VariableSensitivity  `json:"sensitivity,omitempty"`
	Destination VariableDestination  `json:"destination,omitempty"`
}

// CreateProjectRequest represents a request to create a new project
type CreateProjectRequest struct {
	Name                 string     `json:"name"`
	Description          string     `json:"description,omitempty"`
	Labels               []string   `json:"labels,omitempty"`
	CronExecutionPattern string     `json:"cronExecutionPattern,omitempty"`
	Variables            []Variable `json:"variables,omitempty"`
	ParentID             string     `json:"parentId,omitempty"`
}

// UpdateProjectRequest represents a request to update a project
type UpdateProjectRequest struct {
	Name                 string     `json:"name,omitempty"`
	Description          string     `json:"description,omitempty"`
	Labels               []string   `json:"labels,omitempty"`
	CronExecutionPattern string     `json:"cronExecutionPattern,omitempty"`
	Variables            []Variable `json:"variables,omitempty"`
}

// Project represents a Firefly project
type Project struct {
	ID                   string     `json:"id"`
	AccountID            string     `json:"accountId"`
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	Labels               []string   `json:"labels"`
	CronExecutionPattern string     `json:"cronExecutionPattern"`
	Variables            []Variable `json:"variables"`
	MembersCount         int        `json:"membersCount"`
	WorkspaceCount       int        `json:"workspaceCount"`
	ParentID             string     `json:"parentId"`
}

// Member represents a project member
type Member struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(req CreateProjectRequest) (*Project, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodPost, "/v2/runners/projects", req)
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
		return nil, fmt.Errorf("failed to create project: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &project, nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(id string) (*Project, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodGet, fmt.Sprintf("/v2/runners/projects/%s", id), nil)
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
		return nil, fmt.Errorf("failed to get project: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &project, nil
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(id string, req UpdateProjectRequest) (*Project, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodPatch, fmt.Sprintf("/v2/runners/projects/%s", id), req)
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
		return nil, fmt.Errorf("failed to update project: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &project, nil
}

// DeleteProject deletes a project and all its associated workspaces
func (s *ProjectService) DeleteProject(id string) error {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodDelete, fmt.Sprintf("/v2/runners/projects/%s", id), nil)
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
		return fmt.Errorf("failed to delete project: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	return nil
}

// ListProjectMembers retrieves all members of a project
func (s *ProjectService) ListProjectMembers(projectID string) ([]Member, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodGet, fmt.Sprintf("/v2/runners/projects/%s/members", projectID), nil)
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
		return nil, fmt.Errorf("failed to list project members: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var members []Member
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return members, nil
}

// AddProjectMembers adds new members to a project
func (s *ProjectService) AddProjectMembers(projectID string, members []Member) ([]Member, error) {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodPost, fmt.Sprintf("/v2/runners/projects/%s/members", projectID), members)
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
		return nil, fmt.Errorf("failed to add project members: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response - try to decode as array first, fall back to success message
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Try to decode as array of members first
	var addedMembers []Member
	if err := json.Unmarshal(bodyBytes, &addedMembers); err == nil && len(addedMembers) > 0 {
		return addedMembers, nil
	}

	// Fall back to success message format
	var successResp struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(bodyBytes, &successResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Return the original members that were sent (API doesn't return member details)
	return members, nil
}

// RemoveProjectMembers removes members from a project
func (s *ProjectService) RemoveProjectMembers(projectID string, userIDs []string) error {
	// Create the request
	httpReq, err := s.client.newRequest(http.MethodDelete, fmt.Sprintf("/v2/runners/projects/%s/members", projectID), userIDs)
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
		return fmt.Errorf("failed to remove project members: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	return nil
}

// GetProjectMember retrieves a specific member of a project
func (s *ProjectService) GetProjectMember(projectID, userID string) (*Member, error) {
	members, err := s.ListProjectMembers(projectID)
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		if member.UserID == userID {
			return &member, nil
		}
	}

	return nil, fmt.Errorf("member with user ID %s not found in project %s", userID, projectID)
}

// AddProjectMember adds a single member to a project
func (s *ProjectService) AddProjectMember(projectID string, member Member) (*Member, error) {
	members := []Member{member}
	addedMembers, err := s.AddProjectMembers(projectID, members)
	if err != nil {
		return nil, err
	}

	if len(addedMembers) == 0 {
		return nil, fmt.Errorf("no members were added")
	}

	return &addedMembers[0], nil
}

// RemoveProjectMember removes a single member from a project
func (s *ProjectService) RemoveProjectMember(projectID, userID string) error {
	return s.RemoveProjectMembers(projectID, []string{userID})
}

// UpdateProjectMember updates a member's role in a project
func (s *ProjectService) UpdateProjectMember(projectID string, member Member) (*Member, error) {
	// First remove the member
	if err := s.RemoveProjectMember(projectID, member.UserID); err != nil {
		return nil, fmt.Errorf("failed to remove member for update: %w", err)
	}

	// Then add with the new role
	return s.AddProjectMember(projectID, member)
}

// ProjectsListResponse represents the response from listing projects
type ProjectsListResponse struct {
	Data       []Project `json:"data"`
	TotalCount int       `json:"totalCount"`
}

// ListProjects retrieves all projects with pagination support
func (s *ProjectService) ListProjects(pageSize, offset int, searchQuery string) (*ProjectsListResponse, error) {
	// Create the request URL with query parameters
	url := fmt.Sprintf("/v2/runners/projects/list?pageSize=%d&offset=%d", pageSize, offset)
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
		return nil, fmt.Errorf("failed to list projects: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response
	var projectsResp ProjectsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&projectsResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &projectsResp, nil
}