package azure

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

type AzureProviderConfig struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	TenantID       types.String `tfsdk:"tenant_id"`
}

type AzureConfig struct {
	provider   *AzureProviderConfig
	authConfig *auth.Config
}

func Initialize(ctx context.Context, providerConfig *AzureProviderConfig) (*AzureConfig, error) {
	azureConfig := AzureConfig{}

	azureConfig.provider = providerConfig

	authConfig := auth.Config{
		Environment:            environments.Global,
		TenantID:               providerConfig.TenantID.Value,
		ClientID:               providerConfig.ClientID.Value,
		ClientSecret:           providerConfig.ClientSecret.Value,
		EnableClientSecretAuth: true,
	}

	azureConfig.authConfig = &authConfig

	log.Printf("[debug] Azure Config Created")
	return &azureConfig, nil
}
