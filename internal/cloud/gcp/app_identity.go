// Package app_identity implements the massdriver.AppIdentity for GCP
// https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-go
package gcp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/oauth2"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

// TODO: NewService(creds ... are they in ctx w/ GCP?) -> Service to pass to Create()
// Pretty sure we'll need an oauth token like this (https://github.com/massdriver-cloud/satellite/blob/5e5cbba01d2563e7eb2316d3a9b71e007e109a75/src/handler/dns_zone/gcp.go#L20)
// but haven't seen tf provider client auth yet...
func NewIAMService(ctx context.Context, tokenSource oauth2.TokenSource) (*iam.Service, error) {
	service, err := iam.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	return service, nil
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

func (c *GCPConfig) CreateAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamClient, serviceErr := NewIAMService(ctx, c.tokenSource)
	if serviceErr != nil {
		return diag.FromErr(serviceErr)
	}

	request := &iam.CreateServiceAccountRequest{
		AccountId: d.Get("name").(string),
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: d.Get("name").(string),
		},
	}

	projectId := c.project

	serviceAccountOutput, doErr := iamClient.Projects.ServiceAccounts.Create("projects/"+projectId, request).Do()
	if doErr != nil {
		return diag.FromErr(doErr)
	}

	d.SetId(serviceAccountOutput.UniqueId)

	return nil
}

func (c *GCPConfig) DeleteAppIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamClient, serviceErr := NewIAMService(ctx, c.tokenSource)
	if serviceErr != nil {
		return diag.FromErr(serviceErr)
	}

	name := "projects/" + c.project + "/serviceAccounts/" + d.Id()

	_, doErr := iamClient.Projects.ServiceAccounts.Delete(name).Do()
	if doErr != nil {
		return diag.FromErr(doErr)
	}

	return nil
}
