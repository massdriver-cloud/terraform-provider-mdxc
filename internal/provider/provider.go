package provider

import (
	"context"
	"terraform-provider-mdxc/internal/mdxc"

	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &MDXCProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MDXCProvider{}
	}
}

type MDXCProvider struct {
	Configured bool
	Client     *mdxc.MDXCClient
}

var awsProviderSchema = tfsdk.Attribute{
	Optional: true,
	Validators: []tfsdk.AttributeValidator{
		schemavalidator.ConflictsWith(
			path.MatchRoot("azure"),
			path.MatchRoot("gcp"),
		),
	},
	Description: "Credentials for AWS Cloud",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"role_arn": {
			Optional:    true,
			Description: "ARN of AWS Role to assume",
			Type:        types.StringType,
			Sensitive:   true,
			//RequiredWith: []string{"aws.0.external_id", "aws.0.region"},
		},
		"external_id": {
			Optional:    true,
			Description: "A unique identifier that might be required when you assume a role in another account.",
			Type:        types.StringType,
			Sensitive:   true,
			//RequiredWith: []string{"aws.0.role_arn", "aws.0.region"},
		},
		"region": {
			Optional:    true,
			Description: "The region where AWS operations will take place.",
			Type:        types.StringType,
			//RequiredWith: []string{"aws.0.role_arn", "aws.0.external_id"},
		},
	}),
}

var azureProviderSchema = tfsdk.Attribute{
	Optional: true,
	Validators: []tfsdk.AttributeValidator{
		schemavalidator.ConflictsWith(
			path.MatchRoot("aws"),
			path.MatchRoot("gcp"),
		),
	},
	Description: "Credentials for Azure Cloud. See how to authenticate through Service Principal in the [Azure docs](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_secret#creating-a-service-principal)",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"subscription_id": {
			Optional:    true,
			Description: "Azure Subscription ID",
			Type:        types.StringType,
			Sensitive:   true,
			//RequiredWith: []string{"azure.0.client_id", "azure.0.client_secret", "azure.0.tenant_id"},
		},
		"client_id": {
			Optional:    true,
			Description: "Azure Client ID",
			Type:        types.StringType,
			//RequiredWith: []string{"azure.0.subscription_id", "azure.0.client_secret", "azure.0.tenant_id"},
		},
		"client_secret": {
			Optional:    true,
			Description: "Azure Client Secret",
			Type:        types.StringType,
			Sensitive:   true,
			//RequiredWith: []string{"azure.0.subscription_id", "azure.0.client_id", "azure.0.tenant_id"},
		},
		"tenant_id": {
			Optional:    true,
			Description: "Azure Tenant ID",
			Type:        types.StringType,
			Sensitive:   true,
			//RequiredWith: []string{"azure.0.subscription_id", "azure.0.client_id", "azure.0.client_secret"},
		},
	}),
}

var gcpProviderSchema = tfsdk.Attribute{
	Optional: true,
	Validators: []tfsdk.AttributeValidator{
		schemavalidator.ConflictsWith(
			path.MatchRoot("aws"),
			path.MatchRoot("azure"),
		),
	},
	Description: "Credentials for Google Cloud. See how to authenticate through Service Principals in the [Google docs](https://cloud.google.com/compute/docs/authentication)",
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"credentials": {
			Optional:    true,
			Description: "Either the path to or the contents of a service account key file in JSON format.",
			Type:        types.StringType,
			Sensitive:   true,
			//RequiredWith: []string{"gcp.0.project"},
		},
		"project": {
			Optional:    true,
			Description: "The GCP project to manage resources in.",
			Type:        types.StringType,
			//RequiredWith: []string{"gcp.0.credentials"},
		},
	}),
}

func (p *MDXCProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Terraform provider to resource across multiple clouds.",
		Attributes: map[string]tfsdk.Attribute{
			"aws":   awsProviderSchema,
			"azure": azureProviderSchema,
			"gcp":   gcpProviderSchema,
		},
	}, nil
}

func (p *MDXCProvider) GetResources(_ context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"mdxc_application_identity":   ApplicationIdentityType{},
		"mdxc_application_permission": ApplicationPermissionType{},
	}, nil
}

// GetDataSources - Defines Provider data sources
func (p *MDXCProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{}, nil
}

func (p *MDXCProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config mdxc.MDXCProviderConfig
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Configuring Provider.")

	mdxcClient, factoryErr := mdxc.MDXCClientFactory(ctx, &config)
	if factoryErr != nil {
		resp.Diagnostics.AddError(
			"Error configuring credentials.",
			factoryErr.Error(),
		)
		return
	}

	p.Configured = true
	p.Client = mdxcClient
}
