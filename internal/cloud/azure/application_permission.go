package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type ApplicationPermissionConfig struct {
	ID                 string
	RoleName           string
	ServicePrincipalID string
	Scope              string
}

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
	raClient := msgraph.NewRoleAssignmentsClient(c.provider.TenantID.Value)
	raClient.BaseClient.Authorizer = authorizer
	return raClient, nil
}

type RoleDefinitionsClient interface {
	Get(ctx context.Context, id string, query odata.Query) (*msgraph.UnifiedRoleDefinition, int, error)
	List(ctx context.Context, query odata.Query) (*[]msgraph.UnifiedRoleDefinition, int, error)
}

func (c *AzureConfig) NewRoleDefinitionsClient(ctx context.Context) (RoleDefinitionsClient, error) {
	_, err := azidentity.NewClientSecretCredential(
		c.provider.TenantID.Value,
		c.provider.ClientID.Value,
		c.provider.ClientSecret.Value, nil,
	)
	if err != nil {
		return nil, err
	}
	authorizer, authorizerErr := c.authConfig.NewAuthorizer(ctx, c.authConfig.Environment.Storage)
	if authorizerErr != nil {
		return nil, authorizerErr
	}
	rdClient := msgraph.NewRoleDefinitionsClient(c.provider.TenantID.Value)
	rdClient.BaseClient.Authorizer = authorizer
	return rdClient, nil
}

func CreateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, raClient RoleAssignmentsClient, rdClient RoleDefinitionsClient) error {

	roleDefinitions, _, err := rdClient.List(ctx, odata.Query{ //})
		Filter: fmt.Sprintf("displayName eq '%s'", config.RoleName),
		//Filter: fmt.Sprintf("roleName eq '%s'", config.RoleName),
	})
	if err != nil {
		return fmt.Errorf("loading Role Definition List: %+v", err)
	}
	for _, rd := range *roleDefinitions {
		str := fmt.Sprintf("\nDisplayName: %v\nID: %v\n\n", *rd.DisplayName, *rd.ID)
		tflog.Debug(ctx, str)
	}
	// if len(*roleDefinitions) != 1 {
	// 	return fmt.Errorf("loading Role Definition List: could not find role '%s'", config.RoleName)
	// }
	roleDefinitionId := (*roleDefinitions)[0].ID
	// properties := authorization.RoleAssignmentCreateParameters{
	// 	RoleAssignmentProperties: &authorization.RoleAssignmentProperties{
	// 		RoleDefinitionID: &roleDefinitionId,
	// 		PrincipalID:      sp.ID,
	// 		PrincipalType:    authorization.ServicePrincipal,
	// 	},
	// }
	ra, _, createErr := raClient.Create(ctx, msgraph.UnifiedRoleAssignment{
		DirectoryScopeId: &config.Scope,
		PrincipalId:      &config.ServicePrincipalID,
		RoleDefinitionId: roleDefinitionId,
	})
	// ctx,
	// policy.Scope,
	// uuid.NewString(), // this is a GUID for the role assignment to ensure uniqueness we can probably be more careful about storing this id in state in the provider
	// properties)
	if createErr != nil {
		return createErr
	}

	config.ID = *ra.ID

	return nil
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
