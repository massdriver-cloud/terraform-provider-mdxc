// Package app_identity implements the massdriver.AppIdentity for AWS
package aws

import (
	"context"

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
func NewIAMService(cfg *aws.Config) *iam.Client {
	client := iam.NewFromConfig(*cfg)
	return client
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

func (c *AWSClient) CreateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamClient := NewIAMService(c.config)

	assumeRolePolicy, assumeErr := structure.NormalizeJsonString(`
	{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"sts:AssumeRole"
				],
				"Principal": {
					"Service": [
						"ec2.amazonaws.com"
					]
				}
			}
		]
	}
	`)
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

func (c *AWSClient) DeleteAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamClient := NewIAMService(c.config)

	roleInput := iam.DeleteRoleInput{
		RoleName: aws.String(d.Id()),
	}

	_, roleErr := iamClient.DeleteRole(ctx, &roleInput)
	if roleErr != nil {
		return diag.FromErr(roleErr)
	}

	return nil
}
