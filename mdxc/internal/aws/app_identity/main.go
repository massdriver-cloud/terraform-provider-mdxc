package app_identity

import (
	"context"
	"terraform-provider-mdxc/mdxc/internal/massdriver"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// take a client & massdriver.AppIdentityInput{}

type IAMCreateRoleAPI interface {
	CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns []func(*iam.Options)) (*iam.CreateRoleOutput, error)
}

// TODO: return (*iam.CreateRoleOutput, error) -> AppIdentity{}
func Create(ctx context.Context, api IAMCreateRoleAPI, input *massdriver.AppIdentityInput) (*iam.CreateRoleOutput, error) {

	roleInput := iam.CreateRoleInput{
		RoleName: input.Name,
	}

	roleOutput, err := api.CreateRole(ctx, &roleInput, []func(*iam.Options){})

	// TODO: cast to AppIdentityOutput
	return roleOutput, err
}
