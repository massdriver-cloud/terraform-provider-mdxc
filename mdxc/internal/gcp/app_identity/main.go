// Package app_identity implements the massdriver.AppIdentity for GCP
// https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-go
package app_identity

import (
	"context"
	"fmt"
	"terraform-provider-mdxc/mdxc/internal/massdriver"

	"google.golang.org/api/iam/v1"
)

const rolesIamServiceAccountUser = "roles/iam.serviceAccountUser"

// TODO: NewService(creds ... are they in ctx w/ GCP?) -> Service to pass to Create()
// Pretty sure we'll need an oauth token like this (https://github.com/massdriver-cloud/satellite/blob/5e5cbba01d2563e7eb2316d3a9b71e007e109a75/src/handler/dns_zone/gcp.go#L20)
// but haven't seen tf provider client auth yet...
func NewService(ctx context.Context) (*iam.Service, error) {
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	return service, nil
}

// BindServiceAccountUserRole binds the serviceAccountUser role to the iam.ServiceAccount
func BindServiceAccountUserRole(ctx context.Context, api string, svcAcct *iam.ServiceAccount) *iam.Binding {
	members := []string{fmt.Sprintf("serviceAccount:%s", svcAcct.Email)}

	// @@HERE -> add the policy for _this_ SA to access rolesIamServiceAccountUser, which will result in an iam.Binding
	// https://pkg.go.dev/google.golang.org/api/iam/v1#Binding
	binding := iam.Binding{Role: rolesIamServiceAccountUser, Members: members}

	return &binding
}

// CreateServiceAccount makes the service account identity for an app
func CreateServiceAccount(ctx context.Context, serviceAcctApi *iam.ProjectsServiceAccountsService, input *massdriver.AppIdentityInput) (*iam.ServiceAccount, error) {
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

// Create a massdriver AppIdentity in GCP.
func Create(ctx context.Context, api *iam.Service, input *massdriver.AppIdentityInput) (*massdriver.AppIdentityOutput, error) {
	// We will need to apply a number of operations for GCP. We should use a backoff library and
	// Dave's checkpointing idea to handle failures during the 'transaction'
	// TODO: api enablement (we should be non-authoritative)
	svcAcct, _ := CreateServiceAccount(ctx, api.Projects.ServiceAccounts, input)
	// TODO: binding, _ := BindServiceAccountUserRole(ctx, api, svcAcct)

	// TODO func CreateServiceAccountIAMMember()

	return &massdriver.AppIdentityOutput{
		GcpServiceAccount: iam.ServiceAccount{Email: svcAcct.Email},
	}, nil
}
