package gcp

import (
	"context"
	"fmt"
	thirdparty "terraform-provider-mdxc/internal/cloud/gcp/thirdparty/terraform-google-provider"
	"time"

	"github.com/hashicorp/errwrap"
	"golang.org/x/oauth2"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

type ApplicationPermissionConfig struct {
	ID               string
	ServiceAccountID string
	Project          string
	Role             string
	Condition        string
}

type GCPResourceManagerIface interface {
	GetIamPolicy(resourceName string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *cloudresourcemanager.ProjectsGetIamPolicyCall
	SetIamPolicy(resourceName string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *cloudresourcemanager.ProjectsSetIamPolicyCall
}

func gcpResourceManagerClientFactory(ctx context.Context, tokenSource oauth2.TokenSource) (GCPResourceManagerIface, error) {
	service, err := cloudresourcemanager.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("cloudresourcemanager.NewService: %v", err)
	}

	return service.Projects, nil
}

func CreateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) error {
	return readModifyWriteWithBackoff(ctx, config, client, addToPolicy)
}

func ReadApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) error {
	return nil
}

func UpdateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) error {
	return nil
}

func DeleteApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) error {
	return readModifyWriteWithBackoff(ctx, config, client, removeFromPolicy)
}

// https://github.com/hashicorp/terraform-provider-google/blob/2c3be0cf1f9c56231817a2e876fa63b1afdb46e2/google/iam.go#L103
func readModifyWriteWithBackoff(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface, modifyFunc func(ctx context.Context, config *ApplicationPermissionConfig, policy *cloudresourcemanager.Policy) error) error {
	backoff := time.Second

	for {
		policy, err := getProjectIamPolicy(ctx, client, config.Project)
		if err != nil {
			return err
		}

		errModify := modifyFunc(ctx, config, policy)
		if errModify != nil {
			return errModify
		}

		errSave := saveProjectIamPolicy(ctx, client, config.Project, policy)
		if errSave == nil {
			// TODO: fetch again I think?
			// https://github.com/hashicorp/terraform-provider-google/blob/2c3be0cf1f9c56231817a2e876fa63b1afdb46e2/google/iam.go#L103
			break
		}
		if thirdparty.IsConflictError(errSave) {
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return errwrap.Wrapf("Error applying IAM policy to %s: Too many conflicts.  Latest error: {{err}}", err)
			}
			continue
		}
		if errSave != nil {
			return errSave
		}
	}

	config.ID = fmt.Sprintf("%s-%s", config.ServiceAccountID, config.Role)

	return nil
}

// https://cloud.google.com/iam/docs/reference/rest/v1/projects.serviceAccounts/create
func getProjectIamPolicy(ctx context.Context, service GCPResourceManagerIface, projectId string) (*cloudresourcemanager.Policy, error) {
	getCall := service.GetIamPolicy(projectId, &cloudresourcemanager.GetIamPolicyRequest{
		Options: &cloudresourcemanager.GetPolicyOptions{
			// policies with any conditional role bindings must specify version 3.
			// https://cloud.google.com/iam/docs/policies#versions
			RequestedPolicyVersion: 3,
		},
	})
	policy, errDo := getCall.Do()
	if errDo != nil {
		return nil, errDo
	}
	// The "RequestedPolicyVersion" above isn't guaranteeing version 3, so we force it here
	policy.Version = 3
	return policy, nil
}

func saveProjectIamPolicy(ctx context.Context, service GCPResourceManagerIface, projectId string, policy *cloudresourcemanager.Policy) error {
	saveCall := service.SetIamPolicy(projectId, &cloudresourcemanager.SetIamPolicyRequest{
		Policy: policy,
	})
	policy, errDo := saveCall.Do()
	if errDo != nil {
		return errDo
	}
	return nil
}

func addToPolicy(ctx context.Context, config *ApplicationPermissionConfig, policy *cloudresourcemanager.Policy) error {
	role := config.Role
	member := config.ServiceAccountID

	policy.Bindings = thirdparty.AddBinding(policy.Bindings, &cloudresourcemanager.Binding{
		Role: role,
		Condition: &cloudresourcemanager.Expr{
			Expression: config.Condition,
		},
		Members: []string{
			fmt.Sprintf("serviceAccount:%s", member),
		},
	})

	return nil
}

func removeFromPolicy(ctx context.Context, config *ApplicationPermissionConfig, policy *cloudresourcemanager.Policy) error {
	role := config.Role
	member := config.ServiceAccountID

	policy.Bindings = thirdparty.RemoveBinding(policy.Bindings, &cloudresourcemanager.Binding{
		Role: role,
		Condition: &cloudresourcemanager.Expr{
			Expression: config.Condition,
		},
		Members: []string{
			fmt.Sprintf("serviceAccount:%s", member),
		},
	})

	return nil
}
