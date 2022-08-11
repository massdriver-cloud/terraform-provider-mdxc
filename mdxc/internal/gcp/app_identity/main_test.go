package app_identity_test

import (
	"context"
	"terraform-provider-mdxc/mdxc/internal/gcp/app_identity"
	"terraform-provider-mdxc/mdxc/internal/massdriver"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
)

type mockProjectsServiceAccountsCreateCall iam.ProjectsServiceAccountsCreateCall

func (m *mockProjectsServiceAccountsCreateCall) Do(opts ...googleapi.CallOption) (*iam.ServiceAccount, error) {
	return &iam.ServiceAccount{}, nil
}

type mockCreateServiceAccountAPI func(projectName string, createserviceaccountrequest *iam.CreateServiceAccountRequest) *mockProjectsServiceAccountsCreateCall

func (m mockCreateServiceAccountAPI) Create(projectName string, createserviceaccountrequest *iam.CreateServiceAccountRequest) *mockProjectsServiceAccountsCreateCall {
	return &mockProjectsServiceAccountsCreateCall{}
}

func CreateServiceAccount(t *testing.T) {
	m := mockCreateServiceAccountAPI(func(projectName string, createserviceaccountrequest *iam.CreateServiceAccountRequest) *mockProjectsServiceAccountsCreateCall {
		t.Helper()
		// if params.RoleName == nil {
		// 	t.Fatal("expect role name to not be nil")
		// }

		// arn := fmt.Sprintf("arn:aws:iam::account:role/%s", *params.RoleName)
		// role := &types.Role{Arn: aws.String(arn)}
		// return &iam.CreateRoleOutput{Role: role}, nil
		mockCreateCall := mockProjectsServiceAccountsCreateCall{}
		return &mockCreateCall
	})

	appIdentityInput := massdriver.AppIdentityInput{
		Name: aws.String("test"),
	}

	serviceAccount, _ := app_identity.CreateServiceAccount(context.TODO(), m, &appIdentityInput)
	got := serviceAccount.Email
	want := "test@PROJECT_ID.iam.gserviceaccount.com"

	if want != got {
		t.Errorf("expect %v, got %v", want, got)
	}
}
