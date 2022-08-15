package azure

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AzureProviderConfig struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	TenantID       types.String `tfsdk:"tenant_id"`
}

type AzureConfig struct {
	provider *AzureProviderConfig
}

func Initialize(ctx context.Context, providerConfig *AzureProviderConfig) (*AzureConfig, error) {
	azureConfig := AzureConfig{}

	azureConfig.provider = providerConfig

	log.Printf("[debug] Azure Config Created")
	return &azureConfig, nil
}
