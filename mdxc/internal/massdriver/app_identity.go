package massdriver

import (
	awsTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	gcpTypes "google.golang.org/api/iam/v1"
)

type AppIdentityInput struct {
	Name *string
}

type AppIdentityOutput struct {
	AwsIamRole        awsTypes.Role
	GcpServiceAccount gcpTypes.ServiceAccount
}

// TBD: module or struct???
// app_identity.Create(...) -> Cloud agnostic create
