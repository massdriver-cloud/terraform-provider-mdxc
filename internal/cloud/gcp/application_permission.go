package gcp

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2"
	"google.golang.org/api/cloudresourcemanager/v1"
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
	Roles            []Role
	Member           string
}

type GCPIAMResponse struct {
	Email string
}

type GCPResourceManagerIface interface {
	GetIamPolicy(project string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *cloudresourcemanager.ProjectsGetIamPolicyCall
	SetIamPolicy(project string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *cloudresourcemanager.ProjectsSetIamPolicyCall
}

// func NewClient() *GCPConfig {
// 	return &GCPConfig{
// 		NewResourceManagerService: resourceManagerClientFactory,
// 	}
// }

func resourceManagerClientFactory(ctx context.Context, tokenSource oauth2.TokenSource) (GCPResourceManagerIface, error) {
	service, err := cloudresourcemanager.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("cloudresourcemanager.NewService: %v", err)
	}

	return service.Projects, nil
}

func CreateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) (GCPIAMResponse, error) {
	response := GCPIAMResponse{}
	project := "md-wbeebe-0808-example-apps"
	projectPolicy, err := getProjectIamPolicy(client, project)
	if err != nil {
		return response, err
	}

	for _, role := range config.Roles {
		log.Printf("[perms] adding role %s", role.Role)
		AddToPolicy(role.Role, config.Member, projectPolicy)
	}

	if errSave := saveProjectIamPolicy(client, project, projectPolicy); errSave != nil {
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
	projectPolicy, err := getProjectIamPolicy(client, config.Project)
	if err != nil {
		return response, err
	}

	for _, role := range config.Roles {
		RemoveFromPolicy(role.Role, config.Member, projectPolicy)
	}

	if errSave := saveProjectIamPolicy(client, config.Project, projectPolicy); errSave != nil {
		return response, errSave
	}

	return response, nil
}

func getProjectIamPolicy(service GCPResourceManagerIface, projectId string) (*cloudresourcemanager.Policy, error) {
	id := "projects/" + projectId
	getCall := service.GetIamPolicy(id, &cloudresourcemanager.GetIamPolicyRequest{})
	policy, errDo := getCall.Do()
	if errDo != nil {
		return nil, errDo
	}
	log.Printf("[debug] got iam policy")

	return policy, nil
}

func saveProjectIamPolicy(service GCPResourceManagerIface, projectId string, policy *cloudresourcemanager.Policy) error {
	id := "projects/" + projectId

	saveCall := service.SetIamPolicy(id, &cloudresourcemanager.SetIamPolicyRequest{
		Policy: policy,
	})
	policy, errDo := saveCall.Do()
	if errDo != nil {
		return errDo
	}
	return nil
}

// good thing to test
func AddToPolicy(role string, member string, policy *cloudresourcemanager.Policy) error {
	addedToExisting := false
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			log.Printf("[debug] adding to existing")
			// TODO: dedupe members
			binding.Members = append(binding.Members, fmt.Sprintf("serviceAccount:%s", member))
			addedToExisting = true
		}
	}
	if !addedToExisting {
		log.Printf("[debug] adding new policy")
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

func addWorkloadIdentityRole() {

}
