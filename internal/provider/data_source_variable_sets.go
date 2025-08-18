package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
)

var (
	_ datasource.DataSource              = &variableSetsDataSource{}
	_ datasource.DataSourceWithConfigure = &variableSetsDataSource{}
)

func NewVariableSetsDataSource() datasource.DataSource {
	return &variableSetsDataSource{}
}

type variableSetsDataSource struct {
	client *client.Client
}

type VariableSetsDataSourceModel struct {
	SearchQuery  types.String               `tfsdk:"search_query"`
	VariableSets []VariableSetDataModel     `tfsdk:"variable_sets"`
}

type VariableSetDataModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Labels      types.List   `tfsdk:"labels"`
	Parents     types.List   `tfsdk:"parents"`
	Version     types.Int64  `tfsdk:"version"`
}

func (d *variableSetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable_sets"
}

func (d *variableSetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of Firefly variable sets",
		Attributes: map[string]schema.Attribute{
			"search_query": schema.StringAttribute{
				Description: "Optional search query to filter variable sets",
				Optional:    true,
			},
			"variable_sets": schema.ListNestedAttribute{
				Description: "List of variable sets",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the variable set",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the variable set",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the variable set",
							Computed:    true,
						},
						"labels": schema.ListAttribute{
							Description: "Labels assigned to the variable set",
							Computed:    true,
							ElementType: types.StringType,
						},
						"parents": schema.ListAttribute{
							Description: "Parent variable set IDs",
							Computed:    true,
							ElementType: types.StringType,
						},
						"version": schema.Int64Attribute{
							Description: "Version number of the variable set",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *variableSetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *variableSetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VariableSetsDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	searchQuery := ""
	if !data.SearchQuery.IsNull() {
		searchQuery = data.SearchQuery.ValueString()
	}

	tflog.Debug(ctx, "Reading variable sets", map[string]interface{}{
		"search_query": searchQuery,
	})

	variableSets, err := d.client.VariableSets.ListVariableSets(100, 0, searchQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Variable Sets",
			fmt.Sprintf("Could not read variable sets: %s", err),
		)
		return
	}

	// Map response to model
	sets := make([]VariableSetDataModel, len(variableSets))
	for i, vs := range variableSets {
		// Convert labels
		var labelsList types.List
		if len(vs.Labels) > 0 {
			labelValues := make([]types.String, len(vs.Labels))
			for j, label := range vs.Labels {
				labelValues[j] = types.StringValue(label)
			}
			labelsList = types.ListValueMust(types.StringType, labelListToValues(labelValues))
		} else {
			labelsList = types.ListValueMust(types.StringType, []attr.Value{})
		}

		// Convert parents
		var parentsList types.List
		if len(vs.Parents) > 0 {
			parentValues := make([]types.String, len(vs.Parents))
			for j, parent := range vs.Parents {
				parentValues[j] = types.StringValue(parent)
			}
			parentsList = types.ListValueMust(types.StringType, labelListToValues(parentValues))
		} else {
			parentsList = types.ListValueMust(types.StringType, []attr.Value{})
		}

		sets[i] = VariableSetDataModel{
			ID:          types.StringValue(vs.ID),
			Name:        types.StringValue(vs.Name),
			Description: types.StringValue(vs.Description),
			Labels:      labelsList,
			Parents:     parentsList,
			Version:     types.Int64Value(int64(vs.Version)),
		}
	}

	data.VariableSets = sets

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}