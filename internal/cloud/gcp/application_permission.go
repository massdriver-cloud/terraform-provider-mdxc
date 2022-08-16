package gcp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

type ApplicationPermissionConfig struct {
	ID               string
	ServiceAccountID string
	Project          string
	Role             string
	Condition        string
	Member           string
}

type GCPIAMResponse struct {
	Email string
}

type GCPResourceManagerIface interface {
	GetIamPolicy(resourceName string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *cloudresourcemanager.ProjectsGetIamPolicyCall
	SetIamPolicy(resourceName string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *cloudresourcemanager.ProjectsSetIamPolicyCall
}

func resourceManagerClientFactory(ctx context.Context, tokenSource oauth2.TokenSource) (GCPResourceManagerIface, error) {
	service, err := cloudresourcemanager.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("cloudresourcemanager.NewService: %v", err)
	}

	return service.Projects, nil
}

func CreateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) (GCPIAMResponse, error) {
	response := GCPIAMResponse{}

	if errDo := readModifyWriteWithBackoff(ctx, config, client, AddToPolicy); errDo != nil {
		return response, errDo
	}

	return response, nil
}

func ReadApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) (GCPIAMResponse, error) {
	response := GCPIAMResponse{}
	return response, nil
}

func UpdateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) (GCPIAMResponse, error) {
	response := GCPIAMResponse{}
	return response, nil
}

func DeleteApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) (GCPIAMResponse, error) {
	response := GCPIAMResponse{}

	if errDo := readModifyWriteWithBackoff(ctx, config, client, RemoveFromPolicy); errDo != nil {
		return response, errDo
	}

	return response, nil
}

func readModifyWriteWithBackoff(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface, modifyFunc func(ctx context.Context, role string, member string, policy *cloudresourcemanager.Policy) error) error {
	backoff := time.Second

	for {
		projectPolicy, err := getProjectIamPolicy(ctx, client, config.Project)
		if err != nil {
			return err
		}

		AddToPolicy(ctx, config.Role, config.Member, projectPolicy)

		errSave := saveProjectIamPolicy(ctx, client, config.Project, projectPolicy)
		if errSave == nil {
			// TODO: fetch again I think?
			// https://github.com/hashicorp/terraform-provider-google/blob/2c3be0cf1f9c56231817a2e876fa63b1afdb46e2/google/iam.go#L103
			break
		}
		if isConflictError(errSave) {
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return errwrap.Wrapf(fmt.Sprintf("Error applying IAM policy to %s: Too many conflicts.  Latest error: {{err}}", "create permission"), err)
			}
			continue
		}
	}

	return nil
}

// https://cloud.google.com/iam/docs/reference/rest/v1/projects.serviceAccounts/create
func getProjectIamPolicy(ctx context.Context, service GCPResourceManagerIface, projectId string) (*cloudresourcemanager.Policy, error) {
	tflog.Debug(ctx, "getProjectIamPolicy "+projectId)
	getCall := service.GetIamPolicy(projectId, &cloudresourcemanager.GetIamPolicyRequest{})
	policy, errDo := getCall.Do()
	if errDo != nil {
		return nil, errDo
	}
	tflog.Debug(ctx, "got iam policy")

	return policy, nil
}

func saveProjectIamPolicy(ctx context.Context, service GCPResourceManagerIface, projectId string, policy *cloudresourcemanager.Policy) error {
	tflog.Debug(ctx, "saveProjectIamPolicy "+projectId)
	saveCall := service.SetIamPolicy(projectId, &cloudresourcemanager.SetIamPolicyRequest{
		Policy: policy,
	})
	policy, errDo := saveCall.Do()
	if errDo != nil {
		return errDo
	}
	return nil
}

// good thing to test
func AddToPolicy(ctx context.Context, role string, member string, policy *cloudresourcemanager.Policy) error {
	addedToExisting := false
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			tflog.Debug(ctx, "adding to existing")
			// TODO: dedupe members
			binding.Members = append(binding.Members, fmt.Sprintf("serviceAccount:%s", member))
			addedToExisting = true
		}
	}
	if !addedToExisting {
		tflog.Debug(ctx, "adding new existing")
		policy.Bindings = append(policy.Bindings, &cloudresourcemanager.Binding{
			Role: role,
			Members: []string{
				fmt.Sprintf("serviceAccount:%s", member),
			},
		})
	}
	return nil
}

// good thing to test
func RemoveFromPolicy(ctx context.Context, role string, member string, policy *cloudresourcemanager.Policy) error {
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			membersToKeep := []string{}
			for _, existingMember := range binding.Members {
				if existingMember != fmt.Sprintf("serviceAccount:%s", member) {
					membersToKeep = append(membersToKeep, existingMember)
				}
			}
			binding.Members = membersToKeep
		}
	}
	return nil
}

func addWorkloadIdentityRole(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) (GCPIAMResponse, error) {
	// namespace := "default"
	// namePrefix := "example-apps"
	// k8sEmail := fmt.Sprintf("%s.svc.id.goog[%s/%s]", config.Project, namespace, namePrefix)
	response := GCPIAMResponse{}
	// projectPolicy, err := getProjectIamPolicy(ctx, client, config.Project)
	// if err != nil {
	// 	return response, err
	// }
	// AddToPolicy(ctx, "roles/iam.workloadIdentityUser", k8sEmail, projectPolicy)
	// saveProjectIamPolicy(ctx, client, config.Project, projectPolicy)

	return response, nil
}

func GetServiceAccountIamPolicy(projectId string, serviceId string) (*iam.Policy, error) {
	// ctx := context.Background()
	// iamService, err := iam.NewService(ctx, option.WithCredentialsFile(jsonPath))
	// if err != nil {
	// 	return nil, err
	// }
	// getCall := iamService.Projects.ServiceAccounts.GetIamPolicy("projects/" + projectId + "/serviceAccounts/" + serviceId)
	// iamPolicy, errGet := getCall.Do()
	// if errGet != nil {
	// 	return nil, errGet
	// }

	return nil, nil
}
