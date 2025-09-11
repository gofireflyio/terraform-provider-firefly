package provider

import (
	"context"
	"fmt"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &GovernanceInsightsDataSource{}

// NewGovernanceInsightsDataSource creates a new governance insights data source
func NewGovernanceInsightsDataSource() datasource.DataSource {
	return &GovernanceInsightsDataSource{}
}

// GovernanceInsightsDataSource defines the data source implementation
type GovernanceInsightsDataSource struct {
	client *client.Client
}

func (d *GovernanceInsightsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_governance_insights"
}

func (d *GovernanceInsightsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for retrieving Firefly governance insights",
		
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The data source identifier",
				Computed:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "Search query string for filtering insights",
				Optional:            true,
			},
			"labels": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of labels to filter insights",
				Optional:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Category filter for insights",
				Optional:            true,
			},
			"insights": schema.ListNestedAttribute{
				MarkdownDescription: "List of governance insights",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the governance insight",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the governance insight",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the governance insight",
							Computed:            true,
						},
						"code": schema.StringAttribute{
							MarkdownDescription: "The Rego code for the insight rule",
							Computed:            true,
						},
						"type": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of resource types this insight applies to",
							Computed:            true,
						},
						"provider_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of provider IDs this insight applies to",
							Computed:            true,
						},
						"labels": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of labels for categorizing the insight",
							Computed:            true,
						},
						"severity": schema.StringAttribute{
							MarkdownDescription: "The severity level of the insight",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							MarkdownDescription: "The category of the insight",
							Computed:            true,
						},
						"frameworks": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of compliance frameworks this insight relates to",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *GovernanceInsightsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider is not configured
	if req.ProviderData == nil {
		return
	}
	
	client, ok := req.ProviderData.(*client.Client)
	
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		
		return
	}
	
	d.client = client
}

func (d *GovernanceInsightsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GovernanceInsightsDataSourceModel
	
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	// Build request
	listReq := &client.GovernanceInsightListRequest{
		Page:     1,
		PageSize: 100, // Get first 100 insights
		OnlyAvailableProviders: true,
	}
	
	// Add query filter if provided
	if !data.Query.IsNull() && !data.Query.IsUnknown() {
		listReq.Query = data.Query.ValueString()
	}
	
	// Add labels filter if provided
	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		var labels []string
		diags := data.Labels.ElementsAs(ctx, &labels, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		listReq.Labels = labels
	}
	
	// Add category filter if provided
	if !data.Category.IsNull() && !data.Category.IsUnknown() {
		listReq.Category = data.Category.ValueString()
	}
	
	tflog.Debug(ctx, "Reading governance insights", map[string]interface{}{
		"query":    listReq.Query,
		"labels":   listReq.Labels,
		"category": listReq.Category,
	})
	
	// Get insights
	insightsResp, err := d.client.GovernanceInsights.List(listReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading governance insights",
			fmt.Sprintf("Could not read governance insights: %s", err),
		)
		return
	}
	
	// Map insights to data source model
	data.Insights = make([]GovernanceInsightDataSourceModel, len(insightsResp.Data))
	for i, insight := range insightsResp.Data {
		insightModel := &GovernanceInsightDataSourceModel{}
		
		// Map each insight
		insightModel.ID = types.StringValue(insight.ID)
		insightModel.Name = types.StringValue(insight.Name)
		
		if insight.Description != "" {
			insightModel.Description = types.StringValue(insight.Description)
		} else {
			insightModel.Description = types.StringNull()
		}
		
		insightModel.Code = types.StringValue(insight.Code)
		
		// Convert arrays to lists
		typeList, diags := types.ListValueFrom(ctx, types.StringType, insight.Type)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		insightModel.Type = typeList
		
		providerList, diags := types.ListValueFrom(ctx, types.StringType, insight.ProviderIDs)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		insightModel.ProviderIDs = providerList
		
		if len(insight.Labels) > 0 {
			labelsList, diags := types.ListValueFrom(ctx, types.StringType, insight.Labels)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			insightModel.Labels = labelsList
		} else {
			insightModel.Labels = types.ListNull(types.StringType)
		}
		
		if insight.Severity > 0 {
			insightModel.Severity = types.StringValue(client.SeverityToString(insight.Severity))
		} else {
			insightModel.Severity = types.StringNull()
		}
		
		if insight.Category != "" {
			insightModel.Category = types.StringValue(insight.Category)
		} else {
			insightModel.Category = types.StringNull()
		}
		
		if len(insight.Frameworks) > 0 {
			frameworksList, diags := types.ListValueFrom(ctx, types.StringType, insight.Frameworks)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			insightModel.Frameworks = frameworksList
		} else {
			insightModel.Frameworks = types.ListNull(types.StringType)
		}
		
		data.Insights[i] = *insightModel
	}
	
	// Generate unique ID for the data source
	data.ID = types.StringValue("governance-insights")
	
	tflog.Debug(ctx, "Read governance insights", map[string]interface{}{
		"count": len(data.Insights),
	})
	
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}