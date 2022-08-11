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

// TODO: API Enablement
// TODO: google_project_iam_binding
// TODO: google_service_account_iam_member
func Create(ctx context.Context, api string, input *massdriver.AppIdentityInput) (*massdriver.AppIdentityOutput, error) {
	// request := &iam.CreateServiceAccountRequest{
	// 	AccountId: *input.Name,
	// 	ServiceAccount: &iam.ServiceAccount{
	// 		DisplayName: *input.Name,
	// 	},
	// }

	//TODO: projectId must come from tfland
	// projectId := "foo"

	// account, err := api.Projects.ServiceAccounts.Create("projects/"+projectId, request).Do()
	// if err != nil {
	// 	return nil, fmt.Errorf("Projects.ServiceAccounts.Create: %v", err)
	// }
	// fmt.Fprintf(w, "Created service account: %v", account)
	// return account, nil

	return &massdriver.AppIdentityOutput{
		GcpServiceAccount: iam.ServiceAccount{Email: "test@PROJECT_ID.iam.gserviceaccount.com"},
	}, nil
}
