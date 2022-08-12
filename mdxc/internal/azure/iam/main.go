package main

// Credentials for authing requests to azure inferred from following environment variables:
// export AZURE_TENANT_ID="<active_directory_tenant_id"
// export AZURE_CLIENT_ID="<service_principal_appid>"
// export AZURE_CLIENT_SECRET="<service_principal_password>"
// export AZURE_SUBSCRIPTION_ID="<subscription_id>"
// you can grab all of this info from an azure service principal artifact in massdriver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	// "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization"
	azauth "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/google/uuid"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/msgraph"
)

// read env vars
var (
	AzureSubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	AzureTenantID       = os.Getenv("AZURE_TENANT_ID")
	AzureClientID       = os.Getenv("AZURE_CLIENT_ID")
	AzureClientSecret   = os.Getenv("AZURE_CLIENT_SECRET")
)

// ApplicationCreate creates an azure cloud applicaiton resource to
// represent the kubenertes app
func ApplicationCreate(ctx context.Context, authorizer auth.Authorizer, name string) (*msgraph.Application, error) {
	// TODO need to add create if not exist logic as this will create many apps with the same display name
	ac := msgraph.NewApplicationsClient(AzureTenantID)
	ac.BaseClient.Authorizer = authorizer
	app, _, err := ac.Create(ctx, msgraph.Application{
		DisplayName: &name,
	})
	return app, err
}

func ApplicationDelete(ctx context.Context, authorizer auth.Authorizer, app *msgraph.Application) error {
	// TODO need to add create if not exist logic as this will create many apps with the same display name
	ac := msgraph.NewApplicationsClient(AzureTenantID)
	ac.BaseClient.Authorizer = authorizer
	_, err := ac.Delete(ctx, *app.ID)
	return err
}

// ServicePrincipalCreate creates an azure identity for this application
// that can be assigned access policies on azure cloud resources
func ServicePrincipalCreate(ctx context.Context, authorizer auth.Authorizer, app *msgraph.Application) (*msgraph.ServicePrincipal, error) {
	c := msgraph.NewServicePrincipalsClient(AzureTenantID)
	c.BaseClient.Authorizer = authorizer
	sp, _, err := c.Create(ctx, msgraph.ServicePrincipal{AppId: app.AppId})
	return sp, err
}

func ServicePrincipalDelete(ctx context.Context, authorizer auth.Authorizer, sp *msgraph.ServicePrincipal) error {
	c := msgraph.NewServicePrincipalsClient(AzureTenantID)
	c.BaseClient.Authorizer = authorizer
	_, err := c.Delete(ctx, *sp.ID)
	return err
}

// ServicePrincipalPasswordCreate creates an long lived password credential for this application
func ServicePrincipalPasswordCreate(ctx context.Context, authorizer auth.Authorizer, sp *msgraph.ServicePrincipal) (*msgraph.PasswordCredential, error) {
	c := msgraph.NewServicePrincipalsClient(AzureTenantID)
	c.BaseClient.Authorizer = authorizer
	newCredential, _, err := c.AddPassword(ctx, *sp.ID, msgraph.PasswordCredential{
		KeyId: sp.AppId,
	})
	return newCredential, err
}

func ServicePrincipalPasswordDelete(ctx context.Context, authorizer auth.Authorizer, sp *msgraph.ServicePrincipal) error {
	c := msgraph.NewServicePrincipalsClient(AzureTenantID)
	c.BaseClient.Authorizer = authorizer
	_, err := c.RemovePassword(ctx, *sp.ID, *sp.AppId)
	return err
}

type Policy struct {
	Scope              string `json:"scope"`
	RoleDefinitionName string `json:"roleDefinitionName`
}

// AddAccessPoliciesToServicePrincipal adds the access policies
// from the connection to this azure service principal to give
// it accesss to the azure cloud resources the app needs to connect to
func AddAccessPoliciesToServicePrincipal(ctx context.Context, sp *msgraph.ServicePrincipal, policy Policy) error {
	authorizer, err := azauth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}
	roleAssClient := authorization.NewRoleAssignmentsClient(AzureSubscriptionID)
	roleAssClient.Authorizer = authorizer
	roleDefClient := authorization.NewRoleDefinitionsClient(AzureSubscriptionID)
	roleDefClient.Authorizer = authorizer
	roleDefinitions, err := roleDefClient.List(ctx, policy.Scope, fmt.Sprintf("roleName eq '%s'", policy.RoleDefinitionName))
	if err != nil {
		return fmt.Errorf("loading Role Definition List: %+v", err)
	}
	if len(roleDefinitions.Values()) != 1 {
		return fmt.Errorf("loading Role Definition List: could not find role '%s'", policy.RoleDefinitionName)
	}
	roleDefinitionId := *roleDefinitions.Values()[0].ID
	properties := authorization.RoleAssignmentCreateParameters{
		RoleAssignmentProperties: &authorization.RoleAssignmentProperties{
			RoleDefinitionID: &roleDefinitionId,
			PrincipalID:      sp.ID,
			PrincipalType:   authorization.ServicePrincipal,
		},
	}
	ra, createErr := roleAssClient.Create(
		ctx,
		policy.Scope,
		uuid.NewString(), // this is a GUID for the role assignment to ensure uniqueness we can probably be more careful about storing this id in state in the provider
		properties)
	if createErr != nil {
		return createErr
	}
	logObject("role assignment", ra)
	return nil
}

