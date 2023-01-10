package mdxc

import (
	"context"
	"terraform-provider-mdxc/internal/cloud/aws"
	"terraform-provider-mdxc/internal/cloud/azure"
	"terraform-provider-mdxc/internal/cloud/gcp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AWSApplicationIdentityInputData struct {
	AssumeRolePolicy types.String `tfsdk:"assume_role_policy"`
}
type GCPApplicationIdentityInputData struct {
	Kubernetes *GCPKubernetesIdentityInputData `tfsdk:"kubernetes"`
}
type GCPKubernetesIdentityInputData struct {
	Namespace          types.String `tfsdk:"namespace"`
	ServiceAccountName types.String `tfsdk:"service_account_name"`
}
type AzureApplicationIdentityInputData struct {
	Placeholder types.String `tfsdk:"placeholder"`
}

type AWSApplicationIdentityOutputData struct {
	IAMRoleARN types.String `tfsdk:"iam_role_arn"`
}
type AzureApplicationIdentityOutputData struct {
	ApplicationID            types.String `tfsdk:"application_id"`
	ServicePrincipalID       types.String `tfsdk:"service_principal_id"`
	ServicePrincipalClientID types.String `tfsdk:"service_principal_client_id"`
	ServicePrincipalSecret   types.String `tfsdk:"service_principal_secret"`
}
type GCPApplicationIdentityOutputData struct {
	ServiceAccountEmail types.String `tfsdk:"service_account_email"`
}

type ApplicationIdentityData struct {
	Id          types.String                        `tfsdk:"id"`
	Name        types.String                        `tfsdk:"name"`
	Cloud       types.String                        `tfsdk:"cloud"`
	AWSInput    *AWSApplicationIdentityInputData    `tfsdk:"aws_configuration"`
	AzureInput  *AzureApplicationIdentityInputData  `tfsdk:"azure_configuration"`
	GCPInput    *GCPApplicationIdentityInputData    `tfsdk:"gcp_configuration"`
	AWSOutput   *AWSApplicationIdentityOutputData   `tfsdk:"aws_application_identity"`
	AzureOutput *AzureApplicationIdentityOutputData `tfsdk:"azure_application_identity"`
	GCPOutput   *GCPApplicationIdentityOutputData   `tfsdk:"gcp_application_identity"`
}

func (c *MDXCClient) CreateApplicationIdentity(ctx context.Context, d *ApplicationIdentityData) diag.Diagnostics {
	switch c.Cloud {
	case "aws":
		return runApplicationIdentityFunctionAWS(aws.CreateApplicationIdentity, ctx, d, c.AWSConfig)
	case "azure":
		return runApplicationIdentityFunctionAzure(azure.CreateApplicationIdentity, ctx, d, c.AzureConfig)
	case "gcp":
		return runApplicationIdentityFunctionGCP(gcp.CreateApplicationIdentity, ctx, d, c.GCPConfig)
	}
	return diag.Diagnostics{diag.NewErrorDiagnostic("Cloud not supported", "Provider does not support specified cloud: "+c.Cloud)}
}

func (c *MDXCClient) ReadApplicationIdentity(ctx context.Context, d *ApplicationIdentityData) diag.Diagnostics {
	switch c.Cloud {
	case "aws":
		return runApplicationIdentityFunctionAWS(aws.ReadApplicationIdentity, ctx, d, c.AWSConfig)
	case "azure":
		return runApplicationIdentityFunctionAzure(azure.ReadApplicationIdentity, ctx, d, c.AzureConfig)
	case "gcp":
		return runApplicationIdentityFunctionGCP(gcp.ReadApplicationIdentity, ctx, d, c.GCPConfig)
	}
	return diag.Diagnostics{diag.NewErrorDiagnostic("Cloud not supported", "Provider does not support specified cloud: "+c.Cloud)}
}

func (c *MDXCClient) UpdateApplicationIdentity(ctx context.Context, d *ApplicationIdentityData) diag.Diagnostics {
	switch c.Cloud {
	case "aws":
		return runApplicationIdentityFunctionAWS(aws.UpdateApplicationIdentity, ctx, d, c.AWSConfig)
	case "azure":
		return runApplicationIdentityFunctionAzure(azure.UpdateApplicationIdentity, ctx, d, c.AzureConfig)
	case "gcp":
		return runApplicationIdentityFunctionGCP(gcp.UpdateApplicationIdentity, ctx, d, c.GCPConfig)
	}
	return diag.Diagnostics{diag.NewErrorDiagnostic("Cloud not supported", "Provider does not support specified cloud: "+c.Cloud)}
}

func (c *MDXCClient) DeleteApplicationIdentity(ctx context.Context, d *ApplicationIdentityData) diag.Diagnostics {
	switch c.Cloud {
	case "aws":
		return runApplicationIdentityFunctionAWS(aws.DeleteApplicationIdentity, ctx, d, c.AWSConfig)
	case "azure":
		return runApplicationIdentityFunctionAzure(azure.DeleteApplicationIdentity, ctx, d, c.AzureConfig)
	case "gcp":
		return runApplicationIdentityFunctionGCP(gcp.DeleteApplicationIdentity, ctx, d, c.GCPConfig)
	}
	return diag.Diagnostics{diag.NewErrorDiagnostic("Cloud not supported", "Provider does not support specified cloud: "+c.Cloud)}
}

// -------------- AWS --------------
type applicationIdentityFunctionAWS func(context.Context, *aws.ApplicationIdentityConfig, aws.IAMClient) error

