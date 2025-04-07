package client

// This file contains shared data models used across multiple services

// RunnerTypeEnum represents the type of CI/CD runner
type RunnerTypeEnum string

const (
	RunnerTypeGithubActions    RunnerTypeEnum = "github-actions"
	RunnerTypeGitlabPipelines  RunnerTypeEnum = "gitlab-pipelines"
	RunnerTypeBitbucketPipelines RunnerTypeEnum = "bitbucket-pipelines"
	RunnerTypeAzurePipelines   RunnerTypeEnum = "azure-pipelines"
	RunnerTypeJenkins          RunnerTypeEnum = "jenkins"
	RunnerTypeSemaphore        RunnerTypeEnum = "semaphore"
	RunnerTypeAtlantis         RunnerTypeEnum = "atlantis"
	RunnerTypeEnv0             RunnerTypeEnum = "env0"
	RunnerTypeUnrecognized     RunnerTypeEnum = "unrecognized"
)

// WorkspaceRunStatusEnum represents the status of a workspace run
type WorkspaceRunStatusEnum string

const (
	WorkspaceRunStatusInitError     WorkspaceRunStatusEnum = "init_error"
	WorkspaceRunStatusPlanning      WorkspaceRunStatusEnum = "planning"
	WorkspaceRunStatusPlanSuccess   WorkspaceRunStatusEnum = "plan_success"
	WorkspaceRunStatusPlanError     WorkspaceRunStatusEnum = "plan_error"
	WorkspaceRunStatusApplying      WorkspaceRunStatusEnum = "applying"
	WorkspaceRunStatusApplyError    WorkspaceRunStatusEnum = "apply_error"
	WorkspaceRunStatusApplySuccess  WorkspaceRunStatusEnum = "apply_success"
	WorkspaceRunStatusBlocked       WorkspaceRunStatusEnum = "blocked"
	WorkspaceRunStatusAcknowledged  WorkspaceRunStatusEnum = "Acknowledged"
	WorkspaceRunStatusPlanNoChanges WorkspaceRunStatusEnum = "plan_no_changes"
	WorkspaceRunStatusApplyNoChanges WorkspaceRunStatusEnum = "apply_no_changes"
)

// VcsTypeEnum represents the type of version control system
type VcsTypeEnum string

const (
	VcsTypeGithub       VcsTypeEnum = "github"
	VcsTypeGitlab       VcsTypeEnum = "gitlab"
	VcsTypeBitbucket    VcsTypeEnum = "bitbucket"
	VcsTypeCodecommit   VcsTypeEnum = "codecommit"
	VcsTypeAzureDevops  VcsTypeEnum = "azuredevops"
)

// ResourceActionEnum represents the type of action performed on a resource
type ResourceActionEnum string

const (
	ResourceActionNoOp   ResourceActionEnum = "no-op"
	ResourceActionCreate ResourceActionEnum = "create"
	ResourceActionRead   ResourceActionEnum = "read"
	ResourceActionUpdate ResourceActionEnum = "update"
	ResourceActionDelete ResourceActionEnum = "delete"
)

// TaskTypeEnum represents the type of task in the workflow
type TaskTypeEnum string

const (
	TaskTypePostPlan  TaskTypeEnum = "post-plan"
	TaskTypePostApply TaskTypeEnum = "post-apply"
)

// GuardrailTypeEnum represents the type of the guardrail rule
type GuardrailTypeEnum string

const (
	GuardrailTypePolicy   GuardrailTypeEnum = "policy"
	GuardrailTypeCost     GuardrailTypeEnum = "cost"
	GuardrailTypeResource GuardrailTypeEnum = "resource"
	GuardrailTypeTag      GuardrailTypeEnum = "tag"
)

// TagEnforcementModeEnum represents the mode of tag enforcement
type TagEnforcementModeEnum string

const (
	TagEnforcementModeRequiredTags  TagEnforcementModeEnum = "requiredTags"
	TagEnforcementModeAnyTags       TagEnforcementModeEnum = "anyTags"
	TagEnforcementModeRequiredValues TagEnforcementModeEnum = "requiredValues"
)

// SortConfig represents a sort configuration
type SortConfig struct {
	Order int `json:"order,omitempty"`
}

// WorkspaceSortField represents the field to sort workspaces by
type WorkspaceSortField struct {
	LastRunTime     *SortConfig `json:"lastRunTime,omitempty"`
	WorkspaceName   *SortConfig `json:"workspaceName,omitempty"`
	Repository      *SortConfig `json:"repository,omitempty"`
	TerraformVersion *SortConfig `json:"terraformVersion,omitempty"`
	LastRunStatus   *SortConfig `json:"lastRunStatus,omitempty"`
	CITool          *SortConfig `json:"ciTool,omitempty"`
}

// WorkspaceRunSortField represents the field to sort workspace runs by
type WorkspaceRunSortField struct {
	UpdatedAt  *SortConfig `json:"updatedAt,omitempty"`
	Title      *SortConfig `json:"title,omitempty"`
	CommitID   *SortConfig `json:"commitId,omitempty"`
	Branch     *SortConfig `json:"branch,omitempty"`
	Status     *SortConfig `json:"status,omitempty"`
	CITool     *SortConfig `json:"ciTool,omitempty"`
	VCSType    *SortConfig `json:"vcsType,omitempty"`
	Repository *SortConfig `json:"repository,omitempty"`
}

// GuardrailSortField represents the field to sort guardrail rules by
type GuardrailSortField struct {
	CreatedAt  *SortConfig `json:"createdAt,omitempty"`
	AccountID  *SortConfig `json:"accountId,omitempty"`
	CreatedBy  *SortConfig `json:"createdBy,omitempty"`
	Name       *SortConfig `json:"name,omitempty"`
	Severity   *SortConfig `json:"severity,omitempty"`
}