// This uses the "official" azure-sdk-for-go to create a new role definition but was giving errors
// "Assignments to roles with DataActions and NotDataActions are not supported on API version '2015-07-01'. The minimum required API version for this operations is '2018-01-01-preview'
// AddAccessPoliciesToServicePrincipal adds the access policies
// from the connection to this azure service principal to give
// it accesss to the azure cloud resources the app needs to connect to
// func AddAccessPoliciesToServicePrincipal(ctx context.Context, azCreds *azidentity.DefaultAzureCredential, sp *msgraph.ServicePrincipal, policy Policy) error {
// 	roleAssClient, err := armauthorization.NewRoleAssignmentsClient(AzureSubscriptionID, azCreds, nil)
// 	if err != nil {
// 		return err
// 	}
// 	roleDefClient, err := armauthorization.NewRoleDefinitionsClient(azCreds, nil)
// 	if err != nil {
// 		return err
// 	}
// 	pager := roleDefClient.NewListPager(policy.Scope,
// 		&armauthorization.RoleDefinitionsClientListOptions{Filter: stringPtr(fmt.Sprintf("roleName eq '%s'", policy.RoleDefinitionName))})
// 	roleDefs := []*armauthorization.RoleDefinition{}
// 	for pager.More() {
// 		nextResult, err := pager.NextPage(ctx)
// 		if err != nil {
// 			log.Fatalf("failed to advance page: %v", err)
// 		}
// 		roleDefs = append(roleDefs, nextResult.Value...)
// 	}
// 	if len(roleDefs) != 1 {
// 		return fmt.Errorf("loading Role Definition List: could not find role '%s'", policy.RoleDefinitionName)
// 	}
// 	roleDefinitionId := *roleDefs[0].ID
// 	ra, createErr := roleAssClient.Create(ctx,
// 		policy.Scope,
// 		uuid.New().String(), // this is a UUID for the role assignment must be unique should probably be more careful to use this from state in provider
// 		armauthorization.RoleAssignmentCreateParameters{
// 			Properties: &armauthorization.RoleAssignmentProperties{
// 				PrincipalID:      sp.ID,
// 				RoleDefinitionID: &roleDefinitionId,
// 			},
// 		},

// 	)

// 	if createErr != nil {
// 		return createErr
// 	}
// 	logObject("role assignment", ra)
// 	return nil
// }

// TODO pivot to this in the future not yet really stable on azure side for now going to use long lived credential.
// // FederatedIdentityCredentialCreate creates the trust relationship
// // between the KSA used to run the app and the service principal
// // which has access to the azure cloud resources
// func FederatedIdentityCredentialCreate(ctx context.Context, azCreds *azidentity.DefaultAzureCredential, principal interface{}) error {
//   // TODO
//   return nil
// }

func logObject(name string, obj interface{}) {
	b, _ := json.MarshalIndent(obj, "", "  ")
	log.Default().Printf("successfully created %s: %v\n", name, string(b))
}

func main() {
	// Check for subscription id info this would come from service principal artifact
	if len(AzureSubscriptionID) == 0 {
		log.Fatalf("AZURE_SUBSCRIPTION_ID is not set")
	}

	log.Default().Printf("createing azure applicaiton identity resources in subscription id %v\n", AzureSubscriptionID)

	// Create default credentials from environment variables
	appName := "foo"
	ctx := context.Background()

	environment := environments.Global

	authConfig := &auth.Config{
		Environment:            environment,
		TenantID:               AzureTenantID,
		ClientID:               AzureClientID,
		ClientSecret:           AzureClientSecret,
		EnableClientSecretAuth: true,
	}

	// b, _ := json.MarshalIndent(authConfig, "", "  ")
	// log.Default().Printf("authConfig: %v\n", string(b))

	authorizer, err := authConfig.NewAuthorizer(ctx, environment.MsGraph)
	if err != nil {
		log.Fatalf("=== %v", err)
	}
	if err != nil {
		log.Fatalf("failed to obtain an azurecredential: %v", err)
	}

	// All the above will be replaced with wiring into mx provider

	app, err := ApplicationCreate(ctx, authorizer, appName)
	if err != nil {
		log.Fatalf("failed to create application: %v", err)
	}
	logObject("application", app)

	sp, err := ServicePrincipalCreate(ctx, authorizer, app)
	if err != nil {
		log.Fatalf("failed to create service principal: %v", err)
	}
	logObject("service principal", sp)
	spPass, err := ServicePrincipalPasswordCreate(ctx, authorizer, sp)
	log.Default().Printf("service principal password: %#v\n", spPass)
	if err != nil {
		log.Fatalf("failed to create service principal password: %v", err)
	}
	logObject("service principal password", spPass)

	// these environment variables can be set in the pods that need access to cloud services via this service principal
	// azure has not yet implemented stable workload identity so we are going to use the long lived credential for now
	azCredEnv := map[string]string{
		"AZURE_CLIENT_ID":     *sp.AppId,
		"AZURE_CLIENT_SECRET": *spPass.SecretText,
		"AZURE_TENANT_ID":     AzureTenantID,
	}
	log.Default().Printf("success! created azure applicaiton identity resources you can use them with this env: %#v", azCredEnv)

	policies := []Policy{
		{
			Scope:              fmt.Sprintf("/subscriptions/%s", AzureSubscriptionID),
			RoleDefinitionName: "Storage Blob Data Contributor",
		},
	}
	for _, p := range policies {
		if err := AddAccessPoliciesToServicePrincipal(ctx, sp, p); err != nil {
			log.Fatalf("failed to add access policy: %v", err)
		}
	}

}

func stringPtr(s string) *string {
	return &s
}
