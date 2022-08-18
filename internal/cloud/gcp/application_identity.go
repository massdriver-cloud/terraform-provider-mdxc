// Package app_identity implements the massdriver.AppIdentity for GCP
// https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-go
package gcp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

// // Create a massdriver AppIdentity in GCP.
// func Create(ctx context.Context, api *iam.Service, input *massdriver.AppIdentityInput) (*massdriver.AppIdentityOutput, error) {
// 	// We will need to apply a number of operations for GCP. We should use a backoff library and
// 	// Dave's checkpointing idea to handle failures during the 'transaction'
// 	// TODO: api enablement (we should be non-authoritative)
// 	svcAcct, _ := CreateServiceAccount(ctx, api.Projects.ServiceAccounts, input)
// 	// TODO func CreateProjectIAMBinding()
// 	// TODO func CreateServiceAccountIAMMember()

type ApplicationIdentityConfig struct {
	ID                  string
	Project             string
	Name                string
	ServiceAccountEmail string
}

type GCPIamIface interface {
	Create(name string, createserviceaccountrequest *iam.CreateServiceAccountRequest) *iam.ProjectsServiceAccountsCreateCall
	Get(email string) *iam.ProjectsServiceAccountsGetCall
	Patch(email string, patchserviceaccountrequest *iam.PatchServiceAccountRequest) *iam.ProjectsServiceAccountsPatchCall
	Delete(email string) *iam.ProjectsServiceAccountsDeleteCall
}

func gcpIAMClientFactory(ctx context.Context, tokenSource oauth2.TokenSource) (GCPIamIface, error) {
	service, err := iam.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	return service.Projects.ServiceAccounts, nil
}

func CreateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, iamClient GCPIamIface) error {
	request := &iam.CreateServiceAccountRequest{
		AccountId: config.Name,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: config.Name,
		},
	}

	projectResourceName := fmt.Sprintf("projects/%s", config.Project)
	serviceAccount, doErr := iamClient.Create(projectResourceName, request).Do()
	if doErr != nil {
		return doErr
	}

	config.ID = serviceAccount.UniqueId
	config.Name = serviceAccount.DisplayName
	config.ServiceAccountEmail = serviceAccount.Email

	return nil
}

func ReadApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, iamClient GCPIamIface) error {
	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, config.ID)
	serviceAccount, doErr := iamClient.Get(resourceName).Do()
	if doErr != nil {
		return doErr
	}

	config.ID = serviceAccount.UniqueId
	config.Name = serviceAccount.DisplayName
	config.ServiceAccountEmail = serviceAccount.Email

	return nil
}

func UpdateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, iamClient GCPIamIface) error {
	request := &iam.PatchServiceAccountRequest{
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: config.Name,
		},
	}
	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, config.ID)
	_, doErr := iamClient.Patch(resourceName, request).Do()
	if doErr != nil {
		return doErr
	}
	return nil
}

func DeleteApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, iamClient GCPIamIface) error {
	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, config.ID)

	tflog.Debug(ctx, "------------------------------------------------------------------"+resourceName)

	_, doErr := iamClient.Delete(resourceName).Do()
	return doErr
}
