package azure

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type ApplicationPermissionConfig struct {
	ID                 string
	RoleName           string
	ServicePrincipalID string
	Scope              string
}

// type RoleAssignmentsClient interface {
// 	Create(ctx context.Context, roleAssignment msgraph.UnifiedRoleAssignment) (*msgraph.UnifiedRoleAssignment, int, error)
// 	Get(ctx context.Context, id string, query odata.Query) (*msgraph.UnifiedRoleAssignment, int, error)
// 	Delete(ctx context.Context, id string) (int, error)
// }

// func (c *AzureConfig) NewRoleAssignmentsClient(ctx context.Context) (RoleAssignmentsClient, error) {
// 	authorizer, authorizerErr := c.authConfig.NewAuthorizer(ctx, c.authConfig.Environment.MsGraph)
// 	if authorizerErr != nil {
// 		return nil, authorizerErr
// 	}
// 	raClient := msgraph.NewRoleAssignmentsClient(c.Provider.TenantID.Value)
// 	raClient.BaseClient.Authorizer = authorizer
// 	return raClient, nil
// }

// type RoleDefinitionsClient interface {
// 	Get(ctx context.Context, id string, query odata.Query) (*msgraph.UnifiedRoleDefinition, int, error)
// 	List(ctx context.Context, query odata.Query) (*[]msgraph.UnifiedRoleDefinition, int, error)
// }

// func (c *AzureConfig) NewRoleDefinitionsClient(ctx context.Context) (RoleDefinitionsClient, error) {
// 	authorizer, authorizerErr := c.authConfig.NewAuthorizer(ctx, c.authConfig.Environment.ResourceManager)
// 	if authorizerErr != nil {
// 		return nil, authorizerErr
// 	}
// 	rdClient := msgraph.NewRoleDefinitionsClient(c.Provider.TenantID.Value)
// 	rdClient.BaseClient.Authorizer = authorizer
// 	return rdClient, nil
// }

type RoleAssignmentsClient interface {
	Create(ctx context.Context, roleAssignment msgraph.UnifiedRoleAssignment) (*msgraph.UnifiedRoleAssignment, int, error)
	Get(ctx context.Context, id string, query odata.Query) (*msgraph.UnifiedRoleAssignment, int, error)
	Delete(ctx context.Context, id string) (int, error)
}

func (c *AzureConfig) NewRoleAssignmentsClient(ctx context.Context) (RoleAssignmentsClient, error) {
	authorizer, authorizerErr := c.authConfig.NewAuthorizer(ctx, c.authConfig.Environment.MsGraph)
	if authorizerErr != nil {
		return nil, authorizerErr
	}
	raClient := msgraph.NewRoleAssignmentsClient(c.Provider.TenantID.Value)
	raClient.BaseClient.Authorizer = authorizer
	return raClient, nil
}

type RoleDefinitionsClient interface {
	Get(ctx context.Context, scope string, roleDefinitionID string, options *armauthorization.RoleDefinitionsClientGetOptions) (armauthorization.RoleDefinitionsClientGetResponse, error)
	NewListPager(scope string, options *armauthorization.RoleDefinitionsClientListOptions) *runtime.Pager[armauthorization.RoleDefinitionsClientListResponse]
	//List(ctx context.Context, query odata.Query) (*[]msgraph.UnifiedRoleDefinition, int, error)
}

func (c *AzureConfig) NewRoleDefinitionsClient(ctx context.Context) (RoleDefinitionsClient, error) {
	creds, credsErr := azidentity.NewClientSecretCredential(
		c.Provider.TenantID.Value,
		c.Provider.ClientID.Value,
		c.Provider.ClientSecret.Value,
		&azidentity.ClientSecretCredentialOptions{})
	if credsErr != nil {
		return nil, credsErr
	}
	//rdClient := authorization.NewRoleDefinitionsClient()//NewRoleDefinitionsClient(c.Provider.SubscriptionID.Value)
	rdClient, clientErr := armauthorization.NewRoleDefinitionsClient(creds, &arm.ClientOptions{})
	if clientErr != nil {
		return nil, clientErr
	}
	// authorizer, authorizerErr := c.authConfig.NewAuthorizer(ctx, c.authConfig.Environment.ResourceManager)
	// if authorizerErr != nil {
	// 	return nil, authorizerErr
	// }
	//rdClient := msgraph.NewRoleDefinitionsClient(c.Provider.TenantID.Value)
	return rdClient, nil
}

func CreateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, raClient RoleAssignmentsClient, rdClient RoleDefinitionsClient) error {

	filter := fmt.Sprintf("displayName eq '%s'", config.RoleName)
	rdPager := rdClient.NewListPager(config.Scope,
		&armauthorization.RoleDefinitionsClientListOptions{
			Filter: &filter})
	//roleDefinitions, _, err := rdClient.List(ctx, odata.Query{ //})
	//Filter: fmt.Sprintf("displayName eq '%s'", config.RoleName),
	//Filter: fmt.Sprintf("roleName eq '%s'", config.RoleName),
	//})
	// if err != nil {
	// 	return fmt.Errorf("loading Role Definition List: %+v", err)
	// }
	for rdPager.More() {
		nextResult, nextErr := rdPager.NextPage(ctx)
		if nextErr != nil {
			return nextErr
		}
		for _, v := range nextResult.Value {
			log.Default().Printf("\nDisplayName: %v\t\t\tID: %v\n\n", *v.Properties.RoleName, *v.ID)
		}
	}
	return fmt.Errorf("lol")
	// for _, rd := range *roleDefinitions {
	// 	log.Default().Printf("\nDisplayName: %v\t\t\tID: %v\n\n", *rd.DisplayName, *rd.ID)
	// }
	// if len(*roleDefinitions) != 1 {
	// 	return fmt.Errorf("loading Role Definition List: could not find role '%s'", config.RoleName)
	// }
	// roleDefinitionId := (*roleDefinitions)[0].ID
	// properties := authorization.RoleAssignmentCreateParameters{
	// 	RoleAssignmentProperties: &authorization.RoleAssignmentProperties{
	// 		RoleDefinitionID: &roleDefinitionId,
	// 		PrincipalID:      sp.ID,
	// 		PrincipalType:    authorization.ServicePrincipal,
	// 	},
	// }
	// ra, _, createErr := raClient.Create(ctx, msgraph.UnifiedRoleAssignment{
	// 	DirectoryScopeId: &config.Scope,
	// 	PrincipalId:      &config.ServicePrincipalID,
	// 	RoleDefinitionId: roleDefinitionId,
	// })
	// ctx,
	// policy.Scope,
	// uuid.NewString(), // this is a GUID for the role assignment to ensure uniqueness we can probably be more careful about storing this id in state in the provider
	// properties)
	// if createErr != nil {
	// 	return createErr
	// }

	// config.ID = *ra.ID

}

func ReadApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, raClient RoleAssignmentsClient, rdClient RoleDefinitionsClient) error {
	return nil
}

func UpdateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, raClient RoleAssignmentsClient, rdClient RoleDefinitionsClient) error {
	return nil
}

func DeleteApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, raClient RoleAssignmentsClient, rdClient RoleDefinitionsClient) error {

	// input := iam.DetachRolePolicyInput{
	// 	RoleName:  &config.RoleName,
	// 	PolicyArn: &config.PolicyARN,
	// }

	// _, deleteErr := client.DetachRolePolicy(ctx, &input)
	// if deleteErr != nil {
	// 	return deleteErr
	// }

	return nil
}
