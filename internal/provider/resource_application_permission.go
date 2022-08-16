package provider

import (
	"context"
	"terraform-provider-mdxc/internal/mdxc"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = ApplicationPermissionType{}
var _ resource.Resource = ApplicationPermission{}
var _ resource.ResourceWithImportState = ApplicationPermission{}

type ApplicationPermissionType struct{}

var awsApplicationPermissionInputs = tfsdk.Attribute{
	Optional:    true,
	Description: "AWS IAM Role Configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"policy_arn": {
			Type:        types.StringType,
			Required:    true,
			Description: "AWS IAM policy ARN to associate with the application identity",
			PlanModifiers: tfsdk.AttributePlanModifiers{
				resource.RequiresReplace(),
			},
		},
	}),
}

var azureApplicationPermissionInputs = tfsdk.Attribute{
	Optional:    true,
	Description: "Azure IAM Role Configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"role_name": {
			Type:        types.StringType,
			Required:    true,
			Description: "The Azure built-in IAM role to bind to the application identity",
			PlanModifiers: tfsdk.AttributePlanModifiers{
				resource.RequiresReplace(),
			},
		},
		"scope": {
			Type:        types.StringType,
			Required:    true,
			Description: "The scope at which the Role Assignment applies to",
			PlanModifiers: tfsdk.AttributePlanModifiers{
				resource.RequiresReplace(),
			},
		},
	}),
}

var gcpApplicationPermissionInputs = tfsdk.Attribute{
	Optional:    true,
	Description: "Azure IAM Role Configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"role": {
			Type:        types.StringType,
			Required:    true,
			Description: "The GCP role to bind to the application identity",
			PlanModifiers: tfsdk.AttributePlanModifiers{
				resource.RequiresReplace(),
			},
		},
		"condition": {
			Type:        types.StringType,
			Required:    true,
			Description: "An IAM Condition for a given role binding",
			PlanModifiers: tfsdk.AttributePlanModifiers{
				resource.RequiresReplace(),
			},
		},
	}),
}

func (t ApplicationPermissionType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"aws_configuration":   awsApplicationPermissionInputs,
			"azure_configuration": azureApplicationPermissionInputs,
			"gcp_configuration":   gcpApplicationPermissionInputs,
		},
	}, nil
}

func (t ApplicationPermissionType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	return ApplicationPermission{
		provider: *(in.(*MDXCProvider)),
	}, diag.Diagnostics{}
}

type ApplicationPermission struct {
	provider MDXCProvider
}

func (r ApplicationPermission) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

func (r ApplicationPermission) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

func (r ApplicationPermission) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (r ApplicationPermission) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func (r ApplicationPermission) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
