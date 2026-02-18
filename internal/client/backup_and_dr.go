package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// BackupAndDrService provides access to the Backup & DR API methods
type BackupAndDrService struct {
	client *Client
}

// PolicyCreateRequest represents a backup policy for API requests (POST)
type PolicyCreateRequest struct {
	PolicyName          string         `json:"policy_name"`
	IntegrationID       string         `json:"integration_id"`
	Region              string         `json:"region"`
	ProviderType        string         `json:"provider_type"`
	Schedule            ScheduleConfig `json:"schedule"`
	Description         string         `json:"description,omitempty"`
	Scope               []ScopeConfig  `json:"scope,omitempty"`
	NotificationID      string         `json:"notificationId,omitempty"`
	VCS                 *VCSConfig     `json:"vcs,omitempty"`
	RestoreInstructions string         `json:"restore_instructions,omitempty"`
	BackupOnSave        bool           `json:"backup_on_save,omitempty"`
}

// PolicyUpdateRequest represents a backup policy update for API requests (PUT)
// All fields are optional to support partial updates
type PolicyUpdateRequest struct {
	PolicyName          *string         `json:"policy_name,omitempty"`
	IntegrationID       *string         `json:"integration_id,omitempty"`
	Region              *string         `json:"region,omitempty"`
	ProviderType        *string         `json:"provider_type,omitempty"`
	Schedule            *ScheduleConfig `json:"schedule,omitempty"`
	Description         *string         `json:"description,omitempty"`
	Scope               []ScopeConfig   `json:"scope,omitempty"`
	NotificationID      *string         `json:"notificationId,omitempty"`
	VCS                 *VCSConfig      `json:"vcs,omitempty"`
	RestoreInstructions *string         `json:"restore_instructions,omitempty"`
	// Note: backup_on_save is only for creation, not updates
}

// PolicyResponse represents a backup policy from API responses
type PolicyResponse struct {
	PolicyID             string         `json:"policy_id"`
	AccountID            string         `json:"account_id"`
	PolicyName           string         `json:"policy_name"`
	IntegrationID        string         `json:"integration_id"`
	Region               string         `json:"region"`
	ProviderType         string         `json:"provider_type"`
	Schedule             ScheduleConfig `json:"schedule"`
	Description          string         `json:"description,omitempty"`
	Scope                []ScopeConfig  `json:"scope,omitempty"`
	NotificationID       string         `json:"notificationId,omitempty"`
	VCS                  *VCSConfig     `json:"vcs,omitempty"`
	RestoreInstructions  string         `json:"restore_instructions,omitempty"`
	BackupOnSave         bool           `json:"backup_on_save"`
	Status               string         `json:"status"`
	SnapshotsCount       int            `json:"snapshots_count"`
	LastBackupSnapshotID string         `json:"last_backup_snapshot_id,omitempty"`
	LastBackupTime       string         `json:"last_backup_time,omitempty"`
	LastBackupStatus     string         `json:"last_backup_status,omitempty"`
	NextBackupTime       string         `json:"next_backup_time,omitempty"`
	CreatedAt            string         `json:"created_at"`
	UpdatedAt            string         `json:"updated_at"`
}

// ScheduleConfig represents the backup schedule configuration
type ScheduleConfig struct {
	Frequency           string   `json:"frequency"`
	Hour                int      `json:"hour,omitempty"`
	Minute              int      `json:"minute,omitempty"`
	DaysOfWeek          []string `json:"days_of_week,omitempty"`
	MonthlyScheduleType string   `json:"monthly_schedule_type,omitempty"`
	DayOfMonth          int      `json:"day_of_month,omitempty"`
	WeekdayOrdinal      string   `json:"weekday_ordinal,omitempty"`
	WeekdayName         string   `json:"weekday_name,omitempty"`
	CronExpression      string   `json:"cron_expression,omitempty"`
}

// ScopeConfig represents a resource scope configuration
type ScopeConfig struct {
	Type  string   `json:"type"`
	Value []string `json:"value"`
}

