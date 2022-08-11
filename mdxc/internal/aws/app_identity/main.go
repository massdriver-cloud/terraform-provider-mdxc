// Package app_identity implements the massdriver.AppIdentity for AWS
package app_identity

import (
	"context"
	"terraform-provider-mdxc/mdxc/internal/massdriver"

	"github.com/aws/aws-sdk-go-v2/config"
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

// Create an Amazon IAM service client
func getClient() (*iam.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := iam.NewFromConfig(cfg)

	return client, err
}
