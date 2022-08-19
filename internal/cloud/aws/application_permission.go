package aws

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type ApplicationPermissionConfig struct {
	ID        string
	RoleARN   string
	PolicyARN string
}

func CreateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client IAMClient) error {

	roleName := getResourceNameFromARN(config.RoleARN)

	roleInput := iam.AttachRolePolicyInput{
		RoleName:  &roleName,
		PolicyArn: &config.PolicyARN,
	}

	_, attachErr := client.AttachRolePolicy(ctx, &roleInput)
	if attachErr != nil {
		return attachErr
	}

	config.ID = fmt.Sprintf("%s#%s", config.RoleARN, config.PolicyARN)

	return nil
}

func ReadApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client IAMClient) error {
	return nil
}

func UpdateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client IAMClient) error {
	return nil
}

func DeleteApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client IAMClient) error {
	roleName := getResourceNameFromARN(config.RoleARN)
	input := iam.DetachRolePolicyInput{
		RoleName:  &roleName,
		PolicyArn: &config.PolicyARN,
	}

	_, deleteErr := client.DetachRolePolicy(ctx, &input)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func getResourceNameFromARN(arn string) string {
	var nameRegex = regexp.MustCompile(`[^:/]*$`)
	return nameRegex.FindString(arn)
}
