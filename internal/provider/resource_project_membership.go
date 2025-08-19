package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofireflyio/terraform-provider-firefly/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProjectMembershipResource{}
var _ resource.ResourceWithImportState = &ProjectMembershipResource{}

func NewProjectMembershipResource() resource.Resource {
	return &ProjectMembershipResource{}
}

// ProjectMembershipResource defines the resource implementation.
type ProjectMembershipResource struct {
	client *client.Client
}

// ProjectMembershipResourceModel describes the resource data model.
type ProjectMembershipResourceModel struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	UserID    types.String `tfsdk:"user_id"`
	Email     types.String `tfsdk:"email"`
	Role      types.String `tfsdk:"role"`
}

func (r *ProjectMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_membership"
}

func (r *ProjectMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Manages membership of a user in a Firefly project. This resource allows you to add users to projects with specific roles.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the membership (format: project_id:user_id)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the user to add to the project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The email address of the user",
			},
			"role": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The role of the user in the project (e.g., 'admin', 'member', 'viewer')",
			},
		},
	}
}

func (r *ProjectMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider is not configured.
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

func (r *ProjectMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the member object
	member := client.Member{
		UserID: data.UserID.ValueString(),
		Email:  data.Email.ValueString(),
		Role:   data.Role.ValueString(),
	}

	// Add member to project
	tflog.Debug(ctx, "Adding member to project", map[string]interface{}{
		"project_id": data.ProjectID.ValueString(),
		"user_id":    member.UserID,
		"email":      member.Email,
		"role":       member.Role,
	})

	addedMember, err := r.client.Projects.AddProjectMember(data.ProjectID.ValueString(), member)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add member to project, got error: %s", err))
		return
	}

	// Generate composite ID
	data.ID = types.StringValue(fmt.Sprintf("%s:%s", data.ProjectID.ValueString(), data.UserID.ValueString()))
	
	// Update data from response
	data.Email = types.StringValue(addedMember.Email)
	data.Role = types.StringValue(addedMember.Role)

	tflog.Trace(ctx, "Created project membership resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the member from the project
	member, err := r.client.Projects.GetProjectMember(data.ProjectID.ValueString(), data.UserID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// Member has been removed outside of Terraform
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project member, got error: %s", err))
		return
	}

	// Update the model with the latest data
	data.Email = types.StringValue(member.Email)
	data.Role = types.StringValue(member.Role)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create updated member object
	member := client.Member{
		UserID: data.UserID.ValueString(),
		Email:  data.Email.ValueString(),
		Role:   data.Role.ValueString(),
	}

	// Update the member (remove and re-add with new role)
	tflog.Debug(ctx, "Updating project member", map[string]interface{}{
		"project_id": data.ProjectID.ValueString(),
		"user_id":    member.UserID,
		"new_role":   member.Role,
	})

	updatedMember, err := r.client.Projects.UpdateProjectMember(data.ProjectID.ValueString(), member)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update project member, got error: %s", err))
		return
	}

	// Update data from response
	data.Email = types.StringValue(updatedMember.Email)
	data.Role = types.StringValue(updatedMember.Role)

	tflog.Trace(ctx, "Updated project membership resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove member from project
	tflog.Debug(ctx, "Removing member from project", map[string]interface{}{
		"project_id": data.ProjectID.ValueString(),
		"user_id":    data.UserID.ValueString(),
	})

	err := r.client.Projects.RemoveProjectMember(data.ProjectID.ValueString(), data.UserID.ValueString())
	if err != nil {
		// If member is already gone, don't error
		if !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove member from project, got error: %s", err))
			return
		}
	}

	tflog.Trace(ctx, "Deleted project membership resource")
}

func (r *ProjectMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: project_id:user_id
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: project_id:user_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), parts[1])...)
}