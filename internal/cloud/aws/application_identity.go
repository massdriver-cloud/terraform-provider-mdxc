// Package app_identity implements the massdriver.AppIdentity for AWS
package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

type ApplicationIdentityConfig struct {
	Name             string
	IAMRoleARN       string
	AssumeRolePolicy string
}

func CreateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client IAMClient) error {
	assumeRolePolicy, assumeErr := structure.NormalizeJsonString(config.AssumeRolePolicy)
	if assumeErr != nil {
		return assumeErr
	}

	roleInput := iam.CreateRoleInput{
		AssumeRolePolicyDocument: &assumeRolePolicy,
		RoleName:                 &config.Name,
	}

	output, roleErr := client.CreateRole(ctx, &roleInput)
	if roleErr != nil {
		return roleErr
	}

	config.Name = *output.Role.RoleName
	config.IAMRoleARN = *output.Role.Arn

	return nil
}

func ReadApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client IAMClient) error {
	return nil
}

func UpdateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client IAMClient) error {
	return nil
}

func DeleteApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client IAMClient) error {

	input := iam.DeleteRoleInput{
		RoleName: aws.String(config.Name),
	}

	_, deleteErr := client.DeleteRole(ctx, &input)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
