# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a comprehensive Terraform Provider for Firefly, a SaaS platform for infrastructure management. The provider allows users to manage Firefly resources (projects, workspaces, variable sets, guardrails, labels) using Terraform with full CRUD operations and extensive testing coverage.

**Current Status**: The provider is fully functional but not yet published to the Terraform Registry. Users must run it locally in debug mode for now.

## Architecture
The codebase follows HashiCorp Terraform provider best practices and standards:

- **main.go**: Entry point with version tracking and debug support
- **firefly/**: Core provider implementation using Terraform Plugin Framework v1.15.1
  - `provider.go`: Main provider configuration with authentication
  - `resource_*.go`: Complete resource implementations for all Firefly resources
  - `data_source_*.go`: Data source implementations for resource discovery
- **internal/client/**: Comprehensive API client for Firefly REST API
  - `client.go`: Core HTTP client with automatic authentication and token renewal
  - `project.go`: Projects API with full CRUD operations
  - `variable_set.go`: Variable Sets API with inheritance support
  - `runners_workspace.go`: Runners Workspace API with VCS integration
  - `guardrail.go`: Guardrail API with policy management
  - `workspace.go`: Legacy workspace operations
  - `models.go`: Shared data models and enums

## Key Components
- **Authentication**: Uses access_key/secret_key that authenticate via `/api/v2/login` to get Bearer tokens
- **Resources**: `firefly_workflows_project`, `firefly_workflows_runners_workspace`, `firefly_workflows_variable_set`, `firefly_workflows_guardrail`, `firefly_workspace_labels`
- **Data Sources**: `firefly_workflows_projects`, `firefly_workflows_project`, `firefly_workflows_variable_sets`, `firefly_workflows_variable_set`, `firefly_workspaces`, `firefly_workspace_runs`, `firefly_workflows_guardrails`
- **API Client**: Production-ready HTTP client with comprehensive error handling and automatic token renewal

## Common Commands

### Development
```bash
# Build the provider
make build
# or: go install

# Generate documentation  
make docs
# or: go generate

# Format code
make fmt

# Run linting
make vet

# Run with debug support
make debug
# or: go run main.go -debug
```

### Testing
```bash
# Run unit tests
make test

# Run acceptance tests (creates real resources)
make testacc

# Run specific tests
go test ./internal/client -v -run TestProject
go test ./internal/provider -v -run TestAccProject

# Run all tests
go test ./... -v
```

### Documentation Generation
⚠️ **DO NOT USE**: `github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs` - This tool is not working properly.

Create documentation manually in the `docs/` directory following Terraform Registry format.

## Configuration
Provider expects:
- `access_key` (required, sensitive)
- `secret_key` (required, sensitive) 
- `api_url` (optional, defaults to https://api.gofirefly.io)

Environment variables supported:
- `FIREFLY_ACCESS_KEY`
- `FIREFLY_SECRET_KEY`
- `FIREFLY_API_URL`

## Development Notes
- Uses Terraform Plugin Framework v1.15.1 (modern framework, not legacy SDK v2)
- Go module: `github.com/gofireflyio/terraform-provider-firefly`
- Go version: 1.23.7+ (defined in .go-version)
- Professional build system with Makefile and release automation (.goreleaser.yml)
- Comprehensive testing suite with 49 total tests (21 unit + 28 acceptance tests)
- All compilation issues resolved and production-ready as of 2025-08-18

## Testing the Provider (Local Development)

Since the provider is not yet in the Terraform Registry, use this workflow:

```bash
# 1. Build the provider
go build -o terraform-provider-firefly

# 2. Start in debug mode (keep this running)
go run main.go -debug
# This outputs: TF_REATTACH_PROVIDERS='{"registry.terraform.io/firefly/firefly":{"Protocol":"grpc",...}}'

# 3. In another terminal, set the environment variable and use Terraform
export TF_REATTACH_PROVIDERS='{"registry.terraform.io/firefly/firefly":{"Protocol":"grpc","ProtocolVersion":6,"Pid":12345,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin123"}}}'

# 4. Create terraform configuration with correct source
terraform {
  required_providers {
    firefly = {
      source = "registry.terraform.io/firefly/firefly"  # Required for debug mode
    }
  }
}

provider "firefly" {
  access_key = "your-access-key"
  secret_key = "your-secret-key"
}

# 5. Use terraform normally
terraform init
terraform plan
terraform apply
```

### Working Example
The provider has been tested successfully with these resources:
- ✅ `firefly_workflows_project` - Creates projects (verified working)
- ✅ `firefly_project_membership` - Manages project member assignments (NEW)
- ✅ `firefly_workflows_variable_set` - Manages variable sets  
- ✅ `firefly_workflows_runners_workspace` - Manages runner workspaces (fixed 2025-08-27)
- ✅ `firefly_workflows_guardrail` - Manages governance rules
- ✅ Data sources for all above resources

## Testing Coverage
The provider includes comprehensive test coverage:

### Unit Tests (21 tests)
- **Client Layer**: Full API client testing with mock servers
- **Authentication**: Token handling and renewal
- **CRUD Operations**: All resource types (Projects, Variable Sets, Guardrails)
- **Error Handling**: Network failures, API errors, validation

### Acceptance Tests (28 test scenarios)
- **All Resources**: Complete lifecycle testing (Create, Read, Update, Delete, Import)
- **Complex Scenarios**: Variable inheritance, resource relationships, VCS integration
- **Data Sources**: Resource discovery and cross-references
- **Real-World Examples**: Mirrors documentation usage patterns

### Test Commands
```bash
# Unit tests only
make test

# Acceptance tests (requires TF_ACC=1 and valid Firefly credentials)
TF_ACC=1 make testacc

# Specific test suites
go test ./internal/client -v
go test ./internal/provider -v
```

## Registry Status

🚧 **Not Yet Published**: This provider is not published to the Terraform Registry yet.

**Current Usage**: All users must use local development mode with debug provider.

**Publication Checklist**:
- ✅ All resources implemented and tested
- ✅ Comprehensive documentation created  
- ✅ 49 test scenarios passing
- ✅ Professional project structure established
- ⏳ Terraform Registry publication pending

## Recently Fixed Issues

### Runners Workspace Creation (Fixed 2025-08-27)
**Previous Issue**: The `/v2/runners/workspaces` POST endpoint was returning an empty response instead of the created workspace object.

**Resolution**: 
- ✅ API now properly returns the complete workspace object including ID
- ✅ Terraform can fully manage runners workspaces (create, read, update, delete)
- ✅ Full state management is working correctly

**Additional Improvements**:
- ✅ Added workspace name validation to prevent spaces and invalid characters
- ✅ Fixed runner type from "firefly_runners" to "firefly"

## Recent Major Updates
- **2025-08-28**: **CRITICAL FIXES**: Resolved workspace-project relationship and deletion issues
  - ✅ Fixed "storage failed to get project id from workspace id: no records matched" error
  - ✅ Fixed "unknown value for consumed_variable_sets" after apply
  - ✅ Fixed workspace deletion failing due to incorrect status code handling
  - ✅ Added missing ProjectID field to RunnersWorkspace API response struct
  - ✅ Enhanced Read function to preserve project relationships from state
  - ✅ Provider now accepts both 200 and 204 status codes for successful deletion
  - **Root Cause**: API response was missing projectId field, causing state inconsistency
  - **Testing**: Comprehensively tested with real Firefly API using full create/read/update/delete lifecycle
- **2025-08-27**: **RESOLVED**: Fixed all runners workspace creation issues
  - Fixed "invalid runner type" error by changing `"firefly_runners"` to `"firefly"` 
  - Added comprehensive workspace name validation (no spaces, valid characters only)
  - API team fixed the `/v2/runners/workspaces` POST endpoint to return workspace object
  - Full Terraform support for runners workspaces now working
- **2025-08-19**: Added `firefly_project_membership` resource for managing project member assignments
  - New resource allows adding users to projects with specific roles (admin, member, viewer)
  - Enables projects created via Terraform to be visible in the UI by assigning members
  - Includes comprehensive tests and documentation
  - Import/export functionality with format `project_id:user_id`
- **2025-08-19**: Fixed variable set resource to properly populate computed fields after creation/update
  - Issue: `parents` and `version` fields were showing as unknown after apply
  - Solution: Added API call to fetch complete resource after create/update operations
  - All computed fields now properly populated in Terraform state
- **2025-08-19**: Resources renamed to `workflows_*` prefix for clarity
- **2025-08-19**: Added comprehensive documentation structure
- **2025-08-19**: Verified provider functionality with real API testing
- **2025-08-18**: Complete testing suite with 49 comprehensive tests
- **2025-08-18**: Modernized project structure following HashiCorp standards
- **2025-08-18**: Added Makefile, .goreleaser.yml, and development tooling
- **2025-08-18**: Updated to latest Terraform Plugin Framework (v1.15.1)
- **2025-08-18**: Added all major Firefly resources (Projects, Variable Sets, Runners Workspaces)
- **2025-08-14**: Fixed all compilation issues and established stable foundation