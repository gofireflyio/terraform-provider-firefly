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
	_ datasource.DataSource              = &variableSetDataSource{}
	_ datasource.DataSourceWithConfigure = &variableSetDataSource{}
)

func NewVariableSetDataSource() datasource.DataSource {
	return &variableSetDataSource{}
}

type variableSetDataSource struct {
	client *client.Client
}

type VariableSetSingleDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Labels      types.List   `tfsdk:"labels"`
	Parents     types.List   `tfsdk:"parents"`
	Version     types.Int64  `tfsdk:"version"`
}

func (d *variableSetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows_variable_set"
}

func (d *variableSetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Firefly variable set by ID",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the variable set",
				Required:    true,
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
	}
}

func (d *variableSetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *variableSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VariableSetSingleDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	variableSetID := data.ID.ValueString()
	tflog.Debug(ctx, "Reading variable set", map[string]interface{}{"id": variableSetID})

	variableSet, err := d.client.VariableSets.GetVariableSet(variableSetID)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Variable Set", fmt.Sprintf("Could not read variable set ID %s: %s", variableSetID, err))
		return
	}

	// Convert labels
	var labelsList types.List
	if len(variableSet.Labels) > 0 {
		labelValues := make([]types.String, len(variableSet.Labels))
		for i, label := range variableSet.Labels {
			labelValues[i] = types.StringValue(label)
		}
		labelsList = types.ListValueMust(types.StringType, labelListToValues(labelValues))
	} else {
		labelsList = types.ListValueMust(types.StringType, []attr.Value{})
	}

	// Convert parents
	var parentsList types.List
	if len(variableSet.Parents) > 0 {
		parentValues := make([]types.String, len(variableSet.Parents))
		for i, parent := range variableSet.Parents {
			parentValues[i] = types.StringValue(parent)
		}
		parentsList = types.ListValueMust(types.StringType, labelListToValues(parentValues))
	} else {
		parentsList = types.ListValueMust(types.StringType, []attr.Value{})
	}

	data.Name = types.StringValue(variableSet.Name)
	data.Description = types.StringValue(variableSet.Description)
	data.Labels = labelsList
	data.Parents = parentsList
	data.Version = types.Int64Value(int64(variableSet.Version))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}