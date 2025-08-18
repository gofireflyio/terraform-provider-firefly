package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-firefly/internal/client"
)

var (
	_ resource.Resource                = &variableSetResource{}
	_ resource.ResourceWithConfigure   = &variableSetResource{}
	_ resource.ResourceWithImportState = &variableSetResource{}
)

func NewVariableSetResource() resource.Resource {
	return &variableSetResource{}
}

type variableSetResource struct {
	client *client.Client
}

type VariableSetResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Labels      types.List   `tfsdk:"labels"`
	Parents     types.List   `tfsdk:"parents"`
	Variables   types.List   `tfsdk:"variables"`
	Version     types.Int64  `tfsdk:"version"`
}

func (r *variableSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable_set"
}

func (r *variableSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Firefly variable set",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the variable set",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the variable set",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the variable set",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"labels": schema.ListAttribute{
				Description: "Labels to assign to the variable set",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"parents": schema.ListAttribute{
				Description: "Parent variable set IDs",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				Description: "Version number of the variable set",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"variables": schema.ListNestedBlock{
				Description: "Variables in the variable set",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "The variable key",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The variable value",
							Required:    true,
							Sensitive:   true,
						},
						"sensitivity": schema.StringAttribute{
							Description: "The sensitivity of the variable (string or secret)",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("string"),
						},
						"destination": schema.StringAttribute{
							Description: "The destination of the variable (env or iac)",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("env"),
						},
					},
				},
			},
		},
	}
}

func (r *variableSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *variableSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VariableSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert labels and parents
	var labels, parents []string
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		diags = plan.Labels.ElementsAs(ctx, &labels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.Parents.IsNull() && !plan.Parents.IsUnknown() {
		diags = plan.Parents.ElementsAs(ctx, &parents, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert variables
	var variables []client.Variable
	if !plan.Variables.IsNull() && !plan.Variables.IsUnknown() {
		var varModels []ProjectVariableModel
		diags = plan.Variables.ElementsAs(ctx, &varModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, v := range varModels {
			variables = append(variables, client.Variable{
				Key:         v.Key.ValueString(),
				Value:       v.Value.ValueString(),
				Sensitivity: v.Sensitivity.ValueString(),
				Destination: v.Destination.ValueString(),
			})
		}
	}

	createReq := client.CreateVariableSetRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Labels:      labels,
		Parents:     parents,
		Variables:   variables,
	}

	tflog.Debug(ctx, "Creating variable set", map[string]interface{}{"name": createReq.Name})

	createResp, err := r.client.VariableSets.CreateVariableSet(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Variable Set", fmt.Sprintf("Could not create variable set: %s", err))
		return
	}

	plan.ID = types.StringValue(createResp.VariableSetID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *variableSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VariableSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	variableSet, err := r.client.VariableSets.GetVariableSet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Variable Set", fmt.Sprintf("Could not read variable set ID %s: %s", state.ID.ValueString(), err))
		return
	}

	// Update state with API response
	state.Name = types.StringValue(variableSet.Name)
	state.Description = types.StringValue(variableSet.Description)
	state.Version = types.Int64Value(int64(variableSet.Version))

	// Convert labels and parents
	if len(variableSet.Labels) > 0 {
		labelList := make([]types.String, len(variableSet.Labels))
		for i, label := range variableSet.Labels {
			labelList[i] = types.StringValue(label)
		}
		state.Labels = types.ListValueMust(types.StringType, labelListToValues(labelList))
	}

	if len(variableSet.Parents) > 0 {
		parentList := make([]types.String, len(variableSet.Parents))
		for i, parent := range variableSet.Parents {
			parentList[i] = types.StringValue(parent)
		}
		state.Parents = types.ListValueMust(types.StringType, labelListToValues(parentList))
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *variableSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VariableSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert labels and parents  
	var labels, parents []string
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		diags = plan.Labels.ElementsAs(ctx, &labels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.Parents.IsNull() && !plan.Parents.IsUnknown() {
		diags = plan.Parents.ElementsAs(ctx, &parents, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert variables
	var variables []client.Variable
	if !plan.Variables.IsNull() && !plan.Variables.IsUnknown() {
		var varModels []ProjectVariableModel
		diags = plan.Variables.ElementsAs(ctx, &varModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, v := range varModels {
			variables = append(variables, client.Variable{
				Key:         v.Key.ValueString(),
				Value:       v.Value.ValueString(),
				Sensitivity: v.Sensitivity.ValueString(),
				Destination: v.Destination.ValueString(),
			})
		}
	}

	updateReq := client.UpdateVariableSetRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Labels:      labels,
		Parents:     parents,
		Variables:   variables,
	}

	tflog.Debug(ctx, "Updating variable set", map[string]interface{}{"id": plan.ID.ValueString()})

	variableSet, err := r.client.VariableSets.UpdateVariableSet(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Variable Set", fmt.Sprintf("Could not update variable set ID %s: %s", plan.ID.ValueString(), err))
		return
	}

	plan.Version = types.Int64Value(int64(variableSet.Version))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *variableSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VariableSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting variable set", map[string]interface{}{"id": state.ID.ValueString()})

	err := r.client.VariableSets.DeleteVariableSet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Variable Set", fmt.Sprintf("Could not delete variable set ID %s: %s", state.ID.ValueString(), err))
		return
	}
}

func (r *variableSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}