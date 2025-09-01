# firefly_workflows_guardrails (Data Source)

Fetches a list of Firefly guardrails with optional filtering.

## Example Usage

```terraform
# Get all guardrails
data "firefly_workflows_guardrails" "all" {}

# Search for cost-related guardrails
data "firefly_workflows_guardrails" "cost_controls" {
  search_field = "type"
  search_value = "cost"
}

# Output guardrail information
output "cost_guardrails" {
  value = [for gr in data.firefly_workflows_guardrails.cost_controls.guardrails : {
    name     = gr.name
    severity = gr.severity
    enabled  = gr.is_enabled
  }]
}
```

## Schema

### Optional

- `search_field` (String) - Field to search in. Valid values: `name`, `type`
- `search_value` (String) - Value to search for

### Read-Only

- `guardrails` (List of Object) - List of guardrails (see [below for nested schema](#nestedatt--guardrails))

<a id="nestedatt--guardrails"></a>
### Nested Schema for `guardrails`

Read-Only:

- `id` (String) - The unique identifier of the guardrail
- `name` (String) - The name of the guardrail
- `type` (String) - Type of guardrail
- `is_enabled` (Boolean) - Whether the guardrail is enabled
- `severity` (String) - Severity level of the guardrail (Flexible, Strict, Warning)