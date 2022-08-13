package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type AzureConfig struct {
	subscriptionId string
	clientId       string
	clientSecret   string
	tenantId       string
	credentials    *azidentity.ClientSecretCredential
}

func Initialize(ctx context.Context, azureMap map[string]interface{}) (*AzureConfig, error) {
	azureConfig := AzureConfig{}

	if subscriptionId, ok := azureMap["subscription_id"].(string); ok && subscriptionId != "" {
		azureConfig.subscriptionId = subscriptionId
	}
	if clientId, ok := azureMap["client_id"].(string); ok && clientId != "" {
		azureConfig.clientId = clientId
	}
	if clientSecret, ok := azureMap["client_secret"].(string); ok && clientSecret != "" {
		azureConfig.clientSecret = clientSecret
	}
	if tenantId, ok := azureMap["tenant_id"].(string); ok && tenantId != "" {
		azureConfig.tenantId = tenantId
	}

	var clientErr error
	azureConfig.credentials, clientErr = azidentity.NewClientSecretCredential(
		azureConfig.tenantId,
		azureConfig.clientId,
		azureConfig.clientSecret,
		nil,
	)
	if clientErr != nil {
		return nil, clientErr
	}

	log.Printf("[debug] Azure Config Created")
	return &azureConfig, nil
}
