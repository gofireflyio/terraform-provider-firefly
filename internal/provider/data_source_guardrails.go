package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-firefly/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &guardrailsDataSource{}
	_ datasource.DataSourceWithConfigure = &guardrailsDataSource{}
)

// NewGuardrailsDataSource is a helper function to simplify the provider implementation
func NewGuardrailsDataSource() datasource.DataSource {
	return &guardrailsDataSource{}
}

// guardrailsDataSource is the data source implementation
type guardrailsDataSource struct {
	client *client.Client
}

// GuardrailFiltersModel describes the guardrail filters
type GuardrailFiltersModel struct {
	CreatedBy    types.List   `tfsdk:"created_by"`
	Type         types.List   `tfsdk:"type"`
	Labels       types.List   `tfsdk:"labels"`
	Repositories types.List   `tfsdk:"repositories"`
	Workspaces   types.List   `tfsdk:"workspaces"`
	Branches     types.List   `tfsdk:"branches"`
}

// GuardrailDataModel describes a single guardrail rule
type GuardrailDataModel struct {
	ID            types.String `tfsdk:"id"`
	AccountID     types.String `tfsdk:"account_id"`
	CreatedBy     types.String `tfsdk:"created_by"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	IsEnabled     types.Bool   `tfsdk:"is_enabled"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	NotificationID types.String `tfsdk:"notification_id"`
	Severity      types.Int64  `tfsdk:"severity"`
}

// GuardrailsDataSourceModel describes the data source data model
type GuardrailsDataSourceModel struct {
	Guardrails   types.List          `tfsdk:"guardrails"`
	Filters      *GuardrailFiltersModel `tfsdk:"filters"`
	SearchValue  types.String        `tfsdk:"search_value"`
}

// Metadata returns the data source type name
func (d *guardrailsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guardrails"
}

// Schema defines the schema for the data source
func (d *guardrailsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Firefly guardrail rules",
		Attributes: map[string]schema.Attribute{
			"guardrails": schema.ListNestedAttribute{
				Description: "List of guardrail rules",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the guardrail rule",
							Computed:    true,
						},
						"account_id": schema.StringAttribute{
							Description: "Account ID associated with the guardrail rule",
							Computed:    true,
						},
						"created_by": schema.StringAttribute{
							Description: "ID of the user who created the guardrail rule",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the guardrail rule",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the guardrail rule",
							Computed:    true,
						},
						"is_enabled": schema.BoolAttribute{
							Description: "Whether the guardrail rule is enabled",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Timestamp when the guardrail rule was created",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Timestamp when the guardrail rule was last updated",
							Computed:    true,
						},
						"notification_id": schema.StringAttribute{
							Description: "ID of the associated notification",
							Computed:    true,
						},
						"severity": schema.Int64Attribute{
							Description: "Severity level of the guardrail rule",
							Computed:    true,
						},
					},
				},
			},
			"search_value": schema.StringAttribute{
				Description: "Search value to filter guardrails",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"filters": schema.SingleNestedBlock{
				Description: "Filters for guardrail rules",
				Attributes: map[string]schema.Attribute{
					"created_by": schema.ListAttribute{
						Description: "Filter by creator",
						Optional:    true,
						ElementType: types.StringType,
					},
					"type": schema.ListAttribute{
						Description: "Filter by type",
						Optional:    true,
						ElementType: types.StringType,
					},
					"labels": schema.ListAttribute{
						Description: "Filter by labels",
						Optional:    true,
						ElementType: types.StringType,
					},
					"repositories": schema.ListAttribute{
						Description: "Filter by repositories",
						Optional:    true,
						ElementType: types.StringType,
					},
					"workspaces": schema.ListAttribute{
						Description: "Filter by workspaces",
						Optional:    true,
						ElementType: types.StringType,
					},
					"branches": schema.ListAttribute{
						Description: "Filter by branches",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *guardrailsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data
func (d *guardrailsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GuardrailsDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request with filters
	request := &client.ListGuardrailsRequest{}
	
	// Add search value if provided
	if !data.SearchValue.IsNull() {
		request.SearchValue = data.SearchValue.ValueString()
	}
	
	// Add filters if provided
	if data.Filters != nil {
		filters := &client.GuardrailFilters{}
		
		// Add created by filter
		if !data.Filters.CreatedBy.IsNull() {
			var createdBy []string
			diags = data.Filters.CreatedBy.ElementsAs(ctx, &createdBy, false)
			resp.Diagnostics.Append(diags...)
			if createdBy != nil {
				filters.CreatedBy = createdBy
			}
		}
		
		// Add type filter
		if !data.Filters.Type.IsNull() {
			var types []string
			diags = data.Filters.Type.ElementsAs(ctx, &types, false)
			resp.Diagnostics.Append(diags...)
			if types != nil {
				filters.Type = types
			}
		}
		
		// Add labels filter
		if !data.Filters.Labels.IsNull() {
			var labels []string
			diags = data.Filters.Labels.ElementsAs(ctx, &labels, false)
			resp.Diagnostics.Append(diags...)
			if labels != nil {
				filters.Labels = labels
			}
		}
		
		// Add repositories filter
		if !data.Filters.Repositories.IsNull() {
			var repositories []string
			diags = data.Filters.Repositories.ElementsAs(ctx, &repositories, false)
			resp.Diagnostics.Append(diags...)
			if repositories != nil {
				filters.Repositories = repositories
			}
		}
		
		// Add workspaces filter
		if !data.Filters.Workspaces.IsNull() {
			var workspaces []string
			diags = data.Filters.Workspaces.ElementsAs(ctx, &workspaces, false)
			resp.Diagnostics.Append(diags...)
			if workspaces != nil {
				filters.Workspaces = workspaces
			}
		}
		
		// Add branches filter
		if !data.Filters.Branches.IsNull() {
			var branches []string
			diags = data.Filters.Branches.ElementsAs(ctx, &branches, false)
			resp.Diagnostics.Append(diags...)
			if branches != nil {
				filters.Branches = branches
			}
		}
		
		request.Filters = filters
	}
	
	// Get guardrails from API
	guardrails, err := d.client.Guardrails.ListGuardrails(request, 0, 100)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Guardrails",
			fmt.Sprintf("Could not read guardrails: %s", err),
		)
		return
	}
	
	// Map response to model
	var guardrailModels []GuardrailDataModel
	for _, guardrail := range guardrails {
		guardrailModel := GuardrailDataModel{
			ID:             types.StringValue(guardrail.ID),
			AccountID:      types.StringValue(guardrail.AccountID),
			CreatedBy:      types.StringValue(guardrail.CreatedBy),
			Name:           types.StringValue(guardrail.Name),
			Type:           types.StringValue(guardrail.Type),
			IsEnabled:      types.BoolValue(guardrail.IsEnabled),
			CreatedAt:      types.StringValue(guardrail.CreatedAt),
			UpdatedAt:      types.StringValue(guardrail.UpdatedAt),
			NotificationID: types.StringValue(guardrail.NotificationID),
			Severity:       types.Int64Value(int64(guardrail.Severity)),
		}
		
		guardrailModels = append(guardrailModels, guardrailModel)
	}
	
	// Set guardrails in the data model
	guardrailsList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":              types.StringType,
			"account_id":      types.StringType,
			"created_by":      types.StringType,
			"name":            types.StringType,
			"type":            types.StringType,
			"is_enabled":      types.BoolType,
			"created_at":      types.StringType,
			"updated_at":      types.StringType,
			"notification_id": types.StringType,
			"severity":        types.Int64Type,
		},
	}, guardrailModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	data.Guardrails = guardrailsList
	
	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
