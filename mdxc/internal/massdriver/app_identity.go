package massdriver

import "github.com/aws/aws-sdk-go-v2/service/iam/types"

type AppIdentityInput struct {
	Name *string
}

type AppIdentityOutput struct {
	AwsIamRole types.Role
}

// TBD: module or struct???
// app_identity.Create(...) -> Cloud agnostic create
