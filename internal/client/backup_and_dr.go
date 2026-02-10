package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type BackupAndDRService struct {
	client *Client
}

// BackupPolicy represents a backup & DR policy from API responses
type BackupPolicy struct {
	ID               string          `json:"policy_id,omitempty"`
	AccountID        string          `json:"account_id,omitempty"`
	PolicyName       string          `json:"policy_name"`
	Description      string          `json:"description,omitempty"`
	IntegrationID    string          `json:"integration_id"`
	Region           string          `json:"region"`
	ProviderType     string          `json:"provider_type"`
	Schedule         *ScheduleConfig `json:"schedule,omitempty"`
	Scope            []ScopeConfig   `json:"scope,omitempty"`
	NotificationID   string          `json:"notification_id,omitempty"`
	BackupOnSave     *bool           `json:"backup_on_save,omitempty"`
	Status           string          `json:"status,omitempty"`
	SnapshotsCount   int             `json:"snapshots_count,omitempty"`
	LastBackupTime   string          `json:"last_backup_time,omitempty"`
	LastBackupStatus string          `json:"last_backup_status,omitempty"`
	NextBackupTime   string          `json:"next_backup_time,omitempty"`
	VCS              *VCSConfig      `json:"vcs,omitempty"`
	CreatedAt        string          `json:"created_at,omitempty"`
	UpdatedAt        string          `json:"updated_at,omitempty"`
}

type ScheduleConfig struct {
	Frequency           string   `json:"frequency"`
	Hour                *int     `json:"hour,omitempty"`
	Minute              *int     `json:"minute,omitempty"`
	DaysOfWeek          []string `json:"days_of_week,omitempty"`
	MonthlyScheduleType string   `json:"monthly_schedule_type,omitempty"`
	DayOfMonth          *int     `json:"day_of_month,omitempty"`
	WeekdayOrdinal      string   `json:"weekday_ordinal,omitempty"`
	WeekdayName         string   `json:"weekday_name,omitempty"`
	CronExpression      string   `json:"cron_expression,omitempty"`
}

type ScopeConfig struct {
	Type  string   `json:"type"`
	Value []string `json:"value"`
}

type VCSConfig struct {
	ProjectID        string `json:"project_id"`
	VCSIntegrationID string `json:"vcs_integration_id"`
	RepoID           string `json:"repo_id"`
}

type BackupPolicyCreateRequest struct {
	PolicyName     string          `json:"policy_name"`
	Description    string          `json:"description,omitempty"`
	IntegrationID  string          `json:"integration_id"`
	Region         string          `json:"region"`
	ProviderType   string          `json:"provider_type"`
	Schedule       *ScheduleConfig `json:"schedule,omitempty"`
	Scope          []ScopeConfig   `json:"scope,omitempty"`
	NotificationID string          `json:"notification_id,omitempty"`
	BackupOnSave   *bool           `json:"backup_on_save,omitempty"`
	VCS            *VCSConfig      `json:"vcs,omitempty"`
}

type BackupPolicyUpdateRequest struct {
	PolicyName     string          `json:"policy_name"`
	Description    string          `json:"description,omitempty"`
	IntegrationID  string          `json:"integration_id"`
	Region         string          `json:"region"`
	ProviderType   string          `json:"provider_type"`
	Schedule       *ScheduleConfig `json:"schedule,omitempty"`
	Scope          []ScopeConfig   `json:"scope,omitempty"`
	NotificationID string          `json:"notification_id,omitempty"`
	BackupOnSave   *bool           `json:"backup_on_save,omitempty"`
	VCS            *VCSConfig      `json:"vcs,omitempty"`
}

type BackupPolicyListResponse struct {
	Policies []BackupPolicy `json:"policies"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

func (s *BackupAndDRService) Create(req *BackupPolicyCreateRequest) (*BackupPolicy, error) {
	httpReq, err := s.client.newRequest(http.MethodPost, "/v2/backup-and-dr/policies", req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result BackupPolicy
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func (s *BackupAndDRService) Get(policyID string) (*BackupPolicy, error) {
	endpoint := fmt.Sprintf("/v2/backup-and-dr/policies/%s", url.PathEscape(policyID))

	httpReq, err := s.client.newRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("backup policy not found: %s", policyID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result BackupPolicy
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func (s *BackupAndDRService) Update(policyID string, req *BackupPolicyUpdateRequest) (*BackupPolicy, error) {
	endpoint := fmt.Sprintf("/v2/backup-and-dr/policies/%s", url.PathEscape(policyID))

	httpReq, err := s.client.newRequest(http.MethodPut, endpoint, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result BackupPolicy
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func (s *BackupAndDRService) Delete(policyID string, cascade bool) error {
	endpoint := fmt.Sprintf("/v2/backup-and-dr/policies/%s?cascade=%t", url.PathEscape(policyID), cascade)

	httpReq, err := s.client.newRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.doRequest(httpReq)
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

func (s *BackupAndDRService) SetStatus(policyID string, status string) error {
	endpoint := fmt.Sprintf("/v2/backup-and-dr/policies/%s/status", url.PathEscape(policyID))

	body := map[string]string{"status": status}

	httpReq, err := s.client.newRequest(http.MethodPatch, endpoint, body)
	if err != nil {
		return err
	}

	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (s *BackupAndDRService) List(page, pageSize int) (*BackupPolicyListResponse, error) {
	endpoint := fmt.Sprintf("/v2/backup-and-dr/policies?page=%d&pageSize=%d", page, pageSize)

	httpReq, err := s.client.newRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result BackupPolicyListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}
