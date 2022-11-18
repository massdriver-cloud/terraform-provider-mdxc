// Package app_identity implements the massdriver.AppIdentity for GCP
// https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-go
package gcp

import (
	"context"
	"fmt"
	"time"

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
	ID                           string
	Project                      string
	Name                         string
	ServiceAccountEmail          string
	KubernetesNamspace           string
	KubernetesServiceAccountName string
}

type GCPIamIface interface {
	Create(name string, createserviceaccountrequest *iam.CreateServiceAccountRequest) *iam.ProjectsServiceAccountsCreateCall
	Get(id string) *iam.ProjectsServiceAccountsGetCall
	Patch(email string, patchserviceaccountrequest *iam.PatchServiceAccountRequest) *iam.ProjectsServiceAccountsPatchCall
	Delete(email string) *iam.ProjectsServiceAccountsDeleteCall
	GetIamPolicy(resource string) *iam.ProjectsServiceAccountsGetIamPolicyCall
	SetIamPolicy(resource string, setiampolicyrequest *iam.SetIamPolicyRequest) *iam.ProjectsServiceAccountsSetIamPolicyCall
}

func gcpIAMClientFactory(ctx context.Context, tokenSource oauth2.TokenSource) (GCPIamIface, error) {
	service, err := iam.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	return service.Projects.ServiceAccounts, nil
}

func CreateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client GCPIamIface) error {
	request := &iam.CreateServiceAccountRequest{
		AccountId: config.Name,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: config.Name,
		},
	}

	projectResourceName := fmt.Sprintf("projects/%s", config.Project)
	serviceAccount, doErr := client.Create(projectResourceName, request).Do()
	if doErr != nil {
		return doErr
	}
	resourceID := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, serviceAccount.Email)
	err := retry(5, time.Second, func() (operr error) {
		_, errGet := client.Get(resourceID).Do()
		return errGet
	})
	if err != nil {
		return err
	}

	config.ID = serviceAccount.Email
	config.ServiceAccountEmail = serviceAccount.Email
	config.Name = serviceAccount.DisplayName

	if config.KubernetesNamspace != "" {
		if errAddRole := addWorkloadIdentityRole(ctx, config, client); errAddRole != nil {
			return errAddRole
		}
	}

	return nil
}

func ReadApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, iamClient GCPIamIface) error {
	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, config.ID)
	serviceAccount, doErr := iamClient.Get(resourceName).Do()
	if doErr != nil {
		return doErr
	}

	config.Name = serviceAccount.DisplayName

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

func DeleteApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client GCPIamIface) error {
	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, config.ID)
	_, doErr := client.Delete(resourceName).Do()
	return doErr
}

// google_service_account_iam_member
// sets an IAM policy for a GCP service account
func addWorkloadIdentityRole(ctx context.Context, config *ApplicationIdentityConfig, client GCPIamIface) error {
	k8sEmail := fmt.Sprintf("%s.svc.id.goog[%s/%s]", config.Project, config.KubernetesNamspace, config.KubernetesServiceAccountName)
	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", config.Project, config.ID)
	iamPolicy, errGet := client.GetIamPolicy(resourceName).Do()
	if errGet != nil {
		return errGet
	}

	// TODO: test if idempotent
	iamPolicy.Bindings = append(iamPolicy.Bindings, &iam.Binding{
		Role: "roles/iam.workloadIdentityUser",
		Members: []string{
			fmt.Sprintf("serviceAccount:%s", k8sEmail),
		},
	})

	_, errSet := client.SetIamPolicy(resourceName, &iam.SetIamPolicyRequest{
		Policy: iamPolicy,
	}).Do()
	if errSet != nil {
		return errSet
	}

	return nil
}
