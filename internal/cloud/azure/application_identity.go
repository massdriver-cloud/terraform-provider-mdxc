package azure

import (
	"context"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type ApplicationAPI interface {
	Create(ctx context.Context, application msgraph.Application) (*msgraph.Application, int, error)
	Get(ctx context.Context, id string, query odata.Query) (*msgraph.Application, int, error)
	Update(ctx context.Context, application msgraph.Application) (int, error)
	Delete(ctx context.Context, id string) (int, error)
}

func (c *AzureConfig) NewApplicationService(ctx context.Context) (ApplicationAPI, error) {
	authConfig := auth.Config{
		Environment:            environments.Global,
		TenantID:               c.provider.TenantID.Value,
		ClientID:               c.provider.ClientID.Value,
		ClientSecret:           c.provider.ClientSecret.Value,
		EnableClientSecretAuth: true,
	}
	authorizer, authorizerErr := authConfig.NewAuthorizer(ctx, authConfig.Environment.MsGraph)
	if authorizerErr != nil {
		return nil, authorizerErr
	}
	appAPI := msgraph.NewApplicationsClient(c.provider.TenantID.Value)
	appAPI.BaseClient.Authorizer = authorizer
	return appAPI, nil
}

type ApplicationIdentityConfig struct {
	ID   string
	Name string
}

func CreateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, applicationAPI ApplicationAPI) error {
	application, _, err := applicationAPI.Create(ctx, msgraph.Application{
		DisplayName: &config.Name,
	})
	if err != nil {
		return err
	}

	config.ID = *application.ID

	return err
}

func ReadApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, applicationAPI ApplicationAPI) error {
	return nil
}

func UpdateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, applicationAPI ApplicationAPI) error {
	return nil
}

func DeleteApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, applicationAPI ApplicationAPI) error {
	_, err := applicationAPI.Delete(ctx, config.ID)
	return err
}
