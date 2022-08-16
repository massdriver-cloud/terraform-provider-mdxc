package mdxc

import (
	"context"
	"terraform-provider-mdxc/internal/cloud/aws"
	"terraform-provider-mdxc/internal/cloud/azure"
	"terraform-provider-mdxc/internal/cloud/gcp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type AWSApplicationPermissionInputData struct {
	PolicyARN types.String `tfsdk:"policy_arn"`
}
type AzureApplicationPermissionInputData struct {
	RoleName types.String `tfsdk:"role_name"`
	Scope    types.String `tfsdk:"scope"`
}
type GCPApplicationPermissionInputData struct {
	Role      types.String `tfsdk:"role"`
	Condition types.String `tfsdk:"condition"`
}

type ApplicationPermissionData struct {
	Id                    types.String                         `tfsdk:"id"`
	ApplicationIdentityID types.String                         `tfsdk:"application_identity_id"`
	AWSInput              *AWSApplicationPermissionInputData   `tfsdk:"aws_configuration"`
	AzureInput            *AzureApplicationPermissionInputData `tfsdk:"azure_configuration"`
	GCPInput              *GCPApplicationPermissionInputData   `tfsdk:"gcp_configuration"`
}

func (c *MDXCClient) CreateApplicationPermission(ctx context.Context, d *ApplicationPermissionData) diag.Diagnostics {
	switch c.Cloud {
	case "aws":
		return runApplicationPermissionFunctionAWS(aws.CreateApplicationPermission, ctx, d, c.AWSConfig)
	case "azure":
		return runApplicationPermissionFunctionAzure(azure.CreateApplicationPermission, ctx, d, c.AzureConfig)
	case "gcp":
		return runApplicationPermissionFunctionGCP(gcp.CreateApplicationPermission, ctx, d, c.GCPConfig)
	}
	return diag.Diagnostics{diag.NewErrorDiagnostic("Cloud not supported", "Provider does not support specified cloud: "+c.Cloud)}
}

func (c *MDXCClient) ReadApplicationPermission(ctx context.Context, d *ApplicationPermissionData) diag.Diagnostics {
	switch c.Cloud {
	case "aws":
		return runApplicationPermissionFunctionAWS(aws.ReadApplicationPermission, ctx, d, c.AWSConfig)
	case "azure":
		return runApplicationPermissionFunctionAzure(azure.ReadApplicationPermission, ctx, d, c.AzureConfig)
	case "gcp":
		return runApplicationPermissionFunctionGCP(gcp.ReadApplicationPermission, ctx, d, c.GCPConfig)
	}
	return diag.Diagnostics{diag.NewErrorDiagnostic("Cloud not supported", "Provider does not support specified cloud: "+c.Cloud)}
}

func (c *MDXCClient) UpdateApplicationPermission(ctx context.Context, d *ApplicationPermissionData) diag.Diagnostics {
	switch c.Cloud {
	case "aws":
		return runApplicationPermissionFunctionAWS(aws.UpdateApplicationPermission, ctx, d, c.AWSConfig)
	case "azure":
		return runApplicationPermissionFunctionAzure(azure.UpdateApplicationPermission, ctx, d, c.AzureConfig)
	case "gcp":
		return runApplicationPermissionFunctionGCP(gcp.UpdateApplicationPermission, ctx, d, c.GCPConfig)
	}
	return diag.Diagnostics{diag.NewErrorDiagnostic("Cloud not supported", "Provider does not support specified cloud: "+c.Cloud)}
}

func (c *MDXCClient) DeleteApplicationPermission(ctx context.Context, d *ApplicationPermissionData) diag.Diagnostics {
	switch c.Cloud {
	case "aws":
		return runApplicationPermissionFunctionAWS(aws.DeleteApplicationPermission, ctx, d, c.AWSConfig)
	case "azure":
		return runApplicationPermissionFunctionAzure(azure.DeleteApplicationPermission, ctx, d, c.AzureConfig)
	case "gcp":
		return runApplicationPermissionFunctionGCP(gcp.DeleteApplicationPermission, ctx, d, c.GCPConfig)
	}
	return diag.Diagnostics{diag.NewErrorDiagnostic("Cloud not supported", "Provider does not support specified cloud: "+c.Cloud)}
}

// -------------- AWS --------------
type applicationPermissionFunctionAWS func(context.Context, *aws.ApplicationPermissionConfig, aws.IAMClient) error

func convertApplicationPermissionConfigTerraformToAWS(d *ApplicationPermissionData, a *aws.ApplicationPermissionConfig) {
	a.ID = d.Id.Value
	a.RoleName = d.ApplicationIdentityID.Value
	if d.AWSInput != nil {
		a.PolicyARN = d.AWSInput.PolicyARN.Value
	}
}

func convertApplicationPermissionConfigAWSToTerraform(a *aws.ApplicationPermissionConfig, d *ApplicationPermissionData) {
	d.Id = types.String{Value: a.ID}
	d.ApplicationIdentityID = types.String{Value: a.RoleName}
	if d.AWSInput == nil {
		d.AWSInput = &AWSApplicationPermissionInputData{}
	}
	d.AWSInput.PolicyARN = types.String{Value: a.PolicyARN}
}

