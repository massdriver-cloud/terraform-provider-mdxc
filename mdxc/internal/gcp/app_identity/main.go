// Package app_identity implements the massdriver.AppIdentity for GCP
// https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-go
package app_identity

// TODO: SA API Enablement

import (
	"context"
	"fmt"
	"terraform-provider-mdxc/mdxc/internal/massdriver"

	"google.golang.org/api/iam/v1"
)

// TODO: NewService(creds ... are they in ctx w/ GCP?) -> Service to pass to Create()
func NewService(ctx context.Context) (*iam.Service, error) {
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	return service, nil
}

//*iam.ProjectsServiceAccountsService
type IAMServiceAccountAPI interface {
	Create(projectName string, createserviceaccountrequest *iam.CreateServiceAccountRequest) *iam.ProjectsServiceAccountsCreateCall
}

func CreateServiceAccount(ctx context.Context, serviceAcctApi IAMServiceAccountAPI, input *massdriver.AppIdentityInput) (*iam.ServiceAccount, error) {
	request := &iam.CreateServiceAccountRequest{
		AccountId: *input.Name,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: *input.Name,
		},
	}

	//TODO: projectId must come from tfland
	projectId := "foo"

	return serviceAcctApi.Create("projects/"+projectId, request).Do()
}

func Create(ctx context.Context, api *iam.Service, input *massdriver.AppIdentityInput) (*massdriver.AppIdentityOutput, error) {
	svcAcct, _ := CreateServiceAccount(ctx, api.Projects.ServiceAccounts, input)

	// TODO:
	// func CreateProjectIAMBinding()
	// func CreateServiceAccountIAMMember()
	// func APIEnablement
	// func Create() -> calls all of above takes NewService, make the right service for each and passes it in.
	return &massdriver.AppIdentityOutput{
		GcpServiceAccount: iam.ServiceAccount{Email: svcAcct.Email},
	}, nil
}
