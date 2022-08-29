package azure

import (
	"context"
	"fmt"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type ApplicationClient interface {
	Create(ctx context.Context, application msgraph.Application) (*msgraph.Application, int, error)
	Get(ctx context.Context, id string, query odata.Query) (*msgraph.Application, int, error)
	Update(ctx context.Context, application msgraph.Application) (int, error)
	Delete(ctx context.Context, id string) (int, error)
}

func (c *AzureConfig) NewApplicationClient(ctx context.Context) (ApplicationClient, error) {
	authorizer, authorizerErr := c.authConfig.NewAuthorizer(ctx, c.authConfig.Environment.MsGraph)
	if authorizerErr != nil {
		return nil, authorizerErr
	}
	appClient := msgraph.NewApplicationsClient(c.Provider.TenantID.Value)
	appClient.BaseClient.Authorizer = authorizer
	return appClient, nil
}

type ServicePrincipalsClient interface {
	Create(ctx context.Context, servicePrincipal msgraph.ServicePrincipal) (*msgraph.ServicePrincipal, int, error)
	Get(ctx context.Context, id string, query odata.Query) (*msgraph.ServicePrincipal, int, error)
	Update(ctx context.Context, servicePrincipal msgraph.ServicePrincipal) (int, error)
	Delete(ctx context.Context, id string) (int, error)

	AddPassword(ctx context.Context, servicePrincipalId string, passwordCredential msgraph.PasswordCredential) (*msgraph.PasswordCredential, int, error)
	RemovePassword(ctx context.Context, servicePrincipalId string, keyId string) (int, error)
}

func (c *AzureConfig) NewServicePrincipalsClient(ctx context.Context) (ServicePrincipalsClient, error) {
	authorizer, authorizerErr := c.authConfig.NewAuthorizer(ctx, c.authConfig.Environment.MsGraph)
	if authorizerErr != nil {
		return nil, authorizerErr
	}
	servicePrincipalsAPI := msgraph.NewServicePrincipalsClient(c.Provider.TenantID.Value)
	servicePrincipalsAPI.BaseClient.Authorizer = authorizer
	return servicePrincipalsAPI, nil
}

type ApplicationIdentityConfig struct {
	ID                     string
	Name                   string
	ApplicationID          string
	ServicePrincipalID     string
	ServicePrincipalSecret string
}

func CreateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, appClient ApplicationClient, spClient ServicePrincipalsClient) error {
	// The ID is the service principal ID, so technically theres a chance that last time the App was made
	// but the SP was not. We need to check if the Application was made and skip creation if it was.
	if config.ApplicationID == "" {
		app, _, appErr := appClient.Create(ctx, msgraph.Application{
			DisplayName: &config.Name,
		})
		if appErr != nil {
			return appErr
		}
		config.ApplicationID = *app.ID
	}

	// fetch the application to make sure it exists
	app, _, err := appClient.Get(ctx, config.ApplicationID, odata.Query{})
	if err != nil {
		return fmt.Errorf("error retrieving Application with object ID %v: %w", config.ApplicationID, err)
	}

	sp, _, spErr := spClient.Create(ctx, msgraph.ServicePrincipal{
		AppId: app.AppId,
	})
	if spErr != nil {
		return spErr
	}
	config.ServicePrincipalID = *sp.ID
	config.ID = *sp.ID

	// fetch the service principal to make sure it exists
	_, _, spCheckErr := spClient.Get(ctx, config.ServicePrincipalID, odata.Query{})
	if spCheckErr != nil {
		return fmt.Errorf("error retrieving Service Principal with ID %v: %w", config.ServicePrincipalID, spCheckErr)
	}

	// Technically the secret is only needed for Kubernetes until the support workload identity. Maybe we make this on a conditional?
	spSecret, _, secretErr := spClient.AddPassword(ctx, *sp.ID, msgraph.PasswordCredential{})
	if secretErr != nil {
		return secretErr
	}
	config.ServicePrincipalSecret = *spSecret.SecretText

	return nil
}

func ReadApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, appClient ApplicationClient, spClient ServicePrincipalsClient) error {
	return nil
}

func UpdateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, appClient ApplicationClient, spClient ServicePrincipalsClient) error {
	return nil
}

func DeleteApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, appClient ApplicationClient, spClient ServicePrincipalsClient) error {
	_, appErr := appClient.Delete(ctx, config.ApplicationID)
	if appErr != nil {
		return appErr
	}

	config.ID = ""
	config.ApplicationID = ""
	config.ServicePrincipalID = ""
	config.ServicePrincipalSecret = ""

	return nil
}