func runApplicationPermissionFunctionAWS(function applicationPermissionFunctionAWS, ctx context.Context, d *ApplicationPermissionData, config *aws.AWSConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	iamClient := config.NewIAMService()
	cloudApplicationPermissionConfig := aws.ApplicationPermissionConfig{}
	convertApplicationPermissionConfigTerraformToAWS(d, &cloudApplicationPermissionConfig)
	err := function(ctx, &cloudApplicationPermissionConfig, iamClient)
	if err != nil {
		diags.Append(
			diag.NewErrorDiagnostic(err.Error(), ""),
		)
		return diags
	}
	convertApplicationPermissionConfigAWSToTerraform(&cloudApplicationPermissionConfig, d)
	return diags
}

// -------------- Azure --------------
type applicationPermissionFunctionAzure func(context.Context, *azure.ApplicationPermissionConfig, azure.RoleAssignmentsClient, azure.RoleDefinitionsClient) error

func convertApplicationPermissionConfigTerraformToAzure(d *ApplicationPermissionData, a *azure.ApplicationPermissionConfig) {
	a.ID = d.Id.Value
	if d.AzureInput != nil {
		a.RoleName = d.AzureInput.RoleName.Value
		a.Scope = d.AzureInput.Scope.Value
	}
}

func convertApplicationPermissionConfigAzureToTerraform(a *azure.ApplicationPermissionConfig, d *ApplicationPermissionData) {
	d.Id = types.String{Value: a.ID}
	if d.AzureInput == nil {
		d.AzureInput = &AzureApplicationPermissionInputData{}
	}
	d.AzureInput.RoleName = types.String{Value: a.RoleName}
	d.AzureInput.Scope = types.String{Value: a.Scope}
}

func runApplicationPermissionFunctionAzure(function applicationPermissionFunctionAzure, ctx context.Context, d *ApplicationPermissionData, config *azure.AzureConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	raClient, raErr := config.NewRoleAssignmentsClient(ctx)
	if raErr != nil {
		diags.Append(
			diag.NewErrorDiagnostic(raErr.Error(), ""),
		)
		return diags
	}
	rdClient, rdErr := config.NewRoleDefinitionsClient(ctx)
	if rdErr != nil {
		diags.Append(
			diag.NewErrorDiagnostic(rdErr.Error(), ""),
		)
		return diags
	}
	cloudApplicationPermissionConfig := azure.ApplicationPermissionConfig{}
	convertApplicationPermissionConfigTerraformToAzure(d, &cloudApplicationPermissionConfig)
	err := function(ctx, &cloudApplicationPermissionConfig, raClient, rdClient)
	if err != nil {
		diags.Append(
			diag.NewErrorDiagnostic(err.Error(), ""),
		)
		return diags
	}
	convertApplicationPermissionConfigAzureToTerraform(&cloudApplicationPermissionConfig, d)
	return diags
}

// // -------------- GCP --------------
type applicationPermissionFunctionGCP func(context.Context, *gcp.ApplicationPermissionConfig, gcp.GCPResourceManagerIface) (gcp.GCPIAMResponse, error)

func convertApplicationPermissionConfigTerraformToGCP(d *ApplicationPermissionData, a *gcp.ApplicationPermissionConfig, c *gcp.GCPConfig) {
	a.ID = d.Id.Value
	a.Project = c.Provider.Project.Value
	if d.GCPInput != nil {
		a.Role = d.GCPInput.Role.Value
		a.Member = d.ApplicationIdentityID.Value
		a.Condition = d.GCPInput.Condition.Value
	}
}

func convertApplicationPermissionConfigGCPToTerraform(a *gcp.ApplicationPermissionConfig, d *ApplicationPermissionData) {
	d.Id = types.String{Value: a.ID}
	if d.GCPInput == nil {
		d.GCPInput = &GCPApplicationPermissionInputData{}
	}
	d.GCPInput.Role = types.String{Value: a.Role}
	d.ApplicationIdentityID = types.String{Value: a.Member}
	d.GCPInput.Condition = types.String{Value: a.Condition}
}

func runApplicationPermissionFunctionGCP(function applicationPermissionFunctionGCP, ctx context.Context, d *ApplicationPermissionData, config *gcp.GCPConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	iamClient, serviceErr := config.NewResourceManagerService(ctx, config.TokenSource)
	if serviceErr != nil {
		diags.Append(
			diag.NewErrorDiagnostic(serviceErr.Error(), ""),
		)
		return diags
	}

	cloudApplicationPermissionConfig := gcp.ApplicationPermissionConfig{}
	convertApplicationPermissionConfigTerraformToGCP(d, &cloudApplicationPermissionConfig, config)
	response, err := function(ctx, &cloudApplicationPermissionConfig, iamClient)
	if err != nil {
		diags.Append(
			diag.NewErrorDiagnostic(err.Error(), ""),
		)
		return diags
	}
	tflog.Debug(ctx, "Permissions added for"+response.Email)
	convertApplicationPermissionConfigGCPToTerraform(&cloudApplicationPermissionConfig, d)
	return diags
}
