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

type GCPIamIface interface {
	Create(name string, createserviceaccountrequest *iam.CreateServiceAccountRequest) *iam.ProjectsServiceAccountsCreateCall
	Get(name string) *iam.ProjectsServiceAccountsGetCall
	Patch(name string, patchserviceaccountrequest *iam.PatchServiceAccountRequest) *iam.ProjectsServiceAccountsPatchCall
	Delete(name string) *iam.ProjectsServiceAccountsDeleteCall
}

// TODO: NewService(creds ... are they in ctx w/ GCP?) -> Service to pass to Create()
// Pretty sure we'll need an oauth token like this (https://github.com/massdriver-cloud/satellite/blob/5e5cbba01d2563e7eb2316d3a9b71e007e109a75/src/handler/dns_zone/gcp.go#L20)
// but haven't seen tf provider client auth yet...
// func (c *GCPConfig) NewIAMService(ctx context.Context) (*iam.Service, error) {
// 	service, err := iam.NewService(ctx, option.WithTokenSource(c.tokenSource))
// 	if err != nil {
// 		return nil, fmt.Errorf("iam.NewService: %v", err)
// 	}

// 	return service, nil
// }

func gcpIAMClientFactory(ctx context.Context, tokenSource oauth2.TokenSource) (GCPIamIface, error) {
	service, err := iam.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	return service.Projects.ServiceAccounts, nil
}

// func CreateServiceAccount(ctx context.Context, serviceAcctApi *iam.ProjectsServiceAccountsService, input *massdriver.AppIdentityInput) (*iam.ServiceAccount, error) {
// 	request := &iam.CreateServiceAccountRequest{
// 		AccountId: *input.Name,
// 		ServiceAccount: &iam.ServiceAccount{
// 			DisplayName: *input.Name,
// 		},
// 	}

// 	//TODO: projectId must come from tfland
// 	projectId := "foo"

// 	return serviceAcctApi.Create("projects/"+projectId, request).Do()
// }

// // Create a massdriver AppIdentity in GCP.
// func Create(ctx context.Context, api *iam.Service, input *massdriver.AppIdentityInput) (*massdriver.AppIdentityOutput, error) {
// 	// We will need to apply a number of operations for GCP. We should use a backoff library and
// 	// Dave's checkpointing idea to handle failures during the 'transaction'
// 	// TODO: api enablement (we should be non-authoritative)
// 	svcAcct, _ := CreateServiceAccount(ctx, api.Projects.ServiceAccounts, input)
// 	// TODO func CreateProjectIAMBinding()
// 	// TODO func CreateServiceAccountIAMMember()

// 	return &massdriver.AppIdentityOutput{
// 		GcpServiceAccount: iam.ServiceAccount{Email: svcAcct.Email},
// 	}, nil
// }

type ApplicationIdentityConfig struct {
	ID      string
	Project string
	Name    string
}

func CreateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, iamClient GCPIamIface) error {
	request := &iam.CreateServiceAccountRequest{
		AccountId: config.Name,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: config.Name,
		},
	}

	projectResourceName := fmt.Sprintf("projects/%s", config.Project)
	serviceAccountOutput, doErr := iamClient.Create(projectResourceName, request).Do()
	if doErr != nil {
		return doErr
	}

	config.ID = serviceAccountOutput.UniqueId

	return nil
}

func ReadApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, iamClient GCPIamIface) error {
	serviceAccountResourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, config.ID)
	_, doErr := iamClient.Get(serviceAccountResourceName).Do()
	if doErr != nil {
		return doErr
	}

	return nil
}

func UpdateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, iamClient GCPIamIface) error {
	request := &iam.PatchServiceAccountRequest{
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: config.Name,
		},
	}
	serviceAccountResourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, config.ID)
	_, doErr := iamClient.Patch(serviceAccountResourceName, request).Do()
	if doErr != nil {
		return doErr
	}
	return nil
}

func DeleteApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, iamClient GCPIamIface) error {
	serviceAccountResourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, config.ID)

	tflog.Debug(ctx, "------------------------------------------------------------------"+serviceAccountResourceName)

	_, doErr := iamClient.Delete(serviceAccountResourceName).Do()
	return doErr
}