func convertApplicationIdentityConfigTerraformToAWS(d *ApplicationIdentityData, a *aws.ApplicationIdentityConfig) {
	a.Name = d.Name.Value
	if d.AWSInput != nil {
		a.AssumeRolePolicy = d.AWSInput.AssumeRolePolicy.Value
	}
	if d.AWSOutput != nil {
		a.IAMRoleARN = d.AWSOutput.IAMRoleARN.Value
	}
}

func convertApplicationIdentityConfigAWSToTerraform(a *aws.ApplicationIdentityConfig, d *ApplicationIdentityData) {
	d.Id = types.String{Value: a.IAMRoleARN}
	d.Name = types.String{Value: a.Name}
	if d.AWSInput == nil {
		d.AWSInput = &AWSApplicationIdentityInputData{}
	}
	if d.AWSOutput == nil {
		d.AWSOutput = &AWSApplicationIdentityOutputData{}
	}
	d.AWSInput.AssumeRolePolicy = types.String{Value: a.AssumeRolePolicy}
	d.AWSOutput.IAMRoleARN = types.String{Value: a.IAMRoleARN}
}

func runApplicationIdentityFunctionAWS(function applicationIdentityFunctionAWS, ctx context.Context, d *ApplicationIdentityData, config *aws.AWSConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	iamClient := config.NewIAMService()
	cloudApplicationIdentityConfig := aws.ApplicationIdentityConfig{}
	convertApplicationIdentityConfigTerraformToAWS(d, &cloudApplicationIdentityConfig)
	err := function(ctx, &cloudApplicationIdentityConfig, iamClient)
	if err != nil {
		diags.Append(
			diag.NewErrorDiagnostic(err.Error(), ""),
		)
		return diags
	}
	convertApplicationIdentityConfigAWSToTerraform(&cloudApplicationIdentityConfig, d)
	return diags
}

// -------------- Azure --------------
type applicationIdentityFunctionAzure func(context.Context, *azure.ApplicationIdentityConfig, azure.ManagedIdentityClient, azure.FederatedIdentityCredentialClient) error

func convertApplicationIdentityConfigTerraformToAzure(d *ApplicationIdentityData, a *azure.ApplicationIdentityConfig) {
	a.Name = d.Name.Value
	a.ID = d.Id.Value
}

func convertApplicationIdentityConfigAzureToTerraform(a *azure.ApplicationIdentityConfig, d *ApplicationIdentityData) {
	d.Name = types.String{Value: a.Name}
	d.Id = types.String{Value: a.ID}
}

func runApplicationIdentityFunctionAzure(function applicationIdentityFunctionAzure, ctx context.Context, d *ApplicationIdentityData, config *azure.AzureConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := config.NewManagedIdentityClient(ctx, config.Provider)
	if err != nil {
		diags.Append(
			diag.NewErrorDiagnostic(err.Error(), ""),
		)
		return diags
	}
	fedClient, errFed := config.NewFederatedIdentityCredentialsClient(ctx, config.Provider)
	if errFed != nil {
		diags.Append(
			diag.NewErrorDiagnostic(err.Error(), ""),
		)
		return diags
	}

	cloudApplicationIdentityConfig := azure.ApplicationIdentityConfig{}
	convertApplicationIdentityConfigTerraformToAzure(d, &cloudApplicationIdentityConfig)
	errRunFunc := function(ctx, &cloudApplicationIdentityConfig, client, fedClient)
	if errRunFunc != nil {
		diags.Append(
			diag.NewErrorDiagnostic(errRunFunc.Error(), ""),
		)
		return diags
	}
	convertApplicationIdentityConfigAzureToTerraform(&cloudApplicationIdentityConfig, d)
	return diags
}

// -------------- GCP --------------
type applicationIdentityFunctionGCP func(context.Context, *gcp.ApplicationIdentityConfig, gcp.GCPIamIface) error

func convertApplicationIdentityConfigTerraformToGCP(d *ApplicationIdentityData, a *gcp.ApplicationIdentityConfig, c *gcp.GCPConfig) {
	a.ID = d.Id.Value
	a.ServiceAccountEmail = d.Id.Value
	a.Name = d.Name.Value
	a.Project = c.Provider.Project.Value
	if d.GCPInput != nil {
		if d.GCPInput.Kubernetes != nil {
			a.KubernetesNamspace = d.GCPInput.Kubernetes.Namespace.Value
			a.KubernetesServiceAccountName = d.GCPInput.Kubernetes.ServiceAccountName.Value
		}
	}
	if d.GCPOutput != nil {
	}
}

func convertApplicationIdentityConfigGCPToTerraform(a *gcp.ApplicationIdentityConfig, d *ApplicationIdentityData) {
	d.Id = types.String{Value: a.ID}
	d.Name = types.String{Value: a.Name}
	if d.GCPOutput == nil {
		d.GCPOutput = &GCPApplicationIdentityOutputData{}
	}
	d.GCPOutput.ServiceAccountEmail = types.String{Value: a.ID}
}

func runApplicationIdentityFunctionGCP(function applicationIdentityFunctionGCP, ctx context.Context, d *ApplicationIdentityData, config *gcp.GCPConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	iamClient, serviceErr := config.NewIAMService(ctx, config.TokenSource)
	if serviceErr != nil {
		diags.Append(
			diag.NewErrorDiagnostic(serviceErr.Error(), ""),
		)
		return diags
	}
	cloudApplicationIdentityConfig := gcp.ApplicationIdentityConfig{}
	convertApplicationIdentityConfigTerraformToGCP(d, &cloudApplicationIdentityConfig, config)
	err := function(ctx, &cloudApplicationIdentityConfig, iamClient)
	if err != nil {
		diags.Append(
			diag.NewErrorDiagnostic(err.Error(), ""),
		)
		return diags
	}
	convertApplicationIdentityConfigGCPToTerraform(&cloudApplicationIdentityConfig, d)
	return diags
}
