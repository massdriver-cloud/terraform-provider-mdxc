package mdxc

import (
	"context"
	"log"
	"terraform-provider-mdxc/internal/cloud/aws"
	"terraform-provider-mdxc/internal/cloud/azure"
	"terraform-provider-mdxc/internal/cloud/gcp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AWSApplicationIdentityData struct {
	AssumeRolePolicy types.String `tfsdk:"assume_role_policy"`
}

type ApplicationIdentityData struct {
	Id   types.String                `tfsdk:"id"`
	Name types.String                `tfsdk:"name"`
	AWS  *AWSApplicationIdentityData `tfsdk:"aws"`
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
type applicationIdentityFunctionAWS func(context.Context, *aws.ApplicationIdentityConfig, aws.IAMAPI) error

func convertApplicationIdentityConfigTerraformToAWS(d *ApplicationIdentityData, a *aws.ApplicationIdentityConfig) {
	a.AssumeRolePolicy = d.AWS.AssumeRolePolicy.Value
	a.Name = d.Name.Value
}

func convertApplicationIdentityConfigAWSToTerraform(a *aws.ApplicationIdentityConfig, d *ApplicationIdentityData) {
	d.AWS.AssumeRolePolicy = types.String{Value: a.AssumeRolePolicy}
	d.Name = types.String{Value: a.Name}
	d.Id = types.String{Value: a.Name}
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
type applicationIdentityFunctionAzure func(context.Context, *azure.ApplicationIdentityConfig, azure.ApplicationAPI) error

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
	applicationClient, appServiceErr := config.NewApplicationService(ctx)
	if appServiceErr != nil {
		diags.Append(
			diag.NewErrorDiagnostic(appServiceErr.Error(), ""),
		)
		return diags
	}
	cloudApplicationIdentityConfig := azure.ApplicationIdentityConfig{}
	convertApplicationIdentityConfigTerraformToAzure(d, &cloudApplicationIdentityConfig)
	err := function(ctx, &cloudApplicationIdentityConfig, applicationClient)
	if err != nil {
		diags.Append(
			diag.NewErrorDiagnostic(err.Error(), ""),
		)
		return diags
	}
	convertApplicationIdentityConfigAzureToTerraform(&cloudApplicationIdentityConfig, d)
	return diags
}

// -------------- GCP --------------
type applicationIdentityFunctionGCP func(context.Context, *gcp.ApplicationIdentityConfig, gcp.GCPIamIface) error

func convertApplicationIdentityConfigTerraformToGCP(d *ApplicationIdentityData, a *gcp.ApplicationIdentityConfig, c *gcp.GCPConfig) {
	a.Name = d.Name.Value
	a.ID = d.Id.Value
	a.Project = c.Provider.Project.Value
}

func convertApplicationIdentityConfigGCPToTerraform(a *gcp.ApplicationIdentityConfig, d *ApplicationIdentityData) {
	d.Name = types.String{Value: a.Name}
	d.Id = types.String{Value: a.ID}
}

func runApplicationIdentityFunctionGCP(function applicationIdentityFunctionGCP, ctx context.Context, d *ApplicationIdentityData, config *gcp.GCPConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("[debug] NewIAMService")
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
