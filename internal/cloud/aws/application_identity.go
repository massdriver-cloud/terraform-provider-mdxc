// Package app_identity implements the massdriver.AppIdentity for AWS
package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

// type IAMCreateRoleAPI interface {
// 	CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns []func(*iam.Options)) (*iam.CreateRoleOutput, error)
// }

// TODO: call this and inject into Create()
func (c AWSConfig) NewService() (interface{}, error) {
	client := iam.NewFromConfig(*c.config)
	return client, nil
}

// // Create an AWS IAM Role as a massdriver.AppIdentity
// func AppIdentityCreate(ctx context.Context, api IAMCreateRoleAPI, input *massdriver.AppIdentityInput) (*massdriver.AppIdentityOutput, error) {
// 	roleInput := iam.CreateRoleInput{
// 		RoleName: input.Name,
// 	}

// 	roleOutput, err := api.CreateRole(ctx, &roleInput, []func(*iam.Options){})

// 	appIdentityOutput := massdriver.AppIdentityOutput{
// 		AwsIamRole: *roleOutput.Role,
// 	}

// 	return &appIdentityOutput, err
// }

func (c *AWSConfig) CreateApplicationIdentity(ctx context.Context, d *schema.ResourceData, cloudClient interface{}) diag.Diagnostics {
	iamClient := cloudClient.(*iam.Client)

	var applicationIdentityConfig map[string]interface{}

	if awsBlock, ok := d.Get("aws").([]interface{}); ok && len(awsBlock) > 0 && awsBlock[0] != nil {
		applicationIdentityConfig = awsBlock[0].(map[string]interface{})
	} else {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "AWS configuration block not specified",
			},
		}
	}

	log.Printf("--------------------------------------------------------------------------")
	log.Printf("%v", applicationIdentityConfig)

	var assumeRolePolicyFromConfig string
	var assumeRolePolicyFromConfigOk bool
	if assumeRolePolicyFromConfig, assumeRolePolicyFromConfigOk = applicationIdentityConfig["assume_role_policy"].(string); !assumeRolePolicyFromConfigOk || assumeRolePolicyFromConfig == "" {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "assume_role_policy not set",
			},
		}
	}

	assumeRolePolicy, assumeErr := structure.NormalizeJsonString(assumeRolePolicyFromConfig)
	if assumeErr != nil {
		return diag.FromErr(assumeErr)
	}

	roleInput := iam.CreateRoleInput{
		AssumeRolePolicyDocument: &assumeRolePolicy,
		RoleName:                 aws.String(d.Get("name").(string)),
	}

	roleOutput, roleErr := iamClient.CreateRole(ctx, &roleInput)
	if roleErr != nil {
		return diag.FromErr(roleErr)
	}

	roleName := *roleOutput.Role.RoleName
	d.SetId(roleName)

	return nil
}

func (c *AWSConfig) DeleteApplicationIdentity(ctx context.Context, d *schema.ResourceData, cloudClient interface{}) diag.Diagnostics {
	iamClient := cloudClient.(*iam.Client)

	roleInput := iam.DeleteRoleInput{
		RoleName: aws.String(d.Id()),
	}

	_, roleErr := iamClient.DeleteRole(ctx, &roleInput)
	if roleErr != nil {
		return diag.FromErr(roleErr)
	}

	return nil
}
