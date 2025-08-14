# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a Terraform Provider for Firefly, a SaaS platform for infrastructure management. The provider allows users to manage Firefly resources (workspaces, guardrails, labels) using Terraform.

## Architecture
The codebase follows standard Terraform provider patterns:

- **main.go**: Entry point that sets up the plugin server
- **internal/provider/**: Provider implementation using Terraform Plugin Framework
  - `provider.go`: Main provider configuration and setup
  - `resource_*.go`: Resource implementations (guardrails, workspace labels)
  - `data_source_*.go`: Data source implementations (workspaces, workspace runs, guardrails)
- **internal/client/**: API client for Firefly REST API
  - `client.go`: Core HTTP client with authentication
  - `workspace.go`: Workspace-related API methods
  - `guardrail.go`: Guardrail-related API methods
  - `models.go`: Shared data models

## Key Components
- **Authentication**: Uses access_key/secret_key that authenticate via `/api/v2/login` to get Bearer tokens
- **Resources**: `firefly_guardrail`, `firefly_workspace_labels`
- **Data Sources**: `firefly_workspaces`, `firefly_workspace_runs`, `firefly_guardrails`
- **API Client**: Custom HTTP client with automatic token renewal

## Common Commands

### Development
```bash
# Build the provider
go install

# Generate documentation
go generate

# Run with debug support
go run main.go -debug
```

### Testing
```bash
# Run acceptance tests (creates real resources)
make testacc
```

Note: Standard Go commands work as there's no custom Makefile present.

## Configuration
Provider expects:
- `access_key` (required, sensitive)
- `secret_key` (required, sensitive) 
- `api_url` (optional, defaults to https://api.firefly.ai)

Environment variables supported:
- `FIREFLY_ACCESS_KEY`
- `FIREFLY_SECRET_KEY`
- `FIREFLY_API_URL`

## Development Notes
- Uses Terraform Plugin Framework (not legacy SDK v2 for new resources)
- Go module: `terraform-provider-firefly`
- Go version: 1.21+
- No custom build system - standard Go toolchain
- All compilation issues have been resolved as of 2025-08-14

## Testing the Provider
```bash
# Build the provider
go build

# Test in debug mode
./terraform-provider-firefly -debug

# Use with Terraform (set the TF_REATTACH_PROVIDERS env var from debug output)
terraform init
terraform plan
```

## Recent Fixes Applied
- Fixed go.mod filename (was incorrectly named go.mod.go)
- Added missing `attr` imports to data source files
- Added missing fields to `WorkspaceRun` struct 
- Fixed `listToValues()` function return type
- Removed invalid `Optional` fields from `SingleNestedBlock` schemas
- Updated to Plugin Framework (removed Plugin SDK v2 dependency)
- Fixed main.go to use `providerserver` instead of legacy plugin system