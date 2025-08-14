# Firefly API Requirements Documentation

## Overview

This document outlines the specific requirements and quirks discovered while implementing the Terraform provider for the Firefly API.

## API Endpoints

### Authentication
- **Endpoint**: `POST /v2/login`
- **Expected**: JSON object with `accessKey` and `secretKey`
- **Returns**: JSON object with `accessToken`, `expiresAt`, `tokenType`

### Guardrails
- **Base Path**: `/v2/guardrails`
- **List**: `POST /v2/guardrails/search` (with query params for pagination)
- **Create**: `POST /v2/guardrails`
- **Update**: `PATCH /v2/guardrails/{ruleId}`  
- **Delete**: `DELETE /v2/guardrails/{ruleId}`

## API Quirks & Workarounds

### 1. Create Guardrail Response Format
**Issue**: API spec says create returns JSON object, but actually returns just a string (rule ID).

**Spec Says**: 
```json
{
  "ruleId": "507f1f77bcf86cd799439011",
  "notificationId": "507f1f77bcf86cd799439012"
}
```

**Actually Returns**: `"689de93ffbe9fdd796444b49"`

**Workaround**: Provider tries JSON parsing first, falls back to string parsing.

### 2. Required Fields Not Documented as Required

#### CreatedBy Field
- **Issue**: API requires `createdBy` field but spec doesn't clearly indicate this
- **Error**: `"Created By is required"`
- **Workaround**: Provider auto-populates with `"terraform-provider"`

#### Scope Fields
All scope fields are required even when they seem optional:
- `workspaces` (with include/exclude arrays)
- `repositories` (with include/exclude arrays)  
- `branches` (with include/exclude arrays)
- `labels` (with include/exclude arrays)

**Error Messages**:
- `"Repositories is required"`
- `"Branches is required"`  
- `"Labels is required"`

#### Resource Criteria Requirements
For resource-type guardrails, these fields are required:
- `regions` (with include/exclude arrays)
- `asset_types` (with include/exclude arrays)

**Error**: `"Regions is required"`

**Workaround**: Provider auto-populates missing fields with `{"include": ["*"]}` defaults.

## Provider Fixes Applied

### Auto-Population of Required Fields

1. **CreatedBy**: Always set to `"terraform-provider"`
2. **Resource Regions**: Default to `{"include": ["*"]}` if not specified
3. **Resource Asset Types**: Default to `{"include": ["*"]}` if not specified

### Response Parsing

- Create guardrail: Handle both JSON object and string responses
- All other endpoints: Standard JSON parsing

## Recommended User Configuration

### Minimal Guardrail Config
Users must specify ALL scope fields:

```hcl
resource "firefly_guardrail" "example" {
  name        = "Example Guardrail"
  type        = "cost"
  is_enabled  = true
  severity    = 2
  
  scope {
    workspaces {
      include = ["*"]
    }
    repositories {
      include = ["*"] 
    }
    branches {
      include = ["*"]
    }
    labels {
      include = ["*"]
    }
  }
  
  criteria {
    cost {
      threshold_amount = 1000.0
    }
  }
}
```

### Resource Guardrail Config
For resource-type guardrails, regions must be specified:

```hcl
resource "firefly_guardrail" "resource_example" {
  name        = "Resource Guardrail"
  type        = "resource"
  is_enabled  = true
  severity    = 1
  
  scope {
    # ... all scope fields required
  }
  
  criteria {
    resource {
      actions = ["delete"]
      specific_resources = ["aws_instance"]
      
      regions {
        include = ["us-east-1", "eu-west-1"]
      }
      
      asset_types {
        include = ["*"]
      }
    }
  }
}
```

## Testing Results

✅ **Working Features**:
- Authentication with access_key/secret_key
- Create cost guardrails
- Create resource guardrails  
- Create policy guardrails
- Auto-population of required fields
- Terraform state management

⚠️ **Known Issues**:
- API response format differs from spec
- Many fields marked as "optional" are actually required
- Notification ID not returned in create response

## Recommendations for API Team

1. **Fix Response Format**: Make create guardrail return consistent JSON object
2. **Update OpenAPI Spec**: Mark actually required fields as `required: true`
3. **Relax Validation**: Make scope fields truly optional with server-side defaults
4. **Return Notification ID**: Include in create response for better UX

## Provider Status

**Current Status**: ✅ Fully Functional  
**Last Updated**: 2025-08-14  
**Tested With**: Your provided credentials  
**Created Guardrails**: 
- Cost guardrail ID: `689de93ffbe9fdd796444b49`
- Resource guardrail ID: `689de94a221e030ce2455501`