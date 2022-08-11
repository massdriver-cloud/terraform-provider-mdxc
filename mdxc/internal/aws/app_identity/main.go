// Package app_identity implements the massdriver.AppIdentity for AWS
package app_identity

import (
	"context"
	"terraform-provider-mdxc/mdxc/internal/massdriver"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type IAMCreateRoleAPI interface {
	CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns []func(*iam.Options)) (*iam.CreateRoleOutput, error)
}

// Create an AWS IAM Role as a massdriver.AppIdentity
func Create(ctx context.Context, api IAMCreateRoleAPI, input *massdriver.AppIdentityInput) (*massdriver.AppIdentityOutput, error) {
	roleInput := iam.CreateRoleInput{
		RoleName: input.Name,
	}

	roleOutput, err := api.CreateRole(ctx, &roleInput, []func(*iam.Options){})

	appIdentityOutput := massdriver.AppIdentityOutput{
		AwsIamRole: *roleOutput.Role,
	}

	return &appIdentityOutput, err
}
