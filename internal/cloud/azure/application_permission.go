package azure

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization"
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/go-uuid"
	"github.com/manicminer/hamilton/environments"
	"gopkg.in/retry.v1"
)

type ApplicationPermissionConfig struct {
	ID                 string
	RoleName           string
	ServicePrincipalID string
	Scope              string
}

func getAzureResourceManagerAuthorizer(ctx context.Context, c *AzureConfig) (autorest.Authorizer, error) {
	builder := authentication.Builder{
		SubscriptionID:           c.Provider.SubscriptionID.Value,
		ClientID:                 c.Provider.ClientID.Value,
		ClientSecret:             c.Provider.ClientSecret.Value,
		TenantID:                 c.Provider.TenantID.Value,
		Environment:              "public",
		SupportsClientSecretAuth: true,
		UseMicrosoftGraph:        true,
	}
	config, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("building AzureRM Client: %s", err)
	}

	sender := sender.BuildSender("AzureRM")

	env, err := authentication.AzureEnvironmentByNameFromEndpoint(ctx, config.MetadataHost, config.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to find environment %q from endpoint %q: %+v", config.Environment, config.MetadataHost, err)
	}

	environment, err := environments.EnvironmentFromString(config.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to find environment %q from endpoint %q: %+v", config.Environment, config.MetadataHost, err)
	}

	oathConfig, err := config.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("building OAuth Config: %+v", err)
	}

	authorizer, err := config.GetMSALToken(ctx, environments.ResourceManagerPublic, sender, oathConfig, string(environment.ResourceManager.Endpoint))
	if err != nil {
		return nil, fmt.Errorf("unable to get MSAL authorization token for resource manager API: %+v", err)
	}
	return authorizer, nil
}

type RoleAssignmentsClient interface {
	Create(ctx context.Context, scope string, roleAssignmentName string, parameters authorization.RoleAssignmentCreateParameters) (result authorization.RoleAssignment, err error)
	GetByID(ctx context.Context, roleID string, tenantID string) (result authorization.RoleAssignment, err error)
	Delete(ctx context.Context, scope string, roleAssignmentName string, tenantID string) (result authorization.RoleAssignment, err error)
}

func (c *AzureConfig) NewRoleAssignmentsClient(ctx context.Context) (RoleAssignmentsClient, error) {
	raClient := authorization.NewRoleAssignmentsClient(c.Provider.SubscriptionID.Value)

	authorizer, err := getAzureResourceManagerAuthorizer(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("error creating RoleDefinitionsClient: %+v", err)
	}

	raClient.Authorizer = authorizer
	return raClient, nil
}

type RoleDefinitionsClient interface {
	GetByID(ctx context.Context, roleID string) (result authorization.RoleDefinition, err error)
	List(ctx context.Context, scope string, filter string) (result authorization.RoleDefinitionListResultPage, err error)
}

func (c *AzureConfig) NewRoleDefinitionsClient(ctx context.Context) (RoleDefinitionsClient, error) {
	rdClient := authorization.NewRoleDefinitionsClient(c.Provider.SubscriptionID.Value)

	authorizer, err := getAzureResourceManagerAuthorizer(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("error creating RoleDefinitionsClient: %+v", err)
	}

	rdClient.Authorizer = authorizer
	return rdClient, nil
}

func CreateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, raClient RoleAssignmentsClient, rdClient RoleDefinitionsClient) error {

	roleDefinitions, err := rdClient.List(ctx, config.Scope, fmt.Sprintf("roleName eq '%s'", config.RoleName))
	if err != nil {
		return fmt.Errorf("loading Role Definition List: %+v", err)
	}
	if len(roleDefinitions.Values()) != 1 {
		return fmt.Errorf("loading Role Definition List: could not find role '%s'", config.RoleName)
	}
	if roleDefinitions.Values()[0].ID == nil {
		return fmt.Errorf("loading Role Definition List: values[0].ID is nil '%s'", config.RoleName)
	}

	defId := *roleDefinitions.Values()[0].ID
	role, err := rdClient.GetByID(ctx, defId)
	if err != nil {
		return fmt.Errorf("getting Role Definition by ID %s: %+v", defId, err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
	}

	parameters := authorization.RoleAssignmentCreateParameters{
		RoleAssignmentProperties: &authorization.RoleAssignmentProperties{
			RoleDefinitionID: role.ID,
			PrincipalID:      &config.ServicePrincipalID,
		},
	}

	resp, createErr := createRoleAssignment(ctx, config.Scope, uuid, parameters, raClient)
	if createErr != nil {
		return fmt.Errorf("error creating role assignment. Response: %+v Error: %+v", resp.Body, createErr)
	}

	config.ID = *resp.ID

	return nil
}

func ReadApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, raClient RoleAssignmentsClient, rdClient RoleDefinitionsClient) error {
	return nil
}

func UpdateApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, raClient RoleAssignmentsClient, rdClient RoleDefinitionsClient) error {
	return nil
}

func DeleteApplicationPermission(ctx context.Context, config *ApplicationPermissionConfig, raClient RoleAssignmentsClient, rdClient RoleDefinitionsClient) error {

	id, err := parseRoleAssignmentId(config.ID)
	if err != nil {
		return err
	}

	_, err = raClient.Delete(ctx, id.scope, id.name, "") //id.tenantId)
	if err != nil {
		return fmt.Errorf("deletion of Role Assignment %q returned an error: %w", config.ID, err)
	}

	return nil
}

type roleAssignmentId struct {
	scope    string
	name     string
	tenantId string
}

func parseRoleAssignmentId(input string) (*roleAssignmentId, error) {
	segments := strings.Split(input, "/providers/Microsoft.Authorization/roleAssignments/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected Role Assignment ID to be in the format `{scope}/providers/Microsoft.Authorization/roleAssignments/{name}` but got %q", input)
	}

	// /{scope}/providers/Microsoft.Authorization/roleAssignments/{roleAssignmentName}
	// Tenant ID only required when going cross-tenant
	id := roleAssignmentId{
		scope:    strings.TrimPrefix(segments[0], "/"),
		name:     segments[1],
		tenantId: "",
	}
	return &id, nil
}

func createRoleAssignment(ctx context.Context, scope string, roleAssignmentName string, parameters authorization.RoleAssignmentCreateParameters, raClient RoleAssignmentsClient) (authorization.RoleAssignment, error) {
	attempts := retry.Regular{
		Total: 4 * time.Minute,
		Delay: 5 * time.Second,
	}
	var resp authorization.RoleAssignment
	var createErr error
	for attempt := attempts.Start(nil); attempt.Next(); {
		resp, createErr = raClient.Create(ctx, scope, roleAssignmentName, parameters)
		if createErr != nil {
			if responseErrorIsRetryable(createErr) {
				log.Printf("[debug] Experienced retryable error creating service principal. Trying again. error: %s", createErr.Error())
				continue
			} else if responseWasStatusCode(resp.Response, 400) && strings.Contains(createErr.Error(), "PrincipalNotFound") {
				log.Printf("[debug] Service principal not found. Could still be registering. Trying again. error: %s", createErr.Error())
				continue
			}
		}
		if resp.ID == nil {
			return resp, fmt.Errorf("creation of Role Assignment %q did not return an id value", roleAssignmentName)
		}
		return resp, createErr
	}
	return resp, createErr
}
