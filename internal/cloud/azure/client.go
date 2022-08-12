package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AzureConfig struct {
	subscriptionId string
	clientId       string
	clientSecret   string
	tenantId       string
	credentials    *azidentity.ClientSecretCredential
}

func Initialize(ctx context.Context, d *schema.ResourceData, azureMap map[string]interface{}) (*AzureConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
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

	// log.Printf("[debug] Converting Azure values to config")
	// builder := authentication.Builder{
	// 	SubscriptionID: azureConfig.subscriptionId,
	// 	ClientID:       azureConfig.clientId,
	// 	ClientSecret:   azureConfig.clientSecret,
	// 	TenantID:       azureConfig.tenantId,

	// 	Environment:                    "public",
	// 	MetadataHost:                   "",
	// 	SupportsOIDCAuth:               false,
	// 	SupportsManagedServiceIdentity: false,

	// 	// Feature Toggles
	// 	SupportsClientCertAuth:   true,
	// 	SupportsClientSecretAuth: true,
	// 	SupportsAzureCliToken:    true,
	// 	SupportsAuxiliaryTenants: false,

	// 	// Doc Links
	// 	ClientSecretDocsLink: "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_secret",

	// 	// Use MSAL
	// 	UseMicrosoftGraph: true,
	// }

	// config, err := builder.Build()
	// if err != nil {
	// 	return nil, diag.Errorf("building AzureRM Client: %s", err)
	// }
	// terraformVersion := d.TerraformVersion
	// if terraformVersion == "" {
	// 	// Terraform 0.12 introduced this field to the protocol
	// 	// We can therefore assume that if it's missing it's 0.10 or 0.11
	// 	terraformVersion = "0.11+compatible"
	// }
	// clientBuilder := clients.ClientBuilder{
	// 	AuthConfig:                  config,
	// 	SkipProviderRegistration:    false,
	// 	TerraformVersion:            terraformVersion,
	// 	DisableCorrelationRequestID: false,
	// 	DisableTerraformPartnerID:   false,
	// 	StorageUseAzureAD:           false,
	// }

	var clientErr error
	azureConfig.credentials, clientErr = azidentity.NewClientSecretCredential(
		azureConfig.tenantId,
		azureConfig.clientId,
		azureConfig.clientSecret,
		nil,
	)
	if clientErr != nil {
		return nil, diag.FromErr(clientErr)
	}

	log.Printf("[debug] Azure Config Created")
	return &azureConfig, diags
}
