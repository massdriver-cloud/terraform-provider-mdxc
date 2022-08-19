package provider

import (
	"context"
	"terraform-provider-mdxc/internal/mdxc"

	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = ResourceApplicationPermissionType{}
var _ resource.Resource = ResourceApplicationPermission{}
var _ resource.ResourceWithImportState = ResourceApplicationPermission{}

type ResourceApplicationPermissionType struct{}

func (t ResourceApplicationPermissionType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "A cross-cloud application Permission resource (AWS IAM Role, GCP Service Account, Azure Application)",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Cloud specific identifier of the application Permission",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"application_identity_id": {
				Type:        types.StringType,
				Description: "The ID of the Application Permission resource",
				Required:    true,
			},
			"permission": {
				Required:    true,
				Description: "Permission definition to assign application identity",
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"policy_arn": {
						Type:        types.StringType,
						Optional:    true,
						Description: "AWS IAM policy ARN to associate with the application identity",
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.RequiresReplace(),
						},
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("role_name"),
								path.MatchRelative().AtParent().AtName("scope"),
								path.MatchRelative().AtParent().AtName("condition"),
							),
						},
					},
					"role_name": {
						Type:        types.StringType,
						Optional:    true,
						Description: "The Azure or GCP built-in IAM role to bind to the application identity",
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.RequiresReplace(),
						},
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("policy_arn"),
							),
						},
					},
					"scope": {
						Type:        types.StringType,
						Optional:    true,
						Description: "The scope at which the Azure Role Assignment applies to",
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.RequiresReplace(),
						},

						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("policy_arn"),
								path.MatchRelative().AtParent().AtName("condition"),
							),
							schemavalidator.AlsoRequires(
								path.MatchRelative().AtParent().AtName("role_name"),
							),
						},
					},
					"condition": {
						Type:        types.StringType,
						Optional:    true,
						Description: "An GCP IAM Condition for a given role binding",
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.RequiresReplace(),
						},
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("policy_arn"),
								path.MatchRelative().AtParent().AtName("scope"),
							),
							schemavalidator.AlsoRequires(
								path.MatchRelative().AtParent().AtName("role_name"),
							),
						},
					},
				}),
			},
		},
	}, nil
}

func (t ResourceApplicationPermissionType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	return ResourceApplicationPermission{
		provider: *(in.(*MDXCProvider)),
	}, diag.Diagnostics{}
}

type ResourceApplicationPermission struct {
	provider MDXCProvider
}

func (r ResourceApplicationPermission) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdxc.ApplicationPermissionData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.provider.Client.CreateApplicationPermission(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created application Permission")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceApplicationPermission) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdxc.ApplicationPermissionData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.provider.Client.ReadApplicationPermission(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceApplicationPermission) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdxc.ApplicationPermissionData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.provider.Client.UpdateApplicationPermission(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceApplicationPermission) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdxc.ApplicationPermissionData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.provider.Client.DeleteApplicationPermission(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceApplicationPermission) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
