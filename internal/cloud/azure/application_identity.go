package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
)

type ManagedIdentityClient interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, parameters armmsi.Identity, options *armmsi.UserAssignedIdentitiesClientCreateOrUpdateOptions) (armmsi.UserAssignedIdentitiesClientCreateOrUpdateResponse, error)
	Get(ctx context.Context, resourceGroupName string, resourceName string, options *armmsi.UserAssignedIdentitiesClientGetOptions) (armmsi.UserAssignedIdentitiesClientGetResponse, error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string, options *armmsi.UserAssignedIdentitiesClientDeleteOptions) (armmsi.UserAssignedIdentitiesClientDeleteResponse, error)
}

type FederatedIdentityCredentialClient interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, federatedIdentityCredentialResourceName string, parameters armmsi.FederatedIdentityCredential, options *armmsi.FederatedIdentityCredentialsClientCreateOrUpdateOptions) (armmsi.FederatedIdentityCredentialsClientCreateOrUpdateResponse, error)
}

type ApplicationIdentityConfig struct {
	// READ-ONLY; Fully qualified resource ID for the resource.
	// /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID                 string
	Name               string
	KubernetesNamspace string
}

func newManagedIdentityClientFactory(ctx context.Context, config *AzureProviderConfig) (ManagedIdentityClient, error) {
	cred, err := azidentity.NewClientSecretCredential(config.TenantID.Value, config.ClientID.Value, config.ClientSecret.Value, nil)
	if err != nil {
		return nil, err
	}

	client, errClient := armmsi.NewUserAssignedIdentitiesClient(config.SubscriptionID.Value, cred, nil)
	if errClient != nil {
		return nil, errClient
	}

	return client, nil
}

func newFederatedIdentityCredentialClientFactory(ctx context.Context, config *AzureProviderConfig) (FederatedIdentityCredentialClient, error) {
	cred, err := azidentity.NewClientSecretCredential(config.TenantID.Value, config.ClientID.Value, config.ClientSecret.Value, nil)
	if err != nil {
		return nil, err
	}

	client, errClient := armmsi.NewFederatedIdentityCredentialsClient(config.SubscriptionID.Value, cred, nil)
	if errClient != nil {
		return nil, errClient
	}

	return client, nil
}

func CreateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client ManagedIdentityClient, fedClient FederatedIdentityCredentialClient) error {
	res, _ := client.CreateOrUpdate(ctx,
		"rgName",
		config.Name,
		armmsi.Identity{
			Location: to.Ptr("eastus"),
			// Tags: map[string]*string{
			// 	"key1": to.Ptr("value1"),
			// 	"key2": to.Ptr("value2"),
			// },
		},
		nil,
	)

	config.ID = *res.ID

	if config.KubernetesNamspace != "" {
		if errAddRole := addWorkloadIdentityRole(ctx, config, fedClient); errAddRole != nil {
			return errAddRole
		}
	}
	return nil
}

func ReadApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client ManagedIdentityClient, fedClient FederatedIdentityCredentialClient) error {
	return nil
}

func UpdateApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client ManagedIdentityClient, fedClient FederatedIdentityCredentialClient) error {
	return nil
}

func DeleteApplicationIdentity(ctx context.Context, config *ApplicationIdentityConfig, client ManagedIdentityClient, fedClient FederatedIdentityCredentialClient) error {
	// We can use the name or get the name from the ID, which is the full resource ID
	_, err := client.Delete(ctx, "resource-group", config.Name, nil)
	if err != nil {
		return err
	}

	config.ID = ""

	return nil
}

func addWorkloadIdentityRole(ctx context.Context, config *ApplicationIdentityConfig, client FederatedIdentityCredentialClient) error {
	res, err := client.CreateOrUpdate(ctx,
		"rgName",
		"resourceName",
		"ficResourceName",
		armmsi.FederatedIdentityCredential{
			Properties: &armmsi.FederatedIdentityCredentialProperties{
				Audiences: []*string{
					to.Ptr("api://AzureADTokenExchange")},
				// TODO: need the issuer url from the k8s cluster
				Issuer: to.Ptr("https://oidc.prod-aks.azure.com/IssuerGUID"),
				// k8s service account
				Subject: to.Ptr("system:serviceaccount:ns:svcaccount"),
			},
		},
		nil)
	if err != nil {
		return err
	}
	config.KubernetesNamspace = *res.Name
	return nil
}
