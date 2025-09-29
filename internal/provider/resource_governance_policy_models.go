package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GovernancePolicyResourceModel represents the resource model for a governance policy
type GovernancePolicyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Code        types.String `tfsdk:"code"`
	Type        types.List   `tfsdk:"type"`
	ProviderIDs types.List   `tfsdk:"provider_ids"`
	Labels      types.List   `tfsdk:"labels"`
	Severity    types.String `tfsdk:"severity"`
	Category    types.String `tfsdk:"category"`
	Frameworks  types.List   `tfsdk:"frameworks"`
}

// GovernancePolicyDataSourceModel represents the data source model for a governance policy
type GovernancePolicyDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Code        types.String `tfsdk:"code"`
	Type        types.List   `tfsdk:"type"`
	ProviderIDs types.List   `tfsdk:"provider_ids"`
	Labels      types.List   `tfsdk:"labels"`
	Severity    types.String `tfsdk:"severity"`
	Category    types.String `tfsdk:"category"`
	Frameworks  types.List   `tfsdk:"frameworks"`
}

// GovernancePoliciesDataSourceModel represents the data source model for listing governance policies
type GovernancePoliciesDataSourceModel struct {
	ID       types.String                       `tfsdk:"id"`
	Query    types.String                       `tfsdk:"query"`
	Labels   types.List                         `tfsdk:"labels"`
	Category types.String                       `tfsdk:"category"`
	Policies []GovernancePolicyDataSourceModel `tfsdk:"policies"`
}