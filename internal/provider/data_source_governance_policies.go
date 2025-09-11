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
var _ datasource.DataSource = &GovernancePoliciesDataSource{}

// NewGovernancePoliciesDataSource creates a new governance policies data source
func NewGovernancePoliciesDataSource() datasource.DataSource {
	return &GovernancePoliciesDataSource{}
}

// GovernancePoliciesDataSource defines the data source implementation
type GovernancePoliciesDataSource struct {
	client *client.Client
}

func (d *GovernancePoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_governance_policies"
}

func (d *GovernancePoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for retrieving Firefly governance policies",
		
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The data source identifier",
				Computed:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "Search query string for filtering policies",
				Optional:            true,
			},
			"labels": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of labels to filter policies",
				Optional:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Category filter for policies",
				Optional:            true,
			},
			"policies": schema.ListNestedAttribute{
				MarkdownDescription: "List of governance policies",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the governance policy",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the governance policy",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the governance policy",
							Computed:            true,
						},
						"code": schema.StringAttribute{
							MarkdownDescription: "The Rego code for the policy rule",
							Computed:            true,
						},
						"type": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of resource types this policy applies to",
							Computed:            true,
						},
						"provider_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of provider IDs this policy applies to",
							Computed:            true,
						},
						"labels": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of labels for categorizing the policy",
							Computed:            true,
						},
						"severity": schema.StringAttribute{
							MarkdownDescription: "The severity level of the policy",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							MarkdownDescription: "The category of the policy",
							Computed:            true,
						},
						"frameworks": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of compliance frameworks this policy relates to",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *GovernancePoliciesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GovernancePoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GovernancePoliciesDataSourceModel
	
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	// Build request
	listReq := &client.GovernancePolicyListRequest{
		Page:     1,
		PageSize: 100, // Get first 100 policies
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
	
	tflog.Debug(ctx, "Reading governance policies", map[string]interface{}{
		"query":    listReq.Query,
		"labels":   listReq.Labels,
		"category": listReq.Category,
	})
	
	// Get policies
	policiesResp, err := d.client.GovernancePolicies.List(listReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading governance policies",
			fmt.Sprintf("Could not read governance policies: %s", err),
		)
		return
	}
	
	// Map policies to data source model
	data.Policies = make([]GovernancePolicyDataSourceModel, len(policiesResp.Data))
	for i, policy := range policiesResp.Data {
		policyModel := &GovernancePolicyDataSourceModel{}
		
		// Map each policy
		policyModel.ID = types.StringValue(policy.ID)
		policyModel.Name = types.StringValue(policy.Name)
		
		if policy.Description != "" {
			policyModel.Description = types.StringValue(policy.Description)
		} else {
			policyModel.Description = types.StringNull()
		}
		
		policyModel.Code = types.StringValue(policy.Code)
		
		// Convert arrays to lists
		typeList, diags := types.ListValueFrom(ctx, types.StringType, policy.Type)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		policyModel.Type = typeList
		
		providerList, diags := types.ListValueFrom(ctx, types.StringType, policy.ProviderIDs)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		policyModel.ProviderIDs = providerList
		
		if len(policy.Labels) > 0 {
			labelsList, diags := types.ListValueFrom(ctx, types.StringType, policy.Labels)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			policyModel.Labels = labelsList
		} else {
			policyModel.Labels = types.ListNull(types.StringType)
		}
		
		if policy.Severity > 0 {
			policyModel.Severity = types.StringValue(client.SeverityToString(policy.Severity))
		} else {
			policyModel.Severity = types.StringNull()
		}
		
		if policy.Category != "" {
			policyModel.Category = types.StringValue(policy.Category)
		} else {
			policyModel.Category = types.StringNull()
		}
		
		if len(policy.Frameworks) > 0 {
			frameworksList, diags := types.ListValueFrom(ctx, types.StringType, policy.Frameworks)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			policyModel.Frameworks = frameworksList
		} else {
			policyModel.Frameworks = types.ListNull(types.StringType)
		}
		
		data.Policies[i] = *policyModel
	}
	
	// Generate unique ID for the data source
	data.ID = types.StringValue("governance-policies")
	
	tflog.Debug(ctx, "Read governance policies", map[string]interface{}{
		"count": len(data.Policies),
	})
	
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}