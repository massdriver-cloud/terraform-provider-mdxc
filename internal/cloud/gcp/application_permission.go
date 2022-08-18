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
func readModifyWriteWithBackoff(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface, modifyFunc func(ctx context.Context, role string, member string, policy *cloudresourcemanager.Policy) error) error {
	backoff := time.Second

	for {
		projectPolicy, err := getProjectIamPolicy(ctx, client, config.Project)
		if err != nil {
			return err
		}

		modifyFunc(ctx, config.Role, config.Member, projectPolicy)

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

		// TODO: retry on not-found SA
		// retry in the case that a service account is not found. This can happen when a service account is deleted
		// out of band.
		// if isServiceAccountNotFoundError, _ := iamServiceAccountNotFound(err); isServiceAccountNotFoundError {
		// 	// calling a retryable function within a retry loop is not
		// 	// strictly the _best_ idea, but this error only happens in
		// 	// high-traffic projects anyways
		// 	currentPolicy, rerr := iamPolicyReadWithRetry(updater)
		// 	if rerr != nil {
		// 		if p.Etag != currentPolicy.Etag {
		// 			// not matching indicates that there is a new state to attempt to apply
		// 			// log.Printf("current and old etag did not match for %s, retrying", updater.DescribeResource())
		// 			time.Sleep(backoff)
		// 			backoff = backoff * 2
		// 			continue
		// 		}

		// 		// log.Printf("current and old etag matched for %s, not retrying", updater.DescribeResource())
		// 	} else {
		// 		// if the error is non-nil, just fall through and return the base error
		// 		// log.Printf("[DEBUG]: error checking etag for policy %s. error: %v", updater.DescribeResource(), rerr)
		// 	}
		// }
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

func addToPolicy(ctx context.Context, role string, member string, policy *cloudresourcemanager.Policy) error {
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

func removeFromPolicy(ctx context.Context, role string, member string, policy *cloudresourcemanager.Policy) error {
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

func addWorkloadIdentityRole(ctx context.Context, config *ApplicationPermissionConfig, client GCPResourceManagerIface) error {
	// namespace := "default"
	// namePrefix := "example-apps"
	// k8sEmail := fmt.Sprintf("%s.svc.id.goog[%s/%s]", config.Project, namespace, namePrefix)
	// projectPolicy, err := getProjectIamPolicy(ctx, client, config.Project)
	// if err != nil {
	// 	return response, err
	// }
	// AddToPolicy(ctx, "roles/iam.workloadIdentityUser", k8sEmail, projectPolicy)
	// saveProjectIamPolicy(ctx, client, config.Project, projectPolicy)

	return nil
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