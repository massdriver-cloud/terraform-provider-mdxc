// Package app_identity implements the massdriver.AppIdentity for AWS
package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

type IAMAPI interface {
	CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
	DeleteRole(ctx context.Context, params *iam.DeleteRoleInput, optFns ...func(*iam.Options)) (*iam.DeleteRoleOutput, error)
}

func (c AWSConfig) NewIAMService() IAMAPI {
	client := iam.NewFromConfig(*c.config)
	return client
}

type ApplicationIdentityConfig struct {
	Name             string
	AssumeRolePolicy string
}

func CreateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client IAMAPI) error {
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

	return nil
}

func ReadApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client IAMAPI) error {
	return nil
}

func UpdateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client IAMAPI) error {
	return nil
}

func DeleteApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client IAMAPI) error {

	input := iam.DeleteRoleInput{
		RoleName: aws.String(config.Name),
	}

	_, deleteErr := client.DeleteRole(ctx, &input)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
