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
var _ provider.ResourceType = ResourceApplicationIdentityType{}
var _ resource.Resource = ResourceApplicationIdentity{}
var _ resource.ResourceWithImportState = ResourceApplicationIdentity{}

type ResourceApplicationIdentityType struct{}

// https://github.com/hashicorp/terraform-provider-aws/pull/23060
// https://github.com/hashicorp/terraform-provider-aws/blob/main/internal/verify/json.go
var awsApplicationIdentityInputs = tfsdk.Attribute{
	Optional:    true,
	Description: "AWS IAM role configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"assume_role_policy": {
			Type:        types.StringType,
			Description: "The AWS IAM role assume role policy. Required if provisioning into AWS",
			Required:    true,
		},
	}),
}

var azureApplicationIdentityInputs = tfsdk.Attribute{
	Optional:    true,
	Description: "Azure application and service principal configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"placeholder": {
			Type:     types.StringType,
			Optional: true,
		},
	}),
}

var gcpApplicationIdentityInputs = tfsdk.Attribute{
	Optional:    true,
	Description: "GCP service account configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"kubernetes": {
			Optional:    true,
			Description: "Kubernetes configuration",
			Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
				"namespace": {
					Type:     types.StringType,
					Required: true,
				},
				"service_account_name": {
					Type:     types.StringType,
					Required: true,
				},
			}),
		},
	}),
}

var awsApplicationIdentityOutputs = tfsdk.Attribute{
	Computed:    true,
	Description: "AWS IAM role configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"iam_role_arn": {
			Type:     types.StringType,
			Computed: true,
		},
	}),
}

var azureApplicationIdentityOutputs = tfsdk.Attribute{
	Computed:    true,
	Description: "Azure application and service principal configuration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"application_id": {
			Type:     types.StringType,
			Computed: true,
		},
		"service_principal_id": {
			Type:     types.StringType,
			Computed: true,
		},
		"service_principal_client_id": {
			Type:     types.StringType,
			Computed: true,
		},
		"service_principal_secret": {
			Type:     types.StringType,
			Computed: true,
		},
	}),
}

var gcpApplicationIdentityOutputs = tfsdk.Attribute{
	Computed:    true,
	Description: "GCP service account onfiguration",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"service_account_email": {
			Type:     types.StringType,
			Computed: true,
		},
	}),
}

func (t ResourceApplicationIdentityType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"cloud": {
				Type:                types.StringType,
				MarkdownDescription: "The cloud the application identity was provisioned into (value will be `aws`, `azure` or `gcp`)",
				Computed:            true,
			},
			"aws_configuration":          awsApplicationIdentityInputs,
			"azure_configuration":        azureApplicationIdentityInputs,
			"gcp_configuration":          gcpApplicationIdentityInputs,
			"aws_application_identity":   awsApplicationIdentityOutputs,
			"azure_application_identity": azureApplicationIdentityOutputs,
			"gcp_application_identity":   gcpApplicationIdentityOutputs,
		},
	}, nil
}

func (t ResourceApplicationIdentityType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	return ResourceApplicationIdentity{
		provider: *(in.(*MDXCProvider)),
	}, diag.Diagnostics{}
}

type ResourceApplicationIdentity struct {
	provider MDXCProvider
}

func (r ResourceApplicationIdentity) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

func (r ResourceApplicationIdentity) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

func (r ResourceApplicationIdentity) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (r ResourceApplicationIdentity) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func (r ResourceApplicationIdentity) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
