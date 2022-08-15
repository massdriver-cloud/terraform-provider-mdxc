package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type ApplicationPermissionConfig struct {
	ID        string
	RoleName  string
	PolicyARN string
}

func CreateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client IAMAPI) error {

	roleInput := iam.AttachRolePolicyInput{
		RoleName:  &config.RoleName,
		PolicyArn: &config.PolicyARN,
	}

	_, attachErr := client.AttachRolePolicy(ctx, &roleInput)
	if attachErr != nil {
		return attachErr
	}

	config.ID = fmt.Sprintf("%s/%s", config.RoleName, config.PolicyARN)

	return nil
}

func ReadApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client IAMAPI) error {
	return nil
}

func UpdateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client IAMAPI) error {
	return nil
}

func DeleteApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client IAMAPI) error {

	input := iam.DetachRolePolicyInput{
		RoleName:  &config.RoleName,
		PolicyArn: &config.PolicyARN,
	}

	_, deleteErr := client.DetachRolePolicy(ctx, &input)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
