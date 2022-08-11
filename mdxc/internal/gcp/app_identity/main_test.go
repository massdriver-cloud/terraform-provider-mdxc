package app_identity_test

import (
	"context"
	"terraform-provider-mdxc/mdxc/internal/gcp/app_identity"
	"terraform-provider-mdxc/mdxc/internal/massdriver"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"google.golang.org/api/iam/v1"
)

type mockCreateServiceAccountAPI struct{}

func (m mockCreateServiceAccountAPI) Create(projectName string, createserviceaccountrequest *iam.CreateServiceAccountRequest) *iam.ProjectsServiceAccountsCreateCall {

	// This needs to be an interface that implments -> Do(opts ...googleapi.CallOption) (*iam.ServiceAccount, error)
	createCall := iam.ProjectsServiceAccountsCreateCall{}

	return &createCall
}

func TestCreateServiceAccount(t *testing.T) {
	appIdentityInput := massdriver.AppIdentityInput{
		Name: aws.String("test"),
	}

	m := &mockCreateServiceAccountAPI{}

	serviceAccount, _ := app_identity.CreateServiceAccount(context.TODO(), m, &appIdentityInput)
	got := serviceAccount.Email
	want := "test@PROJECT_ID.iam.gserviceaccount.com"

	if want != got {
		t.Errorf("expect %v, got %v", want, got)
	}
}
