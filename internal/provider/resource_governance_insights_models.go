package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GovernanceInsightResourceModel represents the resource model for a governance insight
type GovernanceInsightResourceModel struct {
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

// GovernanceInsightDataSourceModel represents the data source model for a governance insight
type GovernanceInsightDataSourceModel struct {
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

// GovernanceInsightsDataSourceModel represents the data source model for listing governance insights
type GovernanceInsightsDataSourceModel struct {
	ID       types.String                        `tfsdk:"id"`
	Query    types.String                        `tfsdk:"query"`
	Labels   types.List                          `tfsdk:"labels"`
	Category types.String                        `tfsdk:"category"`
	Insights []GovernanceInsightDataSourceModel `tfsdk:"insights"`
}