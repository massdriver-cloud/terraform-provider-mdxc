package gcp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

type Role struct {
	Role      string
	Condition string
}

type ApplicationPermissionConfig struct {
	ID               string
	ServiceAccountID string
	Project          string
	Role             string
	Condition        string
	Member           string
	// TODO: remove
	// Roles []Role
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
	projectPolicy, err := getProjectIamPolicy(ctx, client, config.Project)
	if err != nil {
		return response, err
	}

	tflog.Debug(ctx, "adding role "+config.Role)
	AddToPolicy(ctx, config.Role, config.Member, projectPolicy)

	if errSave := saveProjectIamPolicy(ctx, client, config.Project, projectPolicy); errSave != nil {
		return response, errSave
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
	projectPolicy, err := getProjectIamPolicy(ctx, client, config.Project)
	if err != nil {
		return response, err
	}

	RemoveFromPolicy(config.Role, config.Member, projectPolicy)

	if errSave := saveProjectIamPolicy(ctx, client, config.Project, projectPolicy); errSave != nil {
		return response, errSave
	}

	return response, nil
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
func RemoveFromPolicy(role string, member string, policy *cloudresourcemanager.Policy) error {
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
	projectId := ""
	namespace := "default"
	namePrefix := "example-apps"
	k8sEmail := fmt.Sprintf("%s.svc.id.goog[%s/%s]", projectId, namespace, namePrefix)
	response := GCPIAMResponse{}
	projectPolicy, err := getProjectIamPolicy(ctx, client, config.Project)
	if err != nil {
		return response, err
	}
	AddToPolicy(ctx, "roles/iam.workloadIdentityUser", k8sEmail, projectPolicy)
	saveProjectIamPolicy(ctx, client, config.Project, projectPolicy)

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