// VCSConfig represents VCS integration configuration
type VCSConfig struct {
	ProjectID        string `json:"project_id,omitempty"`
	VCSIntegrationID string `json:"vcs_integration_id,omitempty"`
	RepoID           string `json:"repo_id,omitempty"`
}

// PolicyListFilters represents filters for listing policies
type PolicyListFilters struct {
	Status        string
	IntegrationID string
	Region        string
	ProviderType  string
}

// Pagination represents pagination metadata
type Pagination struct {
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Total    int  `json:"total"`
	HasNext  bool `json:"has_next"`
	HasPrev  bool `json:"has_prev"`
}

// FacetValue represents a single facet value with its count
type FacetValue struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// Facet represents a filterable field with its available values and counts
type Facet struct {
	Field      string       `json:"field"`
	Size       int          `json:"size"`
	Pagination Pagination   `json:"pagination"`
	Values     []FacetValue `json:"values"`
}

// PolicyListResponse represents the response from listing policies
type PolicyListResponse struct {
	Data       []PolicyResponse `json:"data"`
	Pagination Pagination       `json:"pagination"`
	Facets     []Facet          `json:"facets,omitempty"`
}

// ConvertCreateToUpdate converts a PolicyCreateRequest to PolicyUpdateRequest
// This is useful when the same data model is used for both create and update operations
func ConvertCreateToUpdate(create *PolicyCreateRequest) *PolicyUpdateRequest {
	if create == nil {
		return nil
	}

	update := &PolicyUpdateRequest{
		PolicyName:          &create.PolicyName,
		IntegrationID:       &create.IntegrationID,
		Region:              &create.Region,
		ProviderType:        &create.ProviderType,
		Schedule:            &create.Schedule,
		RestoreInstructions: &create.RestoreInstructions,
	}

	if create.Description != "" {
		update.Description = &create.Description
	}

	if create.NotificationID != "" {
		update.NotificationID = &create.NotificationID
	}

	if len(create.Scope) > 0 {
		update.Scope = create.Scope
	}

	if create.VCS != nil {
		update.VCS = create.VCS
	}

	return update
}

// Create creates a new backup policy
func (s *BackupAndDrService) Create(policy *PolicyCreateRequest) (*PolicyResponse, error) {
	endpoint := "/v2/backup-and-dr/policies"

	req, err := s.client.newRequest("POST", endpoint, policy)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result PolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// Get retrieves a specific backup policy by ID
func (s *BackupAndDrService) Get(policyID string) (*PolicyResponse, error) {
	endpoint := fmt.Sprintf("/v2/backup-and-dr/policies/%s", url.PathEscape(policyID))

	req, err := s.client.newRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("policy not found: %s", policyID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result PolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// Update updates an existing backup policy
func (s *BackupAndDrService) Update(policyID string, policy *PolicyUpdateRequest) (*PolicyResponse, error) {
	endpoint := fmt.Sprintf("/v2/backup-and-dr/policies/%s", url.PathEscape(policyID))

	req, err := s.client.newRequest("PUT", endpoint, policy)
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

	var result PolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// Delete deletes a backup policy
func (s *BackupAndDrService) Delete(policyID string) error {
	endpoint := fmt.Sprintf("/v2/backup-and-dr/policies/%s", url.PathEscape(policyID))

	req, err := s.client.newRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Accept both 200 and 204 for successful deletion
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// List retrieves all backup policies with optional filters
func (s *BackupAndDrService) List(filters *PolicyListFilters) (*PolicyListResponse, error) {
	endpoint := "/v2/backup-and-dr/policies"

	// Build query parameters if filters are provided
	queryParams := url.Values{}
	if filters != nil {
		if filters.Status != "" {
			queryParams.Add("status", filters.Status)
		}
		if filters.IntegrationID != "" {
			queryParams.Add("integration_id", filters.IntegrationID)
		}
		if filters.Region != "" {
			queryParams.Add("region", filters.Region)
		}
		if filters.ProviderType != "" {
			queryParams.Add("provider_type", filters.ProviderType)
		}
	}

	if len(queryParams) > 0 {
		endpoint = endpoint + "?" + queryParams.Encode()
	}

	req, err := s.client.newRequest("GET", endpoint, nil)
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

	var result PolicyListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}
