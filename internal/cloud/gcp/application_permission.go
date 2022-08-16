package gcp

import (
	"context"
	"fmt"

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

func (c *GCPConfig) NewResourceManagerService(ctx context.Context) (*cloudresourcemanager.Service, error) {
	service, err := cloudresourcemanager.NewService(ctx, option.WithTokenSource(c.tokenSource))
	if err != nil {
		return nil, fmt.Errorf("cloudresourcemanager.NewService: %v", err)
	}

	return service, nil
}

func CreateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client *cloudresourcemanager.Service) error {
	projectPolicy, err := getProjectIamPolicy(client, config.Project)
	if err != nil {
		return err
	}

	for _, role := range config.Roles {
		AddToPolicy(role.Role, config.Member, projectPolicy)
	}

	if errSave := saveProjectIamPolicy(client, config.Project, projectPolicy); errSave != nil {
		return errSave
	}

	return nil
}

func ReadApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client *cloudresourcemanager.Service) error {
	return nil
}

func UpdateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client *cloudresourcemanager.Service) error {
	return nil
}

func DeleteApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client *cloudresourcemanager.Service) error {
	projectPolicy, err := getProjectIamPolicy(client, config.Project)
	if err != nil {
		return err
	}

	for _, role := range config.Roles {
		RemoveFromPolicy(role.Role, config.Member, projectPolicy)
	}

	if errSave := saveProjectIamPolicy(client, config.Project, projectPolicy); errSave != nil {
		return errSave
	}

	return nil
}

func getProjectIamPolicy(service *cloudresourcemanager.Service, projectId string) (*cloudresourcemanager.Policy, error) {
	getCall := service.Projects.GetIamPolicy(projectId, &cloudresourcemanager.GetIamPolicyRequest{})
	policy, errDo := getCall.Do()
	if errDo != nil {
		return nil, errDo
	}

	return policy, nil
}

func saveProjectIamPolicy(service *cloudresourcemanager.Service, projectId string, policy *cloudresourcemanager.Policy) error {
	saveCall := service.Projects.SetIamPolicy(projectId, &cloudresourcemanager.SetIamPolicyRequest{
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
			// TODO: dedupe members
			binding.Members = append(binding.Members, fmt.Sprintf("serviceAccount:%s", member))
			addedToExisting = true
		}
	}
	if !addedToExisting {
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
			filteredMembers := []string{}
			for _, existingMember := range binding.Members {
				if existingMember != fmt.Sprintf("serviceAccount:%s", member) {
					filteredMembers = append(filteredMembers, existingMember)
				}
			}
			binding.Members = filteredMembers
		}
	}
	return nil
}
