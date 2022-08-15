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
var _ provider.ResourceType = ApplicationIdentityType{}
var _ resource.Resource = ApplicationIdentity{}
var _ resource.ResourceWithImportState = ApplicationIdentity{}

type ApplicationIdentityType struct{}

var awsApplicationIdentitySchema = tfsdk.Attribute{
	Optional:    true,
	Description: "AWS IAM Role Configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"assume_role_policy": {
			Type:     types.StringType,
			Required: true,
			// ValidateFunc:     validation.StringIsJSON,
			// DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
			// StateFunc: func(v interface{}) string {
			// 	json, _ := structure.NormalizeJsonString(v)
			// 	return json
			// },
		},
	}),
}

func (t ApplicationIdentityType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "A cross-cloud application identity resource (AWS IAM Role, GCP Service Account, Azure Application)",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Cloud specific identifier of the application identity",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"name": {
				Type:        types.StringType,
				Description: "The name of the IAM entity in the respective cloud (AWS IAM Role, GCP Service Account, Azure Application)",
				Required:    true,
			},
			"aws": awsApplicationIdentitySchema,
		},
	}, nil
}

func (t ApplicationIdentityType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	return ApplicationIdentity{
		provider: *(in.(*MDXCProvider)),
	}, diag.Diagnostics{}
}

type ApplicationIdentity struct {
	provider MDXCProvider
}

func (r ApplicationIdentity) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdxc.ApplicationIdentityData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.provider.Client.CreateApplicationIdentity(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created application identity")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ApplicationIdentity) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdxc.ApplicationIdentityData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.provider.Client.ReadApplicationIdentity(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ApplicationIdentity) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdxc.ApplicationIdentityData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.provider.Client.UpdateApplicationIdentity(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ApplicationIdentity) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdxc.ApplicationIdentityData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.provider.Client.DeleteApplicationIdentity(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ApplicationIdentity) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
