package iam

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

func main() {
	namePrefix := "md-beebe-06"
	namespace := "default"
	projectId := "md-wbeebe-0808-example-apps"
	k8sEmail := fmt.Sprintf("%s.svc.id.goog[%s/%s]", projectId, namespace, namePrefix)

	serviceAccount, err := CreateServiceAccount(projectId, namePrefix)
	if err != nil {
		log.Fatal(err)
	}

	if errGet := printProjectPolicy(projectId); errGet != nil {
		log.Fatal(errGet)
	}

	// the roles from app connections w/ policies
	roles := []Role{
		{
			Role: "roles/cloudsql.editor",
			// TODO: support condition
		},
	}
	if errRoles := addProjectRolesToServiceAccount(projectId, serviceAccount, roles); errRoles != nil {
		log.Fatal(errRoles)
	}

	// Allow the Kubernetes service account to impersonate the IAM service account by adding an IAM policy binding between the two service accounts. This binding allows the Kubernetes service account to act as the IAM service account
	if errAddMember := addServiceAccountIamPolicyBinding(projectId, serviceAccount, "roles/iam.workloadIdentityUser", k8sEmail); errAddMember != nil {
		log.Fatal(errAddMember)
	}
}

func printProjectPolicy(projectId string) error {
	log.Println("project policy")
	projectPolicy, err := GetProjectIamPolicy(projectId)
	if err != nil {
		return err
	}
	for _, binding := range projectPolicy.Bindings {
		log.Printf("role: %s", binding.Role)
		for _, member := range binding.Members {
			log.Printf("member: %s", member)
		}
	}
	return nil
}

type Role struct {
	Role           string
	Condition      string
	memberAppended bool
}

func addProjectRolesToServiceAccount(projectId string, serviceAccount *iam.ServiceAccount, roles []Role) error {
	projectPolicy, err := GetProjectIamPolicy(projectId)
	if err != nil {
		return err
	}

	for _, binding := range projectPolicy.Bindings {
		for _, role := range roles {
			if binding.Role == role.Role {
				// TODO: dedupe members
				binding.Members = append(binding.Members, "serviceAccount:"+serviceAccount.Email)
				role.memberAppended = true
			}
		}
	}
	for _, role := range roles {
		if !role.memberAppended {
			projectPolicy.Bindings = append(projectPolicy.Bindings, &cloudresourcemanager.Binding{
				Role: role.Role,
				Members: []string{
					"serviceAccount:" + serviceAccount.Email,
				},
			})
		}
	}

	if errSave := SaveProjectIamPolicy(projectId, projectPolicy); errSave != nil {
		return errSave
	}

	return nil
}

// https://cloud.google.com/sdk/gcloud/reference/iam/service-accounts/add-iam-policy-binding
func addServiceAccountIamPolicyBinding(projectId string, serviceAccount *iam.ServiceAccount, role string, member string) error {
	policy, err := GetServiceAccountIamPolicy(projectId, serviceAccount.UniqueId)
	if err != nil {
		return err
	}
	memberAppended := false
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			binding.Members = append(binding.Members, "serviceAccount:"+member)
			memberAppended = true
		}
	}
	if !memberAppended {
		policy.Bindings = append(policy.Bindings, &iam.Binding{
			Role: role,
			Members: []string{
				"serviceAccount:" + member,
			},
		})
	}

	if errSave := SaveServiceAccountPolicy(policy, serviceAccount); errSave != nil {
		return errSave
	}

	return nil
}

// https://cloud.google.com/iam/docs/reference/rest/v1/projects.serviceAccounts/create?apix_params=%7B%22name%22%3A%22projects%2Fmd-wbeebe-0808-example-apps%22%2C%22resource%22%3A%7B%22accountId%22%3A%22md-name-prefix-1234%22%7D%7D
func CreateServiceAccount(projectId string, accountId string) (*iam.ServiceAccount, error) {
	jsonPath := "./creds.json"
	ctx := context.Background()
	iamService, err := iam.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, err
	}

	createCall := iamService.Projects.ServiceAccounts.Create("projects/"+projectId, &iam.CreateServiceAccountRequest{
		AccountId: accountId,
	})
	serviceAccount, errCreate := createCall.Do()
	if errCreate != nil {
		return nil, errCreate
	}
	log.Printf("service account created %s", serviceAccount.Email)

	return serviceAccount, nil
}

func GetServiceAccountIamPolicy(projectId string, serviceId string) (*iam.Policy, error) {
	jsonPath := "./creds.json"
	ctx := context.Background()
	iamService, err := iam.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, err
	}
	getCall := iamService.Projects.ServiceAccounts.GetIamPolicy("projects/" + projectId + "/serviceAccounts/" + serviceId)
	iamPolicy, errGet := getCall.Do()
	if errGet != nil {
		return nil, errGet
	}

	return iamPolicy, nil
}

func GetProjectIamPolicy(projectId string) (*cloudresourcemanager.Policy, error) {
	jsonPath := "./creds.json"
	ctx := context.Background()
	resourceM, _ := cloudresourcemanager.NewService(ctx, option.WithCredentialsFile(jsonPath))
	getCall := resourceM.Projects.GetIamPolicy(projectId, &cloudresourcemanager.GetIamPolicyRequest{})
	policy, errDo := getCall.Do()
	if errDo != nil {
		return nil, errDo
	}

	return policy, nil
}

func SaveProjectIamPolicy(projectId string, policy *cloudresourcemanager.Policy) error {
	jsonPath := "./creds.json"
	ctx := context.Background()
	resourceM, _ := cloudresourcemanager.NewService(ctx, option.WithCredentialsFile(jsonPath))
	saveCall := resourceM.Projects.SetIamPolicy(projectId, &cloudresourcemanager.SetIamPolicyRequest{
		Policy: policy,
	})
	policy, errDo := saveCall.Do()
	if errDo != nil {
		return errDo
	}
	return nil
}

// https://cloud.google.com/iam/docs/reference/rest/v1/projects.serviceAccounts/setIamPolicy
func SaveServiceAccountPolicy(policy *iam.Policy, serviceAccount *iam.ServiceAccount) error {
	jsonPath := "./creds.json"
	ctx := context.Background()
	iamService, err := iam.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return err
	}

	setCall := iamService.Projects.ServiceAccounts.SetIamPolicy("projects/md-wbeebe-0808-example-apps/serviceAccounts/"+serviceAccount.UniqueId, &iam.SetIamPolicyRequest{
		Policy: policy,
	})
	_, errSet := setCall.Do()
	if errSet != nil {
		return errSet
	}
	return nil
}
